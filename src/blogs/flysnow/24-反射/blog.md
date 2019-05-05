# 前言

和Java语言一样,Go也实现运行时反射,这为我们提供一种可以在运行时操作任意类型对象的能力。比如我们可以查看一个接口变量的具体类型,看看一个结构体有多少字段,如何修改某个字段的值等等。

# TypeOf和ValueOf
func main() {
	u:= User{"张三",20}
	t:=reflect.TypeOf(u)
	fmt.Println(t)
}
type User struct{
	Name string
	Age int
}


# reflect.Value转原始类型
reflect.ValueOf 把任意类型的对象转为了一个reflect.Value

如何逆向转?

u1:=v.Interface().(User)
fmt.Println(u1)


reflect.TypeOf

reflect.Value->reflect.Type
t1:=v.Type()

# 获取类型底层类型
fmt.Println(t.Kind()
)


# 遍历字段和方法
for i:=0;i<t.NumField();i++ {
	fmt.Println(t.Field(i).Name)
}
for i:=0;i<t.NumMethod() ;i++  {
	fmt.Println(t.Method(i).Name)
}

# 修改字段的值

func main(){
    x:=2
    v:=reflect.ValueOf(&x)
    v.Elem().SetInt(100)
    fmt.Println(x)
}


# 动态调用方法

func main() {
	u:=User{"张三",20}
	v:=reflect.ValueOf(u)
	mPrint:=v.MethodByName("Print")
	args:=[]reflect.Value{reflect.ValueOf("前缀")}
	fmt.Println(mPrint.Call(args))
}
type User struct{
	Name string
	Age int
}
func (u User) Print(prfix string){
	fmt.Printf("%s:Name is %s,Age is %d",prfix,u.Name,u.Age)
}

MethodByName方法可以让我们根据一个方法名获取一个方法对象，然后我们构建好该方法需要的参数，最后调用Call就达到了动态调用方法的目的。

获取到的方法我们可以使用IsValid 来判断是否可用（存在）。

这里的参数是一个Value类型的数组，所以需要的参数，我们必须要通过ValueOf函数进行转换。

关于反射基本的介绍到这里就结束了，下一篇再介绍一些高级用法，比如获取字段的tag，常用的比如把一个json字符串转为一个struct就用到了字段的tag。
