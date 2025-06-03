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
		URLs:          []string{"https://www.zttc.cn/"},
		Open:          false,
		UserAgent:     "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:139.0) Gecko/20100101 Firefox/139.0",
		ProxyString:   "", // 如需代理，格式如: "http://127.0.0.1:8080"
		Cookies:       []string{"session=abc123", "user=test"},
		ConfigID:      "",               // 留空将自动生成UUID，也可以手动指定
		MaxFolderSize: 50 * 1024 * 1024, // 50MB限制，设置为0表示无限制
	}

	result2 := goclone.Clone(ctx, config)
	if result2.Error != nil {
		log.Printf("克隆失败: %v", result2.Error)
	} else {
		fmt.Printf("成功克隆 %d 个网站:\n", len(result2.ProjectPaths))
		for i, path := range result2.ProjectPaths {
			fmt.Printf("  %d. %s\n", i+1, path)
		}
		fmt.Printf("配置ID: %s\n", config.ConfigID)
	}
}
