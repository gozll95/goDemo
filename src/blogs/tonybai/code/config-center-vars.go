$tree configcenter
configcenter
├── config
│   └── config.go
└── main.go

//flag-demo/configcenter/config/config.go
package config

import (
    "log"
    "sync"
)

var (
    m  map[string]interface{}
    mu sync.RWMutex
)

func init() {
    m = make(map[string]interface{}, 10)
}

func SetString(k, v string) {
    mu.Lock()
    m[k] = v
    mu.Unlock()
}

func SetInt(k string, i int) {
    mu.Lock()
    m[k] = i
    mu.Unlock()
}

func GetString(key string) string {
    mu.RLock()
    defer mu.RUnlock()
    v, ok := m[key]
    if !ok {
        return ""
    }
    return v.(string)
}

func GetInt(key string) int {
    mu.RLock()
    defer mu.RUnlock()
    v, ok := m[key]
    if !ok {
        return 0
    }
    return v.(int)
}

func Dump() {
    log.Println(m)
}

// flag-demo/configcenter/main.go

package main

import (
    "flag"
    "fmt"
    "time"

    "./config"
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

    // inject flag vars to config center
    config.SetString("endpoints", endpoints)
    config.SetString("user", user)
    config.SetString("password", password)

    time.Sleep(5 * time.Second)
}