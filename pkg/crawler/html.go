package crawler

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

// HTMLExtractorFromResponse 从colly响应中提取HTML内容
func HTMLExtractorFromResponse(link string, projectPath string, bodyData []byte) {
	fmt.Println("从响应提取HTML --> ", link)
	fmt.Println("项目路径 --> ", projectPath)

	fmt.Printf("HTML内容长度: %d 字节\n", len(bodyData))

	if len(bodyData) == 0 {
		fmt.Println("警告: HTML内容为空!")
		return
	}

	// 创建或打开index.html文件
	f, err := os.OpenFile(projectPath+"/"+"index.html", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		fmt.Printf("创建文件失败: %v\n", err)
		return
	}
	defer f.Close()

	written, err := f.Write(bodyData)
	if err != nil {
		fmt.Printf("写入文件失败: %v\n", err)
		return
	}

	fmt.Printf("成功写入 %d 字节到文件\n", written)
}

// HTMLExtractor ...
func HTMLExtractor(link string, projectPath string) {
	fmt.Println("Extracting --> ", link)
	fmt.Println("Project path --> ", projectPath)

	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	// get the html body
	resp, err := http.Get(link)
	if err != nil {
		fmt.Printf("HTTP请求失败: %v\n", err)
		return
	}

	// Close the body once everything else is compled
	defer resp.Body.Close()

	fmt.Printf("HTTP状态码: %d\n", resp.StatusCode)
	fmt.Printf("Content-Type: %s\n", resp.Header.Get("Content-Type"))

	// get the project name and path we use the path to
	f, err := os.OpenFile(projectPath+"/"+"index.html", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		fmt.Printf("创建文件失败: %v\n", err)
		return
	}
	defer f.Close()

	htmlData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("读取响应内容失败: %v\n", err)
		return
	}

	fmt.Printf("HTML内容长度: %d 字节\n", len(htmlData))

	if len(htmlData) == 0 {
		fmt.Println("警告: HTML内容为空!")
		return
	}

	written, err := f.Write(htmlData)
	if err != nil {
		fmt.Printf("写入文件失败: %v\n", err)
		return
	}

	fmt.Printf("成功写入 %d 字节到文件\n", written)
}
