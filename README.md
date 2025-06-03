# go-website-clone

goclone-dev/goclone的优化版本，支持函数式API调用，具备智能文件夹管理和大小限制功能。

## 🚀 主要特性

- ✅ **函数式API**: 可直接在Go代码中调用，无需命令行
- ✅ **智能文件夹管理**: 使用UUID自动生成唯一项目文件夹
- ✅ **大小限制控制**: 支持文件夹大小限制，防止过度下载
- ✅ **多URL批量克隆**: 一次配置克隆多个网站
- ✅ **代理支持**: 支持HTTP和SOCKS5代理
- ✅ **Cookie管理**: 支持预设cookie进行认证
- ✅ **实时监控**: 提供详细的下载进度和大小监控
- ✅ **跨平台兼容**: 智能处理文件名，支持Windows/Linux/macOS

## 📦 安装

```bash
go get github.com/z-bool/go-website-clone
```

## 🏗️ 项目结构

```
go-website-clone/
├── pkg/
│   ├── goclone/          # 主要API包
│   │   └── goclone.go    # 核心函数和配置
│   ├── crawler/          # 智能爬虫模块
│   │   ├── collector.go  # 资源收集器
│   │   ├── crawler.go    # 爬虫控制器
│   │   └── extractor.go  # 文件提取器
│   ├── file/             # 文件管理模块
│   │   └── write.go      # 文件写入和大小管理
│   ├── html/             # HTML处理模块
│   ├── parser/           # URL解析模块
│   └── server/           # 本地服务器模块
├── example/
│   └── main.go           # 完整使用示例
└── go.mod                # Go模块文件
```

## ⚙️ 核心配置

### Config 配置结构体

```go
type Config struct {
    URLs          []string  // 要克隆的网站URL列表
    UserAgent     string    // 自定义用户代理
    ProxyString   string    // 代理连接字符串
    Cookies       []string  // 预设的cookie列表
    ConfigID      string    // 配置ID（UUID），用作文件夹名称
    MaxFolderSize int64     // 文件夹最大大小限制（字节）
}
```

### CloneResult 结果结构体

```go
type CloneResult struct {
    Success      bool      // 是否成功
    ProjectPaths []string  // 生成的项目路径列表
    FirstProject string    // 第一个项目路径
    Error        error     // 错误信息
}
```

## 🔧 使用方法

### 导入包

```go
import (
    "context"
    "github.com/z-bool/go-website-clone/pkg/goclone"
)
```

### 基础使用 - 快速克隆

```go
package main

import (
    "context"
    "fmt"
    "log"
    "github.com/z-bool/go-website-clone/pkg/goclone"
)

func main() {
    ctx := context.Background()
    
    // 快速克隆单个网站
    result := goclone.QuickClone(ctx, "https://example.com")
    if result.Error != nil {
        log.Printf("克隆失败: %v", result.Error)
    } else {
        fmt.Printf("成功克隆到: %s\n", result.FirstProject)
    }
}
```

### 完整配置使用

```go
package main

import (
    "context"
    "fmt"
    "log"
    "github.com/z-bool/go-website-clone/pkg/goclone"
)

func main() {
    ctx := context.Background()
    
    config := &goclone.Config{
        URLs: []string{
            "https://example.com",
            "https://httpbin.org",
        },
        UserAgent:     "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
        ProxyString:   "http://127.0.0.1:8080", // 可选：设置代理
        Cookies:       []string{"session=abc123", "user=test"}, // 可选：预设cookies
        ConfigID:      "", // 留空自动生成UUID，或手动指定
        MaxFolderSize: 100 * 1024 * 1024, // 100MB限制，设置为0表示无限制
    }

    result := goclone.Clone(ctx, config)
    if result.Error != nil {
        log.Printf("克隆失败: %v", result.Error)
    } else {
        fmt.Printf("成功克隆 %d 个网站:\n", len(result.ProjectPaths))
        for i, path := range result.ProjectPaths {
            fmt.Printf("  %d. %s\n", i+1, path)
        }
        fmt.Printf("配置ID: %s\n", config.ConfigID)
    }
}
```

## 🎯 高级功能

### 1. 智能文件夹管理

使用UUID作为文件夹名称，确保每次克隆都有唯一标识：

```go
config := &goclone.Config{
    URLs:     []string{"https://example.com"},
    ConfigID: "", // 自动生成UUID如：a1b2c3d4-e5f6-7890-1234-567890abcdef
}

result := goclone.Clone(ctx, config)
fmt.Printf("项目保存在: %s\n", result.FirstProject)
fmt.Printf("使用的ConfigID: %s\n", config.ConfigID)
```

### 2. 文件夹大小限制

防止下载过大文件，保护系统资源：

