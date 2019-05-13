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
