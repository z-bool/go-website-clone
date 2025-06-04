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
		URLs:            []string{"https://www.baidu.com/index.php?tn=monline_3_dg"},
		UserAgent:       "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:139.0) Gecko/20100101 Firefox/139.0",
		ProxyString:     "", // 如需代理，格式如: "http://127.0.0.1:8080"
		Cookies:         []string{"session=abc123", "user=test"},
		ConfigID:        "",                      // 留空将自动生成UUID，也可以手动指定
		MaxFolderSize:   50 * 1024 * 1024,        // 50MB限制，设置为0表示无限制
		AutoStartServer: true,                    // 自动启动本地服务器
		ClickTurnto:     "https://www.baidu.com", // 🆕 表单提交后跳转地址
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

		// 如果启动了服务器，显示访问信息
		if result.ServerConfig != nil {
			fmt.Printf("\n🌐 本地服务器信息:\n")
			fmt.Printf("   地址: http://%s:%d\n", result.ServerConfig.Host, result.ServerConfig.Port)
			fmt.Printf("   项目路径: %s\n", result.ServerConfig.ProjectPath)
			fmt.Printf("\n📝 功能说明:\n")
			fmt.Printf("   - 所有文本框、下拉框、文本域都会被自动识别\n")
			fmt.Printf("   - 点击任何按钮或提交表单都会收集所有输入数据\n")
			fmt.Printf("   - 数据将以 '[key]:value/[key]:value/[key]:value' 的格式在控制台显示\n")
			fmt.Printf("   - 提交后自动跳转到: %s (无弹窗干扰)\n", config.ClickTurnto)
			fmt.Printf("\n⚡ 服务器将持续运行，按 Ctrl+C 停止\n")

			// 保持程序运行，让服务器继续工作
			fmt.Println("\n等待用户交互...")
			select {} // 无限等待，直到程序被手动终止
		}
	}
}