```go
config := &goclone.Config{
    URLs:          []string{"https://example.com"},
    MaxFolderSize: 50 * 1024 * 1024, // 50MB限制
}

// 下载过程中会实时监控文件夹大小
// 超过限制时自动跳过后续文件
```

常用大小限制设置：
- `10 * 1024 * 1024`    // 10MB
- `50 * 1024 * 1024`    // 50MB  
- `100 * 1024 * 1024`   // 100MB
- `0`                   // 无限制

### 3. 代理配置

支持多种代理协议：

```go
config := &goclone.Config{
    URLs:        []string{"https://example.com"},
    ProxyString: "http://127.0.0.1:8080",     // HTTP代理
    // ProxyString: "socks5://127.0.0.1:1080", // SOCKS5代理
    // ProxyString: "https://user:pass@proxy.com:8080", // 认证代理
}
```

### 4. Cookie管理

支持复杂的cookie设置：

```go
config := &goclone.Config{
    URLs: []string{"https://example.com"},
    Cookies: []string{
        "session=abc123; domain=example.com; path=/",
        "user=test; secure; httponly",
        "theme=dark; max-age=3600",
    },
}
```

### 5. 自定义User-Agent

```go
config := &goclone.Config{
    URLs:      []string{"https://example.com"},
    UserAgent: "Mozilla/5.0 (compatible; MyBot/1.0; +http://mybot.com)",
}
```

## 🏃‍♂️ 运行项目

### 方式1: 运行示例代码
```bash
go run ./example/main.go
```

### 方式2: 构建可执行文件
```bash
# 构建
go build -o goclone.exe ./example/

# 运行
./goclone.exe
```

### 方式3: 直接在你的项目中使用
```bash
go mod init your-project
go get github.com/z-bool/go-website-clone
```

## 📊 输出示例

```
开始克隆 1 个URL，配置ID: a1b2c3d4-e5f6-7890-1234-567890abcdef
当前文件夹大小: 0 字节 (限制: 52428800 字节)
正在处理第 1 个URL: https://example.com
获取完整响应: https://example.com/ (目标URL: https://example.com/)
保存主页面HTML: https://example.com/
Css found --> css/style.css
Extracting --> https://example.com/css/style.css
Js found --> js/main.js
Extracting --> https://example.com/js/main.js
Img found --> images/logo.png
Extracting --> https://example.com/images/logo.png
最终文件夹大小: 2456789 字节 (限制: 52428800 字节)
URL https://example.com/ 克隆完成，项目路径: /path/to/a1b2c3d4-e5f6-7890-1234-567890abcdef
所有URL克隆完成
成功克隆 1 个网站:
  1. /path/to/a1b2c3d4-e5f6-7890-1234-567890abcdef
配置ID: a1b2c3d4-e5f6-7890-1234-567890abcdef
```

## 🔥 版本亮点

### v2.0 新功能
- ✨ **UUID文件夹命名**: 使用ConfigID字段自动生成唯一文件夹名
- ✨ **智能大小控制**: MaxFolderSize字段实现文件夹大小限制
- ✨ **实时监控**: 下载过程中实时显示文件夹大小变化
- ✨ **Windows兼容**: 智能处理特殊字符文件名
- ✨ **接口优化**: 通过CrawlConfig接口实现更好的代码结构

### 架构改进
- 🏗️ **模块化设计**: 各功能模块职责清晰，便于维护
- 🔧 **接口解耦**: 使用接口设计避免循环依赖
- 🛡️ **错误处理**: 更完善的错误处理和中文错误信息
- 📝 **日志优化**: 详细的进度日志和状态反馈

## ⚠️ 注意事项

1. **合规使用**: 遵守目标网站的robots.txt和服务条款
2. **网络环境**: 确保网络连接稳定，目标网站可访问
3. **资源限制**: 合理设置MaxFolderSize，避免磁盘空间不足
4. **URL格式**: 确保提供有效的URL格式
5. **权限问题**: 确保程序有文件写入权限
6. **代理设置**: 代理配置错误可能导致连接失败

## 🆘 常见问题

### Q: 如何设置合适的文件夹大小限制？
A: 根据实际需求设置，一般网站推荐50-100MB，大型网站可设置500MB或更多。

### Q: 支持哪些文件类型？
A: 自动下载HTML、CSS、JS、图片文件（jpg、png、gif、svg等）。

### Q: 如何处理需要登录的网站？
A: 使用Cookies字段设置登录后的cookie信息。

### Q: 代理不生效怎么办？
A: 检查代理地址格式，确保代理服务器可用，格式如`http://host:port`。

## 📄 许可证

本项目遵循 [License](LICENSE) 许可证。

## 🤝 贡献

欢迎提交Issue和Pull Request来改进项目！

---
⭐ 如果这个项目对你有帮助，请给它一个Star！ 