package main

import (
	"sync"
)

var cache = struct {
	sync.Mutex
	mapping map[string]string
}{
	mapping: make(map[string]string),
}

func Lookup(key string) string {
	cache.Lock()
	v := cache.mapping[key]
	cache.Unlock()
	return v
}

/*

起初我查阅了相关问题解决方案。大致就是多线程操作map数据结构一定要加锁。否则肯定要出现这个错误。我查看我的代码，我认为我写的map结构都加了锁，附加锁方式：


通用锁
type Demo struct {
  Data map[string]string
  Lock sync.Mutex
}

func (d Demo) Get(k string) string{
  d.Lock.Lock()
  defer d.Lock.UnLock()
  return d.Data[k]
}

func (d Demo) Set(k,v string) {
  d.Lock.Lock()
  defer d.Lock.UnLock()
  d.Data[k]=v
}

读写锁
type Demo struct {
  Data map[string]string
  Lock sync.RwMutex
}

func (d Demo) Get(k string) string{
  d.Lock.RLock()
  defer d.Lock.RUnlock()
  return d.Data[k]
}

func (d Demo) Set(k,v string) {
  d.Lock.Lock()
  defer d.Lock.UnLock()
  d.Data[k]=v
}

*/
