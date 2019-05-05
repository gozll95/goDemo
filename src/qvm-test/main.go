package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/astaxie/beego"
	"github.com/zhu/qvm/server/utils/multi_tag"
)

func init() {
	beego.SetLogFuncCall(true)
}

type Test struct 
	QVMT                      string                         `ason:"T" json:"t"`
	MasterSlaveBackendServers []MasterSlaveBackendServerType `ason:"MasterSlaveBackendServers" json:"master_slave_backend_servers"`
}

type B struct{
	Func()
	 []MasterSlaveBackendServerType


}
type MasterSlaveBackendServerType struct {
	ServerId   string       `json:"server_id"`   //后端服务器名称 ID，为 ECS 实例 ID。
	Port       int          `json:"port"`        //后端服务器使用的端口。
	Weight     int          `json:"weight"`      //权重
	ServerType MSServerType `json:"server_type"` //取值为 Master 或 Slave。默认值为 Master。
}

type MSServerType string

const (
	MasterServerType MSServerType = "Master"
	SlaveServerType  MSServerType = "Slave"
)

func ToQuery(params interface{}, tag string) (url.Values, error) {
	var (
		paramsMap map[string]interface{}
		paramUrl  url.Values
	)

	paramUrl = url.Values{}

	body, err := multi_tag.Marshal(params, tag)
	if err != nil {
		beego.Error(err)
		return nil, err
	}
	fmt.Printf("=====\n")
	fmt.Println(string(body))
	fmt.Printf("=====\n")

	err = json.Unmarshal(body, &paramsMap)
	if err != nil {
		beego.Error(err)
		return nil, err
	}
	beego.Info(paramsMap)

	for k, v := range paramsMap {
		if arr, ok := v.([]interface{}); ok {
			body, err := multi_tag.Marshal(arr, tag)
			if err != nil {
				beego.Error(err)
				return nil, err
			}
			paramUrl.Set(k, string(body))

			// queryMap, err := convertArrayToMap(arr, k)
			// if err != nil {
			// 	beego.Error(err)
			// 	return nil, err
			// }
			// addToQuery(paramUrl, queryMap)

		} else {
			paramUrl.Set(k, fmt.Sprint(v))
		}
	}

	return paramUrl, nil
}

func convertArrayToMap(arrays []interface{}, prefix string) (map[string]interface{}, error) {
	queryMap := map[string]interface{}{}

	// for i, arr := range arrays {

	// 	v, ok := arr.(map[string]interface{})
	// 	if !ok {
	// 		return nil, fmt.Errorf("%v is not type of map[string]interface{}", arr)
	// 	}
	// 	for key, val := range v {
	// 		var buf bytes.Buffer
	// 		buf.WriteString(prefix)
	// 		buf.WriteString(".")
	// 		buf.WriteString(fmt.Sprint(i + 1))
	// 		buf.WriteString(".")
	// 		buf.WriteString(key)
	// 		queryMap[buf.String()] = val
	// 	}
	// }

	for i, arr := range arrays {
		v, ok := arr.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("%v is not type of map[string]interface{}", arr)
		}
		for key, val := range v {
			var buf bytes.Buffer
			buf.WriteString(prefix)
			buf.WriteString(".")
			buf.WriteString(fmt.Sprint(i + 1))
			buf.WriteString(".")
			buf.WriteString(key)
			queryMap[buf.String()] = val
		}
	}

	return queryMap, nil
}

func addToQuery(dst url.Values, src map[string]interface{}) {
	for k, v := range src {
		dst.Set(k, fmt.Sprint(v))
	}
}

func main() {
	test := Test{
		QVMT: "test",
	}
	a := MasterSlaveBackendServerType{
		ServerId:   "i-2ze0o0vnaqm1zcueycfw",
		Port:       8080,
		Weight:     100,
		ServerType: MasterServerType,
	}
	b := MasterSlaveBackendServerType{
		ServerId:   "i-2ze0dsjl58fv8wbx7evg",
		Port:       8080,
		Weight:     100,
		ServerType: SlaveServerType,
	}
	test.MasterSlaveBackendServers = append(test.MasterSlaveBackendServers, a)
	test.MasterSlaveBackendServers = append(test.MasterSlaveBackendServers, b)

	//fmt.Println(string(c))

	utlValues, err := ToQuery(test, "ason")
	if err != nil {
		beego.Error(err)
	}
	fmt.Println(utlValues)
}
