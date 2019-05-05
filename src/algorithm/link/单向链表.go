package main

import "fmt"

// 单向链表的"节点"
type Node struct {
	pval string
	next *Node
}

// 单向链表的"表头"
var pHead = &Node{
	pval: "origin head",
	next: nil,
}

// 创建节点，p为节点值
func CreateNode(pval string) *Node {
	var node *Node
	node = new(Node)
	node.pval = pval
	node.next = nil
	return node
}

// 销毁单向链表
func DestroySingleLink() {
	var pNode *Node
	for {
		if pHead == nil {
			break
		}
		pNode = pHead
		fmt.Println("free ", pNode.pval)
		pHead = pHead.next
	}
}

// 将pval插入到链表的表头位置
func Push(pval string) *Node {
	var node *Node
	node = CreateNode(pval)
	node.next = pHead
	pHead = node
	return pHead
}

// 删除链表的表头
func Pop() string {
	var node *Node

	if pHead == nil {
		fmt.Println("remove failed! link is empty!")
		return ""
	}

	var pret string
	pret = pHead.pval

	node = pHead
	pHead = pHead.next
	fmt.Println("free ", node.pval)

	return pret
}

// 返回链表的表头节点的值
func Peek() string {
	if pHead == nil {
		fmt.Println("peek failed! link is empty!")
		return ""
	}

	return pHead.pval
}

// 返回链表中节点的个数
func Size() int {
	var node *Node
	node = pHead

	var count int
	for {
		if node == nil {
			break
		}
		node = node.next
		count++
	}
	return count
}

// 判断链表是否为空
func IsEmpty() bool {
	return pHead == nil
}

// 打印“栈”
func PrintSingleLink() {
	if IsEmpty() {
		fmt.Println("stack is Empty\n")
		return
	}

	fmt.Printf("stack size()=%d\n", Size())

	var node *Node
	node = pHead

	for {
		if node == nil {
			break
		}
		fmt.Printf("%v\n", node.pval)
		node = node.next
	}
}

func main() {
	// 将 A,B,C 依次推入栈中
	Push("A")
	Push("B")
	Push("C")

	// 打印
	PrintSingleLink()

	// 将“栈顶元素”赋值给tmp，并删除“栈顶元素”
	tmp := Pop()
	fmt.Println(tmp)
	PrintSingleLink()

	// 只将“栈顶”赋值给tmp，不删除该元素.
	tmp = Peek()
	fmt.Println(tmp)

	Push("D")
	PrintSingleLink()

	// 销毁栈
	DestroySingleLink()
}
