package graph

import "fmt"

/*
 * 邻接表
 * 非并发安全
 */
type Graph struct {
	vexnum int
	edgnum int
	vexs   []*VNode
}

// 邻接表中表的顶点
type VNode struct {
	data      interface{}
	FirstEdge *ENode
}

// 邻接表中表对应的链表的顶点
type ENode struct {
	ivex int
	next *ENode
}

/*
 * 返回data在邻接表中的位置
 */
func (pg *Graph) GetPosition(data interface{}) (int, error) {
	for i := 0; i < pg.vexnum; i++ {
		if pg.vexs[i].data == data {
			return i, nil
		}

	}
	return -1, GetPositionErr
}

/*
 * 将node链接到list的末尾
 */
func (list *ENode) LinkLast(node *ENode) {
	var p *ENode
	p = list

	for {
		if p.next == nil {
			break
		}
		p = p.next
	}
	p.next = node
}

/*
 * 创建邻接表
 */
func CreateGraph(vexs []interface{}, edges [][]interface{}) (*Graph, error) {
	var (
		vlen   = len(vexs)
		elen   = len(edges)
		pG     = new(Graph)
		c1, c2 interface{}
		p1, p2 int
		node1  *ENode
		err    error
	)

	// 初始化"顶点数"和"边数"
	pG.vexnum = vlen
	pG.edgnum = elen
	pG.vexs = make([]*VNode, pG.vexnum)

	// 初始化"邻接表"的顶点
	for i := 0; i < pG.vexnum; i++ {
		pG.vexs[i] = new(VNode)
		pG.vexs[i].data = vexs[i]
	}

	// 初始化"邻接表"的边
	for i := 0; i < pG.edgnum; i++ {
		// 读取边的起始顶点和结束顶点
		c1 = edges[i][0]
		c2 = edges[i][1]

		p1, err = pG.GetPosition(c1)
		if err != nil {
			return nil, err
		}
		p2, err = pG.GetPosition(c2)
		if err != nil {
			return nil, err
		}

		// 初始化node1
		node1 = new(ENode)
		node1.ivex = p2

		// 将node1链接到"p1所在链表的末尾"
		if pG.vexs[p1].FirstEdge == nil {
			pG.vexs[p1].FirstEdge = node1
		} else {
			pG.vexs[p1].FirstEdge.LinkLast(node1)
		}
	}
	return pG, nil
}

/*
 * 打印邻接表图
 */

func (pG *Graph) Print() {
	var node *ENode

	fmt.Println("List Graph:")

	for i := 0; i < pG.vexnum; i++ {
		fmt.Printf("%v(%v): ", i, pG.vexs[i].data)
		node = pG.vexs[i].FirstEdge
		for {
			if node == nil {
				break
			}
			fmt.Printf("%v(%v) ", node.ivex, pG.vexs[node.ivex].data)
			node = node.next
		}
		fmt.Println()
	}
}

/*
 * DFS搜索邻接表
 */
func (pG *Graph) DFSTraverse() []interface{} {
	var (
		items   []interface{}
		visited = make([]bool, pG.vexnum)
	)

	for i := 0; i < pG.vexnum; i++ {
		if !visited[i] {
			DFS(pG, i, visited, &items)
		}
	}
	return items
}

/*
 * 深度优先搜索遍历图的递归实现
 */
func DFS(pG *Graph, i int, visited []bool, items *[]interface{}) {
	var (
		node *ENode
	)
	visited[i] = true

	*items = append(*items, pG.vexs[i].data)
	node = pG.vexs[i].FirstEdge
	for {
		if node == nil {
			break
		}
		if !visited[node.ivex] {
			DFS(pG, node.ivex, visited, items)
		}
		node = node.next
	}
}

/*
 * BFS搜索邻接表
 */
func (pG *Graph) BFS() []interface{} {
	var (
		items      []interface{}
		head, rear int
		queue      = make([]int, pG.vexnum)
		visited    = make([]bool, pG.vexnum)
		node       *ENode
		j, k       int
	)

	for i := 0; i < pG.vexnum; i++ {
		if !visited[i] {
			visited[i] = true
			items = append(items, pG.vexs[i].data)
			queue[rear] = i // 入队列
			rear++
		}
		for {
			if head == rear {
				break
			}
			// 出队列
			j = queue[head]
			head++
			node = pG.vexs[j].FirstEdge
			for {
				if node == nil {
					break
				}
				k = node.ivex
				if !visited[k] {
					visited[k] = true
					items = append(items, pG.vexs[k].data)
					queue[rear] = k
					rear++
				}
				node = node.next
			}
		}
	}
	return items
}

