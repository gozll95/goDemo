package main

import (
	"encoding/json"
	"fmt"
	"github.com/koron/go-dproxy"
)

func main() {
	var v interface{}

	json.Unmarshal([]byte(`{
      "name": "dproxy",
  "strings": [ "a", "b", "c" ],
  "numbers": [ 1,2,3 ],
  "mixs": ["a", 1, "b"],

  "tags": {
    "key-1": "value-1",
    "key-2": 100,
    "key-3": [4,5,6]
  }
}`), &v)

	// s == "dproxy", got a string.
	s, _ := dproxy.New(v).M("name").String()

	fmt.Printf("1. name is %s\n", s)

	// err: not matched types: expected=int64 actual=string: name
	_, err := dproxy.New(v).M("name").Int64()
	fmt.Printf("2. err is %s\n", err)

	// can be chained. n == 3, got a int64
	n, _ := dproxy.New(v).M("numbers").A(2).Int64()
	fmt.Printf("3. n is %d\n", n)

	// got value-1
	s, _ = dproxy.New(v).M("tags").M("key-1").String()
	fmt.Printf("4. key-1 is %s\n", s)

	// err.Error() == "not found: data.kustom", wrong query can be verified.
	_, err = dproxy.New(v).M("data").M("kustom").String()
	fmt.Printf("5. err is %s\n", err)

	// n == 5
	n, err = dproxy.Pointer(v, "/tags/key-3/1").Int64()
	fmt.Printf("6. n is %d\n", n)

}
