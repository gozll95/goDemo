package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type TagName struct {
	Name string `json:"nm"`
	Book string `json:"b"`
	user string
}

func main() {
	enc := json.NewEncoder(os.Stdout)

	tags := []TagName{
		{"hello", "Joson", "weter"},
		{"world", "Nano", "jober"},
	}
	if e := enc.Encode(tags); e != nil {
		fmt.Println("e ", e)
	}
}

/*
编码和解码流

json 包提供 Decoder 和 Encoder 类型来支持常用 JSON 数据流读写。NewDecoder 和 NewEncoder 函数分别封装了 io.Reader 和 io.Writer 接口。

func NewDecoder(r io.Reader) *Decoder
func NewEncoder(w io.Writer) *Encoder
要想把 JSON 直接写入文件，可以使用 json.NewEncoder 初始化文件（或者任何实现 io.Writer 的类型），并调用 Encode()；反过来与其对应的是使用 json.Decoder 和 Decode() 函数：

func NewDecoder(r io.Reader) *Decoder
func (dec *Decoder) Decode(v interface{}) error
来看下接口是如何对实现进行抽象的：数据结构可以是任何类型，只要其实现了某种接口，目标或源数据要能够被编码就必须实现 io.Writer 或 io.Reader 接口。由于 Go 语言中到处都实现了 Reader 和 Writer，因此 Encoder 和 Decoder 可被应用的场景非常广泛，例如读取或写入 HTTP 连接、websockets 或文件。

*/
