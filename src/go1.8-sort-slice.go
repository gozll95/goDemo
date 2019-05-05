package main

import (
	"fmt"
	"sort"
)

type Lang struct {
	Name string
	Rank int
}

func main() {
	langs := []Lang{
		{"rust", 2},
		{"go", 1},
		{"swift", 3},
	}
	sort.Slice(langs, func(i, j int) bool { return langs[i].Rank < langs[j].Rank })
	fmt.Printf("%v\n", langs)

}

//[{go 1} {rust 2} {swift 3}]

/*
实现sort，需要三要素：Len、Swap和Less。在1.8之前，我们通过实现sort.Interface实现了这三个要素；而在1.8版本里，Slice函数通过reflect获取到swap和length，通过结合闭包实现的less参数让Less要素也具备了。我们从下面sort.Slice的源码可以看出这一点：
// $GOROOT/src/sort/sort.go
... ...
func Slice(slice interface{}, less func(i, j int) bool) {
    rv := reflect.ValueOf(slice)
    swap := reflect.Swapper(slice)
    length := rv.Len()
    quickSort_func(lessSwap{less, swap}, 0, length, maxDepth(length))
}
*/