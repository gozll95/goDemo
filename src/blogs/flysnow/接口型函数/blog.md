***接口型函数***，指的是***用函数实现接口，这样在调用的时候就会非常简便，我称这种函数，为接口型函数，这种方式使用于只有一个函数的接口。***

我们以迭代一个map为例，演示这一技巧，这种方式有点类似于groovy中Map的each方法一样，也是Gradle里each闭包。

# 原始接口实现

```
type Handler interface {
	Do(k, v interface{})
}
func Each(m map[interface{}]interface{}, h Handler) {
	if m != nil && len(m) > 0 {
		for k, v := range m {
			h.Do(k, v)
		}
	}
}
```

首先定义一个Handler接口，只有一个Do方法，接收k，v两个参数，这就是一个接口了，我们后面会实现他，具体做什么由我们的实现决定。

然后我们定义了一个Each函数，这个函数的功能，就是迭代传递过来的map参数，然后把map的每个key和value值传递给Handler的Do方法，去做具体的事情，可以是输出，也可以是计算，具体由这个Handler的实现来决定，这也是面向接口编程。

现在我们就以新学期开学，大家自我介绍为例，演示使用我们刚刚定义的Each方法和Handler接口。这里我们假设有三个学生，分别为：张三，李四和王五，他们每个人都要介绍自己的名字和年龄。

```
type welcome string
func (w welcome) Do(k, v interface{}) {
	fmt.Printf("%s,我叫%s,今年%d岁\n", w,k, v)
}
func main() {
	persons := make(map[interface{}]interface{})
	persons["张三"] = 20
	persons["李四"] = 23
	persons["王五"] = 26
	var w welcome = "大家好"
	Each(persons, w)
}
```

以上实现，我们定义了一个map来存储学生们，map的key是学生的名字，value是该学生的年龄。welcome是我们新定义的类型，对应基本类型string，该welcome实现了Handler接口，打印出自我介绍。

# 接口型函数出场

以上实现，主要有***两点不太好***：

- 因为必须要实现Handler接口，Do这个方法名不能修改，不能定义一个更有意义的名字
- ***必须要新定义一个类型，才可以实现Handler接口***，才能使用Each函数

首先我们先解决第一个问题，根据我们具体做的事情定义一个更有意义的方法名，比如例子中是自我介绍，那么使用selfInfo要比Do这个干巴巴的方法要好的多。

如果调用者改了方法名，那么就不能实现Handler接口，还要使用Each方法怎么办？那就是由提供Each函数的负责提供Handler的实现，我们添加代码如下：

```
type HandlerFunc func(k, v interface{})
func (f HandlerFunc) Do(k, v interface{}){
	f(k,v)
}
```

以上代码，我们定义了一个新的类型HandlerFunc，它是一个func(k, v interface{})类型，然后这个新的HandlerFunc实现了Handler接口，Do方法的实现是调用HandlerFunc本身，因为HandlerFunc类型的变量就是一个方法。

现在我们使用这种方式实现同样的效果。

```
type welcome string
func (w welcome) selfInfo(k, v interface{}) {
	fmt.Printf("%s,我叫%s,今年%d岁\n", w,k, v)
}
func main() {
	persons := make(map[interface{}]interface{})
	persons["张三"] = 20
	persons["李四"] = 23
	persons["王五"] = 26
	var w welcome = "大家好"
	Each(persons, HandlerFunc(w.selfInfo))
}
```

还是差不多原来的实现，只是把方法名Do改为selfInfo。HandlerFunc(w.selfInfo)不是方法的调用，而是转型，因为selfInfo和HandlerFunc是同一种类型，所以可以强制转型。转型后，因为HandlerFunc实现了Handler接口，所以我们就可以继续使用原来的Each方法了。

# 进一步重构

现在解决了命名的问题，***但是每次强制转型不太好，我们继续重构，可以采用新定义一个函数的方式，帮助调用者强制转型***。

```
func EachFunc(m map[interface{}]interface{}, f func(k, v interface{})) {
	Each(m,HandlerFunc(f))
}
type welcome string
func (w welcome) selfInfo(k, v interface{}) {
	fmt.Printf("%s,我叫%s,今年%d岁\n", w,k, v)
}
func main() {
	persons := make(map[interface{}]interface{})
	persons["张三"] = 20
	persons["李四"] = 23
	persons["王五"] = 26
	var w welcome = "大家好"
	EachFunc(persons, w.selfInfo)
}
```

新增了一个EachFunc函数，帮助调用者强制转型，调用者就不用自己做了。

现在我们发现EachFunc函数接收的是一个func(k, v interface{})类型的函数，没有必要实现Handler接口了，所以我们新的类型可以去掉不用了。


```
func selfInfo(k, v interface{}) {
	fmt.Printf("大家好,我叫%s,今年%d岁\n", k, v)
}
func main() {
	persons := make(map[interface{}]interface{})
	persons["张三"] = 20
	persons["李四"] = 23
	persons["王五"] = 26
	EachFunc(persons, selfInfo)
}
```

去掉了自定义类型welcome之后，整个代码更简洁，可读性更好。我们的方法含义都是：
- 让这学生自我介绍
- 让这些学生起立
- 让这些学生早读
- 让这些学生…
都是这种默认，方法处理，更符合自然语言规则。


# 延伸
以上关于函数型接口就写完了，如果我们仔细留意，发现和我们自己平时使用的http.Handle方法非常像，其实接口http.Handler就是这么实现的。

```
type Handler interface {
	ServeHTTP(ResponseWriter, *Request)
}
func Handle(pattern string, handler Handler) {
	DefaultServeMux.Handle(pattern, handler)
}
func HandleFunc(pattern string, handler func(ResponseWriter, *Request)) {
	DefaultServeMux.HandleFunc(pattern, handler)
}
```

这是一种非常好的技巧，提供两种函数，既可以以接口的方式使用，也可以以方法的方式，对应我们例子中的Each和EachFunc这两个函数，灵活方便。


# 最后，附上完整的源代码：
```
package main
import (
	"fmt"
)
type Handler interface {
	Do(k, v interface{})
}
type HandlerFunc func(k, v interface{})
func (f HandlerFunc) Do(k, v interface{}) {
	f(k, v)
}
func Each(m map[interface{}]interface{}, h Handler) {
	if m != nil && len(m) > 0 {
		for k, v := range m {
			h.Do(k, v)
		}
	}
}
func EachFunc(m map[interface{}]interface{}, f func(k, v interface{})) {
	Each(m, HandlerFunc(f))
}
func selfInfo(k, v interface{}) {
	fmt.Printf("大家好,我叫%s,今年%d岁\n", k, v)
}
func main() {
	persons := make(map[interface{}]interface{})
	persons["张三"] = 20
	persons["李四"] = 23
	persons["王五"] = 26
	EachFunc(persons, selfInfo)
}
```