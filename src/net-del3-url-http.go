package main

//bad
import (
	"bytes"
	"fmt"
	"io/ioutil"
	//"log"
	"net/http"
	"net/url"
)

func main() {
	fmt.Println("hello1")
	u := "http://127.0.0.1:8889/api/v1/cluster/office/user"
	data := url.Values{}
	data.Set("username", "hzxxxxxxxxxx")
	data.Set("cluster", "idc")
	b := bytes.NewBufferString(data.Encode())
	fmt.Println("hello2")
	req, err := http.NewRequest("DELETE", u, b)
	fmt.Println("hello3")
	req.Header.Set("Content-Type", "application/json")
	fmt.Println("hello4")
	if err != nil {
		fmt.Println(err)
		fmt.Println("hello5")
		//log.Fatal(err)
	}
	c := &http.Client{}
	resp, err := c.Do(req)
	defer resp.Body.Close()

	//r, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		//log.Fatal(err)
	}

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))
	/*
		url := "http://127.0.0.1:8889/api/v1/cluster/office/user"
		fmt.Println("URL:>", url)

		var jsonStr = []byte(`{"username":"hzxxxxxxxxxx", "cluster":"idc"}`)
		//buf1 := bytes.New
		req, err := http.NewRequest("DELETE", url, bytes.NewBuffer(jsonStr))
		//req.Header.Set("X-Custom-Header", "myvalue")
		//req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

		fmt.Println("response Status:", resp.Status)
		fmt.Println("response Headers:", resp.Header)
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Println("response Body:", string(body))
	*/
}
