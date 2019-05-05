接口

type S struct{
	i int
}

func(p *S)Get()int {
	return p.i
}

func(p *S)Put(v int){
	p.i=v
}

==>
type I interface{
	Get() int 
	Put(int)
}

==>
func f(p I){
	fmt.Println(p.Get())
	p.Put(1)
}


\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\
判断是否接口是否是特定的类型
方法一:

func f(p I){
	switch t:=p.(type){
		case *S:
		case *R:
		case S:
		case R:
		default:
	}
}

方法二:
common,ok

if t,ok:=something.(I);ok{
	//对于某些实现了接口I的t是其所拥有的类型
}



确定一个变量实现了某个接口，可以使用：
t:=something.(I)



\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\
可以在任意类型上定义方法(除了非本地类型，包括内建类型:int类型不能有方法)。然而可以新建一个拥有方法的整数类型。例如：

type Foo int 

func(self Foo)Emit(){
	fmt.Printf("%v",self)
}

type Emitter interface{
	Emit()
}


\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\
func sort(i []interface{}){
	switch i.(type){
		case string:
			//...
		case int:
			//...
	}
	return /*...*/
}


但是如果用 sort([]int{1, 4, 5}) 调用这个函数，会失败：cannot use i (type []int) as type []interface in function argument

##这是因为 Go 不能简单的将其转换为接口的 slice。转换到接口是容易的，但是转换到 slice 的开销就高了。
##简单来说 ：Go 不能（隐式）转换为 slice。



\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\
在接口中列出接口
type Interface interface{
	sort.Interface 
	Push(x interface{})
	Pop() interface{}
}