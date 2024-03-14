package main

import (
	"fmt"
	"github.com/duke-git/lancet/v2/slice"
	"github.com/spf13/cast"
	"github.com/tidwall/gjson"
	"math"
	"time"
)

const BuyURL = "https://api.buff.market/api/market/goods/buy"

func SmallerOrEqual(a, b float64) bool {
	return math.Max(a, b) == b || math.Abs(a-b) < 0.00001
}

func groupRun() {
	URL := fmt.Sprintf("https://api.buff.market/api/market/goods?category_group=knife,hands,rifle,pistol,smg,shotgun,machinegun&rarity=ancient_weapon,legendary_weapon,mythical_weapon,ancient,ancient_character&exterior=wearcategory0,wearcategory1,wearcategory2&quality=normal,strange,unusual&min_price=1.39&game=csgo&page_num=1&page_size=100&_d=%d",
		time.Now().Unix())
	resp, err := client.R().Get(URL)
	if err != nil {
		logger.Error().Str("URL", URL).Str("error", err.Error()).Msg("get error")
		return
	}

	updateCsrfToken(resp)

	respStr := resp.String()

	if gjson.Get(respStr, "code").Str == "OK" {
	} else {
		logger.Error().Str("URL", URL).Str("resp", respStr).Msg("code is not OK")
		return
	}

	for _, item := range gjson.Get(respStr, "data.items").Array() {
		id := item.Get("id").Int()

		values, found := groupPriceInfo.Get(id)
		if found {
		} else {
			continue
		}

		marketHashName := item.Get("market_hash_name").Str
		sellMinPrice := cast.ToFloat64(item.Get("sell_min_price").Str)

		if sellMinPrice <= 0 {
			logger.Error().Float64("sellMinPrice", sellMinPrice).Str("marketHashName", marketHashName).
				Msg("sellMinPrice error")
			continue
		}

		for _, value := range values {
			if value.PaintwearOpen == 1 && SmallerOrEqual(sellMinPrice, value.PaintwearBuyPrice) {
				query(sellMinPrice, value)
			}
		}
	}

}

func query(sellMinPrice float64, csItem GroupItem) {
	URL := fmt.Sprintf("https://api.buff.market/api/market/goods/sell_order?game=csgo&page_num=1&page_size=10&goods_id=%d&sort_by=default",
		csItem.GoodsID)
	resp, err := client.R().Get(URL)
	if err != nil {
		logger.Error().Str("URL", URL).Str("error", err.Error()).Msg("get error")
		return
	}

	updateCsrfToken(resp)

	respStr := resp.String()

	if gjson.Get(respStr, "code").Str == "OK" {
	} else {
		logger.Error().Str("URL", URL).Str("resp", respStr).Msg("code is not OK")
		return
	}

	for _, item := range gjson.Get(respStr, "data.items").Array() {
		price := cast.ToFloat64(item.Get("price").Str)
		if price > 0 && price <= sellMinPrice {
		} else {
			continue
		}

		paintWear := cast.ToFloat64(item.Get("asset_info.paintwear").Str)

		id := item.Get("id").Str
		logger.Info().Str("id", id).Int64("GoodsID", csItem.GoodsID).
			Float64("paintWear", paintWear).Msg("item check")
		if csItem.PaintwearOpen == 1 &&
			csItem.PaintwearBuyPrice <= price &&
			csItem.PaintwearLow <= paintWear &&
			paintWear <= csItem.PaintwearHigh {

			if slice.Contain(sellIDList, id) {
				return
			}

			msg := fmt.Sprintf("buff market ID:%d 中文名:%s 当前价格:%.2f 阈值最高价:%.2f 磨损度:%.2f",
				csItem.GoodsID, csItem.Name, price, csItem.PaintwearBuyPrice, paintWear)

			sellIDList = append(sellIDList, id)
			feishuNotify(msg)

			// buy(id, price)
		}
	}
}

type Buy struct {
	Game        string  `json:"game"`
	SellOrderID string  `json:"sell_order_id"`
	Price       float64 `json:"price"`
	PayMethod   int     `json:"pay_method"`
}

func buy(SellOrderID string, price float64) {

	b := Buy{
		Game:        "csgo",
		SellOrderID: SellOrderID,
		Price:       price,
		PayMethod:   12,
	}

	resp, err := client.R().SetBodyJsonMarshal(b).Post(BuyURL)

	if err != nil {
		logger.Error().Str("URL", BuyURL).Str("error", err.Error()).Msg("get error")
		return
	}

	updateCsrfToken(resp)

	respStr := resp.String()

	if gjson.Get(respStr, "code").Str == "OK" {
	} else {
		logger.Error().Str("resp", respStr).Msg("code is not OK")
		return
	}

}
