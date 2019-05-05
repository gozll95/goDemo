package main

import (
	"fmt"
	"test1/graph"
)

func main() {
	// martix := graph.NewMartix()

	// martix.AddEdge("A", "B")
	// martix.AddEdge("B", "C")
	// martix.AddEdge("B", "E")
	// martix.AddEdge("B", "F")
	// martix.AddEdge("C", "E")
	// martix.AddEdge("D", "C")
	// martix.AddEdge("E", "B")
	// martix.AddEdge("E", "D")
	// martix.AddEdge("F", "G")

	// var vexs []interface{}
	// vexs = append(vexs, "A", "B", "C", "D", "E", "F", "G")

	vexs := []interface{}{"A", "B", "C", "D", "E", "F", "G"}
	var edges [][]interface{}
	tmp1 := []interface{}{"A", "G"}
	edges = append(edges, tmp1)
	tmp1 = []interface{}{"B", "A"}
	edges = append(edges, tmp1)
	tmp1 = []interface{}{"B", "D"}
	edges = append(edges, tmp1)
	tmp1 = []interface{}{"B", "F"}
	edges = append(edges, tmp1)
	tmp1 = []interface{}{"C", "E"}
	edges = append(edges, tmp1)
	tmp1 = []interface{}{"D", "C"}
	edges = append(edges, tmp1)
	tmp1 = []interface{}{"E", "B"}
	edges = append(edges, tmp1)
	tmp1 = []interface{}{"E", "D"}
	edges = append(edges, tmp1)
	tmp1 = []interface{}{"F", "G"}
	edges = append(edges, tmp1)

	pG, err := graph.CreateGraph(vexs, edges)
	if err != nil {
		panic(err)
	}

	pG.Print()

	//pG.Test()

	fmt.Println(pG.DFSTraverse())
	fmt.Println(pG.BFS())

	items, err := pG.TopologicalSort()
	if err != nil {
		panic(err)
	}
	fmt.Println(items)

	levels, err := pG.TopologicalSortWithBFS()
	if err != nil {
		panic(err)
	}
	fmt.Println(levels)
}

/*
== List Graph:
0(A): 6(G)
1(B): 0(A) 3(D)
2(C): 5(F) 6(G)
3(D): 4(E) 5(F)
4(E):
5(F):
6(G):
== DFS: A G B D E F C
== BFS: A G B D E F C
== TopSort: B C A D G E F
*/
