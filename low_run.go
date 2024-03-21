package main

import (
	"fmt"
	"github.com/duke-git/lancet/v2/slice"
	"github.com/spf13/cast"
	"github.com/tidwall/gjson"
	"math"
	"time"
)

func lowRun() {

	URL := fmt.Sprintf("https://api.buff.market/api/market/sell_order/low_price?game=csgo&page_num=1&page_size=50&category_group=knife,hands,rifle,pistol,smg,shotgun,machinegun&exterior=wearcategory0,wearcategory1,wearcategory2,wearcategory3,wearcategory4,wearcategoryna&quality=normal,strange,unusual,unusual_strange&min_price=5&sort_by=discount.desc&_d=%d",
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
		id := item.Get("id").Str
		if slice.Contain(sellIDList, id) {
			return
		}

		goodsID := item.Get("goods_id").Int()
		values, found := lowPriceInfo.Get(goodsID)
		if found {
		} else {
			continue
		}

		marketHashName := item.Get("market_hash_name").Str
		price := cast.ToFloat64(item.Get("price").Str)
		discount := math.Abs(item.Get("discount").Float())

		for _, value := range values {
			if discount >= value.Discount {
			} else {
				continue
			}
			var msg string
			switch {
			case value.maxOpen == 1 && price <= value.MAX:
				msg = fmt.Sprintf("buff market ID:%d marketHashName:%s 中文名:%s 当前折扣:%.2f 阈值折扣:%.2f 当前价格:%.2f 阈值最高价:%.2f",
					value.GoodsID, marketHashName, value.Name, discount, value.Discount, price, value.MAX)
			case value.maxOpen == 0:
				msg = fmt.Sprintf("buff market ID:%d marketHashName:%s 中文名:%s 当前折扣:%.2f 阈值折扣:%.2f 当前价格:%.2f 阈值最高价:%.2f",
					value.GoodsID, marketHashName, value.Name, discount, value.Discount, price, value.MAX)

			}
			sellIDList = append(sellIDList, id)
			feishuNotify(msg)
		}
	}
}
