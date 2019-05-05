//http://docs.ruanjiadeng.com/gopl-zh/ch11/ch11-04.html
package main 

import "testing"

func BenchmarkIsPalindrome(b *testing.B) {
    for i := 0; i < b.N; i++ {
        IsPalindrome("A man, a plan, a canal: Panama")
    }
}


$ cd $GOPATH/src/gopl.io/ch11/word2
$ go test -bench=.
PASS
BenchmarkIsPalindrome-8 1000000                1035 ns/op
ok      gopl.io/ch11/word2      2.179s


基准测试名的数字后缀部分, 这里是8, 表示运行时对应的 GOMAXPROCS 的值, 这对于一些和并发相关的基准测试是重要的信息.

报告显示每次调用 IsPalindrome 函数花费 1.035微秒, 是执行 1,000,000 次的平均时间. 因为基准测试驱动器并不知道每个基准测试函数运行所花的时候, 它会尝试在真正运行基准测试前先尝试用较小的 N 运行测试来估算基准测试函数所需要的时间, 然后推断一个较大的时间保证稳定的测量结果.

循环在基准测试函数内实现, 而不是放在基准测试框架内实现, 这样可以让每个基准测试函数有机会在循环启动前执行初始化代码, 这样并不会显著影响每次迭代的平均运行时间. 如果还是担心初始化代码部分对测量时间带来干扰, 那么可以通过 testing.B 参数的方法来临时关闭或重置计时器, 不过这些一般很少会用到.



如这个例子所示, 快的程序往往是有很少的内存分配. -benchmem 命令行标志参数将在报告中包含内存的分配数据统计. 我们可以比较优化前后内存的分配情况:

$ go test -bench=. -benchmem
PASS
BenchmarkIsPalindrome    1000000   1026 ns/op    304 B/op  4 allocs/op



func benchmark(b *testing.B, size int) { /* ... */ }
func Benchmark10(b *testing.B)         { benchmark(b, 10) }
func Benchmark100(b *testing.B)        { benchmark(b, 100) }
func Benchmark1000(b *testing.B)       { benchmark(b, 1000) }