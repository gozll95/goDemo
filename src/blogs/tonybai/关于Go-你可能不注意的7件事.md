# 一、源文件字符集和字符集编码

Go源码文件默认采用***Unicode字符集***，***Unicode码点(code point)和内存中字节序列（byte sequence）的变换实现使用了UTF-8：一种变长多字节编码，同时也是一种事实字符集编码标准，为Linux、MacOSX 上的默认字符集编码，因此使用Linux或MacOSX进行Go程序开发，你会省去很多字符集转换方面的烦恼。***但如果你是在Windows上使用 默认编辑器编辑Go源码文本，当你编译以下代码时会遇到编译错误：
//hello.go
package main
import "fmt"
func main() {
    fmt.Println("中国人")
}
$ go build hello.go
# command-line-arguments
hello.go:6 illegal UTF-8 sequence d6 d0
hello.go:6 illegal UTF-8 sequence b9
hello.go:6 illegal UTF-8 sequence fa c8
hello.go:6 illegal UTF-8 sequence cb 22
hello.go:6 newline in string
hello.go:7 syntax error: unexpected }, expected )
这是因为Windows默认采用的是CP936字符集编码，也就是GBK编码，“中国人”三个字的内存字节序列为：
“d0d6    fab9    cbc8    000a” （通过iconv转换，然后用od -x查看）
这个字节序列并非utf-8字节序列，Go编译器因此无法识别。要想通过编译，需要将该源文件转换为UTF-8编码格式。
字符集编码对字符和字符串字面值(Literal)影响最大，在Go中对于字符串我们可以有三种写法：

- 1) **字面值**
var s = "中国人"

- 2) **码点表示法**
var s1 = "\u4e2d\u56fd\u4eba"
or
var s2 = "\U00004e2d\U000056fd\U00004eba"

- 3) **字节序列表示法（二进制表示法）**
var s3 = "\xe4\xb8\xad\xe5\x9b\xbd\xe4\xba\xba"

这三种表示法中，***除字面值转换为字节序列存储时根据编辑器保存的源码文件编码格式之外，其他两种均不受编码格式影响。*** 我们可以通过逐字节输出来查 看字节序列的内容：
    fmt.Println("s byte sequence:")
    for i := 0; i < len(s); i++ {
        fmt.Printf("0x%x ", s[i])
    }
    fmt.Println("")


# 三、Method Set

很多人纠结于method定义时receiver的类型（T or *T），个人觉得有两点考虑：
# 1) 效率
   Go方法调用receiver是以传值的形式传入方法中的。如果类型size较大，以value形式传入消耗较大，这时指针类型就是首选。
# 2) 是否赋值给interface变量、以什么形式赋值
   就像本节所描述的，由于T和*T的Method Set可能不同，我们在设计Method receiver type时需要考虑在interface赋值时通过对Method set的校验。


# 四、Method Type、Method Expression、Method Value

Go中没有class，方法与对象通过receiver联系在一起，我们可以为任何非builtin类型定义method：
```
type T struct {
    a int
}
func (t T) Get() int       { return t.a }
func (t *T) Set(a int) int { t.a = a; return t.a }
```

在C++等OO语言中，对象在调用方法时，编译器会自动在方法的第一个参数中传入this/self指针，而对于Go来 说，receiver也是同样道理，将T的method转换为普通function定义：

```
func Get(t T) int       { return t.a }
func Set(t *T, a int) int { t.a = a; return t.a }
```

这种function形式被称为***Method Type***,也可以被称为***Method的signature***。

Method的一般使用方式如下:
```
var t T
t.Get()
t.Set(1)
```

不过我们也可以像普通function那样使用它,根据上面的Method Type定义:

```
var t T
T.Get(t)
(*T).Set(&t, 1)
```

这种以直接以***类型名T调用方法M***的表达方法称为***Method Expression***。

类型T只能调用T的Method Set中的方法；同理*T只能调用*T的Method Set中的方法。上述例子中T的Method Set中只有Get，因此T.Get是合法的。但T.Set则不合法：

T.Set(2) //invalid method expression T.Set (needs pointer receiver: (*T).Set)
我们只能使用(*T).Set(&t, 11)。
这样看来Method Expression有些类似于C++中的static方法(以该类的某个对象实例作为第一个参数)。

