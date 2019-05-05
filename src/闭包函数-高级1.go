package main

import "fmt"

func main() {
	a := test()
	a()
	a()
}

type Token struct {
	TokenExpiry string
}

func test() func() string {
	token := &Token{}
	return func() string {
		fmt.Println("call test token is ", token.TokenExpiry)
		//*token = *tk
		token = &Token{
			TokenExpiry: "bb",
		}
		return token.TokenExpiry
	}
}
