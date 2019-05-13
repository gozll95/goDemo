package main

import (
	"fmt"
	"strings"

	"strconv"

	mapset "github.com/deckarep/golang-set"
)

var qvm = `
1315564408
1369076646
1369077400
1369080337
1377016252
1378262280
1378262392
1380260887
1380266017
1380268845
1380271660
1380274142
1380274312
1380274395
1380285359
1380286579
1380286754
1380294735
1380306818
1380317514
1380324559
1380341292
1380349130
1380355448
1380377315
1380380799
1380390564
1380418374
1380420092
1380424979
1380427967
1380428167
1380430036
1380430123
1380445203
1380448964
1380451873
1380459707
1380498824
1380499432
1380505636
1380582158
1380590444
1380618401
1380631309
1380677317
1380679757
1380680703
1380681123
1380681293
1380681295
1380681302
1380686566
1380689019
1380697642
1380707881
1380711717
1380717415
1380727834
1380736904
1380764156
1380766663
1380768810
1380772162
1380774047
1380799708
1380802653
1380803382
1380811226
1380817118
1380817119
1380817121
1380817123
1380818363
1380819072
1380824098
1380826584
1380844963
1380854529
1380856508
1380862110
1380863522
1380867443
1380867950
1380872920
1380885366
1380886626
1380891407
1380908198
1380914741
1380923455
1380926737
1380943660
1380950452
1380968631
1380972796
1380976107
1380976211
1380980443
1380989730
1380998657
1381001019
1381004756
1381022383
1381028894
1381065481
1381075616
1381080309
1381084496
1381092231
1381104896
1381107605
1381109675
1381113289
1381117325
1381140017
1381168308
1381178650
1381193450
1381221254
1381223778
1381227081
1381227805
1381239387
1381255870
1381258175
1381293539
1381299323
1381299873
1381308487
1381315228
1381322374
1381328752
1381331111
1381344641
1381345542
1381347444
1381354173
1381360700
1381364985
1381368642
1381400613
1381401458
1381403820
1381403989
1381405935
1381407673
1381407894
1381412126
1381412647
1381417113
1381417608
1381424989
1381426231
1381429112
1381431772
1381435086
1381435128
1381439022
1381448150
1381453144
1381457787
1381461989
1381477372
1381477495
1381478212
1381478745
1381482477
1381484440
1381484473
1381485636
1381486517
1381487032
1381493809
1381498844
1381499796
1381500198
1381502457
1381503510
1381504202
1381505137
1381508953
1381520590
1381525584
1381530452
1381530794
1381531470
1381531504
1381538527
1381538772
1381538925
1381539452
1381539723
1381540923
1381541291
1381543301
1381545165
1381545370
1381545908
1381545992
1381546613
1381546758
1381549536
1381549669
1381549734
1381551890
1381553081
1381553542
1381554295
1381554337
1381555489
1381559622
1381561081
1381561692
1381563899
1381567101
1381567606
1381569214
1381572338
1381573393
1381573512
1381574398
1381574860
1381574945
1381575465
1381575809
1381577545
1381577875
1381578022
1381579179
1381579378
1381581661
1381587172
1381588233
1381588463
1381588535
1381589862
`

