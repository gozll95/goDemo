package graph

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Createt_Martix(t *testing.T) {
	martix := NewMartix()

	vexs := []interface{}{"A", "B", "C", "D", "E", "F", "G"}
	err := martix.AddVexs(vexs...)
	assert.Nil(t, err)

	err = martix.AddEdge("A", "G")
	assert.Nil(t, err)
	err = martix.AddEdge("B", "A")
	assert.Nil(t, err)
	err = martix.AddEdge("B", "D")
	assert.Nil(t, err)

	assert.Equal(t, len(martix.Edges()), 3)
	assert.Equal(t, len(martix.Vexs()), 7)
}

func Test_Martix_OverSize(t *testing.T) {
	martix := NewMartix()

	vexs := []interface{}{}
	for i := 0; i < 1000; i++ {
		vexs = append(vexs, i)
	}
	err := martix.AddVexs(vexs...)
	assert.Equal(t, err, MartixOverSize)
}

func Test_Martix_Create_Graph(t *testing.T) {
	martix := NewMartix()

	vexs := []interface{}{"A", "B", "C", "D", "E", "F", "G"}
	err := martix.AddVexs(vexs...)
	assert.Nil(t, err)

	martix.AddEdge("A", "G")
	martix.AddEdge("B", "A")
	martix.AddEdge("B", "D")
	martix.AddEdge("C", "F")
	martix.AddEdge("C", "G")
	martix.AddEdge("D", "E")
	martix.AddEdge("D", "F")

	pG, err := martix.CreateGraph()
	assert.Nil(t, err)

	pG.Print()

	a := pG.AllInsVex("A")
	assert.Equal(t, a[0].(string), "B")

	b := pG.AllInsVex("B")
	assert.Equal(t, len(b), 0)

	c := pG.AllInsVex("C")
	assert.Equal(t, len(c), 0)

	d := pG.AllInsVex("D")
	assert.Equal(t, d[0].(string), "B")

	e := pG.AllInsVex("E")
	assert.Equal(t, e[0].(string), "D")

	f := pG.AllInsVex("F")
	assert.Equal(t, true, f[0].(string) == "D" || f[0].(string) == "C")
	assert.Equal(t, true, f[1].(string) == "D" || f[1].(string) == "C")

	g := pG.AllInsVex("G")
	assert.Equal(t, true, g[0].(string) == "A" || g[0].(string) == "C")
	assert.Equal(t, true, g[1].(string) == "A" || g[1].(string) == "C")
}

func Test_Graph_DFS(t *testing.T) {
	martix := NewMartix()

	vexs := []interface{}{"A", "B", "C", "D", "E", "F", "G"}
	err := martix.AddVexs(vexs...)
	assert.Nil(t, err)

	err = martix.AddEdge("A", "B")
	assert.Nil(t, err)
	err = martix.AddEdge("B", "C")
	assert.Nil(t, err)
	err = martix.AddEdge("B", "E")
	assert.Nil(t, err)
	err = martix.AddEdge("B", "F")
	assert.Nil(t, err)
	err = martix.AddEdge("C", "E")
	assert.Nil(t, err)
	err = martix.AddEdge("D", "C")
	assert.Nil(t, err)
	err = martix.AddEdge("E", "B")
	assert.Nil(t, err)
	err = martix.AddEdge("E", "D")
	assert.Nil(t, err)
	err = martix.AddEdge("F", "G")
	assert.Nil(t, err)

	pG, err := martix.CreateGraph()
	assert.Nil(t, err)

	want := []string{"A", "B", "C", "E", "D", "F", "G"}

	item := pG.DFSTraverse()
	for i, v := range item {
		assert.Equal(t, v.(string), want[i])
	}
}

func Test_Graph_BFS(t *testing.T) {
	martix := NewMartix()

	vexs := []interface{}{"A", "B", "C", "D", "E", "F", "G"}
	err := martix.AddVexs(vexs...)
	assert.Nil(t, err)

	err = martix.AddEdge("A", "B")
	assert.Nil(t, err)
	err = martix.AddEdge("B", "C")
	assert.Nil(t, err)
	err = martix.AddEdge("B", "E")
	assert.Nil(t, err)
	err = martix.AddEdge("B", "F")
	assert.Nil(t, err)
	err = martix.AddEdge("C", "E")
	assert.Nil(t, err)
	err = martix.AddEdge("D", "C")
	assert.Nil(t, err)
	err = martix.AddEdge("E", "B")
	assert.Nil(t, err)
	err = martix.AddEdge("E", "D")
	assert.Nil(t, err)
	err = martix.AddEdge("F", "G")
	assert.Nil(t, err)

	pG, err := martix.CreateGraph()
	assert.Nil(t, err)

	want := []string{"A", "B", "C", "E", "F", "D", "G"}

	item := pG.BFS()
	for i, v := range item {
		assert.Equal(t, v.(string), want[i])
	}
}

