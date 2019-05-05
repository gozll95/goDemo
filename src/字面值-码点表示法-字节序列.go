package main

import (
	"fmt"
)

func main() {
	// 字面值表示法
	var s = "中国人"
	// 码点表示法
	var s2 = "\U00004e2d\U000056fd\U00004eba"
	// 字节序列表示法(二进制表示法)
	var s3 = "\xe4\xb8\xad\xe5\x9b\xbd\xe4\xba\xba"

	fmt.Println("s byte sequence:")
	for i := 0; i < len(s); i++ {
		fmt.Printf("0x%x", s[i])
	}
	fmt.Println("")

	fmt.Println("s byte sequence:")
	for i := 0; i < len(s2); i++ {
		fmt.Printf("0x%x", s[i])
	}
	fmt.Println("")

	fmt.Println("s byte sequence:")
	for i := 0; i < len(s3); i++ {
		fmt.Printf("0x%x", s[i])
	}

	fmt.Println("")
}

/*
s byte sequence:
0xe40xb80xad0xe50x9b0xbd0xe40xba0xba
s byte sequence:
0xe40xb80xad0xe50x9b0xbd0xe40xba0xba
s byte sequence:
0xe40xb80xad0xe50x9b0xbd0xe40xba0xba
*/
