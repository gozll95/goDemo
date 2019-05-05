package main

import (
	"encoding/json"
	"fmt"
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
	str := `{"servers":[{"serverName":"Shanghai_VPN","serverIP":"127.0.0.1"},{"serverName":"Beijing_VPN","serverIP":"127.0.0.2"}]}`
	json.Unmarshal([]byte(str), &s)
	test(str,&s)
	fmt.Printf("origin:%p\n",&s)
	fmt.Println(s)
	
}


func test(str string, s interface{}){
	fmt.Printf("%T\n",s)
	fmt.Printf("%p\n",s)
	fmt.Printf("%T\n",&s)
	fmt.Printf("%p\n",&s)
	
	//beego.Error(s)
	json.Unmarshal([]byte(str), &s)
}
/*

{"Servers":[{"ServerName":"Shanghai_VPN","ServerIP":"127.0.0.1"},{"ServerName":"Beijing_VPN","ServerIP":"127.0.0.2"}]}

*/
