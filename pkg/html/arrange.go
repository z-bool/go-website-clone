package html

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// arrange 重构HTML文件中的链接，将外部资源链接改为本地路径
func arrange(projectDir string) error {
	indexfile := projectDir + "/index.html"

	// 读取整个HTML文件
	input, err := ioutil.ReadFile(indexfile)
	if err != nil {
		return fmt.Errorf("读取HTML文件失败: %w", err)
	}

	// 使用goquery解析整个HTML文档
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(input)))
	if err != nil {
		return fmt.Errorf("解析HTML文档失败: %w", err)
	}

	// 替换CSS链接
	doc.Find("link[rel='stylesheet']").Each(func(i int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if exists && !strings.HasPrefix(href, "css/") {
			file := filepath.Base(href)
			// 移除查询参数
			if idx := strings.Index(file, "?"); idx != -1 {
				file = file[:idx]
			}
			s.SetAttr("href", "css/"+file)
		}
	})

	// 替换JS链接
	doc.Find("script[src]").Each(func(i int, s *goquery.Selection) {
		src, exists := s.Attr("src")
		if exists && !strings.HasPrefix(src, "js/") {
			file := filepath.Base(src)
			// 移除查询参数
			if idx := strings.Index(file, "?"); idx != -1 {
				file = file[:idx]
			}
			s.SetAttr("src", "js/"+file)
		}
	})

	// 替换图片链接
	doc.Find("img[src]").Each(func(i int, s *goquery.Selection) {
		src, exists := s.Attr("src")
		if exists && !strings.HasPrefix(src, "imgs/") && !strings.HasPrefix(src, "data:") {
			file := filepath.Base(src)
			// 移除查询参数
			if idx := strings.Index(file, "?"); idx != -1 {
				file = file[:idx]
			}
			s.SetAttr("src", "imgs/"+file)
		}
	})

	// 获取修改后的HTML
	html, err := doc.Html()
	if err != nil {
		return fmt.Errorf("生成HTML失败: %w", err)
	}

	// 写回文件
	return ioutil.WriteFile(indexfile, []byte(html), 0777)
}

var reSrc = regexp.MustCompile(`src\s*=\s*"(.+?)"`)
