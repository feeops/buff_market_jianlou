package main

import (
	"fmt"
	"github.com/tidwall/gjson"
)

var (
	marketSteamID string
	marketBuffID  string
)

func checkMarket() {

	accountURL := "https://api.buff.market/account/api/user/info"
	client.SetCommonHeaders(getHeaders)
	resp, _ := client.R().Get(accountURL)

	respStr := resp.String()
	marketBuffID = gjson.Get(respStr, "data.id").String()

	if len(marketBuffID) == 0 {
		fmt.Printf("状态码: %d 响应数据: %s\n", resp.StatusCode, resp.String())
		logger.Error().Str("resp", respStr).Msg("buffID is empty")
		fmt.Printf("获取不到buff.market账号ID，或重试执行程序，如果无效，请重新导出cookies，错误消息:%s\n", resp.String())
		waitExit()
	}

	nickname := gjson.Get(respStr, "data.nickname").Str
	email := gjson.Get(respStr, "data.email").Str
	accountInfo := fmt.Sprintf("关联的buff账号信息 昵称: %s 邮箱: %s", nickname, email)
	fmt.Println(accountInfo)

	if gjson.Get(respStr, "data.allow_buyer_bargain").Bool() {
		fmt.Println("允许买家还价功能处于开启状态，建议去https://buff.market/account/profile?game=csgo中关闭")
		waitExit()
	}

	marketSteamID = gjson.Get(respStr, "data.steamid").String()
	if len(marketSteamID) == 0 {
		fmt.Println("steamid没有关联，建议去https://buff.market/account/profile?game=csgo中关联")
		waitExit()
	}

	if len(gjson.Get(respStr, "data.trade_url").String()) == 0 {
		fmt.Println("steam交易链接没有绑定，建议先去https://buff.market/account/profile?game=csgo中绑定")
		waitExit()
	}

	if gjson.Get(respStr, "data.steam_api_key_state").Int() != 2 {
		fmt.Println("steam API状态不正常，建议先去https://buff.market/account/profile?game=csgo中确认")
		waitExit()
	}

	if err := auth(marketBuffID); err != nil {
		fmt.Println(err.Error())
		waitExit()
	}
}
