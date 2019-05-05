package main

import (
	"fmt"
	"os"
	"sort"
)

func main() {
	m := map[int]string{1: "a", 2: "b", 3: "c"}
	s := make([]int, len(m))
	i := 0
	for k, _ := range m {
		s[i] = k
		i++
	}
	sort.Ints(s)
	fmt.Println(s)
	fmt.Fprintf(os.Stderr, "排序后的map的key为%v", s)
}
