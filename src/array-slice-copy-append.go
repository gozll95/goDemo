package main

import "fmt"

func test1() {
	s := make([]int, 5)
	ss := make([]int, 5)
	s1 := []int{1, 2, 3, 4, 5}
	copy(s, s1)
	copy(ss, s1)

	fmt.Println("s1 == ", s1)
	fmt.Println("s(copy from s1) == ", s)
	fmt.Println("ss(copy from s1) == ", ss)

	fmt.Println("set s[0] 11")

	s[0] = 11

	fmt.Println("s1 == ", s1)
	fmt.Println("s(copy from s1) == ", s)
	fmt.Println("ss(copy from s1) == ", ss)

	fmt.Println("set ss[0] 22")

	ss[0] = 22

	fmt.Println("s1 == ", s1)
	fmt.Println("s(copy from s1) == ", s)
	fmt.Println("ss(copy from s1) == ", ss)

	/*
		s1 ==  [1 2 3 4 5]
		s(copy from s1) ==  [1 2 3 4 5]
		ss(copy from s1) ==  [1 2 3 4 5]
		set s[0] 11
		s1 ==  [1 2 3 4 5]
		s(copy from s1) ==  [11 2 3 4 5]
		ss(copy from s1) ==  [1 2 3 4 5]
		set ss[0] 22
		s1 ==  [1 2 3 4 5]
		s(copy from s1) ==  [11 2 3 4 5]
		ss(copy from s1) ==  [22 2 3 4 5]
	*/
}
func test2() {

	s := [5]int{1, 2, 3, 4, 5}
	s1 := s[1:3]
	s2 := s[0:2]

	copy(s2, s1)

	fmt.Println("s: ", s)
	fmt.Println("s1: ", s1)
	fmt.Println("s2: ", s2)

	s1[0] = 9

	fmt.Println("s: ", s)
	fmt.Println("s1: ", s1)
	fmt.Println("s2: ", s2)

	/*
		s:  [2 3 3 4 5]
		s1:  [3 3]
		s2:  [2 3]
		s:  [2 9 3 4 5]
		s1:  [9 3]
		s2:  [2 9]
	*/
}

func test3() {
	s := [5]int{1, 2, 3, 4, 5}
	s1 := s[1:3]
	s2 := make([]int, 2) // make已经给s2分配好底层array的空间，并且用0填补array

	fmt.Println("s: ", s)
	fmt.Println("s1: ", s1)
	fmt.Println("s2: ", s2)

	copy(s2, s1)

	fmt.Println("s: ", s)
	fmt.Println("s1: ", s1)
	fmt.Println("s2: ", s2)

	s2[0] = 9
	s1[0] = 99

	fmt.Println("s: ", s)
	fmt.Println("s1: ", s1)
	fmt.Println("s2: ", s2)

	/*
		s:  [1 2 3 4 5]
		s1:  [2 3]
		s2:  [0 0]
		s:  [1 2 3 4 5]
		s1:  [2 3]
		s2:  [2 3]
		s:  [1 99 3 4 5]
		s1:  [99 3]
		s2:  [9 3]
	*/
}

func test4() {
	s := [5]int{1, 2, 3, 4, 5}
	s1 := s[1:3]
	s3 := append(s1, 9)

	fmt.Println("s: ", s)
	fmt.Println("s1: ", s1)
	fmt.Println("s3: ", s3)

	s3[0] = 99

	fmt.Println("s: ", s)
	fmt.Println("s1: ", s1)
	fmt.Println("s3: ", s3)

	arr := [5]int{1, 2, 3, 4, 5}
	arr1 := arr[2:]
	arr2 := append(arr1, 9)

	fmt.Println("arr: ", arr)
	fmt.Println("arr1: ", arr1)
	fmt.Println("arr2: ", arr2)

	arr2[0] = 99

	fmt.Println("arr: ", arr)
	fmt.Println("arr1: ", arr1)
	fmt.Println("arr2: ", arr2)

	// 对切片的数组其实就是对切片对应的底层数组的修改

	//APPEND总结：s2 := append(s1, *)是切片s1上记录的切片信息复制给s2，如果s1指向的底层array长度不够，append的过程会发生如下操作：内存中不仅新开辟一块区域存储append后的切片信息，而且需要新开辟一块区域存储底层array（复制原来的array至这块新array中），最后再append新数据进新array中，这样，s2指向新array。反之，s2和s1指向同一个array，append的结果是内存中新开辟一个区域存储新切片信息。

	// s: 1 2 3 4 5
	// s1:  2 3
	// s3 := append(s1, 9)
	// s的长度是够的
	// 那么
	// s3与s1将都指向s
	// s: 1 2 3 4 5
	// s3:  2 3 9
	// s1:  2 3

	// s[3]=99
	// s: 1 99 3 9 5
	// s3:  99 3 9
	// s1:  99 3

	// array: 1 2 3 4 5
	// arr1:    2 3 4 5
	// arr2 := append(arr1, 9)
	// 发现arr1指向的底层array的长度不够
	// 1.开辟新区域 首先将array复制成array2
	// array2: 1 2 3 4 5
	// 再加
	// array2: 1 2 3 4 5 9
	// arr2:     2 3 4 5 9

	// arr2[0]=99
	// arr2:     99 3 4 5 9

}

