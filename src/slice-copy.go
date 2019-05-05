http://hsl46346.blog.163.com/blog/static/1776405020125402646586/


package main

import "fmt"

func main() {
 s := [5]int{1, 2, 3, 4, 5}

 s1 := s[1:3]
 s2 := s[0:2]
 fmt.Println(s1, s2)

 //  s1和s2指向同一个底层数组，copy只是数据上的变化，而没有影响到各个切片的指向位置！
 copy(s2, s1)
 fmt.Println(s, s1, s2)

 s1[0] = 9
 fmt.Println(s, s1, s2)
}

output:
[2 3] [1 2] 
[2 3 3 4 5] [3 3] [2 3] 
[2 9 3 4 5] [9 3] [2 9]


package main

import "fmt"

func main() {
s := [5]int{1, 2, 3, 4, 5}

s1 := s[1:3]  

s2 := make([]int, 2) // make已经给s2分配好底层array的空间，并且用0填补array  
fmt.Println(s2)  
 
copy(s2, s1)  
fmt.Println(s2)
   
s2[0] = 9  
s1[0] = 99  

fmt.Println(s1, s2)
}

output:
[0 0] 
[2 3] 
[99 3] [9 3]


COPY总结：copy(s2, s1)过程只是将切片s1指向底层array的数据copy至切片s2指向底层array的数据上

-------------------------------------------------------------------

package main

import "fmt"

func main() {
 s := [5]int{1, 2, 3, 4, 5}

 s1 := s[1:3]
 s3 := append(s1, 9)

 s3[0] = 99
 fmt.Println(s, s1, s3)

 arr := [5]int{1, 2, 3, 4, 5}
 arr1 := arr[2:]

 arr2 := append(arr1, 9)
 arr2[0] = 99

 fmt.Println(arr, arr1, arr2)
}

output:
[1 99 3 9 5] [99 3] [99 3 9]
[1 2 3 4 5] [3 4 5] [99 4 5 9]

APPEND总结：s2 := append(s1, *)是切片s1上记录的切片信息复制给s2，如果s1指向的底层array长度不够，append的过程会发生如下操作：内存中不仅新开辟一块区域存储append后的切片信息，而且需要新开辟一块区域存储底层array（复制原来的array至这块新array中），最后再append新数据进新array中，这样，s2指向新array。反之，s2和s1指向同一个array，append的结果是内存中新开辟一个区域存储新切片信息。


参考官网：http://golang.org/doc/articles/slices_usage_and_internals.html

----------------------------------------------------------------

添加2个有趣的例子：
package main  

import "fmt" 

func main() {  
a1 := make([]int, 10, 10)  
a := a1[:5]  b := a // 不管是直接赋值，还是函数式传值，都是新建一个切片数据，与原切片指向同一个底层array
c := append(a, 1)  
c[0] = 9      

fmt.Println(a, b, c, a1) 
}

output：
[9 0 0 0 0] [9 0 0 0 0] [9 0 0 0 0 1] [9 0 0 0 0 1 0 0 0 0]


package main  

import "fmt"

func main() {  
slice := make([]int, 5, 10)  
newSlice := append(slice, 1)  
test(newSlice)    
fmt.Println("main slice: ", slice, len(slice), cap(slice))  
fmt.Println("main newSlice: ", newSlice, len(newSlice), cap(newSlice)) 
}  

func test(test []int) {  
test = append(test, 2)  
test[0] = 10  
fmt.Println("test: ", test, len(test), cap(test))   
}

output:
test:  [10 0 0 0 0 1 2] 7 10 
main slice:  [10 0 0 0 0] 5 10 
main newSlice:  [10 0 0 0 0 1] 6 10