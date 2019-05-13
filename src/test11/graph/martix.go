package graph

import (
	mapset "github.com/deckarep/golang-set"
)

/*
 * 最多有100条边(边即依赖关系)
 */
const MAXSIZE int = 100

/*
 * 矩阵
 * 非并发安全
 */
type Martix struct {
	edges   [][]interface{}
	vexs    []interface{}
	edgeSet mapset.Set
	vexSet  mapset.Set
}

type Edge struct {
	v1 interface{}
	v2 interface{}
}

func NewMartix() *Martix {
	return &Martix{
		vexSet:  mapset.NewSet(),
		edgeSet: mapset.NewSet(),
	}
}

/*
 * 添加元素
 */
func (m *Martix) AddVex(vex interface{}) (err error) {
	if vex == nil {
		return MartixVexMustNotNil
	}
	if len(m.vexSet.ToSlice()) == MAXSIZE {
		return MartixOverSize
	}
	m.vexSet.Add(vex)
	return
}

/*
 * 添加依赖关系,vex2依赖vex1
 */
func (m *Martix) AddEdge(vex1, vex2 interface{}) (err error) {
	if vex1 == nil || vex2 == nil {
		return MartixVexMustNotNil
	}
	if len(m.edgeSet.ToSlice()) == MAXSIZE {
		return MartixOverSize
	}

	m.edgeSet.Add(Edge{vex1, vex2})
	m.vexSet.Add(vex1)
	m.vexSet.Add(vex2)

	return nil
}

func (m *Martix) Vexs() []interface{} {
	return m.vexSet.ToSlice()
}

func (m *Martix) Edges() [][]interface{} {
	for _, v := range m.edgeSet.ToSlice() {
		e := make([]interface{}, 2)
		e[0] = v.(Edge).v1
		e[1] = v.(Edge).v2
		m.edges = append(m.edges, e)
	}
	return m.edges
}

/*
 * 根据martix创建邻接表
 */
func (m *Martix) CreateGraph() (*Graph, error) {
	vexs := m.Vexs()
	edges := m.Edges()

	return CreateGraph(vexs, edges)
}
