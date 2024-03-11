package main

import (
	"fmt"
	"os"
	"path/filepath"
)

var (
	hostname string
	dirPath  string
)

func getDirPath() {
	absPath, err := filepath.Abs(os.Args[0])
	if err != nil {
		fmt.Println("获取绝对路径失败:", err)
		return
	}

	// 获取程序所在的目录路径
	dirPath = filepath.Dir(absPath)

}

func main() {
	readConfig()
	excelMap()
	hostname, _ = os.Hostname()
	getDirPath()
	setClient()
	checkMarket()
	run()
}
