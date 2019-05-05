package main

import (
	"fmt"
)

func main() {
	s := make([]int, 1, 3)

	fmt.Printf("%p %v \n", s, s[0])
	// fmt.Println(s[1]) 越界
	s = append(s, 2)
	fmt.Printf("%p %v \n", s, s[1]) // 注意这里的指针和上面的指针一样，因为没有超出他的容量

	s = append(s, 3, 4)
	fmt.Printf("%p  \n", s) // 这里指针变化。前面追加的时候，已经超过了容量，会重新生成一个新的slice
	s1 := s[1:]
	fmt.Println(s1)
	s1[0] = 1123
	fmt.Println("after changed :: ", s1, " \t s0 : ", s)

	idx := 2
	value := 44

	//在任意位置插入任意数据
	s2 := make([]int, len(s[:idx]))

	copy(s2, s[:idx])

	s = append(append(s2, value), s[idx:]...)

	fmt.Println("after changed :: ", s1, " \t s0 : ", s)

	//  s2 := make([]T, len(s[:idx]))
	// copy(s2, s[:idx])
	// s = append(append(s2, value ),s[idx:]...))

}

// 0xc42000a200 0
// 0xc42000a200 2
// 0xc420016120
// [2 3 4]
// after changed ::  [1123 3 4]     s0 :  [0 1123 3 4]
// after changed ::  [1123 3 4]     s0 :  [0 1123 44 3 4]

// 切分会生成一个共享源数据的slice。
