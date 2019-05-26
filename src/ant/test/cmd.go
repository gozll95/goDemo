req, err := http.ReadRequest(bufio.NewReader(strings.NewReader(v)))


go test -v -race

// 竞态条件fix
// 1、使用channel
// 2、使用Mutex
// 3、使用atomic

go test -v -run=^$ -bench=.

go test -v -run=^$ -bench=^BenchmarkHi$ -benchtime=2s -cpuprofile=prof.cpu

go tool pprof step2.test prof.cpu

top –cum

list handleHi

go test -v -run=^$ -bench=^BenchmarkHi$ -benchtime=2s -memprofile=prof.mem

go tool pprof –alloc_space step3.test prof.mem


go test -bench=. -memprofile=prof.mem | tee mem.2

go test -bench=. -memprofile=prof.mem | tee mem.3

benchcmp step3/mem.3 step4/mem.4

Go type 在 runtime的内存占用:

A Go interface is 2 words of memory: (type, pointer).
A Go string is 2 words of memory: (base pointer, length)
A Go slice is 3 words of memory: (base pointer, length, capacity)