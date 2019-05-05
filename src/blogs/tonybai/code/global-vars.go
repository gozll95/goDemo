$tree globalvars
globalvars
├── config
│   └── config.go
├── etcd
│   └── etcd.go
└── main.go

// flag-demo/globalvars/config/config.go

package config

var (
    Endpoints string
    User      string
    Password  string
)

// flag-demo/globalvars/etcd/etcd.go
package etcd

import (
    "fmt"

    "../config"
)

func EtcdProxy() {
    fmt.Println(config.Endpoints, config.User, config.Password)
    //... ....
}

// flag-demo/globalvars/main.go
package main

import (
    "flag"
    "fmt"
    "time"

    "./config"
    "./etcd"
)

func init() {
    flag.StringVar(&config.Endpoints, "endpoints", "127.0.0.1:2379", "comma-separated list of etcdv3 endpoints")
    flag.StringVar(&config.User, "user", "", "etcdv3 client user")
    flag.StringVar(&config.Password, "password", "", "etcdv3 client password")
}

.... ...

func main() {
    flag.Usage = usage
    flag.Parse()

    go etcd.EtcdProxy()

    time.Sleep(5 * time.Second)
}