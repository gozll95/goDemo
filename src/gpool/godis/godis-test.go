package main

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"gotest/godis"
	"time"
)

var (
	Godispool *godis.GodisPool
)

func main() {
	Godispool = &godis.GodisPool{
		MaxIdle:       5,
		MaxActive:     6,
		ZkServerList:  []string{"127.0.0.1:2181"}, //10.20.30.91
		ZkDir:         "/serverlist",
		IdleTimeout:   (time.Duration(3) * time.Second),
		ZkConnTimeout: 3 * time.Second,
		Wait:          true,
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if _, err := c.Do("PING"); err != nil {
				return err
			}
			return nil
		},
	}
	Godispool.InitPool()

	gConn := Godispool.Get()

	_, sErr := gConn.Do("set", "aaa", "aaaaaaaaa")

	if sErr != nil {
		fmt.Println("set err -: ", sErr.Error())
	}

	rs, err := redis.String(gConn.Do("get", "aaa"))

	if err != nil {
		fmt.Println("get err:", err.Error())
	}

	fmt.Println("rs--", rs)
	//
	//time.Sleep(time.Second)
	//
	//CConn := Godispool.Get()
	//
	//fmt.Println("----222222222222222-----")
	//rs2, err2 := redis.String(gConn.Do("get", "aaa"))
	//
	//if err2 != nil {
	//	fmt.Println("get err:", err.Error())
	//}
	//
	//fmt.Println("rs----------------------", rs2)
	//
	//defer CConn.Close()
	defer gConn.Close()
	for i := 0; i < 50; i++ {
		go testg()
		//time.Sleep(500 * time.Millisecond)
	}

	time.Sleep(1000 * time.Second)

}

func testg() {
	gConn := Godispool.Get()
	defer gConn.Close()

	rs, err := redis.String(gConn.Do("get", "aaa"))

	if err != nil {
		fmt.Println("get err:", err.Error())
	}

	fmt.Println("rs----------------------", rs)

	time.Sleep(time.Second)
	return

}
