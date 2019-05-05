/*
https://blog.csdn.net/webxscan/article/details/78818070
使用redis实现信号量
我这里用了redis的hash、list、string实现的网络信号量。
利用redis hash存token的时间;
list用来实现堵塞及非阻塞模式;
利用redis string的setnx来实现分布式锁，锁的目的是用来原子性写一致。
*/
package go_redis_semaphore

import (
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/garyburd/redigo/redis"
)

const (
	PREFIX_KYE = "redis_semaphore_"
	VERSION    = "v0.1"
)

type Semaphore struct {
	Limit           int
	InitLockTimeout int

	ScanLock     *sync.RWMutex
	ScanInterval int // 每 多少秒 扫描一次
	ScanTimeout  int // 多少算超时
	LastScanTs   time.Time

	RedisClient *redis.Pool

	NameSpace       string
	QueueName       string
	LockName        string
	TokenTsHashName string

	Tokens []string
}

// NameSpace 推荐使用 自增的version版本号，不然会出现更改limit出现阈值不准的问题.
func NewRedisSemaphore(redis_client *redis.Pool, limit int, namespace string) *Semaphore {
	return &Semaphore{
		Limit:           limit,
		ScanLock:        new(sync.RWMutex),
		InitLockTimeout: 30,
		ScanInterval:    5,
		ScanTimeout:     3,
		RedisClient:     redis_client,
		NameSpace:       namespace,
		QueueName:       namespace + "_" + "queue",
		LockName:        namespace + "_" + "lock",
		TokenTsHashName: namespace + "_" + "hash",
		LastScanTs:      time.Now(),
	}
}

func (s *Semaphore) Init() {
	rc := s.RedisClient.Get()
	defer rc.Close()

	var tmp_token string

	ok, _ := s.TryLock(0)
	if !ok {
		fmt.Println("lock failed")
		return
	}

	// clean old token list
	// to do: pipeline
	rc.Do("DEL", s.QueueName)
	for i := 1; i <= s.Limit; i++ {
		tmp_token = fmt.Sprintf("token_seq_%d", i)
		s.Push(tmp_token)
		s.Tokens = append(s.Tokens, tmp_token)
	}

	// del lock
	// rc.Do("DEL", s.LockName)
}

func (s *Semaphore) ScanIsContinue() bool {
	now := time.Now()
	if int(now.Sub(s.LastScanTs).Seconds()) < s.ScanInterval {
		return false
	}
	return true
}

func (s *Semaphore) ScanTimeoutToken() []string {
	rc := s.RedisClient.Get()
	defer rc.Close()

	expire_tokens := []string{}

	s.ScanLock.Lock()

	if !s.ScanIsContinue() {
		s.ScanLock.Unlock()
		return expire_tokens
	}

	res, _ := redis.StringMap(rc.Do("HGETALL", s.TokenTsHashName))

	for token, ts_s := range res {
		ts, _ := strconv.Atoi(ts_s)
		diff_ts := time.Now().Sub(time.Unix(int64(ts), 0))
		if int(diff_ts.Seconds()) > s.ScanTimeout {
			expire_tokens = append(expire_tokens, token)
		}
	}
	s.LastScanTs = time.Now()
	s.ScanLock.Unlock()

	return expire_tokens
}

func (s *Semaphore) TryLock(timeout int) (bool, error) {
	rc := s.RedisClient.Get()
	defer rc.Close()

	var err error

	if timeout == 0 {
		_, err = redis.String(rc.Do("SET", s.LockName, "locked", "NX"))
	} else {
		_, err = redis.String(rc.Do("SET", s.LockName, "locked", "EX", s.InitLockTimeout, "NX"))
	}

	if err == redis.ErrNil {
		return false, nil
	}

	if err != nil {
		return false, err
	}
	return true, nil
}

func (s *Semaphore) Acquire(timeout int) (string, error) {
	var token string
	var err error

	s.ScanTimeoutToken()
	if timeout > 0 {
		token, err = s.PopBlock(timeout)
	} else {
		token, err = s.Pop()
	}
	return token, err
}

func (s *Semaphore) Release(token string) {
	s.Push(token)
}

func (s *Semaphore) Pop() (string, error) {
	rc := s.RedisClient.Get()
	defer rc.Close()

	res, err := redis.String(rc.Do("LPOP", s.QueueName))
	// 允许队列为空值
	if err == redis.ErrNil {
		err = nil
	}

	rc.Do("HSET", s.TokenTsHashName, res, time.Now().Unix())
	return res, err
}

func (s *Semaphore) Push(body string) (int, error) {
	rc := s.RedisClient.Get()
	defer rc.Close()

	// to do: pipeline
	res, err := redis.Int(rc.Do("RPUSH", s.QueueName, body))
	rc.Do("HDEL", s.TokenTsHashName, body)
	return res, err
}

func (s *Semaphore) PopBlock(timeout int) (string, error) {
	rc := s.RedisClient.Get()
	defer rc.Close()

	// refer: https://gowalker.org/github.com/BPing/Golib/cache/mredis#RedisPool_BLPop
	res_map, err := redis.StringMap(rc.Do("BLPOP", s.QueueName, timeout))
	// 允许队列为空值
	if err == redis.ErrNil {
		err = nil
	}

	res, ok := res_map[s.QueueName]
	if res != "" {
		rc.Do("HSET", s.TokenTsHashName, res, time.Now().Unix())
	}

	if !ok {
		return "", err
	}

	return res, err
}

/*
用法:
limiter := go_redis_semaphore.NewRedisSemaphore(redis连接池对象, 信号量数目, "redis key的名字")

// 初始化redis
limiter.Init()

// 非阻塞拿锁
token, _ := limiter.Acquire(0)

// 超时5s阻塞拿锁
token, _ := limiter.Acquire(5)

// 释放锁
limiter.Release(token)





package main
// xiaorui.cc

import (
	"fmt"
	"github.com/rfyiamcool/go_redis_semaphore"
)

func main() {
	fmt.Println("实例化redis连接池")
	redis_client_config := go_redis_semaphore.RedisConfType{
		RedisPw:          "",
		RedisHost:        "127.0.0.1:6379",
		RedisDb:          0,
		RedisMaxActive:   100,
		RedisMaxIdle:     100,
		RedisIdleTimeOut: 1000,
	}
	redis_client := go_redis_semaphore.NewRedisPool(redis_client_config)

	fmt.Println("实例化 redis Semaphore")
	limiter := go_redis_semaphore.NewRedisSemaphore(redis_client, 2, "love")
	limiter.Init()

	fmt.Println("非阻塞拿锁")
	token, _ := limiter.Acquire(0)

	fmt.Println("释放锁")
	limiter.Release(token)
	fmt.Println(limiter.ScanTimeoutToken())
	fmt.Println("end")
}
*/
