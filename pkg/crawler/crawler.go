package crawler

import (
	"context"
	"net/http/cookiejar"
)

// CrawlConfig 爬取配置接口
type CrawlConfig interface {
	GetProxyString() string
	GetUserAgent() string
	GetMaxFolderSize() int64
}

// Crawl asks the necessary crawlers for collecting links for building the web page
func Crawl(ctx context.Context, site string, projectPath string, cookieJar *cookiejar.Jar, proxyString string, userAgent string) error {
	// searches for css, js, and images within a given link
	return Collector(ctx, site, projectPath, cookieJar, proxyString, userAgent)
}

// CrawlWithConfig 使用配置对象进行爬取，支持大小检查
func CrawlWithConfig(ctx context.Context, site string, projectPath string, cookieJar *cookiejar.Jar, config CrawlConfig) error {
	return CollectorWithSizeLimit(ctx, site, projectPath, cookieJar, config.GetProxyString(), config.GetUserAgent(), config.GetMaxFolderSize())
}
