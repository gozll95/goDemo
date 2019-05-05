package main

import(
    "fmt"
)

func main(){
    var a int 
    a = 10
    fmt.Println(a)

    var b *int 
    fmt.Printf("%p\n",&b)
    fmt.Printf("%p\n",b)
    fmt.Printf("%p\n",a)
    fmt.Printf("%p\n",&a)
}

10
0xc042004030
0x0
%!p<int=10>
0xc04200e098



close:主要用来关闭channel
len: 用来求长度,比如 string/array/slice/map/channel
new: 用来分配内存,主要用来分配值类型,比如int/struct,返回的是指针
make: 用来分配内存,主要用来分配引用类型,比如chan/map/slice

new和make的区别
new([]int)
    - > nil,0,0  []int 
        ptr len cap 
make([]int,0)
    - > ptr 0 0 []int 
        [0]int len cap




defer用途:
1.当函数返回时,执行defer语句,因此,可以用来做资源清理
2.多个defer语句，按先进后出的方式执行
3.defer语句中的变量,在defer声明时就决定了。


递归的设计原则:
1)一个大的问题能够分解成相似的小问题
2)定义好出口条件


闭包:
一个函数,与其相关的引用环境组合而成的实体。
// 变量 ，对变量 进行相同的作用力
// 外面相当于全局变量!!!

/*
package main 

import(
    "fmt"
)

func main(){
    var f=Adder()
    fmt.Print(f(1),"-")
    fmt.Print(f(20),"-")
    fmt.Print(f(300))
}

func Adder()func(int)int{
    var x int 
    return func(delta int)int{
        x+=delta
    }
}
*/

package main 

import(
    "fmt"
    "strings"
)

func makeSuffixFunc(suffix string)func(string){
    return func(name string)string{
        if !string.HasSuffix(name,suffix){
            return name+suffix
        }
        return name
    }
}

func main(){
    func1:=makeSuffixFunc(".bmp")
    func2:=makeSuffixFunc(".jpg")
    fmt.Println(func1("test"))
    fmt.Println(func2("test"))
}


#gorilla/websocket 用来做聊天室 websocket是长连接



1.实现一个冒泡排序


// 38, 1, 4, 5, 40
// 1, 38, 4, 5, 10
// 1, 4, 38, 5, 10
// 1, 4, 5, 38, 10
// 1, 4, 5, 10, 38

func bubble_sort(a []int) {

	for i := len(a)-1; i > 0; i-- {
		for j := 0; j < i; j++ {
			if a[j] > a[j+1] {
				a[j], a[j+1] = a[j+1], a[j]
			}
		}
	}
}


2.实现一个选择排序
每次选择最小的元素,放到第一个位置
下次从第二个元素开始找,放到第二个位置

//38, 1, 4, 5, 10
//1, 38, 4, 5, 10
//1, 4, 38, 5, 10
func select_sort(a []int) {
	for i := 0; i < len(a)-1; i++ {
		for j := i+1; j < len(a);j++ {
			if a[j] < a[i] {
				a[j], a[i] = a[i], a[j]
			}
		}
	}
}

3.插入排序
//38, 1, 4, 5, 10
//38
//1, 38, 
//1, 4, 38,
//1, 4, 5, 38,
//1, 4, 5, 10, 38

手上拿着5开始插,
从38往前走,找到第一个比它小的的

func insert_sort(a []int) {
	for i := 1; i < len(a); i++ {
		for j := i;j > 0;j-- {
			if a[j] < a[j-1] {
				a[j], a[j-1] = a[j-1], a[j]
			} else {
				break
			}
		}
	}
}


4.快速排序

左边的数都比38小,右边的数都比38大


38跟50比较,50>38,并且50已经在38右边了,所以50不动
// 38,100,4,5,10,50
38跟10比较,10小于38但是10在38右边所以交换位置
// 10,100,4,5,38,50
38跟100比较,100比38大但是100在38的左边所以交换位置
// 10,38,4,5,100,50
调整4
// 10,4,38,5,100,50
调整5
// 10,4,5,38,100,50
38左边数字都比38小了,右边数字都比38大-->那么38这一个数已经排好序了
以38这个元素为一个点,就划分成两个子数组
两个子数组再如法炮制-->递归
直到数组的长度=1就是有序了



//50, 100, 4, 5, 10, 50
//10, 100, 4, 5, 38, 50
//10, 38, 4, 5, 100, 50
//10, 4, 38, 5, 100, 50
//10, 4, 5, 38, 100, 50
//5, 4, 10,
//4, 5, 10, 38, 

left,right标识数组的长度
if left>=right 就是 一个 元素了


有两种:从前往后扫
      从后往前扫

func partion(a []int, left, right int) int {
	var i = left
	var j = right
	for i < j {
		for j > i && a[j] > a[left] {
			j--
		}
		a[j], a[left] = a[left], a[j]
		for i < j && a[i] < a[left] {
			i++
		}
		a[left], a[i] = a[i], a[left]
		fmt.Println(i)
	}
	return i
}


func qsort(a []int, left, right int) {
	if left >= right {
		return
	}

    // partion 找 中间数字的算法
	mid := partion(a, left, right)
	qsort(a, left, mid-1)
	qsort(a, mid+1, right)
}


