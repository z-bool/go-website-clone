# GoClone 函数式API使用指南

## 概述

本项目已从命令行工具重构为可直接调用的函数库。您可以在Go代码中直接导入并使用这些函数，而无需通过命令行调用。

## 项目结构

```
goclone/
├── pkg/
│   ├── goclone/          # 主要API包
│   │   └── goclone.go    # 核心函数
│   ├── crawler/          # 爬虫功能
│   ├── file/             # 文件处理
│   ├── html/             # HTML处理
│   ├── parser/           # URL解析
│   └── server/           # 本地服务器
├── example/
│   └── main.go           # 使用示例
├── main.go               # 项目入口点
└── go.mod                # Go模块文件
```

## 主要功能

### 1. 配置结构体

```go
type Config struct {
    URLs        []string  // 要克隆的网站URL列表
    UserAgent   string    // 自定义用户代理
    ProxyString string    // 代理连接字符串
    Cookies     []string  // 预设的cookie列表
}
```

### 2. 结果结构体

```go
type CloneResult struct {
    Success      bool      // 是否成功
    ProjectPaths []string  // 生成的项目路径列表
    FirstProject string    // 第一个项目路径（用于服务器或打开）
    Error        error     // 错误信息
}
```

## 使用方法

### 本地导入

```go
import (
    "context"
    "github.com/z-bool/go-website-clone/pkg/goclone"
)
```

### 快速克隆单个网站

```go
ctx := context.Background()
result := goclone.QuickClone(ctx, "https://example.com")
if result.Error != nil {
    log.Printf("克隆失败: %v", result.Error)
} else {
    fmt.Printf("成功克隆到: %s\n", result.FirstProject)
}
```

### 完整配置克隆

```go
config := &goclone.Config{
    URLs:        []string{"https://example.com", "https://httpbin.org"},
    UserAgent:   "GoClone/1.0",
    ProxyString: "http://127.0.0.1:8080", // 可选的代理设置
    Cookies:     []string{"session=abc123", "user=test"},
}

result := goclone.Clone(ctx, config)
if result.Error != nil {
    log.Printf("克隆失败: %v", result.Error)
} else {
    fmt.Printf("成功克隆 %d 个网站\n", len(result.ProjectPaths))
    for _, path := range result.ProjectPaths {
        fmt.Printf("项目路径: %s\n", path)
    }
}
```

## 运行项目

### 方式1: 直接运行主程序
```bash
go run main.go
```

### 方式2: 构建并运行
```bash
go build -o goclone.exe
./goclone.exe
```

### 方式3: 运行示例
```bash
go run ./example
```

## 主要改进

1. **去除命令行依赖**: 不再依赖cobra等命令行框架
2. **函数式调用**: 可以直接在Go代码中调用
3. **结构化配置**: 使用Config结构体替代命令行参数
4. **详细结果**: 返回CloneResult结构体，包含成功状态和详细信息
5. **错误处理**: 更好的错误处理和中文错误信息
6. **多URL支持**: 支持一次克隆多个网站
7. **便捷函数**: 提供QuickClone等便捷函数
8. **本地包结构**: 使用pkg/goclone包组织代码

## 代理设置

支持HTTP和SOCKS5代理：

```go
config := &goclone.Config{
    URLs:        []string{"https://example.com"},
    ProxyString: "http://127.0.0.1:8080",    // HTTP代理
    // ProxyString: "socks5://127.0.0.1:1080", // SOCKS5代理
}
```

## Cookie设置

支持预设cookies：

```go
config := &goclone.Config{
    URLs:    []string{"https://example.com"},
    Cookies: []string{
        "session=abc123; domain=example.com",
        "user=test; path=/",
    },
}
```

## 注意事项

1. 确保提供有效的URL
2. 网络连接正常
3. 目标网站允许爬取
4. 遵守目标网站的robots.txt和服务条款
5. 使用 `go mod init goclone` 初始化项目
6. 导入路径为 `github.com/z-bool/go-website-clone/pkg/goclone` 