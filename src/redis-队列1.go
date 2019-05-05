//http://blog.csdn.net/xcl168/article/details/44590457

//Redis做后台任务队列
//author: Xiong Chuan Liang
//date: 2015-3-25

package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/garyburd/redigo/redis"
)

func main() {

	r, err := newRedisPool("192.168.20.25:6379", "")
	if err != nil {
		fmt.Println(err)
		return
	}

	//将job放入队列
	if err := r.Enqueue(); err != nil {
		fmt.Println(err)
	}

	//依次取出两个Job
	r.GetJob()
	r.GetJob()
}

type RedisPool struct {
	pool *redis.Pool
}

func newRedisPool(server, password string) (*RedisPool, error) {

	if server == "" {
		server = ":6379"
	}

	pool := &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", server)
			if err != nil {
				return nil, err
			}

			if password != "" {
				if _, err := c.Do("AUTH", password); err != nil {
					c.Close()
					return nil, err
				}
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}

	return &RedisPool{pool}, nil
}

type Job struct {
	Class string        `json:"Class"`
	Args  []interface{} `json:"Args"`
}

//模拟客户端
func (r *RedisPool) Enqueue() error {

	c := r.pool.Get()
	defer c.Close()

	j := &Job{}
	j.Class = "mail"
	j.Args = append(j.Args, "xcl_168@aliyun.com", "", "body", 2, true)

	j2 := &Job{}
	j2.Class = "Log"
	j2.Args = append(j2.Args, "ccc.log", "ddd.log", []int{222, 333})

	for _, v := range []*Job{j, j2} {
		b, err := json.Marshal(v)
		if err != nil {
			return err
		}

		_, err = c.Do("rpush", "queue", b)
		if err != nil {
			return err
		}
	}

	fmt.Println("[Enqueue()] succeed!")

	return nil
}

//模拟Job Server
func (r *RedisPool) GetJob() error {
	count, err := r.QueuedJobCount()
	if err != nil || count == 0 {
		return errors.New("暂无Job.")
	}
	fmt.Println("[GetJob()] Jobs count:", count)

	c := r.pool.Get()
	defer c.Close()

	for i := 0; i < int(count); i++ {
		reply, err := c.Do("LPOP", "queue")
		if err != nil {
			return err
		}

		var j Job
		decoder := json.NewDecoder(bytes.NewReader(reply.([]byte)))
		if err := decoder.Decode(&j); err != nil {
			return err
		}

		fmt.Println("[GetJob()] ", j.Class, " : ", j.Args)
	}
	return nil
}

func (r *RedisPool) QueuedJobCount() (int, error) {
	c := r.pool.Get()
	defer c.Close()

	lenqueue, err := c.Do("llen", "queue")
	if err != nil {
		return 0, err
	}

	count, ok := lenqueue.(int64)
	if !ok {
		return 0, errors.New("类型转换错误!")
	}
	return int(count), nil
}

/*
运行结果:

[Enqueue()] succeed!
[GetJob()] Jobs count: 2
[GetJob()]  mail  :  [xcl_168@aliyun.com  body 2 true]
[GetJob()]  Log  :  [ccc.log ddd.log [222 333]]
[root@xclos src]#

*/