func Test_Graph_TopoSort(t *testing.T) {
	martix := NewMartix()

	vexs := []interface{}{"A", "B", "C", "D", "E", "F", "G"}
	err := martix.AddVexs(vexs...)
	assert.Nil(t, err)

	err = martix.AddEdge("A", "G")
	assert.Nil(t, err)
	err = martix.AddEdge("B", "A")
	assert.Nil(t, err)
	err = martix.AddEdge("B", "D")
	assert.Nil(t, err)
	err = martix.AddEdge("C", "F")
	assert.Nil(t, err)
	err = martix.AddEdge("C", "G")
	assert.Nil(t, err)
	err = martix.AddEdge("D", "E")
	assert.Nil(t, err)
	err = martix.AddEdge("D", "F")
	assert.Nil(t, err)

	pG, err := martix.CreateGraph()
	assert.Nil(t, err)

	want := []string{"B", "C", "A", "D", "G", "E", "F"}

	item, err := pG.TopologicalSort()
	assert.Nil(t, err)
	for i, v := range item {
		assert.Equal(t, v.(string), want[i])
	}
}

func Test_Graph_TopologicalSortWithBFS(t *testing.T) {
	martix := NewMartix()

	vexs := []interface{}{"A", "B", "C", "D", "E", "F", "G"}
	err := martix.AddVexs(vexs...)
	assert.Nil(t, err)

	err = martix.AddEdge("A", "G")
	assert.Nil(t, err)
	err = martix.AddEdge("B", "A")
	assert.Nil(t, err)
	err = martix.AddEdge("B", "D")
	assert.Nil(t, err)
	err = martix.AddEdge("C", "F")
	assert.Nil(t, err)
	err = martix.AddEdge("C", "G")
	assert.Nil(t, err)
	err = martix.AddEdge("D", "E")
	assert.Nil(t, err)
	err = martix.AddEdge("D", "F")
	assert.Nil(t, err)

	pG, err := martix.CreateGraph()
	assert.Nil(t, err)

	items, err := pG.TopologicalSortWithBFS()
	assert.Nil(t, err)
	for i, item := range items {
		if i == 0 {
			for _, v := range item {
				assert.Equal(t, true, v == "B" || v == "C")
			}
		}
		if i == 1 {
			for _, v := range item {
				assert.Equal(t, true, v == "A" || v == "D")
			}
		}
		if i == 2 {
			for _, v := range item {
				assert.Equal(t, true, v == "G" || v == "E" || v == "F")
			}
		}
	}

}

func Test_Graph_Cycle(t *testing.T) {
	martix := NewMartix()

	vexs := []interface{}{"A", "B", "C", "D", "E", "F", "G"}
	err := martix.AddVexs(vexs...)
	assert.Nil(t, err)

	err = martix.AddEdge("A", "B")
	assert.Nil(t, err)
	err = martix.AddEdge("B", "C")
	assert.Nil(t, err)
	err = martix.AddEdge("B", "E")
	assert.Nil(t, err)
	err = martix.AddEdge("B", "F")
	assert.Nil(t, err)
	err = martix.AddEdge("C", "E")
	assert.Nil(t, err)
	err = martix.AddEdge("D", "C")
	assert.Nil(t, err)
	err = martix.AddEdge("E", "B")
	assert.Nil(t, err)
	err = martix.AddEdge("E", "D")
	assert.Nil(t, err)
	err = martix.AddEdge("F", "G")
	assert.Nil(t, err)

	pG, err := martix.CreateGraph()
	assert.Nil(t, err)

	_, err = pG.TopologicalSort()
	assert.Equal(t, err, GraphCycleErr)

	_, err = pG.TopologicalSortWithBFS()
	assert.Equal(t, err, GraphCycleErr)
}
