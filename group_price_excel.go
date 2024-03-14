package main

import (
	"fmt"
	"github.com/bingoohuang/xlsx"
	"github.com/rafos/go-multimap/slicemultimap"
)

var (
	groupPriceInfo = slicemultimap.New[int64, GroupItem]()
)

type GroupItem struct {
	GoodsID int64 `title:"goodsId"`
	// GeneralBuyPrice float64 `title:"generalBuyPrice,omitempty"`

	PaintwearLow      float64 `title:"paintwearLow"`
	PaintwearHigh     float64 `title:"paintwearHigh"`
	PaintwearBuyPrice float64 `title:"paintwearBuyPrice"`
	PaintwearOpen     int     `title:"paintwearOpen"`

	//NameTagBuyPrice float64 `title:"nameTagBuyPrice,omitempty"`
	//NameTagOpen     int     `title:"nameTagOpen,omitempty"`

	IsOpen int `title:"isOpen"`
	// Category string `title:"category,omitempty"`
	Name string `title:"name"`
	// SellMinPrice  float64 `title:"sellMinPrice,omitempty"`
	// AutoPrice     int     `title:"autoPrice,omitempty"`
	// MonthMidPrice float64 `title:"monthMidPrice,omitempty"`
	// BuffURL       string  `title:"buffURL,omitempty"`
}

func groupMap() {

	// sheetName := "Sheet1"
	path := "MpriceF.xlsx"

	var GroupItems []GroupItem
	x, _ := xlsx.New(xlsx.WithInputFile(path))
	defer x.Close()

	if err := x.Read(&GroupItems); err != nil {
		fmt.Printf("%s价格解析出错，出错代码:%s\n", path, err.Error())
		panic(err)
	}

	for _, item := range GroupItems {
		if item.IsOpen == 0 {
			continue
		}

		goodsID := item.GoodsID

		groupPriceInfo.Put(goodsID, item)
	}

	fmt.Printf("%s表共有%d条数据\n", path, len(groupPriceInfo.Values()))
}
