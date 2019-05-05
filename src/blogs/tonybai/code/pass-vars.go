$tree parampass
parampass
├── etcd
│   └── etcd.go
└── main.go

// flag-demo/parampass/etcd/etcd.go
package etcd
... ...

func EtcdProxy(endpoints, user, password string) {
    fmt.Println(endpoints, user, password)
}

// flag-demo/parampass/main.go
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

... ...

func main() {
    flag.Usage = usage
    flag.Parse()

    go etcd.EtcdProxy(endpoints, user, password)

    time.Sleep(5 * time.Second)
}
