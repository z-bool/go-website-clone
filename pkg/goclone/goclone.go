package goclone

import (
	"context"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"

	"github.com/google/uuid"
	"github.com/z-bool/go-website-clone/pkg/crawler"
	"github.com/z-bool/go-website-clone/pkg/file"
	"github.com/z-bool/go-website-clone/pkg/html"
	"github.com/z-bool/go-website-clone/pkg/parser"
)

// Config 配置结构体，用于替代命令行参数
type Config struct {
	// URLs 要克隆的网站URL列表
	URLs []string
	// Open 是否自动在默认浏览器中打开项目
	Open bool
	// UserAgent 自定义用户代理
	UserAgent string
	// ProxyString 代理连接字符串
	ProxyString string
	// Cookies 预设的cookie列表
	Cookies []string
	// ConfigID 配置ID，使用UUID标识，用作保存文件夹名称
	ConfigID string
	// MaxFolderSize 文件夹最大大小限制（字节）
	MaxFolderSize int64
}

// GetProxyString 实现CrawlConfig接口
func (c *Config) GetProxyString() string {
	return c.ProxyString
}

// GetUserAgent 实现CrawlConfig接口
func (c *Config) GetUserAgent() string {
	return c.UserAgent
}

// GetMaxFolderSize 实现CrawlConfig接口
func (c *Config) GetMaxFolderSize() int64 {
	return c.MaxFolderSize
}

// CloneResult 克隆结果
type CloneResult struct {
	// Success 是否成功
	Success bool
	// ProjectPaths 生成的项目路径列表
	ProjectPaths []string
	// FirstProject 第一个项目路径（用于服务器或打开）
	FirstProject string
	// Error 错误信息
	Error error
}

// Clone 克隆网站的主函数
func Clone(ctx context.Context, config *Config) *CloneResult {
	result := &CloneResult{
		ProjectPaths: make([]string, 0),
	}

	if len(config.URLs) == 0 {
		result.Error = fmt.Errorf("没有提供要克隆的URL")
		return result
	}

	// 如果ConfigID为空，生成新的UUID
	if config.ConfigID == "" {
		config.ConfigID = uuid.New().String()
	}

	fmt.Printf("开始克隆 %d 个URL，配置ID: %s\n", len(config.URLs), config.ConfigID)

	// 创建cookie jar
	jar, err := cookiejar.New(&cookiejar.Options{})
	if err != nil {
		result.Error = err
		return result
	}

	// 处理cookies
	if err := setupCookies(jar, config.Cookies, config.URLs); err != nil {
		result.Error = err
		return result
	}

	// 处理每个URL
	for i, u := range config.URLs {
		fmt.Printf("正在处理第 %d 个URL: %s\n", i+1, u)
		projectPath, err := cloneURL(ctx, u, jar, config)
		if err != nil {
			result.Error = fmt.Errorf("克隆 %q 失败: %w", u, err)
			return result
		}

		fmt.Printf("URL %s 克隆完成，项目路径: %s\n", u, projectPath)
		result.ProjectPaths = append(result.ProjectPaths, projectPath)
		if result.FirstProject == "" {
			result.FirstProject = projectPath
		}
	}

	fmt.Println("所有URL克隆完成")

	result.Success = true
	return result
}

// cloneURL 克隆单个URL
func cloneURL(ctx context.Context, targetURL string, jar *cookiejar.Jar, config *Config) (string, error) {
	isValid, isValidDomain := parser.ValidateURL(targetURL), parser.ValidateDomain(targetURL)
	if !isValid && !isValidDomain {
		return "", fmt.Errorf("URL %q 无效", targetURL)
	}

	finalURL := targetURL
	if isValidDomain {
		finalURL = parser.CreateURL(targetURL)
	}

	// 使用ConfigID作为项目文件夹名称
	projectPath := file.CreateProjectWithID(config.ConfigID)

	// 执行爬取，传递配置对象以便进行大小检查
	if err := crawler.CrawlWithConfig(ctx, finalURL, projectPath, jar, config); err != nil {
		return "", fmt.Errorf("爬取失败: %w", err)
	}

	// 重构HTML链接
	if err := html.LinkRestructure(projectPath); err != nil {
		return "", fmt.Errorf("重构HTML链接失败: %w", err)
	}

	return projectPath, nil
}

// setupCookies 设置cookies
func setupCookies(jar *cookiejar.Jar, cookies []string, urls []string) error {
	if len(cookies) == 0 {
		return nil
	}

	var cs []*http.Cookie
	cs = make([]*http.Cookie, 0, len(cookies))

	for _, c := range cookies {
		ff := strings.Fields(c)
		for _, f := range ff {
			var k, v string
			if i := strings.IndexByte(f, '='); i >= 0 {
				k, v = f[:i], strings.TrimRight(f[i+1:], ";")
			} else {
				return fmt.Errorf("cookie格式错误，缺少'=' %q", c)
			}
			cs = append(cs, &http.Cookie{Name: k, Value: v})
		}
	}

	for _, urlStr := range urls {
		u, err := url.Parse(urlStr)
		if err != nil {
			return fmt.Errorf("解析URL失败 %q: %w", urlStr, err)
		}
		jar.SetCookies(&url.URL{Scheme: u.Scheme, User: u.User, Host: u.Host}, cs)
	}

	return nil
}

// QuickClone 快速克隆单个网站的便捷函数
func QuickClone(ctx context.Context, url string) *CloneResult {
	config := &Config{
		URLs: []string{url},
	}
	return Clone(ctx, config)
}
