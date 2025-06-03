package crawler

import (
	"context"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"strings"

	"github.com/gocolly/colly/v2"
	"github.com/z-bool/go-website-clone/pkg/file"
)

// Collector searches for css, js, and images within a given link
// TODO improve for better performance
func Collector(ctx context.Context, url string, projectPath string, cookieJar *cookiejar.Jar, proxyString string, userAgent string) error {
	// create a new collector
	c := colly.NewCollector(colly.Async(true))
	setUpCollector(c, ctx, cookieJar, proxyString, userAgent)

	// search for all link tags that have a rel attribute that is equal to stylesheet - CSS
	c.OnHTML("link[rel='stylesheet']", func(e *colly.HTMLElement) {
		// hyperlink reference
		link := e.Attr("href")
		// print css file was found
		fmt.Println("Css found", "-->", link)
		// extraction
		Extractor(e.Request.AbsoluteURL(link), projectPath)
	})

	// search for all script tags with src attribute -- JS
	c.OnHTML("script[src]", func(e *colly.HTMLElement) {
		// src attribute
		link := e.Attr("src")
		// Print link
		fmt.Println("Js found", "-->", link)
		// extraction
		Extractor(e.Request.AbsoluteURL(link), projectPath)
	})

	// serach for all img tags with src attribute -- Images
	c.OnHTML("img[src]", func(e *colly.HTMLElement) {
		// src attribute
		link := e.Attr("src")
		if strings.HasPrefix(link, "data:image") || strings.HasPrefix(link, "blob:") {
			return
		}
		// Print link
		fmt.Println("Img found", "-->", link)
		// extraction
		Extractor(e.Request.AbsoluteURL(link), projectPath)
	})

	// 获取完整的HTML文档 - 改用OnResponse来获取原始HTML
	c.OnResponse(func(r *colly.Response) {
		currentURL := r.Request.URL.String()
		fmt.Printf("获取完整响应: %s (目标URL: %s)\n", currentURL, url)

		// 处理URL匹配，考虑尾部斜杠的差异
		normalizedCurrentURL := strings.TrimSuffix(currentURL, "/")
		normalizedTargetURL := strings.TrimSuffix(url, "/")

		if normalizedCurrentURL == normalizedTargetURL {
			fmt.Printf("保存主页面HTML: %s\n", currentURL)
			fmt.Printf("Content-Type: %s\n", r.Headers.Get("Content-Type"))
			htmlContent := r.Body
			fmt.Printf("HTML内容长度: %d 字节\n", len(htmlContent))

			// 检查内容类型，确保是HTML
			contentType := r.Headers.Get("Content-Type")
			if strings.Contains(strings.ToLower(contentType), "text/html") {
				HTMLExtractorFromResponse(currentURL, projectPath, htmlContent)
			} else {
				fmt.Printf("跳过非HTML内容: %s\n", contentType)
			}
		} else {
			fmt.Printf("URL不匹配，跳过响应: %s vs %s\n", normalizedCurrentURL, normalizedTargetURL)
		}
	})

	// Visit each url and wait for stuff to load :)
	if err := c.Visit(url); err != nil {
		return err
	}
	c.Wait()
	return nil
}

