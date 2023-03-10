package main

import (
	"flag"
	"fmt"
	"github.com/larspensjo/config"
	"log"
	"runtime"
)

var (
	configFile = flag.String("configfile", "config.ini", "General configuration file")
)

//topic list
var TOPIC = make(map[string]string)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	flag.Parse()

	//set config file std
	cfg, err := config.ReadDefault(*configFile)
	if err != nil {
		log.Fatalf("Fail to find", *configFile, err)
	}
	//set config file std End

	//Initialized topic from the configuration
	if cfg.HasSection("topicArr") {
		section, err := cfg.SectionOptions("topicArr")
		if err == nil {
			for _, v := range section {
				options, err := cfg.String("topicArr", v)
				if err == nil {
					TOPIC[v] = options
				}
			}
		}
	}
	//Initialized topic from the configuration END

	fmt.Println(TOPIC)
	fmt.Println(TOPIC["debug"])
}

//参考配置文件
/*
文件名：config.ini

[topicArr]
addr = 192.168.1.100
debug = true
login = LoginRequest

[other]
t1 = 0000337
t2 = admin
*/