/*
 * 邻接表拓扑排序
 * 也可用于判断是否为有向无环图
 */
func (pG *Graph) TopologicalSort() ([]interface{}, error) {
	var (
		j                 int
		index, head, rear int = 0, 0, 0
		queue                 = make([]int, pG.vexnum)
		ins                   = make([]int, pG.vexnum)
		tops                  = make([]interface{}, pG.vexnum)
		num                   = pG.vexnum
		node              *ENode
	)

	// 统计每个顶点的入度数
	for i := 0; i < num; i++ {
		node = pG.vexs[i].FirstEdge
		for {
			if node == nil {
				break
			}
			ins[node.ivex]++
			node = node.next
		}
	}

	// 将所有入度为0的顶点入队列
	for i := 0; i < num; i++ {
		if ins[i] == 0 {
			queue[rear] = i // 入队列
			rear++
		}
	}
	for {
		if head == rear {
			break
		}
		j = queue[head] // 出队列
		head++
		tops[index] = pG.vexs[j].data
		index++

		node = pG.vexs[j].FirstEdge

		// 将与"node"关联的节点的入度减1；
		// 若减1之后，该节点的入度为0；则将该节点添加到队列中。
		for {
			if node == nil {
				break
			}
			ins[node.ivex]--
			if ins[node.ivex] == 0 {
				queue[rear] = node.ivex
				rear++
			}
			node = node.next
		}

	}

	if index != pG.vexnum {
		return tops, GraphCycleErr
	}

	return tops, nil
}

/*
 * 邻接表拓扑排序的变种
 * 输出层级关系的拓扑排序
 * 也可用于判断是否为有向无环图
 */
func (pG *Graph) TopologicalSortWithBFS() ([][]interface{}, error) {
	var (
		ins   = make([]int, pG.vexnum) // 入度表
		node  *ENode
		items [][]interface{}
		tops  []interface{}
	)

	// 统计每个顶点的入度数
	for i := 0; i < pG.vexnum; i++ {
		node = pG.vexs[i].FirstEdge
		for {
			if node == nil {
				break
			}
			ins[node.ivex]++
			node = node.next
		}
	}

	var queue []interface{}
	// 将所有入度为0的顶点入队列
	for i := 0; i < pG.vexnum; i++ {
		if ins[i] == 0 {
			queue = append(queue, pG.vexs[i].data)
		}
	}
	items = append(items, queue)
	tops = append(tops, queue...)

	for {
		var tmpQueue []interface{}
		for _, j := range queue {
			k, err := pG.GetPosition(j)
			if err != nil {
				return items, GetPositionErr
			}
			q := findZeroIns(pG, k, ins)
			tmpQueue = append(tmpQueue, q...)
		}
		queue = tmpQueue

		if len(queue) == 0 {
			break
		}
		items = append(items, queue)
		tops = append(tops, queue...)
	}
	if len(tops) != pG.vexnum {
		return items, GraphCycleErr
	}
	return items, nil
}

/*
 * 寻找当前邻接表中的入度为0的节点
 */
func findZeroIns(pG *Graph, i int, ins []int) []interface{} {
	var node *ENode
	var queue []interface{}

	node = pG.vexs[i].FirstEdge
	for {
		if node == nil {
			break
		}
		ins[node.ivex]--
		if ins[node.ivex] == 0 {
			queue = append(queue, pG.vexs[node.ivex].data)
		}
		node = node.next
	}
	return queue
}

/*
 * 输出指定data的所有依赖
 */
func (pG *Graph) AllInsVex(data interface{}) []interface{} {
	var node *ENode
	var item []interface{}
	for i := 0; i < pG.vexnum; i++ {
		node = pG.vexs[i].FirstEdge
		for {
			if node == nil {
				break
			}
			if pG.vexs[node.ivex].data == data {
				item = append(item, pG.vexs[i].data)
				break
			}
			node = node.next
		}
	}
	return item
}
