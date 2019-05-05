// Go1.9 安全map使用
// sync.map就是1.9版本带的线程安全map，主要有:

// Store
// LoadOrStore
// Load
// Delete
// Range
// 这几种用法

// Store(key, value interface{})
// 说明： 存储一个设置的键值。

// LoadOrStore(key, value interface{}) (actual interface{}, loaded bool)
// 说明： 返回键的现有值(如果存在)，否则存储并返回给定的值，如果是读取则返回true，如果是存储返回false。

// Load(key interface{}) (value interface{}, ok bool)
// 说明： 读取存储在map中的值，如果没有值，则返回nil。OK的结果表示是否在map中找到值。

// Delete(key interface{})
// 说明： 删除键对应的值。

// Range(f func(key, value interface{}) bool)
// 说明： 循环读取map中的值。

// 代码如下：
package main

import (
	"fmt"
	"sync"
)

type userInfo struct {
	Name string
	Age  int
}

var m sync.Map

func main() {

	vv, ok := m.LoadOrStore("1", "one")
	fmt.Println(vv, ok) //one false

	vv, ok = m.Load("1")
	fmt.Println(vv, ok) //one true

	vv, ok = m.LoadOrStore("1", "oneone")
	fmt.Println(vv, ok) //one true

	vv, ok = m.Load("1")
	fmt.Println(vv, ok) //one true

	m.Store("1", "oneone")
	vv, ok = m.Load("1")
	fmt.Println(vv, ok) // oneone true

	m.Store("2", "two")
	m.Range(func(k, v interface{}) bool {
		fmt.Println(k, v)
		return true
	})

	m.Delete("1")
	m.Range(func(k, v interface{}) bool {
		fmt.Println(k, v)
		return true
	})

	map1 := make(map[string]userInfo)
	var user1 userInfo
	user1.Name = "ChamPly"
	user1.Age = 24
	map1["user1"] = user1

	var user2 userInfo
	user2.Name = "Tom"
	user2.Age = 18
	m.Store("map_test", map1)

	mapValue, _ := m.Load("map_test")

	for k, v := range mapValue.(interface{}).(map[string]userInfo) {
		fmt.Println(k, v)
		fmt.Println("name:", v.Name)
	}
}
