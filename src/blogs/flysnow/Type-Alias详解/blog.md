# 安装go 1.9
很多众所周知的原因，大家可能无法下载最新的go 1.9 sdk，如果你没有梯子，可以到我自建的这个镜像网站下载，有很多常用的开发软件，其中就包含最新的go 1.9。镜像地址：http://mirrors.flysnow.org/

# 作用
type alias这个特性的主要目的是用于已经定义的类型，在package之间的移动时的兼容。比如我们有一个导出的类型flysnow.org/lib/T1，现在要迁移到另外一个package中, 比如flysnow.org/lib2/T1中。

没有type alias的时候我们这么做，就会导致其他第三方引用旧的package路径的代码，都要统一修改，不然无法使用。

有了type alias就不一样了，类型T1的实现我们可以迁移到lib2下，同时我们在原来的lib下定义一个lib2下T1的别名，这样第三方的引用就可以不用修改，也可以正常使用，只需要兼容一段时间，再彻底的去掉旧的package里的类型兼容，这样就可以渐进式的重构我们的代码，而不是一刀切。

//package:flysnow.org/lib
type T1=lib2.T1
type alias vs defintion
我们基于一个类型创建一个新类型，称之为defintion；基于一个类型创建一个别名，称之为alias，这就是他们最大的区别。


type MyInt1 int
type MyInt2 = int
第一行代码是基于基本类型int创建了新类型MyInt1，第二行是创建的一个int的类型别名MyInt2，注意类型别名的定义是=。

var i int =0
var i1 MyInt1 = i //error
var i2 MyInt2 = i
fmt.Println(i1,i2)
仔细看这个示例，第二行把一个int类型的变量i,赋值给MyInt1类型的变量i1会被提示编译错误：类型无法转换。但是第三行把int类型的变量i,赋值给MyInt2类型的变量i2就可以，不会提示错误。

从这个例子也可以看出来，这两种定义方式的不同，因为Go是强类型语言，所以类型之间的转换必须强制转换，因为int和MyInt1是不同的类型，所以这里会报编译错误。

但是因为***MyInt2只是int的一个别名，本质上还是一个int类型***，所以可以直接赋值，不会有问题。

#类型方法
每个类型都可以通过接受者的方式，添加属于它自己的方法，我们看下通过type alias的类型是否可以，以及拥有哪些方法。

type MyInt1 int
type MyInt2 = int
func (i MyInt1) m1(){
	fmt.Println("MyInt1.m1")
}
func (i MyInt2) m2(){
	fmt.Println("MyInt2.m2")
}
func main() {
	var i1 MyInt1
	var i2 MyInt2
	i1.m1()
	i2.m2()
}
以上示例代码看着是没有任何问题，但是我们编译的时候会提示：

i2.m2 undefined (type int has no field or method m2)
cannot define new methods on non-local type int
这里面有2个错误，一个是提示类型int没有m2这个方法，所以我们不能调用，因为MyInt2本质上就是int。

第二个错误是***我们不能为int类型添加新方法***，什么意思呢？因为***int是一个非本地类型，所以我们不能为其增加方法***。既然这样，那我们自定义个struct类型试试。

type User struct {
}
type MyUser1 User
type MyUser2 = User
func (i MyUser1) m1(){
	fmt.Println("MyUser1.m1")
}
func (i MyUser2) m2(){
	fmt.Println("MyUser2.m2")
}
//Blog:www.flysnow.org
//Wechat:flysnow_org
func main() {
	var i1 MyUser1
	var i2 MyUser2
	i1.m1()
	i2.m2()
}
换成struct，正常运行。所以***本地定义的类型的别名，还是可以为其添加方法的***。现在我们接着上面的例子，看一个有趣的现象，我在main函数里增加如下代码：

var i User
i.m2()
然后运行，发现，可以正常运行。是不是很奇怪，我们并没有为类型User 定义方法啊，怎么可以调用呢？这就得益于type alias，MyUser2完全等价于User，所以为MyUser2定义方法，等于就为User定义了方法，反之，亦然。

但是对于新定义的类型MyUser1就不行了，因为它完全是个新类型，所以User的方法，MyUser是没有的。这里不再举例，大家自己可以试试。

还有一点需要注意，因为MyUser2完全等价于User，所以User已经有的方法，MyUser2不能再声明，反之亦然，如果定义了会有如下提示：

./main.go:37:6: User.m1 redeclared in this block
	previous declaration at ./main.go:31:6
其实就是重复声明的意思，不能再次重复声明了。

# 接口实现
上面的小结我们可以发现，User和MyUser2是等价的，并且其中一个新增了方法，另外一个也会有。那么基于此推导出，一个实现了某个接口，另外一个也会实现。现在验证一下：

