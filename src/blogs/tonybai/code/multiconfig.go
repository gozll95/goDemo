package main 

import(
	"multiconfi"
)

func main(){
	m:=multiconfig.NewWithPath("config.toml")

	serverConf:=new(Server)

	m.MustLoad(serverConf)

	fmt.Println("After Loading: ")
	fmt.Printf("%+v\n",serverConf)

	
}