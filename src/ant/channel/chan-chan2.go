import "container/heap"


// chan
func HeapPop(h heap.Interface) interface{} {
	var result = make(chan interface{})
	heapPopChan <- heapPopChanMsg{
		h:      h,
		result: result,
	}
	return <-result
}

select{
	case popMsg := <-heapPopChan:
	popMsg.result <- heap.Pop(popMsg.h) 
}

func aa()<-chan heapPopChanMsg{
	a:=make(chan heapPopChanMsg)
	go func(){
		a.result<-heap.Pop(a.h)
	}()
	return a
}


/* 
 * chan-chan, select 外层 chan,赋值里面chan
 * 另外一个程序等里面的chan
 * 好处: 在同一个函数里,注入外层chan,等待里面的chan
 *  func(){
 *	 c<-a
 *	 <-c.result
 *  }
 *  而这个给出c的程序包括了对c的处理过程
 * 
 * 总结: 我:go处理函数,返回外层c
 *      你赋值外层c,等待内层c
 */



