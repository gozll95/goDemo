Nav logo
首页
下载App

搜索
Golang自动化测试介绍
96  caojunxyz 
2016.10.14 22:14* 字数 1038 阅读 1200评论 0喜欢 1
Go语言提供了对包（package）进行自动化测试的支持，带来了很大的工程便利。本文主要介绍相关的常用知识，进一步学习可阅读官方文档。

测试文件必须以"_test.go"结尾，并和被测试包放在同一目录下。测试文件可以有多个，用以组织复杂的测试逻辑。一个典型的目录结构如下：
student/ student.go student_test.go

student.go内容：

package student

import "math/rand"

type Math struct {
    name string
}

func (m *Math) Perm(n int) []int {
    return rand.Perm(n)
}

func (m *Math) Add(a, b int) int {
    return a + b
}

func (m *Math) StupidSum(N int) int {
    ret := 0
    for i := 1; i <= N; i++ {
        ret += i
    }
    return ret
}

func (m *Math) SmartSum(N int) int {
    return N * (N + 1) / 2
}
三类测试
func ExampleXxx() {}
func TestXxx(t *testing.T) {}
func BenchmarkXxx(b *testing.B) {}
三种函数分别用于不同的测试场景，Xxx部分可以省略，但不能以小写字母开头。良好的书写规范可以提高godoc生成文档的可读性，例如ExampleT_M()，T为类型，M为类型方法。

范例测试
ExampleXxx主要用于测试在给定输入的情况下方法的输出是否与预期相符。测试将用例标准输出的结果与// Output: expected result注释的预期结果进行一致性比较。如果运行结果与注释中预期不一致，该用例失败。// Output没有严格的大小写和空格要求，例如写为//ouput也是可以的，不过最好统一。

func ExampleMath_Add() {
    m := &Math{}
    fmt.Println(m.Add(2, 3))
    // Output: 5
}

func ExampleMath_Perm() {
    m := &Math{}
    for n := range m.Perm(4) {
        fmt.Println(n)
    }
    // Output:
    // 0
    // 1
    // 2
    // 3
}
注意，测试用例ExampleMath_Perm在给定输入的情况下，输出集合是确定的，但是顺序是随机的。预期结果可以书写任意一种合理结果，例如上面的预期也可以改为

// Output:
// 3
// 1
// 2
// 0
通用测试
TestXxx类型的测试函数的参数为*testing.T类型，用于管理测试状态和格式化测试日志。测试日志累计到执行结束后才输出到标准输出。可以通过T的相关方法控制测试逻辑。

func (c *T) Log(args ...interface{})
func (c *T) Logf(format string, args ...interface{})
Log和Logf方法用于日志输出，默认只输出错误日志，如果要输出全部日志需要使用-v标识运行go test命令。benchmarks默认输出全部日志。

func (c *T) Fail()
func (c *T) FailNow() 
func TestFail(t *testing.T) {
    t.Log("mark A")
    t.Fail()
    t.Log("mark B")
}

func TestFailNow(t *testing.T) {
    t.Log("mark C")
    t.FailNow()
    t.Log("mark D")
}

func TestOther(t *testing.T) {
    t.Log("mark E")
}

/* Output:
--- FAIL: TestFail (0.00s)
    student_test.go:44: mark A
    student_test.go:46: mark B
=== RUN   TestFailNow
--- FAIL: TestFailNow (0.00s)
    student_test.go:50: mark C
=== RUN   TestOther
--- PASS: TestOther (0.00s)
    student_test.go:56: mark E
*/
Fail标记用例失败，但继续执行当前用例。FailNow标记用例失败并且立即停止执行当前用例，继续执行下一个（默认按书写顺序）用例。

func (c *T) Error(args ...interface{})
func (c *T) Errorf(format string, args ...interface{})
Error等价于Log加Fail，Errorf等价于Logf加Fail。

func (c *T) Skip(args ...interface{})
func (c *T) SkipNow()
func (c *T) Skipf(format string, args ...interface{})
func (c *T) Skipped() bool
SkipNow标记跳过并停止执行该用例，继续执行下一个用例。Skip等价于Log加SkipNow，Skipf等价于Logf加SkipNow，Skipped返还用例是否被跳过。

func (c *T) Parallel()
示意该测试用例和其它并行用例（也调用该方法的）一起并行执行。我们构造一个例子来测试Parallel()方法是否真的是并行执行测试用例：

var counter int32
var wg sync.WaitGroup
var N = 100000

func TestParallel1(t *testing.T) {
    t.Parallel()
    wg.Add(1)
    for i := 1; i <= N; i++ {
        atomic.AddInt32(&counter, 1)
        if i%5555 == 0 {
            t.Log("p1 -->", counter)
        }
    }
    wg.Done()
}

func TestParallel2(t *testing.T) {
    t.Parallel()
    wg.Add(1)
    for i := 1; i <= N; i++ {
        atomic.AddInt32(&counter, 1)
        if i%5555 == 0 {
            t.Log("p2 -->", counter)
        }
    }
    wg.Done()
}

