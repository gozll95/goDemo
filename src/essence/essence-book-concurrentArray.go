# book-concurrentArray 


"并发性优先考虑原子操作,所以我们的model里放的是atomic.Value"

直接看代码吧
type ConcurrentArray interface {
	// Set 用于设置指定索引上的元素值。
	Set(index uint32, elem int) (err error)
	// Get 用于获取指定索引上的元素值。
	Get(index uint32) (elem int, err error)
	// Len 用于获取数组的长度。
	Len() uint32
}

type intArray struct {
	length uint32
	val    atomic.Value
}


// NewConcurrentArray 会创建一个ConcurrentArray类型值。
func NewConcurrentArray(length uint32) ConcurrentArray {
	array := intArray{}
	array.length = length
	array.val.Store(make([]int, array.length))
	return &array
}

"这里是重点"
- "copy" "参考 array-slice-copy-append.go"
- "Copy on write这个思想"
- "考虑如何做并发性测试"

func (array *intArray) Set(index uint32, elem int) (err error) {
	if err = array.checkIndex(index); err != nil {
		return
	}
	if err = array.checkValue(); err != nil {
		return
	}

	// 不要这样做！否则会形成竞态条件！
	// oldArray := array.val.Load().([]int)
	// oldArray[index] = elem
	// array.val.Store(oldArray)

	newArray := make([]int, array.length)
	copy(newArray, array.val.Load().([]int))
	newArray[index] = elem
	array.val.Store(newArray)
	return
}



func (array *intArray) Get(index uint32) (elem int, err error) {
	if err = array.checkIndex(index); err != nil {
		return
	}
	if err = array.checkValue(); err != nil {
		return
	}
	elem = array.val.Load().([]int)[index]
	return
}

func (array *intArray) Len() uint32 {
	return array.length
}