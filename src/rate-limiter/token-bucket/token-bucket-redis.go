package util

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/cache"
	_ "github.com/astaxie/beego/cache/redis"
	"github.com/garyburd/redigo/redis"
)

func CheckRateLimit(ip, request, action string) bool {
	current := int(time.Now().Unix())
	currentStr := strconv.Itoa(current)
	//limit  100次
	//timeset 600秒
	//限制600秒最多访问100次
	limit, timeset := GetRateLimitConfig()
	allowanceStr, timestampStr := LoadAllowance(ip, request, action)
	allowance, _ := strconv.Atoi(allowanceStr)
	timestamp, _ := strconv.Atoi(timestampStr)
	allowance += int(current-timestamp) * limit / timeset
	if allowance > limit {
		allowance = limit
	}

	if allowance < 1 {
		SaveAllowance(ip, request, action, "0", currentStr)
		//返回true 代表速率超过,进行错误输出
		return true
	} else {
		allowanceStr = strconv.Itoa(allowance - 1)
		SaveAllowance(ip, request, action, allowanceStr, currentStr)
		//返回false 代表速率未超过
		return false
	}
}

func LoadAllowance(ip, request, action string) (allowance, timestamp string) {
	redisConfig := getRedisConfig()
	rs, err := cache.NewCache("redis", redisConfig)
	if err != nil {
		fmt.Println(err)
		return
	}
	res, _ := (redis.String(rs.Get(ip+"_"+request), err))
	if len(res) == 0 {
		currentStr := string(time.Now().Unix())
		defaultLimitInt, _ := GetRateLimitConfig()
		defaultLimitStr := strconv.Itoa(defaultLimitInt)
		allowance, timestamp = defaultLimitStr, currentStr
	} else {
		kv := strings.Split(res, "-")
		allowance, timestamp = kv[0], kv[1]
	}
	return

}

func GetRateLimitConfig() (limit, timeset int) {
	limit, _ = beego.AppConfig.Int("rateLimit")
	timeset, _ = beego.AppConfig.Int("rateTimeset")
	return
}

func SaveAllowance(ip, request, action, allowance, current string) {
	redisConfig := getRedisConfig()
	rs, err := cache.NewCache("redis", redisConfig)
	if err != nil {
		fmt.Println(err)
		return
	}
	rs.Put(ip+"_"+request, allowance+"-"+current, 600*time.Second)
}

func getRedisConfig() string {
	env := beego.AppConfig.String("runmode")
	conn := beego.AppConfig.String("redisConn")
	pass := beego.AppConfig.String("redisPass")
	name := beego.AppConfig.String("redisName")
	redisHash := make(map[string]interface{})
	redisHash["conn"] = conn
	redisHash["key"] = name
	if env == "prod" {
		redisHash["password"] = pass
	}
	redisConfig, _ := json.Marshal(redisHash)
	return string(redisConfig)
}

/*
哇,这个是令牌桶
*/
