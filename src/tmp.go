package main

import (
	"fmt"
	"io/ioutil"
	"net/url"
)

func main() {
	u := "http://127.0.0.1:8889/api/v1/cluster/office/user"
	data := url.Values{}
	data.Set("username", "hzxxxxxxxxxx")
	data.Set("cluster", "idc")
}
