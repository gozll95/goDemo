package main 

import(
	"fmt"
	"sync/atomic"
)

func main(){
	var countVal atomic.Value
	a:=1
	//countVal.Store([]int{1,3,5,7})
	countVal.Store(a)
	anotherStore(&countVal)
	fmt.Printf("The count value:%+v\n",countVal.Load())

}

func anotherStore(countVal *atomic.Value){
	//countVal.Store([]int{2,4,6,8})
	countVal.Store(222)
}

//对atomic.Value的参数操作不会影响原来的值。
//需要使用指针变量

//其实,对于sync包中的Mutext、RWMutex和Cond类型,go vet命令同样会检查此类复制问题,其原因也是相似的。
//一个比较彻底的解决方案是,避免直接使用它们,而使用它们的指针值