另外***Method express自身类型就是一个普通function,可以作为右值赋值给一个函数类型的变量***:
```
    f1 := (*T).Set //函数类型：func (t *T, int)int
    f2 := T.Get //函数类型：func(t T)int
    f1(&t, 3)
    fmt.Println(f2(t))
```

Go中还定义了一种与Method有关的语法:如果一个表达式t具有静态类型T,M是T的Method Set中的一个方法,***那么t.M即为Method Value。***注意这里是t.M而不是T.M。
```
    f3 := (&t).Set //函数类型：func(int)int
    f3(4)
    f4 := t.Get//函数类型：func()int   
    fmt.Println(f4())
```

***可以看出，Method value与Method Expression不同之处在于，Method value绑定了T对象实例，它的函数原型并不包含Method Expression函数原型中的第一个参数。完整例子参见：details-in-go/4/methodexpressionandmethodvalue.go。***


# 五、for range"坑"大阅兵

## 1、iteration variable重用
for range的idiomatic的使用方式是使用short variable declaration（:=）形式在for expression中声明iteration variable，但需要注意的是这些variable在每次循环体中都会被重用，而不是重新声明。
```
//details-in-go/5/iterationvariable.go
… …
    var m = [...]int{1, 2, 3, 4, 5}
    for i, v := range m {
        go func() {
            time.Sleep(time.Second * 3)
            fmt.Println(i, v)
        }()
    }
    time.Sleep(time.Second * 10)
… …
```

在我的Mac上，输出结果如下：
```
$go run iterationvariable.go
4 5
4 5
4 5
4 5
4 5
```

各个goroutine中输出的i,v值都是for range循环结束后的i, v最终值，而不是各个goroutine启动时的i, v值。一个可行的fix方法：

```
    for i, v := range m {
        go func(i, v int) {
            time.Sleep(time.Second * 3)
            fmt.Println(i, v)
        }(i, v)
    }
```


## 2.range expression副本参与iteration
range后面接受的表达式的类型包括：array, pointer to array, slice, string, map和channel(有读权限的)。我们以array为例来看一个简单的例子：
```
//details-in-go/5/arrayrangeexpression.go
func arrayRangeExpression() {
    var a = [5]int{1, 2, 3, 4, 5}
    var r [5]int
    fmt.Println("a = ", a)
    for i, v := range a {
        if i == 0 {
            a[1] = 12
            a[2] = 13
        }
        r[i] = v
    }
    fmt.Println("r = ", r)
}


我们期待输出结果：
a =  [1 2 3 4 5]
r =  [1 12 13 4 5]
a =  [1 12 13 4 5]
但实际输出结果却是：
a =  [1 2 3 4 5]
r =  [1 2 3 4 5]
a =  [1 12 13 4 5]
```

我们原以为在第一次iteration，也就是i = 0时，我们对a的修改(a[1] = 12，a[2] = 13)会在第二次、第三次循环中被v取出，但结果却是v取出的依旧是a被修改前的值：2和3。这就是for range的一个不大不小的坑：***range expression副本参与循环***。也就是说在上面这个例子里，***真正参与循环的是a的副本，而不是真正的a，***伪代码如 下：

```
    for i, v := range a' {//a' is copy from a
        if i == 0 {
            a[1] = 12
            a[2] = 13
        }
        r[i] = v
    }
```

***Go中的数组在内部表示为连续的字节序列，虽然长度是Go数组类型的一部分，但长度并不包含的数组的内部表示中，而是由编译器在编译期计算出 来。这个例子中，对range表达式的拷贝，即对一个数组的拷贝，a'则是Go临时分配的连续字节序列，与a完全不是一块内存。因此无论a被 如何修改，其副本a'依旧保持原值，并且参与循环的是a'，因此v从a'中取出的仍旧是a的原值，而非修改后的值。***


我们再来试试***pointer to array***：
```
func pointerToArrayRangeExpression() {
    var a = [5]int{1, 2, 3, 4, 5}
    var r [5]int
    fmt.Println("pointerToArrayRangeExpression result:")
    fmt.Println("a = ", a)
    for i, v := range &a {
        if i == 0 {
            a[1] = 12
            a[2] = 13
        }
        r[i] = v
    }
    fmt.Println("r = ", r)
    fmt.Println("a = ", a)
    fmt.Println("")
}
这回的输出结果如下：
pointerToArrayRangeExpression result:
a =  [1 2 3 4 5]
r =  [1 12 13 4 5]
a =  [1 12 13 4 5]
```

