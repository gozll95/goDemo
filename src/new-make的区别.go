\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\

var p *int

fmt.Printf("%v",p)
//打印nil

var i int 
p=&i 

fmt.Printf("%v",p)


\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\
内建函数 new 本质上说跟其他语言中的同名函数功能一样：new(T) 分配了零值填充的 T 类型的内存空间，并且返回其地址，一个 *T 类型的值。用 Go 的术语说，它返回了一个指针，指向新分配的类型 T 的零值。记住这点非常重要。

这意味着使用者可以用 new 创建一个数据结构的实例并且可以直接工作。如 bytes.Buffer 的文档所述 “Buffer 的零值是一个准备好了的空缓冲。” 类似的，sync.Mutex 也没有明确的构造函数或 Init 方法。取而代之， sync.Mutex 的零值被定义为非锁定的互斥量。

零值是非常有用的。


#new分配零值
type SyncedBuffer struct{
	lock sync.Mutex 
	buffer bytes.Buffer 
}


p:=new(SyncedBuffer)   //Type *SyncedBuffer,已经可以使用
var v SyncedBuffer		//Type SyncedBuffer,同上



#new 分配；make 初始化

上面的两段可以简单总结为：

new(T) 返回 *T 指向一个零值 T
make(T) 返回初始化后的 T
当然 make 仅适用于slice，map 和channel。




\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\



func (n1 *NameAge) doSomething(n2 i n t ) { /* */ }

var n *NameAge
n.doSomething(2)

如果 x 可获取地址，并且 &x 的方法中包含了 m，x.m() 是 (&x).m() 更短的写法。

var n NameAge
n.doSomething(2)

这里 Go 会查找 NameAge 类型的变量n 的方法列表，没有找到就会再查找 *NameAge 类型的方法列表，并且将其转化为 (&n).doSomething(2)。



现在用两种不同的风格创建了两个数据类型。

#type NewMutex Mutex;
#type PrintableMutex struct {Mutex }.
现在 NewMutux 等同于 Mutex，但是它没有任何 Mutex 的方法。换句话说，它的方法是空的。

但是 PrintableMutex 已经从 Mutex 继承了方法集合。如同 [10] 所说：

*PrintableMutex 的方法集合包含了 Lock 和 Unlock 方法，被绑定到其匿名字段 Mutex。


\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\
#从 string 到字节或者 ruin 的 slice。

mystring := "hello this is string"  
byteslice := []byte(mystring)

runeslice := []rune(mystring)

#从字节或者整形的slice到string
b := []byte {'h','e','l','l','o'} // 复合声明
s := s t r i n g (b)
i := []rune {257,1024,65}
r := s t r i n g (i)