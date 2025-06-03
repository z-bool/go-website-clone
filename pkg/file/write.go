package file

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
)

// CreateProject initializes the project directory and returns the path to the project
// TODO make function more modular to obtain different html files
func CreateProject(projectName string) string {
	// current workin directory
	path := currentDirectory()

	// define project path
	projectPath := path + "/" + projectName

	// create base directory
	err := os.MkdirAll(projectPath, 0777)
	check(err)

	// create CSS/JS/Image directories
	createCSS(projectPath)
	createJS(projectPath)
	createIMG(projectPath)

	// main inedx file
	_, err = os.Create(projectPath + "/" + "index.html")
	check(err)
	// project path
	return projectPath
}

// CreateProjectWithID 使用指定的ID创建项目目录并返回项目路径
func CreateProjectWithID(configID string) string {
	// 当前工作目录
	path := currentDirectory()

	// 使用ConfigID定义项目路径
	projectPath := filepath.Join(path, configID)

	// 创建基础目录
	err := os.MkdirAll(projectPath, 0777)
	check(err)

	// 创建CSS/JS/Image目录
	createCSS(projectPath)
	createJS(projectPath)
	createIMG(projectPath)

	// 主index文件
	_, err = os.Create(filepath.Join(projectPath, "index.html"))
	check(err)

	fmt.Printf("项目目录已创建: %s\n", projectPath)
	return projectPath
}

// GetFolderSize 计算文件夹的总大小（字节）
func GetFolderSize(folderPath string) (int64, error) {
	var totalSize int64

	err := filepath.WalkDir(folderPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			info, err := d.Info()
			if err != nil {
				return err
			}
			totalSize += info.Size()
		}
		return nil
	})

	return totalSize, err
}

// CheckFolderSizeLimit 检查文件夹大小是否超过限制
func CheckFolderSizeLimit(folderPath string, maxSize int64) (bool, int64, error) {
	if maxSize <= 0 {
		return true, 0, nil // 没有大小限制
	}

	currentSize, err := GetFolderSize(folderPath)
	if err != nil {
		return false, 0, err
	}

	return currentSize <= maxSize, currentSize, nil
}

// currentDirectory get the current working directory
func currentDirectory() string {
	path, err := os.Getwd()
	check(err)
	return path
}

// createCSS create a css directory in the current path
func createCSS(path string) {
	// create css directory
	err := os.MkdirAll(path+"/"+"css", 0777)
	check(err)
}

// createJS create a JS directory in the current path
func createJS(path string) {
	err := os.MkdirAll(path+"/"+"js", 0777)
	check(err)
}

// createIMG create a image directory in the current path
func createIMG(path string) {
	err := os.MkdirAll(path+"/"+"imgs", 0777)
	check(err)
}

func check(err error) {
	if err != nil {
		log.Println(err)
	}
}
