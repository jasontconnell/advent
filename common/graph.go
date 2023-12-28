package common

import "fmt"

type Edge[V comparable, W Number] interface {
	GetLeft() V
	GetRight() V
	GetWeight() W
}

type Graph[V comparable, W Number] interface {
	AddVertex(v V)
	AddVertices(v ...V)
	AddEdge(v1, v2 V)
	AddWeightedEdge(v1, v2 V, w W)
	GetEdges() []Edge[V, W]
	GetEdgesFrom(v V) []Edge[V, W]

	Print()

	Adjacent(v1, v2 V) bool
	Neighbors(v V) []V
	RemoveVertex(v V)
	RemoveEdge(v1, v2 V)

	GetEdge(v1, v2 V) Edge[V, W]
	GetPaths(v1, v2 V) [][]Edge[V, W]
}

func NewGraph[V comparable]() Graph[V, int] {
	g := new(graph[V, int])
	g.vertices = make(map[V]vertex[V])
	g.originmap = make(map[V][]edge[V, int])
	g.directed = true
	return g
}

func NewWeightedGraph[V comparable, W Number]() Graph[V, W] {
	g := new(graph[V, W])
	g.vertices = make(map[V]vertex[V])
	g.originmap = make(map[V][]edge[V, W])
	return g
}

func NewDirectedGraph[V comparable, W Number]() Graph[V, W] {
	g := new(graph[V, W])
	g.vertices = make(map[V]vertex[V])
	g.originmap = make(map[V][]edge[V, W])
	g.directed = true
	return g
}

type graph[V comparable, W Number] struct {
	vertices  map[V]vertex[V]
	originmap map[V][]edge[V, W]
	directed  bool
}

type vertex[V any] struct {
	data V
}

type edge[V any, W Number] struct {
	weight      W
	left, right vertex[V]
}

func (g *graph[V, W]) AddVertex(v V) {
	g.vertices[v] = vertex[V]{data: v}
}

func (g *graph[V, W]) AddVertices(vs ...V) {
	for _, v := range vs {
		g.AddVertex(v)
	}
}

func (g *graph[V, W]) AddEdge(v1, v2 V) {
	var def W
	g.AddWeightedEdge(v1, v2, def)
}

func (g *graph[V, W]) AddWeightedEdge(v1, v2 V, w W) {
	var v1d, v2d vertex[V]
	if tv, ok := g.vertices[v1]; ok {
		v1d = tv
	} else {
		v1d = vertex[V]{v1}
		g.vertices[v1] = v1d
	}
	if tv, ok := g.vertices[v2]; ok {
		v2d = tv
	} else {
		v2d = vertex[V]{v2}
		g.vertices[v2] = v2d
	}
	newedge := edge[V, W]{weight: w, left: v1d, right: v2d}
	// g.edges = append(g.edges, newedge)
	g.originmap[v1] = append(g.originmap[v1], newedge)
}

func (g graph[V, W]) GetEdges() []Edge[V, W] {
	list := []Edge[V, W]{}
	for _, v := range g.originmap {
		for _, e := range v {
			list = append(list, edge[V, W]{left: e.left, right: e.right, weight: e.weight})
		}
	}
	return list
}

func (g graph[V, W]) GetEdgesFrom(v V) []Edge[V, W] {
	origins := g.originsFrom(v)
	if origins == nil {
		return nil
	}

	list := []Edge[V, W]{}
	for _, e := range origins {
		list = append(list, edge[V, W]{left: e.left, right: e.right, weight: e.weight})
	}
	return list
}

func (g graph[V, W]) Adjacent(v1, v2 V) bool {
	adj := false
	var origins []edge[V, W]
	if list, ok := g.originmap[v1]; ok {
		origins = list
	}
	if list, ok := g.originmap[v2]; ok {
		origins = append(origins, list...)
	}
	for _, e := range origins {
		if e.left.data == v1 && e.right.data == v2 {
			adj = true
			break
		} else if e.left.data == v2 && e.right.data == v1 {
			adj = true
			break
		}
	}
	return adj
}

func (g graph[V, W]) Neighbors(v V) []V {
	if _, ok := g.vertices[v]; !ok {
		return nil
	}
	origins := g.originsFrom(v)
	if origins == nil {
		return nil
	}

	ns := []V{}
	for _, e := range origins {
		ns = append(ns, e.right.data)
	}
	return ns
}

func (g *graph[V, W]) RemoveVertex(v V) {
	if _, ok := g.vertices[v]; !ok {
		return
	}

	origins := g.originsFrom(v)

	for _, ev := range origins {
		rlist, ok := g.originmap[ev.right.data]
		if !ok {
			continue
		}
		for j := len(rlist) - 1; j >= 0; j-- {
			r := rlist[j]
			if r.right.data == v {
				rlist = append(rlist[:j], rlist[j+1:]...)
				break
			}
		}
	}

	delete(g.vertices, v)
	delete(g.originmap, v)
}

func (g *graph[V, W]) RemoveEdge(v1, v2 V) {
	origins := g.originsFrom(v1)
	if origins == nil {
		return
	}

	for i := len(origins) - 1; i >= 0; i-- {
		e := origins[i]
		if e.right.data == v2 {
			origins = append(origins[:i], origins[i+1:]...)
		}
	}
	g.originmap[v1] = origins
}

func (g graph[V, W]) Print() {
	fmt.Println(g.vertices)
	fmt.Println(g.originmap)
}

func (g graph[V, W]) GetEdge(v1, v2 V) Edge[V, W] {
	origins := g.originsFrom(v1)
	if origins == nil {
		return nil
	}

	var found Edge[V, W]
	for i := 0; i < len(origins); i++ {
		e := origins[i]
		if e.right.data == v2 {
			found = e
			break
		}
	}
	return found
}

func (e edge[V, W]) GetLeft() V {
	return e.left.data
}

func (e edge[V, W]) GetRight() V {
	return e.right.data
}

func (e edge[V, W]) GetWeight() W {
	return e.weight
}

func (g graph[V, W]) originsFrom(v V) []edge[V, W] {
	var origins []edge[V, W]
	if list, ok := g.originmap[v]; ok {
		origins = list
	}
	return origins
}
