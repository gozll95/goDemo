 go doc
package queue // import “imooc.com/ccmouse/learngo/queue”

type Queue []int

go doc Queue
type Queue []int 
func(q *Queue)IsEmpty() bool
func(q *Queue)Pop() int
func(q *Queue)Push(v int)

go doc IsEmpty
func(q *Queue)IsEmpty() bool

godoc -http :6060

用//写注释

写示例代码

queue_test.go

func ExampleQueue_Pop(){}