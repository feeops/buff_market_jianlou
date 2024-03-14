package main

import (
	"fmt"
	"github.com/bingoohuang/xlsx"
	"github.com/rafos/go-multimap/slicemultimap"
)

var (
	lowPriceInfo = slicemultimap.New[int64, LowPriceItem]()
)

type LowPriceItem struct {
	GoodsID  int64   `title:"ID"`
	Discount float64 `title:"Discount"`
	tagID    int     `title:"tag_ids"`

	MAX     float64 `title:"MAX"`
	maxOpen int     `title:"maxopen"`

	buffID int `title:"buffid"`

	IsOpen int `title:"isopen"`
	// Category string `title:"category,omitempty"`
	Name     string `title:"name"`
	Category string `title:"category"`
	// SellMinPrice  float64 `title:"sellMinPrice,omitempty"`
	// AutoPrice     int     `title:"autoPrice,omitempty"`
	// MonthMidPrice float64 `title:"monthMidPrice,omitempty"`
	// BuffURL       string  `title:"buffURL,omitempty"`
}

func lowMap() {

	path := "MpriceD.xlsx"

	var LowPriceItems []LowPriceItem
	x, _ := xlsx.New(xlsx.WithInputFile(path))
	defer x.Close()

	if err := x.Read(&LowPriceItems); err != nil {
		fmt.Printf("%s价格解析出错，出错代码:%s\n", path, err.Error())
		panic(err)
	}

	for _, item := range LowPriceItems {
		if item.IsOpen == 0 {
			continue
		}

		// 先跳过有tagID的情况
		if item.tagID > 0 {
			continue
		}

		lowPriceInfo.Put(item.GoodsID, item)

	}

	fmt.Printf("%s表共有%d条数据\n", path, len(lowPriceInfo.Values()))
}
