package common

import (
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

func ParseUrl(link string) *url.URL {
	if link == "" {
		return nil
	}
	u, err := url.Parse(link)
	if err != nil {
		return nil
	}
	return u
}

func SplitAndAddBearer(authTokens string) []string {
	var authTokenList []string
	for _, v := range strings.Split(authTokens, ",") {
		authTokenList = append(authTokenList, "Bearer "+v)
	}
	return authTokenList
}

// GetAbsPathAndGenerate 获取绝对路径并生成文件或文件夹
func GetAbsPathAndGenerate(path string, isFilePath bool, content string) string {
	// 获取绝对路径
	path = GetAbsPath(path)
	if isFilePath {
		//	判断文件是否存在
		if isExist := fileIsExistAndCreat(path, content); isExist {
			return path
		}
	} else {
		//	判断文件夹是否存在
		if isExist := dirIsExistAndMkdir(path, false); isExist {
			return path
		}
	}
	return path
}

// GetAbsPath 获取绝对路径
func GetAbsPath(path string) string {
	if !filepath.IsAbs(path) {
		absPath, err := filepath.Abs(path)
		if err != nil {
			return ""
		}
		return absPath
	}
	return path
}

func dirIsExistAndMkdir(dirPath string, isFile bool) bool {
	// 判断路径是否存在
	_, err := os.Stat(dirPath)
	dir := dirPath
	if err != nil {
		if isFile {
			dir = filepath.Dir(dirPath)
		}
		// 创建路径
		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			return false
		}
	}
	return true
}

func fileIsExistAndCreat(filePath string, content string) bool {
	//判断文件是否存在
	_, err := os.Stat(filePath)
	if err != nil {
		// 判断文件夹是否存在
		if isExist := dirIsExistAndMkdir(filePath, true); !isExist {
			return false
		}
		// 创建文件
		_, err := os.Create(filePath)
		if err != nil {
			return false
		}
		if content != "" {
			//	写入content
			file, _ := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0777)
			_, _ = file.Write([]byte(content))
			defer func(file *os.File) {
				_ = file.Close()
			}(file)
		}
	}
	return true
}
