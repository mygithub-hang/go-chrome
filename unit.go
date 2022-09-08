package go_chrome

import (
	"archive/zip"
	"encoding/json"
	"github.com/mitchellh/go-homedir"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
)

func createDir(filePath string) error {
	filePath = path.Dir(filePath)
	if !IsExist(filePath) {
		err := os.MkdirAll(filePath, os.ModePerm)
		return err
	}
	return nil
}

// getWorkingDirPath 获取执行路径
func getWorkingDirPath() string {
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return dir
}

// IsExist 判断文件是否存在
func IsExist(fileAddr string) bool {
	// 读取文件信息，判断文件是否存在
	_, err := os.Stat(fileAddr)
	if err != nil {
		if os.IsExist(err) { // 根据错误类型进行判断
			return true
		}
		return false
	}
	return true
}

func getPackageConfig() *PackageConf {
	runPath := getWorkingDirPath()
	fileContent, err := os.ReadFile(runPath + "/" + "package.json")
	if err != nil {
		return nil
	}
	jsonData := PackageConf{}
	jsonData.RunBuildPath = runPath
	err = json.Unmarshal(fileContent, &jsonData)
	if err != nil {
		return nil
	}
	return &jsonData
}

func UnPackZip(src, dest string) error {
	reader, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer reader.Close()
	for _, file := range reader.File {
		filePath := path.Join(dest, file.Name)
		if file.FileInfo().IsDir() {
			os.MkdirAll(filePath, os.ModePerm)
		} else {
			if err = os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
				return err
			}
			inFile, err := file.Open()
			if err != nil {
				return err
			}
			defer inFile.Close()
			outFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
			if err != nil {
				return err
			}
			defer outFile.Close()
			_, err = io.Copy(outFile, inFile)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func GetHomePath() string {
	dir, err := homedir.Dir()
	if err != nil {
		log.Fatal(err)
	}
	return dir
}
