package main

import (
	"fmt"
	"github.com/rafos/go-multimap/slicemultimap"

	"github.com/KarlTango/xs"
)

var (
	itemPrice = map[string]float64{}
	itemInfo  = map[string]CsgoItem{}
	matchInfo = slicemultimap.New[int64, CsgoItem]()
)

type CsgoItem struct {
	GoodsID         int64   `xs:"goodsId"`
	GeneralBuyPrice float64 `xs:"generalBuyPrice,omitempty"`
	GeneralOpen     int     `xs:"generalOpen,omitempty"`

	PaintwearLow      float64 `xs:"paintwearLow"`
	PaintwearHigh     float64 `xs:"paintwearHigh"`
	PaintwearBuyPrice float64 `xs:"paintwearBuyPrice"`
	PaintwearOpen     int     `xs:"paintwearOpen"`

	PrintLowValue float64 `xs:"printLowValue,omitempty"`
	PrintBuyPrice float64 `xs:"printBuyPrice,omitempty"`
	PrintOpen     int     `xs:"printOpen,omitempty"`

	NameTagBuyPrice float64 `xs:"nameTagBuyPrice,omitempty"`
	NameTagOpen     int     `xs:"nameTagOpen,omitempty"`

	IsOpen int `xs:"isOpen"`
	// Category string `xs:"category,omitempty"`
	Name string `xs:"name"`
	// SellMinPrice  float64 `xs:"sellMinPrice,omitempty"`
	// AutoPrice     int     `xs:"autoPrice,omitempty"`
	// MonthMidPrice float64 `xs:"monthMidPrice,omitempty"`
	// BuffURL       string  `xs:"buffURL,omitempty"`
}

func excelMap() {

	// sheetName := "Sheet1"
	path := "price.xlsx"

	var CsgoItems = make([]CsgoItem, 0)

	if err := xs.UnmarshalFromFile(path, &CsgoItems); err != nil {
		fmt.Printf("price.xlsx价格解析出错，出错代码:%s\n", err.Error())
		panic(err)
	}

	for _, item := range CsgoItems {
		if item.IsOpen == 0 {
			continue
		}

		goodsID := item.GoodsID
		defaultURL := fmt.Sprintf(
			"https://api.buff.market/api/market/goods/sell_order?game=csgo&page_num=1&page_size=50&goods_id=%d&sort_by=default",
			goodsID)
		if item.GeneralOpen == 1 && item.GeneralBuyPrice > 0 {
			URL := defaultURL
			itemPrice[URL] = item.GeneralBuyPrice
			itemInfo[URL] = item
		}

		if item.PaintwearOpen == 1 && item.PaintwearHigh > 0 && item.PaintwearBuyPrice > 0 {
			URL := fmt.Sprintf("%s&min_paintwear=%.3f&max_paintwear=%.3f",
				defaultURL, item.PaintwearLow, item.PaintwearHigh)
			itemPrice[URL] = item.PaintwearBuyPrice
			itemInfo[URL] = item
		}

		if item.NameTagOpen == 1 && item.NameTagBuyPrice > 0 {
			URL := fmt.Sprintf("%s&name_tag=1",
				defaultURL)
			itemPrice[URL] = item.NameTagBuyPrice
			itemInfo[URL] = item
		}

		matchInfo.Put(item.GoodsID, item)
	}

	fmt.Println(len(itemPrice), len(itemInfo))
}