func main() {
	boItem := mapset.NewSet()
	qvmItem := mapset.NewSet()

	var bo []uint32
	bo = []uint32{1380968631, 1380306818, 1381530794, 1381549734, 1380950452, 1380448964, 1381539723, 1381426231, 1381453144, 1380681123, 1381400613, 1381520590, 1381109675, 1381299323, 1380427967, 1381503510, 1381331111, 1380324559, 1380390564, 1380811226, 1380341292, 1380891407, 1381308487, 1380451873, 1381538925, 1380885366, 1381004756, 1380590444, 1380980443, 1380686566, 1380908198, 1380854529, 1381559622, 1381545165, 1381477372, 1381403989, 1380818363, 1380268845, 1381499796, 1381221254, 1380697642, 1381538772, 1381538527, 1381545992, 1380867443, 1381407894, 1381545370, 1380689019, 1381457787, 1381502457, 1381405935, 1380681302, 1381424989, 1381461989, 1380914741, 1381546613, 1380285359, 1381344641, 1381543301, 1380803382, 1381493809, 1380459707, 1380817123, 1381431772, 1381484473, 1380819072, 1380976211, 1380802653, 1381485636, 1380799708, 1381258175, 1381567606, 1380420092, 1381561081, 1381107605, 1380271660, 1381541291, 1381545908, 1380294735, 1381546758, 1381561692, 1380286579, 1380286754, 1381482477, 1369076646, 1381435128, 1380707881, 1381478745, 1369077400, 1381293539, 1380976107, 1380582158, 1315564408, 1381227805, 1380844963, 1380768810, 1380355448, 1380766663, 1381322374, 1381315228, 1381104896, 1381554295, 1380631309, 1381345542, 1380418374, 1381540923, 1381484440, 1381567101, 1381555489, 1381178650, 1381347444, 1381500198, 1381407673, 1381530452, 1381113289, 1381140017, 1380817118, 1369080337, 1381001019, 1381239387, 1380856508, 1381498844, 1381531470, 1380274142, 1381525584, 1381401458, 1380863522, 1381412126, 1381508953, 1381549669, 1380989730, 1381504202, 1380817121, 1381554337, 1380727834, 1380677317, 1380681293, 1380505636, 1381022383, 1381255870, 1381193450, 1380274312, 1381478212, 1380736904, 1380872920, 1380680703, 1381368642, 1381223778, 1381448150, 1381092231, 1380972796, 1381439022, 1381417608, 1381553081, 1377016252, 1380772162, 1381403820, 1380886626, 1381364985, 1381553542, 1380711717, 1378262392, 1381360700, 1381065481, 1380424979, 1380349130, 1381354173, 1380867950, 1381551890, 1381569214, 1381227081, 1380377315, 1381028894, 1381487032, 1380428167, 1381117325, 1381531504, 1381549536, 1381477495, 1381328752, 1380826584, 1380824098, 1380943660, 1381080309, 1381505137, 1380681295, 1380651806, 1380260887, 1380862110, 1380717415, 1381412647, 1380266017, 1380998657, 1380926737, 1381429112, 1381417113, 1381168308, 1381435086, 1381539452, 1380764156, 1378262280, 1380274395, 1380817119, 1380430123, 1381299873, 1381572338, 1381573512, 1381573393, 1381574398, 1381575465, 1381575809, 1381577545, 1381577875, 1381578022, 1381579378, 1381579179, 1381581661, 1381587172, 1381588463, 1381588535, 1381588233}
	for _, v := range bo {
		boItem.Add(v)
	}
	//fmt.Println("boItem:", boItem)

	qvms := strings.Split(qvm, "\n")
	//	fmt.Println(qvms)

	for _, v := range qvms {
		vv, _ := strconv.Atoi(v)
		qvmItem.Add(uint32(vv))
	}
	//fmt.Println("qvmItem:", qvmItem)

	fmt.Println(boItem.Difference(qvmItem))

	fmt.Println(qvmItem.Difference(boItem))

}

func test() {
	// kide := mapset.NewSet()
	// kide.Add("xiaorui.cc")
	// kide.Add("blog.xiaorui.cc")
	// kide.Add("vps.xiaorui.cc")
	// kide.Add("linode.xiaorui.cc")

	// special := []interface{}{"Biology", "Chemistry"}
	// scienceClasses := mapset.NewSetFromSlice(special)

	// address := mapset.NewSet()
	// address.Add("beijing")
	// address.Add("nanjing")
	// address.Add("shanghai")

	// bonusClasses := mapset.NewSet()
	// bonusClasses.Add("Go Programming")
	// bonusClasses.Add("Python Programming")

	// //一个并集的运算
	// allClasses := kide.Union(scienceClasses).Union(address).Union(bonusClasses)
	// fmt.Println(allClasses)

	// //是否包含"Cookiing"
	// fmt.Println(scienceClasses.Contains("Cooking")) //false

	// //两个集合的差集
	// fmt.Println(allClasses.Difference(scienceClasses)) //Set{Music, Automotive, Go Programming, Python Programming, Cooking, English, Math, Welding}

	// //两个集合的交集
	// fmt.Println(scienceClasses.Intersect(kide)) //Set{Biology}

	// //有多少基数
	// fmt.Println(bonusClasses.Cardinality()) //2

}