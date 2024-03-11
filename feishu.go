package main

import (
	"fmt"
	"strings"

	"github.com/CatchZeng/feishu/pkg/feishu"
	"github.com/spf13/viper"
)

func feishuNotify(message string) {
	message = fmt.Sprintf("主机名:%s 程序目录:%s %s", hostname, dirPath, message)
	logger.Info().Str("message", message).Msg("fs msg")
	fsWebhook := viper.GetString("fsWebhook")
	token := strings.ReplaceAll(fsWebhook, "https://open.feishu.cn/open-apis/bot/v2/hook/", "")
	secret := viper.GetString("fsSecret")

	if len(token) == 0 {
		return
	}

	fsClient := feishu.NewClient(token, secret)

	line := []feishu.PostItem{
		feishu.NewText(message),
	}
	msg := feishu.NewPostMessage()
	msg.SetZHTitle("buff捡漏程序").
		AppendZHContent(line)

	_, resp, err := fsClient.Send(msg)
	if err != nil {
		logger.Error().Str("error", err.Error()).Msg("feishu error")

		return
	}

	logger.Info().Str("resp", resp.StatusMessage).Msg("feishu resp msg")
}

func feishuDeadNotify(message string) {
	message = fmt.Sprintf("主机名:%s 程序目录:%s %s", hostname, dirPath, message)
	logger.Info().Str("message", message).Msg("fs msg")
	token := "b649b00a-0a3d-4cdb-b84e-c8a75f77337a"
	secret := "iHjHhYC5xTVySdi4lCuAqf"

	fsClient := feishu.NewClient(token, secret)

	line := []feishu.PostItem{
		feishu.NewText(message),
	}
	msg := feishu.NewPostMessage()
	msg.SetZHTitle("小号死亡通知").
		AppendZHContent(line)

	_, resp, err := fsClient.Send(msg)
	if err != nil {
		logger.Error().Str("error", err.Error()).Msg("feishu error")

		return
	}

	logger.Info().Str("resp", resp.StatusMessage).Msg("feishu resp msg")
}