// CollectorWithSizeLimit 带大小限制的收集器
func CollectorWithSizeLimit(ctx context.Context, url string, projectPath string, cookieJar *cookiejar.Jar, proxyString string, userAgent string, maxFolderSize int64) error {
	// 在开始下载前检查当前大小
	if maxFolderSize > 0 {
		withinLimit, currentSize, err := file.CheckFolderSizeLimit(projectPath, maxFolderSize)
		if err != nil {
			return fmt.Errorf("检查文件夹大小失败: %w", err)
		}
		if !withinLimit {
			return fmt.Errorf("文件夹大小已超过限制: 当前 %d 字节, 限制 %d 字节", currentSize, maxFolderSize)
		}
		fmt.Printf("当前文件夹大小: %d 字节 (限制: %d 字节)\n", currentSize, maxFolderSize)
	}

	// 创建新的收集器
	c := colly.NewCollector(colly.Async(true))
	setUpCollector(c, ctx, cookieJar, proxyString, userAgent)

	// 在每次下载前检查大小限制
	c.OnHTML("link[rel='stylesheet']", func(e *colly.HTMLElement) {
		if maxFolderSize > 0 {
			if withinLimit, currentSize, err := file.CheckFolderSizeLimit(projectPath, maxFolderSize); err != nil || !withinLimit {
				if err != nil {
					fmt.Printf("检查文件夹大小失败: %v\n", err)
				} else {
					fmt.Printf("跳过CSS文件，文件夹大小超限: %d/%d 字节\n", currentSize, maxFolderSize)
				}
				return
			}
		}
		link := e.Attr("href")
		fmt.Println("Css found", "-->", link)
		Extractor(e.Request.AbsoluteURL(link), projectPath)
	})

	c.OnHTML("script[src]", func(e *colly.HTMLElement) {
		if maxFolderSize > 0 {
			if withinLimit, currentSize, err := file.CheckFolderSizeLimit(projectPath, maxFolderSize); err != nil || !withinLimit {
				if err != nil {
					fmt.Printf("检查文件夹大小失败: %v\n", err)
				} else {
					fmt.Printf("跳过JS文件，文件夹大小超限: %d/%d 字节\n", currentSize, maxFolderSize)
				}
				return
			}
		}
		link := e.Attr("src")
		fmt.Println("Js found", "-->", link)
		Extractor(e.Request.AbsoluteURL(link), projectPath)
	})

	c.OnHTML("img[src]", func(e *colly.HTMLElement) {
		if maxFolderSize > 0 {
			if withinLimit, currentSize, err := file.CheckFolderSizeLimit(projectPath, maxFolderSize); err != nil || !withinLimit {
				if err != nil {
					fmt.Printf("检查文件夹大小失败: %v\n", err)
				} else {
					fmt.Printf("跳过图片文件，文件夹大小超限: %d/%d 字节\n", currentSize, maxFolderSize)
				}
				return
			}
		}
		link := e.Attr("src")
		if strings.HasPrefix(link, "data:image") || strings.HasPrefix(link, "blob:") {
			return
		}
		fmt.Println("Img found", "-->", link)
		Extractor(e.Request.AbsoluteURL(link), projectPath)
	})

	// 获取完整的HTML文档
	c.OnResponse(func(r *colly.Response) {
		currentURL := r.Request.URL.String()
		fmt.Printf("获取完整响应: %s (目标URL: %s)\n", currentURL, url)

		normalizedCurrentURL := strings.TrimSuffix(currentURL, "/")
		normalizedTargetURL := strings.TrimSuffix(url, "/")

		if normalizedCurrentURL == normalizedTargetURL {
			fmt.Printf("保存主页面HTML: %s\n", currentURL)
			fmt.Printf("Content-Type: %s\n", r.Headers.Get("Content-Type"))
			htmlContent := r.Body
			fmt.Printf("HTML内容长度: %d 字节\n", len(htmlContent))

			contentType := r.Headers.Get("Content-Type")
			if strings.Contains(strings.ToLower(contentType), "text/html") {
				HTMLExtractorFromResponse(currentURL, projectPath, htmlContent)
			} else {
				fmt.Printf("跳过非HTML内容: %s\n", contentType)
			}
		} else {
			fmt.Printf("URL不匹配，跳过响应: %s vs %s\n", normalizedCurrentURL, normalizedTargetURL)
		}
	})

	if err := c.Visit(url); err != nil {
		return err
	}
	c.Wait()

	// 最终大小检查和报告
	if maxFolderSize > 0 {
		finalSize, err := file.GetFolderSize(projectPath)
		if err != nil {
			fmt.Printf("无法计算最终文件夹大小: %v\n", err)
		} else {
			fmt.Printf("最终文件夹大小: %d 字节 (限制: %d 字节)\n", finalSize, maxFolderSize)
		}
	}

	return nil
}

type cancelableTransport struct {
	ctx       context.Context
	transport http.RoundTripper
}

func (t cancelableTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if err := t.ctx.Err(); err != nil {
		return nil, err
	}
	return t.transport.RoundTrip(req.WithContext(t.ctx))
}

func setUpCollector(c *colly.Collector, ctx context.Context, cookieJar *cookiejar.Jar, proxyString, userAgent string) {
	if cookieJar != nil {
		c.SetCookieJar(cookieJar)
	}
	if proxyString != "" {
		c.SetProxy(proxyString)
	} else {
		c.WithTransport(cancelableTransport{ctx: ctx, transport: http.DefaultTransport})
	}
	if userAgent != "" {
		c.UserAgent = userAgent
	}
}