我们看到这次r数组的值与最终a被修改后的值一致了。这个例子中我们使用了*[5]int作为range表达式，***其副本依旧是一个指向原数组 a的指针***，因此后续所有循环中均是&a指向的原数组亲自参与的，因此v能从&a指向的原数组中取出a修改后的值。

idiomatic go建议我们尽可能的用***slice替换掉array的使用***，这里用slice能否实现预期的目标呢？我们来试试：
```
func sliceRangeExpression() {
    var a = [5]int{1, 2, 3, 4, 5}
    var r [5]int
    fmt.Println("sliceRangeExpression result:")
    fmt.Println("a = ", a)
    for i, v := range a[:] {
        if i == 0 {
            a[1] = 12
            a[2] = 13
        }
        r[i] = v
    }
    fmt.Println("r = ", r)
    fmt.Println("a = ", a)
    fmt.Println("")
}
pointerToArrayRangeExpression result:
a =  [1 2 3 4 5]
r =  [1 12 13 4 5]
a =  [1 12 13 4 5]
```

显然用slice也能实现预期要求。我们可以分析一下slice是如何做到的。***slice在go的内部表示为一个struct，由(*T, len, cap)组成，其中*T指向slice对应的underlying array的指针，len是slice当前长度，cap为slice的最大容量。当range进行expression复制时，它实际上复制的是一个 slice，也就是那个struct。副本struct中的*T依旧指向原slice对应的array，为此对slice的修改都反映到 underlying array a上去了，v从副本struct中*T指向的underlying array中获取数组元素，也就得到了被修改后的元素值。***


slice与array还有一个不同点，就是其len在运行时可以被改变，而array的len是一个常量，不可改变。那么len变化的 slice对for range有何影响呢？我们继续看一个例子：

```
func sliceLenChangeRangeExpression() {
    var a = []int{1, 2, 3, 4, 5}
    var r = make([]int, 0)
    fmt.Println("sliceLenChangeRangeExpression result:")
    fmt.Println("a = ", a)
    for i, v := range a {
        if i == 0 {
            a = append(a, 6, 7)
        }
        r = append(r, v)
    }
    fmt.Println("r = ", r)
    fmt.Println("a = ", a)
}
输出结果：
a =  [1 2 3 4 5]
r =  [1 2 3 4 5]
a =  [1 2 3 4 5 6 7]
```

在这个例子中，原slice a在for range过程中被附加了两个元素6和7，其len由5增加到7，***但这对于r却没有产生影响。这里的原因就在于a的副本a'的内部表示struct中的 len字段并没有改变，依旧是5，***因此for range只会循环5次，也就只获取a对应的underlying数组的前5个元素。

range的副本行为会带来一些性能上的消耗，尤其是当range expression的类型为数组时，range需要复制整个数组；而当range expression类型为pointer to array或slice时，这个消耗将小得多，仅仅需要复制一个指针或一个slice的内部表示（一个struct）即可。我们可以通过 benchmark test来看一下三种情况的消耗情况对比：
对于元素个数为100的int数组或slice，测试结果如下：
//details-in-go/5/arraybenchmark
go test -bench=.
testing: warning: no tests to run
PASS
BenchmarkArrayRangeLoop-4             20000000           116 ns/op
BenchmarkPointerToArrayRangeLoop-4    20000000            64.5 ns/op
BenchmarkSliceRangeLoop-4             20000000            70.9 ns/op
可以看到range expression类型为slice或pointer to array的性能相近，消耗都近乎是数组类型的1/2。


# 六、select求值
```
func takeARecvChannel() chan int {
	fmt.Println("invoke takeARecvChannel")
	c := make(chan int)
	go func() {
		time.Sleep(3 * time.Second)
		c <- 1
	}()
	return c
}
func getAStorageArr() *[5]int {
	fmt.Println("invoke getAStorageArr")
	var a [5]int
	return &a
}
func takeASendChannel() chan int {
	fmt.Println("invoke takeASendChannel")
	return make(chan int)
}
func getANumToChannel() int {
	fmt.Println("invoke getANumToChannel")
	return 2
}
func main() {
	select {
	//recv channels
	case (getAStorageArr())[0] = <-takeARecvChannel():
		fmt.Println("recv something from a recv channel")
		//send channels
	case takeASendChannel() <- getANumToChannel():
		fmt.Println("send something to a send channel")
	}
}

运行结果：
$go run select.go
invoke takeARecvChannel
invoke takeASendChannel
invoke getANumToChannel
invoke getAStorageArr
recv something from a recv channel
```

