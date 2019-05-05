package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println(time.Now().Unix())
	now := time.Now()
	sendTime, _ := time.Parse("2006-01-02 15:04:05", "2013-06-07 13:51:39")
	fmt.Println(sendTime)
	println(now.Format("2006-01-02T15:04:05+0800"), sendTime.Format("2006-01-02 15:04:05"))
	println(now.Unix(), sendTime.Unix())
	println(now.After(sendTime))

	time.Now().String()

}

func compare() {
	time1 := "2016/09/28-06:40"
	time3 := "2016-03-20"
	//先把时间字符串格式化成相同的时间类型
	t1, err := time.Parse("2016/09/28-06:40", time1)
	t3, err := time.Parse("2006-01-02", time3)

	d, _ := time.ParseDuration("-8h")
	//t2 := t1.Add(d * 7)
	t2 := t1.Add(d)

	fmt.Println(err)
	t3.Before(t2)

	if err == nil && t2.Before(t3) {
		//处理逻辑
		fmt.Println("true")
	} else {
		fmt.Println("false")
	}
	time.Local()
}

func a() {
	withNanos := "2006-01-02 15:04:05"
	t, _ := time.Parse(withNanos, "2013-10-05 18:30:50")
	fmt.Println(t.Year())

	expireTime, _ := time.ParseDuration("-3s")
	fmt.Println(time.Now().Add(expireTime))

}