func test5() {
	a1 := make([]int, 5, 5)
	a := a1[:3] // 不管是直接赋值，还是函数式传值，都是新建一个切片数据，与原切片指向同一个底层array
	b := a
	c := append(a, 1)
	c[0] = 9

	fmt.Println("a1:", a1)
	fmt.Println("a:", a)
	fmt.Println("b:", b)
	fmt.Println("c:", c)

	/*
		a1: [9 0 0 1 0]
		a: [9 0 0]
		b: [9 0 0]
		c: [9 0 0 1]
	*/
}

func test6() {
	slice := make([]int, 5, 8)
	newSlice := append(slice, 1)
	test(newSlice)
	fmt.Println("main slice: ", slice, len(slice), cap(slice))
	fmt.Println("main newSlice: ", newSlice, len(newSlice), cap(newSlice))

	/*
		test:  [10 0 0 0 0 1 2] 7 8
		main slice:  [10 0 0 0 0] 5 8
		main newSlice:  [10 0 0 0 0 1] 6 8
	*/
}

func test(test []int) {
	test = append(test, 2)
	test[0] = 10
	fmt.Println("test: ", test, len(test), cap(test))
}

func test7() {
	old := make([]int, 5)
	old[0] = 1
	new1 := make([]int, 5)
	copy(new1, old)
	new1[1] = 2
	fmt.Println(old)
	fmt.Println(new1)

	new2 := make([]int, 5)
	copy(new2, old)
	new2[2] = 3
	fmt.Println(old)
	fmt.Println(new2)

}

func main() {
	test7()

}

// 结论:
/*
	a1 := make([]int, 5, 5)
	a := a1[:3] // 不管是直接赋值，还是函数式传值，都是新建一个切片数据，与原切片指向同一个底层array


	copy(s2, s1)过程只是将切片s1指向底层array的数据copy至切片s2指向底层array的数据上

	// 对切片的数组其实就是对切片对应的底层数组的修改

	//APPEND总结：s2 := append(s1, *)是切片s1上记录的切片信息复制给s2，如果s1指向的底层array长度不够，append的过程会发生如下操作：内存中不仅新开辟一块区域存储append后的切片信息，而且需要新开辟一块区域存储底层array（复制原来的array至这块新array中），最后再append新数据进新array中，这样，s2指向新array。反之，s2和s1指向同一个array，append的结果是内存中新开辟一个区域存储新切片信息。

	// s: 1 2 3 4 5
	// s1:  2 3
	// s3 := append(s1, 9)
	// s的长度是够的
	// 那么
	// s3与s1将都指向s
	// s: 1 2 3 4 5
	// s3:  2 3 9
	// s1:  2 3

	// s[3]=99
	// s: 1 99 3 9 5
	// s3:  99 3 9
	// s1:  99 3

	// array: 1 2 3 4 5
	// arr1:    2 3 4 5
	// arr2 := append(arr1, 9)
	// 发现arr1指向的底层array的长度不够
	// 1.开辟新区域 首先将array复制成array2
	// array2: 1 2 3 4 5
	// 再加
	// array2: 1 2 3 4 5 9
	// arr2:     2 3 4 5 9

	// arr2[0]=99
	// arr2:     99 3 4 5 9
*/
