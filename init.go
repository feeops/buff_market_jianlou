package main

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/rs/zerolog"
	"github.com/virtuald/go-paniclog"
)

var logger zerolog.Logger

// 最终方案-全兼容.
func getCurrentAbPath() string {
	dir := getCurrentAbPathByExecutable()
	tmpDir, _ := filepath.EvalSymlinks(os.TempDir())

	if strings.Contains(dir, tmpDir) {
		return getCurrentAbPathByCaller()
	}

	return dir
}

// 获取当前执行文件绝对路径.
func getCurrentAbPathByExecutable() string {
	exePath, err := os.Executable()
	if err != nil {
		panic(err)
	}

	res, _ := filepath.EvalSymlinks(filepath.Dir(exePath))

	return res
}

// 获取当前执行文件绝对路径（go run）.
func getCurrentAbPathByCaller() string {
	var abPath string

	_, filename, _, ok := runtime.Caller(0)
	if ok {
		abPath = path.Dir(filename)
	}

	return abPath
}

func init() {
	os.Chdir(getCurrentAbPath())
	fileLog := filepath.Base(fmt.Sprintf("%s.log", os.Args[0]))
	writer, err := os.OpenFile(fileLog, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)

	if err != nil {
		panic(err)
	}

	logger = zerolog.New(writer).With().Timestamp().Logger()

	f, err := os.OpenFile("crash.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("Error creating file:", err)
		os.Exit(1)
	}

	undo, err := paniclog.RedirectStderr(f)
	if err != nil {
		fmt.Println("Error redirecting stderr:", err)
		os.Exit(1)
	}

	if os.Getenv("UNDO_PANICLOG") != "" {
		// demonstrates undoing the stderr redirect
		undo()
	}
}
