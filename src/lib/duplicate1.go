package main

import (
	"fmt"
)

func main() {
	a := []string{"a", "b", "a", "b"}
	b := []string{"a", "d", "c", "e"}
	fmt.Println(judgeDuplicate(a, b))
}

func judgeDuplicate(a, b []string) bool {
	m := make(map[string]map[string]bool)
	for k, va := range a {
		if val, ok := m[va]; ok {
			// m has key va
			if _, ok := val[b[k]]; ok {
				//duplicate
				return true
			}
			val[b[k]] = true
		}
		m[va] = make(map[string]bool)
		m[va][b[k]] = true
	}
	return false
}
