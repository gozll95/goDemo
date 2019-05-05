package main

import (
	"time"
)

type a struct {
	timeourNA  time.Duration //响应超时时间,单位:ns
	Ips        uint32        //每秒载荷量
	durationNS time.Duration //负载持续时间,单位:ns
}
