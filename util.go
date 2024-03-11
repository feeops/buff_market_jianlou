package main

import (
	"fmt"
	"os"
	"time"

	"github.com/imroc/req/v3"
)

var (
	csrfToken string
)

func waitExit() {
	fmt.Println("请按任意键退出,如果没有操作，5分钟后自动退出")
	ch := make(chan string)
	go func() {
		var input string
		fmt.Scanln(&input)
		ch <- input
	}()

	select {
	case <-ch:
		os.Exit(0)
	case <-time.After(5 * time.Minute):
		fmt.Println("超时退出")
		os.Exit(0)
	}

}

func updateCsrfToken(resp *req.Response) {
	for _, ck := range resp.Cookies() {
		if ck.Name == "csrf_token" {
			csrfToken = ck.Value
		}
	}
}
