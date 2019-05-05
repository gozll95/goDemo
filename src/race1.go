# https://studygolang.com/articles/1531

介绍：

竞争条件是最狡诈的、最难以找到的编程错误。通常，在代码被布置到生产环境很久以后，它们才会出现并且造成奇怪的、神秘的错误。尽管Go语言的并发机制使得更容易的编写出干净的并发代码，依然无法避免竞争条件的出现。小心、勤勉以及测试是必须的。工具也可以提供帮助。

我们很高兴的宣布Go1.1包含了一个竞争检测器，一个全新的工具，用于在Go代码中找到竞争条件。该工具当前在Linux,OS X 和Windows平台可用，只要CPU是64位的x86架构。

竞争检测器基于C/C++的ThreadSanitizer 运行时库，该库在Google内部代码基地和Chromium找到许多错误。这个技术在2012年九月集成到Go中，从那时开始，它已经在标准库中检测到42个竞争条件。现在，它已经是我们持续构建过程的一部分，当竞争条件出现时，它会继续捕捉到这些错误。


工作原理
竞争检测器集成在go工具链中。当使用了-race作为命令行参数后，编译器会插桩代码，使得所有代码在访问内存时，会记录访问时间和方法。同时运行时库会观察对共享变量的未同步访问。当这种竞争行为被检测到，就会输出一个警告信息。（查看这里了解算法的具体细节）

由于设计原因，竞争检测器只有在被运行的代码触发时，才能检测到竞争条件，因此在现实的负载条件下运行是非常重要的。但是由于代码插桩，程序会使用10倍的CPU和内存，所以总是运行插桩后的程序是不现实的。矛盾的解决方法之一就是使用插桩后的程序来运行测试。负载测试和集成测试是好的候选，因为它们倾向于检验代码的并发部分。另外的方法是将单个插桩后的程序布置到运行服务器组成的池中，并且给予生产环境的负载。

使用方法:

竞争检测器已经完全集成到Go工具链中，仅仅添加-race标志到命令行就使用了检测器。
$ go test -race mypkg    // 测试包
$ go run -race mysrc.go  // 编译和运行程序
$ go build -race mycmd   // 构建程序
$ go install -race mypkg // 安装程序
获取以下示例程序，就可以尝试使用竞争检测器

$ go get -race code.google.com/p/go.blog/support/racy
$ racy
这里有两个例子，是由竞争检测器捕捉到的而上报产生的话题。
第1个例子: Timer.Reset
第一个例子是一个由竞争检测器发生的真实bug的简化版。它使用一个计时器来随机的时间间隔后打印消息，该间隔在0-1秒之间。该过程重复5秒。代码使用了time.AfterFunc来创建一个计时器，该计时器用于初次打印;随后使用Reset方法来安排下一条消息，每次都重用了该计时器。

 9 func main() {
10     start := time.Now()
11     var t *time.Timer
12     t = time.AfterFunc(randomDuration(), func() {
13         fmt.Println(time.Now().Sub(start)) 
14         t.Reset(randomDuration())
15     })
16     time.Sleep(5 * time.Second)
17 }
18 
19 func randomDuration() time.Duration {
20     return time.Duration(rand.Int63n(1e9))
21 }
看上去代码很合理，但是在特定约束下，程序以一种令人惊奇的方式失败了。

panic: runtime error: invalid memory address or nil pointer dereference
[signal 0xb code=0x1 addr=0x8 pc=0x41e38a]

goroutine 4 [running]:
time.stopTimer(0x8, 0x12fe6b35d9472d96)
    src/pkg/runtime/ztime_linux_amd64.c:35 +0x25
time.(*Timer).Reset(0x0, 0x4e5904f, 0x1)
    src/pkg/time/sleep.go:81 +0x42
main.func·001()
    race.go:14 +0xe3
created by time.goFunc
    src/pkg/time/sleep.go:122 +0x48
发生了什么事情? 运行插桩后的代码将给予更深入的信息。

==================
WARNING: DATA RACE
Read by goroutine 5:
  main.func·001()
     race.go:14 +0x169

Previous write by goroutine 1:
  main.main()
      race.go:15 +0x174

Goroutine 5 (running) created at:
  time.goFunc()
      src/pkg/time/sleep.go:122 +0x56
  timerproc()
     src/pkg/runtime/ztime_linux_amd64.c:181 +0x189
==================
竞争检测器发现了问题:在不同的goroutine中，对t有未同步的读、写操作。如果最初的计时间隔非常小，那么计时器函数有可能在main线程给t赋值前就调用，所以t.Reset就是在值为nil的t上调用。

