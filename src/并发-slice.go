package main

import (
       "fmt"
       "sync"
)

func main() {

       sourceArray := [...]string{"a", "b", "c", "d", "e", "f", "g"}
       fmt.Println(sourceArray)
       slice_1 := sourceArray[:4]
       slice_2 := sourceArray[2:]
       slice_3 := sourceArray[1:]
       slice_4 := slice_3[1:]
       slice_5 := slice_4[1:]

       slice_1[3] = "i"
       fmt.Println(sourceArray)
       fmt.Println(slice_1)
       fmt.Println(slice_2)
       fmt.Println(slice_3)
       fmt.Println(slice_4)
       fmt.Println(slice_5)

       wg := &sync.WaitGroup{}

       wg.Add(5)
       /*slice, map, channel, function, method 皆为go语言中的引用类型; 一个切片slice)是一个隐藏数组的引用,并且对于该切片的切片也是引用同一个数组.*/
       go func( p []string) {

              defer wg.Done()
              p[3] = "w" //p[3] --> sourceArray[3]

       }(slice_1) //在函数调用传参时, 因为slice是引用类型, 所以slice_1与p这两个引用皆指向同一个数组sourceArray
       go func( p []string) {

              defer wg.Done()
              p[1] = "z" //p[1] --> sourceArray[3]

       }(slice_2)
       go func( p []string) {

              defer wg.Done()
              p[2] = "x" //p[2] --> sourceArray[3]

       }(slice_3)
       go func(p []string) {

              defer  wg.Done()
              p[1] = "q" //p[1] --> sourceArray[3] //即使切片的切片,依然指向同一个数组sourceArray

       }(slice_4)
       go func(p []string) {

              defer  wg.Done()
              p[0] = "t" //p[0] --> sourceArray[3] //即使切片的切片的切片等等,依然指向同一个数组sourceArray

       }(slice_5)

       wg.Wait()

       /* 由以上分析可知, 在多线程编程, 协程编程中, 对于共享的变量,资源的读写是一定要串行化的, 比如加锁,或通过channel来实现排他性访问.
       通过上面的测试代码,我们了解到,对于切片(slice)的读写最终都是对切片的隐藏数组的读写, 如果读写的数组索引范围,或是元素重合, 则多个协程
       读写共享元素,在并发情况下,则必然产生竞争,破坏共享元素数据, 所以要保护,要么加锁, 要么用channel将访问排队串行化.
       所以对于slice的使用一定要特别注意3点: （1）切片为引用类型 （2）一个切片是一个隐藏数组的引用 （3）切片再切片等等,依然指向同一个隐藏数组
       注意: 这些结论都是在没有超过切片容量（隐藏数组的长度)的情况下得出的, 因为go语言内置函数append会在切片容量不够时, 分配新的更大的切片及相应的隐藏数组;
       但是这并不影响以上结论,新切片, 旧切片依然各自指向自己的隐藏数组.
       */

       fmt.Println("------")
       fmt.Println(slice_1)
       fmt.Println(slice_2)
       fmt.Println(slice_3)
       fmt.Println(slice_4)
       fmt.Println(slice_5)
}