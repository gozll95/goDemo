package main

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"hash"
	"io"
	"net/http"
	"time"
)

var accessKeyId string = "6MKOqxGiGU4AUk44"
var accessKeySecret string = "ufu7nS8kS59awNihtjSonMETLI0KLy"
var host string = "http://post-test.oss-cn-hangzhou.aliyuncs.com"
var expire_time int64 = 60
var upload_dir string = "user-dir/"

const (
	base64Table = "123QRSTUabcdVWXYZHijKLAWDCABDstEFGuvwxyzGHIJklmnopqr234560178912"
)

var coder = base64.NewEncoding(base64Table)

func base64Encode(src []byte) []byte {
	return []byte(coder.EncodeToString(src))
}

func get_gmt_iso8601(expire_end int64) string {
	var tokenExpire = time.Unix(expire_end, 0).Format("2006-01-02T15:04:05Z")
	return tokenExpire
}

type ConfigStruct struct {
	Expiration string     `json:"expiration"`
	Conditions [][]string `json:"conditions"`
}

type PolicyToken struct {
	AccessKeyId string `json:"accessid"`
	Host        string `json:"host"`
	Expire      int64  `json:"expire"`
	Signature   string `json:"signature"`
	Policy      string `json:"policy"`
	Directory   string `json:"dir"`
}

func get_policy_token() string {
	now := time.Now().Unix()
	expire_end := now + expire_time
	var tokenExpire = get_gmt_iso8601(expire_end)

	//create post policy json
	var config ConfigStruct
	config.Expiration = tokenExpire
	var condition []string
	condition = append(condition, "starts-with")
	condition = append(condition, "$key")
	condition = append(condition, upload_dir)
	config.Conditions = append(config.Conditions, condition)

	//calucate signature
	result, err := json.Marshal(config)
	debyte := base64.StdEncoding.EncodeToString(result)
	h := hmac.New(func() hash.Hash { return sha1.New() }, []byte(accessKeySecret))
	io.WriteString(h, debyte)
	signedStr := base64.StdEncoding.EncodeToString(h.Sum(nil))

	var policyToken PolicyToken
	policyToken.AccessKeyId = accessKeyId
	policyToken.Host = host
	policyToken.Expire = expire_end
	policyToken.Signature = string(signedStr)
	policyToken.Directory = upload_dir
	policyToken.Policy = string(debyte)
	response, err := json.Marshal(policyToken)
	if err != nil {
		fmt.Println("json err:", err)
	}
	return string(response)
}

func hello(w http.ResponseWriter, r *http.Request) {
	response := get_policy_token()
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	io.WriteString(w, response)
}

func main() {
	http.HandleFunc("/", hello)
	http.ListenAndServe(":9999", nil)
}