为了修正这个竞争条件，更改代码使得对t的读写只能从main线程发出。
 9 func main() {
10     start := time.Now()
11     reset := make(chan bool)
12     var t *time.Timer
13     t = time.AfterFunc(randomDuration(), func() {
14         fmt.Println(time.Now().Sub(start))
15         reset <- true
16     })
17     for time.Since(start) < 5*time.Second {
18         <-reset
19         t.Reset(randomDuration())
20     }
21 }
这里main线程完全负责设置、重置定时器t，一个全新的channel以线程安全的方式通讯是否定时器需要重置
另一种更简单但是效率比较低的方法就是避免重用定时器
第2个例子: ioutil.Discard
第二个例子更加微妙

ioutil包中的Discard对象实现了接口io.Writer,但是抛弃了所有写入的数据。可以将其当做/dev/null:用于发送需要读取但不想存储的数据。该对象被广泛使用于io.Copy()，目的是耗尽读取端的数据，如:
io.Copy(ioutil.Discard, reader)
在2011年7年，Go团队注意到这样使用Discard比较低效:Copy函数每次调用时会分配一个大小为32K字节的内部缓存，当Discadrd作为参数时，这个缓存就没用了。我们认为这种将Copy和Discard组合使用，不应该这样浪费。

解决方案很简单。如果指定的写入端实现了ReadFrom方法，那么如下的Copy调用
io.Copy(writer, reader)
 将被托管给可能更高效的调用

writer.ReadFrom(reader)
我们添加了一个ReadFrom()方法到Discard的底层类型，该类型有一个在所有使用者之间共享的内部缓存。我们知道在理论上，存在着竞争条件，但是所有写入该缓存的数据被抛弃，所以我们认为这没有什么大不了的。

当实现了竞争检测器后，它立即检测出这段代码存在问题。我们再次认为这是多此一举，决定将这个竞争条件定义为非真实的。为了避免在构建过程中的该误报，我们实现了一个非竞争版本，该版本只在竞争检测器运行时可用。
但是几个月后，Brad遇到了一个奇怪的bug。在经过数日的调试后，他将问题定位于一个有ioutil.Discrd引发的真实竞争条件上。
这就是已经的竞争发生的代码，Discard是一个devNull,在所有用户间共享缓存。

var blackHole [4096]byte // shared buffer

func (devNull) ReadFrom(r io.Reader) (n int64, err error) {
    readSize := 0
    for {
        readSize, err = r.Read(blackHole[:])
        n += int64(readSize)
        if err != nil {
            if err == io.EOF {
                return n, nil
            }
            return
        }
    }
}
Brad的程序包含一个trackDigestReader类型，其中包装了一个io.Reader，并且记录了已经阅读内容的hash digest
type trackDigestReader struct {
    r io.Reader
    h hash.Hash
}

func (t trackDigestReader) Read(p []byte) (n int, err error) {
    n, err = t.r.Read(p)
    t.h.Write(p[:n])
    return
}
该代码可以用于在读取文件的同时，计算其SHA-1 哈希值:F

tdr := trackDigestReader{r: file, h: sha1.New()}
io.Copy(writer, tdr)
fmt.Printf("File hash: %x", tdr.h.Sum(nil))
某些情况下，数据并无写入，但是依然需要hash该文件，所以使用了Discard

io.Copy(ioutil.Discard, tdr)
但是在这种情况下，blackHole缓存已经不再是一个黑洞，它是一个合法的地方，用于在读取到数据到写入hash.Hash之间保存数据。在多线程并发hash文件时，每个线程共享blackHole缓存，竞争条件就会在读取和hash之间破坏数据。没有错误或者panic发生，但是hash值出错了。多么令人不愉快！

func (t trackDigestReader) Read(p []byte) (n int, err error) {
    // the buffer p is blackHole
    n, err = t.r.Read(p)
    // p可能被其余线程破坏
    // 
    t.h.Write(p[:n])
    return
}
这个bug最终被修正，方法是给予每个ioutil.Discard用户一个独一无二的内部缓存，消除了在共享缓存上的竞争条件。

结论
在检测并发编程的正确性方面，竞争检测器是一个威力强大的工具。它不会误报，要认真对待其发出的警告。但是它依赖于你的测试，必须保证测试覆盖了代码的并发部分，这样检测器才能完成任务。

你还在等待什么?今天就在你的代码上运行"go test -race"

本文来自：CSDN博客

感谢作者：fighterlyt

查看原文：介绍Go竞争检测器