通过例子我们可以看出：
- 1) select执行开始时，首先所有case expression的表达式都会被求值一遍，按语法先后次序。
invoke takeARecvChannel
invoke takeASendChannel
invoke getANumToChannel
例外的是recv channel的位于赋值等号左边的表达式（这里是：(getAStorageArr())[0]）不会被求值。
- 2) 如果选择***要执行的case***是***一个recv channel***，那么***它的赋值等号左边的表达式会被求值***：如例子中当goroutine 3s后向recvchan写入一个int值后，select选择了recv channel执行，此时对=左侧的表达式 (getAStorageArr())[0] 开始求值，输出“invoke getAStorageArr”。
***


# 七、panic的recover过程
Go没有提供“try-catch-finally”这样的异常处理设施，而仅仅提供了panic和recover，其中recover还要结合 defer使用。最初这也是被一些人诟病的点。但和错误码返回值一样，渐渐的大家似乎适应了这些，征讨之声渐稀，即便有也是排在“缺少generics” 之后了。

- 【panicking】
在没有recover的时候，一旦panic发生，panic会按既定顺序结束当前进程，这一过程成为panicking。下面的例子模拟了这一过程：
```
//details-in-go/7/panicking.go
… …
func foo() {
    defer func() {
        fmt.Println("foo defer func invoked")
    }()
    fmt.Println("foo invoked")
    bar()
    fmt.Println("do something after bar in foo")
}
func bar() {
    defer func() {
        fmt.Println("bar defer func invoked")
    }()
    fmt.Println("bar invoked")
    zoo()
    fmt.Println("do something after zoo in bar")
}
func zoo() {
    defer func() {
        fmt.Println("zoo defer func invoked")
    }()
    fmt.Println("zoo invoked")
    panic("runtime exception")
}
func main() {
    foo()
}
执行结果：
$go run panicking.go
foo invoked
bar invoked
zoo invoked
zoo defer func invoked
bar defer func invoked
foo defer func invoked
panic: runtime exception
goroutine 1 [running]:
… …
exit status 2
```

***从结果可以看出：***
    panic在zoo中发生，在zoo真正退出前，zoo中注册的defer函数会被逐一执行(FILO)，由于zoo defer中没有捕捉panic，因此panic被抛向其caller：bar。
    这时对于bar而言，***其函数体中的zoo的调用就好像变成了panic调用似的***，zoo有些类似于“黑客帝国3”中里奥被史密斯(panic)感 染似的，也变成了史密斯(panic)。panic在bar中扩展开来，bar中的defer也没有捕捉和recover panic，因此在bar中的defer func执行完毕后，panic继续抛给bar的caller: foo；
    这时对于foo而言，bar就变成了panic，同理，最终foo将panic抛给了main
    main与上述函数一样，没有recover，直接异常返回，导致进程异常退出。
  
- 【recover】
recover只有在defer函数中调用才能起到recover的作用，这样recover就和defer函数有了紧密联系。我们在zoo的defer函数中捕捉并recover这个panic：
```
//details-in-go/7/recover.go
… …
func zoo() {
    defer func() {
        fmt.Println("zoo defer func1 invoked")
    }()
    defer func() {
        if x := recover(); x != nil {
            log.Printf("recover panic: %v in zoo recover defer func", x)
        }
    }()
    defer func() {
        fmt.Println("zoo defer func2 invoked")
    }()
    fmt.Println("zoo invoked")
    panic("zoo runtime exception")
}
… …
这回的执行结果如下：
$go run recover.go
foo invoked
bar invoked
zoo invoked
zoo defer func2 invoked
2015/09/17 16:28:00 recover panic: zoo runtime exception in zoo recover defer func
zoo defer func1 invoked
do something after zoo in bar
bar defer func invoked
do something after bar in foo
foo defer func invoked
```

由于zoo在defer里恢复了panic，这样在zoo返回后，bar不会感知到任何异常，将按正常逻辑输出函数执行内容，比如：“do something after zoo in bar”,以此类推。
但若如果在zoo defer func中recover panic后，又raise another panic，那么zoo对于bar来说就又会变成panic了。