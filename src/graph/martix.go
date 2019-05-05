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
	//edges   [][]interface{}
	edgeSet mapset.Set
	vexs    []interface{}
	edges   []Edge
}

/*
 *
 * 边(v2依赖v1)
 */
type Edge struct {
	v1 interface{}
	v2 interface{}
}

func NewMartix() *Martix {
	return &Martix{
		edgeSet: mapset.NewSet(),
	}
}

/*
 * 添加点
 * 执行AddVex的顺序即为点搜索的伊始顺序
 */
func (m *Martix) AddVex(vex interface{}) (err error) {
	if vex == nil {
		return MartixVexMustNotNil
	}
	if len(m.vexs) == MAXSIZE {
		return MartixOverSize
	}
	if m.IsExistVex(vex) {
		return VexAlreadyExistInMartix
	}

	m.vexs = append(m.vexs, vex)
	return
}

/*
 * 添加一组点
 * 执行AddVex的顺序即为vexs的顺序
 */
func (m *Martix) AddVexs(vexs ...interface{}) (err error) {
	for _, v := range vexs {
		err = m.AddVex(v)
		if err != nil {
			return
		}
	}
	return
}

/*
 * 点是否已经存在在矩阵中
 */
func (m *Martix) IsExistVex(vex interface{}) bool {
	for _, v := range m.vexs {
		if v == vex {
			return true
		}
	}
	return false
}

/*
 * 边是否已经存在在矩阵中
 */
func (m *Martix) IsExistEdge(edge Edge) bool {
	for _, e := range m.edges {
		if e == edge {
			return true
		}
	}
	return false
}

/*
 * 添加边,vex2依赖vex1
 */
func (m *Martix) AddEdge(vex1, vex2 interface{}) (err error) {
	if vex1 == nil || vex2 == nil {
		return MartixVexMustNotNil
	}
	if len(m.edgeSet.ToSlice()) == MAXSIZE {
		return MartixOverSize
	}
	if !m.IsExistVex(vex1) || !m.IsExistVex(vex2) {
		return VexNotFoundInMartix
	}

	edge := Edge{vex1, vex2}
	if m.IsExistEdge(edge) {
		return EdgeAlreadyExistInMartix
	}

	m.edges = append(m.edges, edge)

	return nil
}

/*
 * 添加一组边
 */
func (m *Martix) AddEdges(edges ...Edge) (err error) {
	for _, edge := range edges {
		err = m.AddEdge(edge.v1, edge.v2)
		if err != nil {
			return
		}
	}
	return
}

/*
 * 返回点
 */
func (m *Martix) Vexs() []interface{} {
	t := make([]interface{}, len(m.vexs))
	copy(t, m.vexs)
	return t
}

/*
 * 返回边
 */
func (m *Martix) Edges() [][]interface{} {
	var edges [][]interface{}

	t := make([]Edge, len(m.edges))
	copy(t, m.edges)

	for _, v := range t {
		vv := make([]interface{}, 2)
		vv[0] = v.v1
		vv[1] = v.v2
		edges = append(edges, vv)
	}
	return edges
}

/*
 * 根据martix创建邻接表
 */
func (m *Martix) CreateGraph() (*Graph, error) {
	vexs := m.Vexs()
	edges := m.Edges()

	return CreateGraph(vexs, edges)
}
