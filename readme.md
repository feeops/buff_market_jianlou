# buff.market捡漏软件

## 入口

按折扣购买入口(只考虑折扣是否符合，不考虑磨损度)，对应MpriceD.xlsx文件
https://buff.market/market/best_deals?category_group=knife,hands,rifle,pistol,smg,shotgun,machinegun&exterior=wearcategory0,wearcategory1,wearcategory2,wearcategory3,wearcategory4,wearcategoryna&quality=normal,strange,unusual,unusual_strange&min_price=5&sort_by=discount.desc&game=csgo
请求URL
https://api.buff.market/api/market/sell_order/low_price?game=csgo&page_num=1&page_size=50&exterior=wearcategory0,wearcategory1,wearcategory2&quality=normal,strange,unusual,unusual_strange&min_price=3.51&sort_by=discount.desc

最新上架入口，需要匹配麿损度，对应MpriceF.xlsx文件
https://buff.market/market/all?category_group=knife,hands,rifle,pistol,smg,shotgun,machinegun&rarity=ancient_weapon,legendary_weapon,mythical_weapon,ancient,ancient_character&exterior=wearcategory0,wearcategory1,wearcategory2&quality=normal,strange,unusual&min_price=1.39&game=csgo

## 备用库

https://github.com/bingoohuang/xlsx