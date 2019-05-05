// 二分查找法

// https://stackoverflow.com/questions/38607733/sorting-a-uint64-slice-in-go
// https://blog.csdn.net/chenbaoke/article/details/42340301

package main

import (
	"fmt"
	"sort"
)

func main() {
	a := []int{1, 2, 3, 5, 4}
	sort.Slice(a, func(i, j int) bool {
		return a[i] < a[j]
	})

	fmt.Println(a)
	d := sort.Search(len(a), func(i int) bool {
		return a[i] >= 100
	})
	fmt.Println(d) //2
}

func GuessingGame() {
	var s string
	fmt.Printf("Pick an integer from 0 to 100.\n")
	answer := sort.Search(100, func(i int) bool {
		fmt.Printf("Is your number <= %d? ", i)
		fmt.Scanf("%s", &s)
		return s != "" && s[0] == 'y'
	})
	fmt.Printf("Your number is %d.\n", answer)
}
