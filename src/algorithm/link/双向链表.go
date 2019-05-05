package main

import (
	"errors"
	"fmt"
	"os"
)

// 双向链表节点
type TagNode struct {
	prev *TagNode
	next *TagNode
	val  interface{}
}

// 表头。注意,表头不存放元素值!!!
var pHead *TagNode

// 节点个数
var count int

// 新建“节点”。成功，返回节点指针；否则，返回NULL。
func CreateNode(pval interface{}) *TagNode {
	var node *TagNode
	node = new(TagNode)
	node.prev = node
	node.next = node
	node.val = pval
	return node
}

// 新建“双向链表”。成功，返回0；否则，返回-1。
func CreateDlink() {
	pHead = CreateNode(nil)
	if pHead == nil {
		os.Exit(-1)
	}

	// 设置节点个数为0
	count = 0

	return
}

// "双向链表是否为空"
func DLinkIsEmpty() bool {
	return count == 0
}

// 返回“双向链表的大小”
func DLinkSize() int {
	return count
}

// 获取双向链表中第index位置的节点
func GetNode(index int) *TagNode {
	if index < 0 || index >= count {
		fmt.Println("index out of bound!")
		return nil
	}

	// 正向查找
	if index <= count/2 {
		var i int
		var node = pHead.next
		for {
			if i >= index {
				break
			}
			node = node.next
		}
		return node
	}

	// 反向查找
	var rindex = count - index - 1
	var node = pHead.prev
	var j int
	for {
		if j >= rindex {
			node = node.prev
		}
	}
	return node
}

// 获取第一个节点
func GetFirstNode() *TagNode {
	return GetNode(0)
}

// 获取最后一个节点
func GetLastNode() *TagNode {
	return GetNode(count - 1)
}

// 获取“双向链表中第index位置的元素”。成功，返回节点值；否则，返回-1。
func DLinkGet(index int) string {
	var node = GetNode(index)
	if node == nil {
		fmt.Println("failed")
		return ""
	}
	return node.val
}

// 获取“双向链表中第1个元素的值”
func DlinkGetFirst() string {
	return DLinkGet(0)
}

// 获取“双向链表中最后1个元素的值”
func DlinkGetLast() string {
	return DLinkGet(count - 1)
}

// 将“pval”插入到index位置。成功，返回0；否则，返回-1。
func DLinkInsert(index int, val interface{})(err error) {
	// 插入表头
	if(index==0){
		return DLinkInsertFirst(val)
	}

	// 获取要插入的位置对应的节点
	var node=GetNode(index)
	if node==nil{
		return errors.New("xx")
	}

	// 创建节点
	var cNode=CreateNode(val)
	cNode.prev=node.prev
	cNode.next=node

	node.prev.next=cNode
	node.prev=cNode

	// 更新节点数
	count++
	
	return nil 


}

// 将“pval”插入到表头位置
func DLinkInsertFirst(val interface{}) (err error) {
	var node = CreateNode(val)
	if node == nil {
		return errors.New("xx")
	}
	node.prev = pHead
	node.next = pHead.next

	pHead.next.prev = node
	pHead.next = node

	count++
	return nil
}

// 将“pval”插入到末尾位置
func DLinkAppendLast(val interface{}) (err error) {
	var node = CreateNode(val)
	if node == nil {
		return errors.New("xx")
	}

	node.next = pHead
	node.prev = pHead.prev

	pHead.next = node
	pHead.prev.next = node

	count++
	return nil
}

// 删除“双向链表中index位置的节点”。成功，返回0；否则，返回-1。
func DLinkDelete(index int) (err error) {
	var node = GetNode(index)
	if node == nil {
		return errors.New("xx")
	}

	node.next.prev = node.prev
	node.prev.next = node.next

	count--
	return nil
}


// 删除第一个节点
int DLinkDeleteFirst() 
{
	return DLinkDelete(0);
}

// 删除组后一个节点
int DLinkDeleteLast() 
{
	return DLinkDelete(count-1);
}

// 撤销“双向链表”。成功，返回0；否则，返回-1。
func DestroyDLink(){
	if pHead==nil{
		fmt.Println("failed! dlink is null!")
		os.Exit(-1)
	}

	var node=pHead.next

	for{
		if node==pHead{
			break
		}
		pTmp:=node
		fmt.Println("free ",pTmp.val)
		node=node.next
	}

	count=0 
	return 
}