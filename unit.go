package go_chrome

import (
	"encoding/json"
	"os"
)

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
