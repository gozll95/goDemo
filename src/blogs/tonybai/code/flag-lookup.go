$tree flaglookup
flaglookup
├── etcd
│   └── etcd.go
└── main.go

// flag-demo/flaglookup/main.go
package main

import (
    "flag"
    "fmt"
    "time"

    "./etcd"
)

var (
    endpoints string
    user      string
    password  string
)

func init() {
    flag.StringVar(&endpoints, "endpoints", "127.0.0.1:2379", "comma-separated list of etcdv3 endpoints")
    flag.StringVar(&user, "user", "", "etcdv3 client user")
    flag.StringVar(&password, "password", "", "etcdv3 client password")
}

......

func main() {
    flag.Usage = usage
    flag.Parse()

    go etcd.EtcdProxy()

    time.Sleep(5 * time.Second)
}

// flag-demo/flaglookup/etcd/etcd.go
package etcd

import (
    "flag"
    "fmt"
)

func EtcdProxy() {
    endpoints := flag.Lookup("endpoints").Value.(flag.Getter).Get().(string)
    user := flag.Lookup("user").Value.(flag.Getter).Get().(string)
    password := flag.Lookup("password").Value.(flag.Getter).Get().(string)

    fmt.Println(endpoints, user, password)
}

运行该程序：
$go run main.go -endpoints 192.168.10.69:2379,10.10.12.36:2378 -user tonybai -password xyz123
192.168.10.69:2379,10.10.12.36:2378 tonybai xyz123
输出与我们的预期是一致的。