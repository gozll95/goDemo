package main

import "fmt"

type Graph struct {
	vexs   []string //定点集合
	vexnum int      //定点数量
	edgnum int      //边数量
	matrix [][]int  //邻接矩阵
}

func getPosition(g *Graph, ch string) int {
	var i int
	for i := 0; i < g.vexnum; i++ {
		if g.vexs[i] == ch {
			return i
		}
	}
	i = -1
	return i
}

func getGraph() *Graph {
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

	var vlen int = len(vexs)
	var elen = len(edges)
	var p1, p2 int

	pG := new(Graph)

	// 初始化"顶点数"和"边数"
	pG.vexnum = vlen
	pG.edgnum = elen

	// 初始化"顶点"
	for i := 0; i < pG.vexnum; i++ {
		pG.vexs = append(pG.vexs, vexs[i])
	}

	// 初始化 pG.matrix
	for i := 0; i < pG.edgnum; i++ {
		a := make([]int, pG.edgnum)
		pG.matrix = append(pG.matrix, a)
	}

	// 初始化边
	for i := 0; i < pG.edgnum; i++ {
		p1 = getPosition(pG, edges[i][0])
		p2 = getPosition(pG, edges[i][1])

		pG.matrix[p1][p2] = 1

	}

	return pG

}

func printGraph(pG *Graph) {
	fmt.Println("Martix Graph:\n")
	for i := 0; i < pG.vexnum; i++ {
		for j := 0; j < pG.vexnum; j++ {
			fmt.Printf("%d ", pG.matrix[i][j])
		}
		fmt.Println("\n")
	}
}

/*
 * 深度优先搜索遍历图
 */
func DFSTraverse(pG *Graph) {
	var visited []bool // 顶点访问标记

	// 初始化所有顶点都没有被访问
	for i := 0; i < pG.vexnum; i++ {
		visited = append(visited, false)
	}

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
	var w int

	visited[i] = true

	fmt.Printf("%v ", pG.vexs[i])

	// 遍历该顶点的所有邻接顶点。若是没有访问过，那么继续往下走
	for w = firstVertex(pG, i); w >= 0; w = nextVertix(pG, i, w) {
		if !visited[w] {
			DFS(pG, w, visited)
		}
	}
}

/*
 * 广度优先搜索（类似于树的层次遍历）
 */
func BFS(pG *Graph) {
	var (
		head int = 0
		rear int = 0
		j    int
		k    int
	)

	var (
		queue   []int
		visited []bool
	)

	for i := 0; i < pG.vexnum; i++ {
		visited = append(visited, false)
	}
	fmt.Println("BFS: ")

	for i := 0; i < pG.vexnum; i++ {
		if !visited[i] {
			visited[i] = true
			fmt.Printf("%v ", pG.vexs[i])
			queue = append(queue, i)
			rear++
		}
		for {
			if head == rear {
				break
			}
			j = queue[head]
			head++
			for k = firstVertex(pG, j); k >= 0; k = nextVertix(pG, j, k) {
				if !visited[k] {
					visited[k] = true
					fmt.Printf("%v ", pG.vexs[k])
					queue = append(queue, k)
					rear++
				}
			}
		}
	}

	fmt.Println("\n")
}

/*
 * 返回顶点v的第一个邻接顶点的索引，失败则返回-1
 */
func firstVertex(pG *Graph, v int) int {
	if v < 0 || v > (pG.vexnum-1) {
		return -1
	}

	for i := 0; i < pG.vexnum; i++ {
		if pG.matrix[v][i] == 1 {
			return i
		}
	}
	return -1
}

func nextVertix(pG *Graph, v, w int) int {
	if v < 0 || v > (pG.vexnum-1) || w < 0 || w > (pG.vexnum-1) {
		return -1
	}

	for i := w + 1; i < pG.vexnum; i++ {
		if pG.matrix[v][i] == 1 {
			return i
		}
	}
	return -1
}

func topologicalSort(pG *Graph) {
	var (
		i, j  int
		index int = 0
		head  int = 0 // 辅助队列的头
		rear  int = 0 // 辅助队列的尾
		num       = pG.vexnum
	)

	var (
		queue []int    // 辅助队列
		ins   []int    // 入度队列
		tops  []string // 拓扑排序结果数组，记录每个节点的排序后的序号。
	)

	for i := 0; i < num; i++ {
		ins = append(ins, 0)
		tops = append(tops, "")
	}

	// 统计每个顶点的入度数
	for i := 0; i < num; i++ {
		for j := 0; j < num; j++ {
			if pG.matrix[i][j] == 1 {
				ins[j]++
			}
		}
	}

	// 将所有入度为0的顶点入队列
	for i := 0; i < num; i++ {
		if ins[i] == 0 {
			queue = append(queue, i)
			rear++
		}
	}
	fmt.Println("queue is ", queue)

	for {
		if head == rear {
			break
		}
		j = queue[head]
		head++
		tops[index] = pG.vexs[j]
		index++

		// 将与之有关的入度-1
		for w := 0; w < num; w++ {
			if pG.matrix[j][w] == 1 {
				ins[w] -= 1
				if ins[w] == 0 {
					queue = append(queue, w)
					rear++
				}
			}
		}
		fmt.Println("queue is ", queue)
	}

	if index != pG.vexnum {
		fmt.Println("Graph has a cycle\n")
		return
	}

	// 打印拓扑排序结果
	fmt.Println("== TopSort: ")
	for i = 0; i < num; i++ {
		fmt.Printf("%v ", tops[i])
	}
	fmt.Println("\n")

}

func main() {
	pG := getGraph()
	printGraph(pG)

	//DFSTraverse(pG)
	//BFS(pG)
	topologicalSort(pG)
}

/*
Martix Graph:

0 0 0 0 0 0 1

1 0 0 1 0 0 0

0 0 0 0 0 1 1

0 0 0 0 1 1 0

0 0 0 0 0 0 0

0 0 0 0 0 0 0

0 0 0 0 0 0 0

queue is  [1 2]
queue is  [1 2 0 3]
queue is  [1 2 0 3]
queue is  [1 2 0 3 6]
queue is  [1 2 0 3 6 4 5]
queue is  [1 2 0 3 6 4 5]
queue is  [1 2 0 3 6 4 5]
queue is  [1 2 0 3 6 4 5]
== TopSort:
B C A D G E F
*/
