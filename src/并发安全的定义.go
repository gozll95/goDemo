/*
如果在并发的情况下，一个函数可以正确的工作的话，那么我们就说这个函数是并发安全的，并发安全的函数不需要额外的
同步工作。
如果这个类型是并发安全的话，那么访问它的方法和操作就都是并发安全的。

所以，只有当文档中明确的说明了其是在并发安全的情况下，你才可以并发的访问它。

for-并发安全

1）将变量局限在单一的一个goroutine内
2）用互斥条件维持更高级别的不变性

相反，导出包级别的函数一般情况下都是并发安全的。由于package级的变量没法被限制在单一的gorouine，所以修改这些变量“必须”使用互斥条件。

deadlock死锁
livelock活锁
resource starvation饿死
*/

package bank

var balance int

func Deposit(amount int) { balance = balance + amount }

//先r再w
//Ar1
//Aw1

func Balance() int { return balance }

//r操作
//Ar2

// 上面两个并发
// Ar1-Ar2-Aw1
// 就会错

//无论任何时候,只要有两个goroutine并发访问同一变量，且至少其中的一个是写操作的时候就会发生数据竞争

/*
如果数据竞争的对象是一个比一个机器字更大的类型，就更麻烦了。
比如interface/string/slice类型都是如此

var x []int
go func(){x=make([]int,10)}()
go func(){x=make([]int,10000)}()
x[9999]=1

最后一个语句中的x的值是未定义的；其可能是nil，或者也可能是一个长度为10的slice，也可能是一个程度为1,000,000的slice。但是回忆一下slice的三个组成部分：指针(pointer)、长度(length)和容量(capacity)。

如果指针是从第一个make调来,而长度从第二个make来，x就变成了一个混合体。
*/

//并发并不是简单的语句交叉执行


避免数据竞争:
1)第一种方法是不要去写变量
2)避免从多个goroutine访问变量
3)第三种避免数据竞争的方法是允许很多goroutine去访问变量，但是在同一个时刻最多只有一个goroutine在访问。这种方式被称为“互斥”