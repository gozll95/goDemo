package main

import (
	"fmt"
	"regexp"
)

const text = `
my email is ccmouse@gmail.com@abc.com
email1 is abc@def.org
email2 is kkk@qqq.com 
`

func main() {
	re := regexp.MustCompile("ccmouse@gmail.com") // must,后面如果不符合正则语法会直接panic
	match := re.FindString(text)
	fmt.Println(match)
	// ccmouse@gmail.com

	re = regexp.MustCompile(
		`[a-zA-Z0-9]+@[a-zA-Z0-9]+\.[a-zA-Z0-9]+`)
	match1 := re.FindAllString(text, -1)
	fmt.Println(match1)
	// [ccmouse@gmail.com abc@def.org kkk@qqq.com]

	re = regexp.MustCompile(
		`([a-zA-Z0-9]+)@([a-zA-Z0-9]+)\.([a-zA-Z0-9]+)`)
	match2 := re.FindAllStringSubmatch(text, -1)
	fmt.Println(match2)
	//[[ccmouse@gmail.com ccmouse gmail com] [abc@def.org abc def org] [kkk@qqq.com kkk qqq com]]
	for _, m := range match2 {
		fmt.Println(m)
		/*
			[ccmouse@gmail.com ccmouse gmail com]
			[abc@def.org abc def org]
			[kkk@qqq.com kkk qqq com]
		*/
	}
}