func TestParallelEnd(t *testing.T) {
    t.Parallel()
    time.After(time.Second)
    // wg.Wait()
    result := 2 * N
    t.Log("result: ", counter, result)
}
因为现在CPU速度很快，两个不同的测试用例启动会有时间差，所以需要将N设得足够大才能看到并行随机效果。if i%5555 == 0这个条件分支的作用是减少日志输出。运行可以看到每次counter的结果是随机的，当把wg.Wait()这条语句打开时，可以看到确定的输出（counter等于result）。

func (c *T) Run(name string, f func(t *T)) bool 
func TestRunSuits(t *testing.T) {
    t.Run("A=1", func(t *testing.T) {
        t.Log("sub test A=1")
    })
    t.Run("A=2", func(t *testing.T) {
        t.Log("sub test A=2")
    })
    t.Run("B=1", func(t *testing.T) {
        t.Log("sub test B=1")
    })
    t.Run("B=2", func(t *testing.T) {
        t.Log("sub test B=2")
    })
}
Run 运行一个名为name的子用例，返回该子用例是否通过。可以通过-run exp正则表达式参数指定要运行的子用例，例如上面例子中可以通过go test -v -run "/A" 运行两个子用例，正则表达式的顶层以 / 开头。子用例的引入方便更好的对测试用例进行组织。

Bechmarks
Benchmarks类型的测试函数的参数为*testing.B类型，通过go test参数可以对测试用例进行正则匹配，可以控制使用的CPU核心数等。看一个例子：

func BenchmarkMath_StupidSum(b *testing.B) {
    m := &Math{}
    for i := 0; i < b.N; i++ {
        m.StupidSum(1, 100)
    }
}

func BenchmarkMath_SmartSum(b *testing.B) {
    m := &Math{}
    for i := 0; i < b.N; i++ {
        m.SmartSum(1, 100)
    }
}

func BenchmarkMath_StupidSumSerial(b *testing.B) {
    m := &Math{}
    for i := 0; i < b.N; i++ {
        m.StupidSum(1, 100)
    }
}

func BenchmarkMath_StupidSumParallel(b *testing.B) {
    m := &Math{}
    for i := 0; i < b.N; i++ {
        b.RunParallel(func(pb *testing.PB) {
            for pb.Next() {
                m.StupidSum(1, 100)
            }
        })
    }
}

func BenchmarkMath_SmartSumSerial(b *testing.B) {
    m := &Math{}
    for i := 0; i < b.N; i++ {
        m.SmartSum(1, 100)
    }
}

func BenchmarkMath_SmartSumParallel(b *testing.B) {
    m := &Math{}
    for i := 0; i < b.N; i++ {
        b.RunParallel(func(pb *testing.PB) {
            for pb.Next() {
                m.SmartSum(1, 100)
            }
        })
    }
}
/*
go test -cpu 4 -benchmem -bench . 运行结果：
BenchmarkMath_StupidSum-4               30000000            57.3 ns/op         0 B/op          0 allocs/op
BenchmarkMath_SmartSum-4                2000000000           0.41 ns/op        0 B/op          0 allocs/op
BenchmarkMath_StupidSumSerial-4         30000000            57.1 ns/op         0 B/op          0 allocs/op
BenchmarkMath_StupidSumParallel-4          10000        259288 ns/op         181 B/op          8 allocs/op
BenchmarkMath_SmartSumSerial-4          2000000000           0.39 ns/op        0 B/op          0 allocs/op
BenchmarkMath_SmartSumParallel-4           50000        441988 ns/op         176 B/op          8 allocs/op
PASS
ok      github.com/caojunxyz/gotest/student 30.222s
*/

/*
go test -cpu 1 -benchmem -bench . 运行结果：
BenchmarkMath_StupidSum             30000000            58.6 ns/op         0 B/op          0 allocs/op
BenchmarkMath_SmartSum              2000000000           0.42 ns/op        0 B/op          0 allocs/op
BenchmarkMath_StupidSumSerial       20000000            57.0 ns/op         0 B/op          0 allocs/op
BenchmarkMath_StupidSumParallel        10000        659949 ns/op          80 B/op          5 allocs/op
BenchmarkMath_SmartSumSerial        2000000000           0.52 ns/op        0 B/op          0 allocs/op
BenchmarkMath_SmartSumParallel         50000        424859 ns/op          80 B/op          5 allocs/op
PASS
ok      github.com/caojunxyz/gotest/student 33.202s
*/

正则匹配：只运行SmartSum相关的测试命令go test -bench *** "SmartSum"
指定使用CPU核心数：_-cpu*** n_
每个测试用例运行n次: -count n
打印内存分配统计信息: -benchmem
代码覆盖率：-cover




//b.ResetTimer()一般用于准备时间比较长的时候重置计时器减少准备时间带来的误差，这里可用可不用