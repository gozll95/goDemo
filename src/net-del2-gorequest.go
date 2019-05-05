package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/datastream/aws"
	"github.com/parnurzeal/gorequest"
	"io/ioutil"
	"os"
)

var (
	//reqMethod     = flag.String("X", "GET", "Request method")
	//reqBodyFile   = flag.String("f", "", "file include Post body")
	config = flag.String("c", "s4.json", "aws signv4 config")
	//reqBodyString = flag.String("b", "", "request body")
	//debug         = flag.String("v", "", "debug")
)

func ReadConfig(file string) (map[string]string, error) {
	configFile, err := os.Open(file)
	config, err := ioutil.ReadAll(configFile)
	if err != nil {
		return nil, err
	}
	configFile.Close()
	setting := make(map[string]string)
	if err = json.Unmarshal(config, &setting); err != nil {
		return nil, err
	}
	return setting, err
}

func main() {
	flag.Parse()
	S4setting, err := ReadConfig(*config)
	var s *sign4.Signature
	s = nil
	if err == nil {
		s = &sign4.Signature{
			AccessKey: S4setting["access_keys"],
			SecretKey: S4setting["secret_keys"],
			Region:    S4setting["region"],
			Service:   S4setting["service"],
		}
	}
	request := gorequest.New()
	superagent := request.Delete("http://127.0.0.1:8889/api/v1/cluster/office/user")
	reqest, _ := superagent.MakeRequest()
	fmt.Println(reqest)
	if s != nil {
		s.SignRequest(reqest)
	}
	resp, body, errs := superagent.Send(`{"username":"hzxxxxxxxxxx", "cluster":"idc"}`).End()

	/*
		resp, body, errs := request.Delete("http://127.0.0.1:8889/api/v1/cluster/office/user").
			Send(`{"username":"zhutest2", "cluster":"idc"}`).
			End()
	*/
	fmt.Println(resp, body, errs)

}