type I interface {
	m2()
}
type User struct {
}
type MyUser2 = User
func (i User) m(){
	fmt.Println("User.m")
}
func (i MyUser2) m2(){
	fmt.Println("MyUser2.m2")
}
func main() {
	var u User
	var u2 MyUser2
	var i1 I =u
	var i2 I =u2
	fmt.Println(i1,i2)
}
定义了一个接口I，从代码上看，只有MyUser2实现了它，但是我们代码的演示中，发现User也实现了接口I，所以这就验证了我们的推到是正确的，返回来如果User实现了某个接口，那么它的type alias也同样会实现这个接口。

以上讲了很多示例都是类型struct的别名，我们看下接口interface的type alias是否也是等价的。

type I interface {
	m()
}
type MyI1 I
type MyI2 = I
type MyInt int
func (i MyInt) m(){
	fmt.Println("MyInt.m")
}
定义了一个接口I，MyI1是基于I的新类型；MyI2是I的类型别名；MyInt实现了接口I。下面进行测试。


//Blog:www.flysnow.org
//Wechat:flysnow_org
func main() {
	//赋值实现类型MyInt
	var i I = MyInt(0)
	var i1 MyI1 = MyInt(0)
	var i2 MyI2 = MyInt(0)
	//接口间相互赋值
	i = i1
	i = i2
	i1 = i2
	i1 = i
	i2 = i
	i2 = i1
}
以上代码运行是正常的，这个是前面讲的具体类型（struct，int等）的type alias不一样，只要实现了接口，就可以相互赋值，管你是新定义的接口MyI1，还是接口的别名MyI2。

# 类型的嵌套
我们都知道type alias的两个类型是等价的，但是他们在类型嵌套时有些不一样。


//Blog:www.flysnow.org
//Wechat:flysnow_org
func main() {
	my:=MyStruct{}
	my.T2.m1()
}
type T1 struct {
}
func (t T1) m1(){
	fmt.Println("T1.m1")
}
type T2 = T1
type MyStruct struct {
	T2
}
示例中T2是T1的别名，但是我们把T2嵌套在MyStruct中，在调用的时候只能通过T2这个名称调用，而不能通过T1，会提示没这个字段的。反过来也一样。

这是因为T1,T2是两个名称，虽然他们等价，但他们是具有两个不同名字的等价类型，所以在类型嵌套的时候，就是两个字段。

当然我们可以把T1,T2同时嵌入到MyStrut中，进行分别调用。


//Blog:www.flysnow.org
//Wechat:flysnow_org
func main() {
	my:=MyStruct{}
	my.T2.m1()
	my.T1.m1()
}
type MyStruct struct {
	T2
	T1
}
以上也是可以正常运行的，证明这是具有两个不同名字的，同种类型的字段。

下面我们做个有趣的实验，把main方法的代码改为如下：

//Blog:www.flysnow.org
//Wechat:flysnow_org
func main() {
	my:=MyStruct{}
	my.m1()
}
猜猜是不是可以正常编译运行呢？答应可能出乎意料，是不能正常编译的，提示如下：

1
./main.go:25:4: ambiguous selector my.m1
其实想想很简单，不知道该调用哪个，太模糊了，匹配不了，不然该用T1的m1,还是T2的m1。这种结果不限于方法，字段也也一样；也不限于type alias，type defintion也是一样的，只要有重复的方法、字段，就会有这种提示，因为不知道该选择哪个。

# 类型循环
type alias的声明，一定要留意类型循环，不要产生了循环，一旦产生，就会编译不通过，那么什么是类型循环呢。假如type T2 = T1,那么T1绝对不能直接、或者间接的引用到T2，一旦有，就会类型循环。

type T2 = *T2
type T2 = MyStruct
type MyStruct struct {
	T1
	T2
}
以上两种定义都是类型循环，我们自己在使用的过程中，要避免这种定义的出现。

byte and rune
这两个类型一个是int8的别名，一个是int32的别名，在Go 1.9之前，他们是这么定义的。

type byte byte
type rune rune
现在Go 1.9有了type alias这个新特性后，他们的定义就变成如下了：

type byte = uint8
type rune = int32
恩，非常很省事和简洁。

# 导出未导出的类型
type alias还有一个功能，可以***导出一个未被导出的类型***。

package lib
//Blog:www.flysnow.org
//Wechat:flysnow_org
type user struct {
	name string
	Email string
}
func (u user) getName() string {
	return u.name
}
func (u user) GetEmail() string {
	return u.Email
}
//把这个user导出为User
type User = user
user本身是一个未导出的类型，不能被其他package访问，但是我们可以通过type User = user，定义一个User，这样这个User就可以被其他package访问了，可以使用user类型导出的字段和方法，示例中是Email字段和GetEmail方法，另外未被导出name字段和getName方法是不能被其他package使用的。

#小结
type alias的定义，本质上是一样的类型，只是起了一个别名，源类型怎么用，别名类型也怎么用，保留源类型的所有方法、字段等。
