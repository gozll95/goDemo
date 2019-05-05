package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type Server struct {
	ServerName string `json:"zhu_server_name"`
	ServerIP   string `json:"server_ip"`
}

type Serverslice struct {
	Servers []Server `json:"servers"`
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
