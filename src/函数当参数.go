//遍历Map,并把遍历的值给回调函数,可以让调用者控制做任何事情

func(sm *SynchronizedMap)Each(cb func(interface{}),interface{}){
	sm.rw.RLock()
	defer sm.rw.RUnlock()
	for k, v := range sm.data {
		cb(k,v)
	}
}

sm.Each(func(k interface{}, v interface{}) {
	fmt.Println(k," is ",v)
})