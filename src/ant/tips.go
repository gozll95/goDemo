## atomic.Value
var GlobalUser atomic.Value
GlobalUser.Store(user)   
GlobalUser.Load().(User)



## 测试用例
func BenchmarkPAtomicGet(b *testing.B) {
	var config atomic.Value
	config.Store(Config{endpoint: "api.example.com"})
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = config.Load().(Config)
		}
	})
}


/*
 * 类型是方法，挂自身实现接口,别的方法经过类型转换成我实现接口
 */

type HandlerFunc func(k, v interface{})
func (f HandlerFunc) Do(k, v interface{}){
	f(k,v)
}

type welcome string
func (w welcome) selfInfo(k, v interface{}) {
	fmt.Printf("%s,我叫%s,今年%d岁\n", w,k, v)
}
func main() {
	persons := make(map[interface{}]interface{})
	persons["张三"] = 20
	persons["李四"] = 23
	persons["王五"] = 26
	var w welcome = "大家好"
	Each(persons, HandlerFunc(w.selfInfo))
}

package main
import (
	"fmt"
)
type Handler interface {
	Do(k, v interface{})
}
type HandlerFunc func(k, v interface{})
func (f HandlerFunc) Do(k, v interface{}) {
	f(k, v)
}
func Each(m map[interface{}]interface{}, h Handler) {
	if m != nil && len(m) > 0 {
		for k, v := range m {
			h.Do(k, v)
		}
	}
}
func EachFunc(m map[interface{}]interface{}, f func(k, v interface{})) {
	Each(m, HandlerFunc(f))
}
func selfInfo(k, v interface{}) {
	fmt.Printf("大家好,我叫%s,今年%d岁\n", k, v)
}
func main() {
	persons := make(map[interface{}]interface{})
	persons["张三"] = 20
	persons["李四"] = 23
	persons["王五"] = 26
	EachFunc(persons, selfInfo)
}