package main

import "fmt"

// 邻接表
type Graph struct {
	vexnum int
	edgnum int
	vexs   []*VNode
}

// 邻接表中表的顶点
type VNode struct {
	data      string
	FirstEdge *ENode
}

// 邻接表中表对应的链表的顶点
type ENode struct {
	ivex int
	next *ENode
}

/*
 * 返回ch在邻接表中的位置
 */
func getPosition(g *Graph, ch string) int {
	for i := 0; i < g.vexnum; i++ {
		if g.vexs[i].data == ch {
			return i
		}

	}
	return -1
}

/*
 * 将node链接到list的末尾
 */
func LinkLast(list, node *ENode) {
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

func CreateGraph() *Graph {
	var vexs []string
	vexs = []string{"A", "B", "C", "D", "E", "F", "G"}

	var edges [][]string
	tmp1 := []string{"A", "G"}
	edges = append(edges, tmp1)
	tmp1 = []string{"B", "A"}
	edges = append(edges, tmp1)
	tmp1 = []string{"B", "D"}
	edges = append(edges, tmp1)
	tmp1 = []string{"C", "F"}
	edges = append(edges, tmp1)
	tmp1 = []string{"C", "G"}
	edges = append(edges, tmp1)
	tmp1 = []string{"D", "E"}
	edges = append(edges, tmp1)
	tmp1 = []string{"D", "F"}
	edges = append(edges, tmp1)

	var (
		vlen   = len(vexs)
		elen   = len(edges)
		pG     = new(Graph)
		c1, c2 string
		p1, p2 int
		node1  *ENode
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

		p1 = getPosition(pG, c1)
		p2 = getPosition(pG, c2)

		// 初始化node1
		node1 = new(ENode)
		node1.ivex = p2

		// 将node1链接到"p1所在链表的末尾"
		if pG.vexs[p1].FirstEdge == nil {
			pG.vexs[p1].FirstEdge = node1
		} else {
			LinkLast(pG.vexs[p1].FirstEdge, node1)
		}
	}
	return pG
}

/*
 * 打印邻接表图
 */
func PrintGraph(pG *Graph) {
	var node *ENode
	fmt.Println("List Graph:\n")

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
		fmt.Println("\n")
	}
}

func DFSTraverse(pG *Graph) {
	var visited = make([]bool, pG.vexnum)

	fmt.Println("DFS: ")
	for i := 0; i < pG.vexnum; i++ {
		if !visited[i] {
			DFS(pG, i, visited)
		}
	}
	fmt.Println("\n")

}

/*
 * 深度优先搜索遍历图的递归实现
 */
func DFS(pG *Graph, i int, visited []bool) {
	var (
		node *ENode
	)
	visited[i] = true
	fmt.Printf("%v", pG.vexs[i].data)
	node = pG.vexs[i].FirstEdge
	for {
		if node == nil {
			break
		}
		if !visited[node.ivex] {
			DFS(pG, node.ivex, visited)
		}
		node = node.next
	}
}

func BFS(pG *Graph) {
	var (
		head, rear int
		queue      = make([]int, pG.vexnum)
		visited    = make([]bool, pG.vexnum)
		node       *ENode
		j, k       int
	)

	fmt.Println("BFS:  ")
	for i := 0; i < pG.vexnum; i++ {
		if !visited[i] {
			visited[i] = true
			fmt.Printf("%v", pG.vexs[i].data)
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
					fmt.Printf("%v", pG.vexs[k].data)
					queue[rear] = k
					rear++
				}
				node = node.next
			}
		}

	}
	fmt.Println("\n")
}

func TopologicalSort(pG *Graph) {
	var (
		j                 int
		index, head, rear int = 0, 0, 0
		queue                 = make([]int, pG.vexnum)
		ins                   = make([]int, pG.vexnum)
		tops                  = make([]string, pG.vexnum)
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
		fmt.Println("Graph has a cycle\n")
		return
	}

	// 打印拓扑排序的结果
	fmt.Println("== TopSort: ")
	for i := 0; i < num; i++ {
		fmt.Printf("%v ", tops[i])
	}
	fmt.Println("\n")
}

func main() {
	pG := CreateGraph()
	PrintGraph(pG)
	DFSTraverse(pG)
	BFS(pG)
	TopologicalSort(pG)
}

/*
List Graph:

0(A): 6(G)

1(B): 0(A) 3(D)

2(C): 5(F) 6(G)

3(D): 4(E) 5(F)

4(E):

5(F):

6(G):

DFS:
AGBDEFC

BFS:
AGBDEFC

== TopSort:
B C A D G E F
*/
