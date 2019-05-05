package main

//Heap
//author:Xiong Chuan Liang
//date:2015-2-4

import (
	"fmt"
)

var(
  heap = []int{100,16,4,8,70,2,36,22,5,12}	
)


func main(){
	fmt.Println("\n数组:")
	Print(heap)

    MakeHeap()
    fmt.Println("\n构建树后:")
	Print(heap)

	fmt.Println("\n增加 90,30,1 :")
	Push(90)
	Push(30)
	Push(1)	
	Print(heap)	
	
	n := Pop()
	fmt.Println("\nPop出最小值(",n,")后:")
	Print(heap)	

	fmt.Println("\nRemove()掉idx为3即值",heap[3-1],"后:")
	Remove(3)	
	Print(heap)	

	fmt.Println("\nHeapSort()后:")
	HeapSort()
	Print(heap)	

}


func Print(arr []int){		
	for _,v := range arr {
		fmt.Printf("%d ",v)
	}
}

//构建堆
func MakeHeap(){	
	n := len(heap)	
	for i := n/2-1 ;i >= 0;i--{		
		down(i, n)
	}
}


//由父节点至子节点依次建堆
//parent      : i
//left child  : 2*i+1
//right child : 2*i+2
func down(i,n int) {
	//构建最小堆,父小于两子节点值
	for {	

		j1 := 2*i + 1
		if j1 >= n || j1 < 0 { // j1 < 0 after int overflow
			break
		}

		//找出两个节点中最小的(less: a<b)
		j := j1 // left child		
		if j2 := j1 + 1; j2 < n && !Less(j1, j2) {
			j = j2 // = 2*i + 2  // right child
		}

		//然后与父节点i比较,如果父节点小于这个子节点最小值,则break,否则swap
		if !Less(j, i) {
			break
		}
		Swap(i, j)
		i = j	
	}
}


//由子节点至父节点依次重建堆
func up(j int)  {
	
	for {
		i := (j - 1) / 2 // parent      
    
		if i == j || !Less(j, i) { 
			//less(子,父) !Less(9,5) == true 
			//父节点小于子节点,符合最小堆条件,break
			break
		}
		//子节点比父节点小,互换
		Swap(i, j)
		j = i
	}
}

func Push(x interface{}){
	heap = append(heap, x.(int))
	up(len(heap)-1)
	return 
}

func Pop() interface{} {
	n := len(heap) - 1
	Swap(0, n)
	down(0, n)

	old :=heap
	n = len(old)
	x := old[n-1]
	heap = old[0 : n-1]
	return x
}

func Remove(i int) interface{} {
	n := len(heap) - 1
	if n != i {
		Swap(i, n)
		down(i, n)
		up(i)
	}
	return Pop()
}

func Less(a,b int)bool{
	return heap[a] < heap[b]


func Swap(a,b int){
	heap[a],heap[b] = heap[b],heap[a]
}

func HeapSort(){
	//升序 Less(heap[a] > heap[b])	//最大堆
	//降序 Less(heap[a] < heap[b])	//最小堆
	for i := len(heap)-1 ;i > 0;i--{	
		//移除顶部元素到数组末尾,然后剩下的重建堆,依次循环
		Swap(0, i)
		down(0, i)
	}
}


/*
输出结果:
-----------------------------------
数组:
100 16 4 8 70 2 36 22 5 12
构建树后:
2 5 4 8 12 100 36 22 16 70
增加 90,30,1 :
1 5 2 8 12 4 36 22 16 70 90 100 30
Pop出最小值( 1 )后:
2 5 4 8 12 30 36 22 16 70 90 100
Remove()掉idx为3即值 4 后:
4 5 8 16 12 30 36 22 100 70 90
HeapSort()后:
100 90 70 36 30 22 16 12 8 5 4

*/