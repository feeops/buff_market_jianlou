package main

import (
	"fmt"
	"github.com/duke-git/lancet/v2/cryptor"
	"github.com/duke-git/lancet/v2/fileutil"
	"github.com/imroc/req/v3"
	"github.com/tidwall/gjson"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

var (
	getHeaders  = map[string]string{}
	client      = req.C()
	postHeaders = map[string]string{}
)

func cookiesTXT() {
	fileStr, _ := fileutil.ReadFileToString("buff_cookies.txt")
	if len(fileStr) == 0 {
		fmt.Println("buff_cookies.txt数据为空，请检查")
		waitExit()
	}

	getHeaders = map[string]string{
		"accept":          "application/json, text/plain, */*",
		"user-agent":      "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36",
		"origin":          "https://buff.market",
		"referer":         "https://buff.market/",
		"accept-language": "zh-CN,zh;q=0.9",
	}

	for _, item := range gjson.Parse(fileStr).Array() {

		if strings.Contains(fileStr, "api.buff.market") {
		} else {
			fmt.Println("cookies文件不对，请确保在https://api.buff.market页面上导出cookies")
			waitExit()
		}

		value := item.Get("value").Str

		// 解决net/http: invalid byte '"' in Cookie.Value; dropping invalid bytes问题
		if strings.Contains(value, `"`) {
			continue
		}

		name := item.Get("name").Str

		if name == "csrf_token" {
			csrfToken = value
		}

		client.SetCommonCookies(&http.Cookie{
			// Domain: domain,
			Name:  name,
			Path:  item.Get("path").Str,
			Value: value,
		})
	}

	client.SetCommonHeaders(getHeaders)
	updatePost()
	fmt.Println("已读取buff_cookies.txt中的数据")

}

func cookiesHAR(content string) {

	for _, entry := range gjson.Get(content, "log.entries").Array() {
		url := entry.Get("request.url").String()
		if strings.HasPrefix(url, "https://api.buff.market") {
		} else {
			continue
		}

		for _, header := range entry.Get("request.headers").Array() {
			name := header.Get("name").String()
			value := header.Get("value").String()
			lowerName := strings.ToLower(name)
			if strings.HasPrefix(name, ":") || lowerName == "cookie" || lowerName == "accept-encoding" {
				continue
			}

			getHeaders[name] = value

		}

		for _, cookie := range entry.Get("response.cookies").Array() {
			name := cookie.Get("name").String()
			value := cookie.Get("value").String()

			if name == "csrf_token" {
				csrfToken = value
			}

			client.SetCommonCookies(&http.Cookie{
				// Domain: domain,
				Name:  name,
				Value: value,
			})
		}

	}

	client.SetCommonHeaders(getHeaders)
	updatePost()

}

func updatePost() {
	for k, v := range getHeaders {
		postHeaders[k] = v
	}

	postHeaders["x-csrftoken"] = csrfToken
	postHeaders["content-type"] = "application/json"
}

func loadCSGO() bool {
	current := getCurrentAbPath()
	files, _ := os.ReadDir(current)

	path, _ := filepath.Abs(current)
	for _, file := range files {
		fullPath := filepath.Join(path, file.Name())
		if strings.HasSuffix(fullPath, ".csgo") {
		} else {
			continue
		}

		fmt.Printf("读取%s文件\n", file.Name())
		fileBytes, _ := os.ReadFile(fullPath)
		key := "CR385gYCr86WdbWAfYz8v8uUZDWNAtmC"
		decrypt := cryptor.AesEcbDecrypt(fileBytes, []byte(key))
		cookiesHAR(string(decrypt))
		return true
	}

	return false
}

func setClient() {

	if loadCSGO() {
		return
	}

	if fileutil.IsExist("cookies.har") {
		content, _ := fileutil.ReadFileToString("cookies.har")
		cookiesHAR(content)
		fmt.Println("已读取cookies.har中的数据")
		return
	}

	if fileutil.IsExist("buff_cookies.txt") {
		cookiesTXT()
	}

}
