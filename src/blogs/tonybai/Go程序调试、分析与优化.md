go test -v -race

go test -v -run=^$ -bench=.

go test -v -run=^$ -bench=^BenchmarkHi$ -benchtime=2s -cpuprofile=prof.cpu

go tool pprof step2.test prof.cpu

(pprof) top –cum

(pprof) list handleHi

从top –cum来看，handleHi消耗cpu较大，而handleHi中，又是MatchString耗时最长。

go test -v -run=^$ -bench=^BenchmarkHi$ -benchtime=3s -cpuprofile=prof.cpu

top –cum 30

go test -v -run=^$ -bench=^BenchmarkHi$ -benchtime=2s -memprofile=prof.mem

go test -v -run=^$ -bench=^BenchmarkHi$ -benchtime=2s -memprofile=prof.mem

go test -bench=Parallel -blockprofile=prof.block
