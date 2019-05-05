package main

import (
	"fmt"
	"github.com/smallnest/goreq"
)

func main() {
	q := `{"username":"xxxxxxxxxx1"}`
	resp, _, err := goreq.New().Delete("http://127.0.0.1:8889/api/v1/cluster/office/user").ContentType("json").SendMapString(q).End()
	fmt.Println(resp)
}
