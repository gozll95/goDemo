// 利用Go的 container/heap 可以很方便的实现堆排序的队列，heap包中的示例代码有一个优先级队列的实现，但是并不是线程安全的，因此，要实现线程安全的优先级队列，需要堆示例代码稍作修改，

//下面为 container/heap 包的示例代码：

// This example demonstrates a priority queue built using the heap interface.
package main

import (
	"container/heap"
	"fmt"
)

// An Item is something we manage in a priority queue.
type Item struct {
	value    string // The value of the item; arbitrary.
	priority int    // The priority of the item in the queue.
	// The index is needed by update and is maintained by the heap.Interface methods.
	index int // The index of the item in the heap.
}

// A PriorityQueue implements heap.Interface and holds Items.
type PriorityQueue []*Item

func (pq PriorityQueue) Len() int { return len(pq) }
func (pq PriorityQueue) Less(i, j int) bool {
	// We want Pop to give us the highest, not lowest, priority so we use greater than here.
	return pq[i].priority > pq[j].priority
}
func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}
func (pq *PriorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*Item)
	item.index = n
	*pq = append(*pq, item)
}
func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	item.index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

// update modifies the priority and value of an Item in the queue.
func (pq *PriorityQueue) update(item *Item, value string, priority int) {
	item.value = value
	item.priority = priority
	heap.Fix(pq, item.index)
}

// This example creates a PriorityQueue with some items, adds and manipulates an item,
// and then removes the items in priority order.
func main() {
	// Some items and their priorities.
	items := map[string]int{
		"banana": 3, "apple": 2, "pear": 4,
	}
	// Create a priority queue, put the items in it, and
	// establish the priority queue (heap) invariants.
	pq := make(PriorityQueue, len(items))
	i := 0
	for value, priority := range items {
		pq[i] = &Item{
			value:    value,
			priority: priority,
			index:    i,
		}
		i++
	}
	heap.Init(&pq)
	// Insert a new item and then modify its priority.
	item := &Item{
		value:    "orange",
		priority: 1,
	}
	heap.Push(&pq, item)
	pq.update(item, item.value, 5)
	// Take the items out; they arrive in decreasing priority order.
	for pq.Len() > 0 {
		item := heap.Pop(&pq).(*Item)
		fmt.Printf("%.2d:%s ", item.priority, item.value)
	}
}

// 以上是网友写的并发安全的优先级队列

package priorityqueue
import (
    "container/heap"
    "sort"
    "sync"
)
// An Item is something we manage in a priority queue.
type Item struct {
    Key      interface{} //The unique key of the item.
    Value    interface{} // The value of the item; arbitrary.
    Priority int    // The priority of the item in the queue.
    // The index is needed by update and is maintained by the heap.Interface methods.
    Index int // The index of the item in the heap.
}
type ItemSlice struct {
    items    []*Item
    itemsMap map[interface{}]*Item
}
func (s ItemSlice) Len() int { return len(s.items) }
func (s ItemSlice) Less(i, j int) bool {
    return s.items[i].Priority < s.items[j].Priority
}
func (s ItemSlice) Swap(i, j int) {
    s.items[i], s.items[j] = s.items[j], s.items[i]
    s.items[i].Index = i
    s.items[j].Index = j
    if s.itemsMap != nil {
        s.itemsMap[s.items[i].Key] = s.items[i]
        s.itemsMap[s.items[j].Key] = s.items[j]
    }
}
func (s *ItemSlice) Push(x interface{}) {
    n := len(s.items)
    item := x.(*Item)
    item.Index = n
    s.items = append(s.items, item)
    s.itemsMap[item.Key] = item
}
func (s *ItemSlice) Pop() interface{} {
    old := s.items
    n := len(old)
    item := old[n-1]
    item.Index = -1 // for safety
    delete(s.itemsMap, item.Key)
    s.items = old[0 : n-1]
    return item
}
// update modifies the priority and value of an Item in the queue.
func (s *ItemSlice) update(key interface{}, value interface{}, priority int) {
    item := s.itemByKey(key)
    if item != nil {
        s.updateItem(item, value, priority)
    }
}
// update modifies the priority and value of an Item in the queue.
func (s *ItemSlice) updateItem(item *Item, value interface{}, priority int) {
    item.Value = value
    item.Priority = priority
    heap.Fix(s, item.Index)
}
func (s *ItemSlice) itemByKey(key interface{}) *Item {
    if item, found := s.itemsMap[key]; found {
        return item
    }
    return nil
}
// A PriorityQueue implements heap.Interface and holds Items.
type PriorityQueue struct {
    slice   ItemSlice
    maxSize int
    mutex   sync.RWMutex
}
func (pq *PriorityQueue) Init(maxSize int) {
    pq.slice.items = make([]*Item, 0, pq.maxSize)
    pq.slice.itemsMap = make(map[interface{}]*Item)
    pq.maxSize = maxSize
}
func (pq PriorityQueue) Len() int {
    pq.mutex.RLock()
    size := pq.slice.Len()
    pq.mutex.RUnlock()
    return size
}
func (pq *PriorityQueue) minItem() *Item {
    len := pq.slice.Len()
    if len == 0 {
        return nil
    }
    return pq.slice.items[0]
}
func (pq *PriorityQueue) MinItem() *Item {
    pq.mutex.RLock()
    defer pq.mutex.RUnlock()
    return pq.minItem()
}
func (pq *PriorityQueue) PushItem(key, value interface{}, priority int) (bPushed bool) {
    pq.mutex.Lock()
    defer pq.mutex.Unlock()
    size := pq.slice.Len()
    item := pq.slice.itemByKey(key)
    if size > 0 && item != nil {
        pq.slice.updateItem(item, value, priority)
        return true
    }
    item = &Item{
        Value:    value,
        Key:      key,
        Priority: priority,
        Index:    -1,
    }
    if pq.maxSize <= 0 || size < pq.maxSize {
        heap.Push(&(pq.slice), item)
        return true
    }
    min := pq.minItem()
    if min.Priority >= priority {
        return false
    }
    heap.Pop(&(pq.slice))
    heap.Push(&(pq.slice), item)
    return true
}
func (pq PriorityQueue) GetQueue() []interface{} {
    items := pq.GetQueueItems()
    values := make([]interface{}, len(items))
    for i := 0; i < len(items); i++ {
        values[i] = items[i].Value
    }
    return values
}
func (pq PriorityQueue) GetQueueItems() []*Item {
    size := pq.Len()
    if size == 0 {
        return []*Item{}
    }
    s := ItemSlice{}
    s.items = make([]*Item, size)
    pq.mutex.RLock()
    for i := 0; i < size; i++ {
        s.items[i] = &Item{
            Value:    pq.slice.items[i].Value,
            Priority: pq.slice.items[i].Priority,
        }
    }
    pq.mutex.RUnlock()
    sort.Sort(sort.Reverse(s))
    return s.items
}