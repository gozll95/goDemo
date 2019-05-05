package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type Server struct {
	ServerName string
	ServerIP   string
}

type Serverslice struct {
	Servers []Server
}

func main() {
	var s Serverslice
	s.Servers = append(s.Servers, Server{ServerName: "Shanghai_VPN", ServerIP: "127.0.0.1"})
	s.Servers = append(s.Servers, Server{ServerName: "Beijing_VPN", ServerIP: "127.0.0.2"})
	b, err := json.Marshal(s)
	if err != nil {
		fmt.Println("json err:", err)
	}
	fmt.Println(string(b))

	var out bytes.Buffer
	err = json.Indent(&out, b, "", "\t")

	if err != nil {
		log.Fatalln(err)
	}

	out.WriteTo(os.Stdout)
}

//json.NewDecoder(configFile).Decode(&setting)

/*
func main() {
	//data, err := json.Marshal(movies)
	/*
		为了生成便于阅读的格式，另一个json.MarshalIndent函数将产生整齐缩进的输出。
		该函数有两个额外的字符串参数用于表示每一行输出的前缀和每一个层级的缩进：
	*/
	data, err := json.MarshalIndent(movies, "", "    ")
	if err != nil {
		log.Fatalf("JSON marshaling failed: %s", err)
	}
	fmt.Printf("%s\n", data)
}
*/