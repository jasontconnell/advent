package common

import (
	"cmp"
	"fmt"
	"sort"
)

type Edge[V Ordered, W Number] interface {
	GetLeft() V
	GetRight() V
	GetWeight() W
}

type Graph[V Ordered, W Number] interface {
	AddVertex(v V)
	AddVertices(v ...V)
	GetVertices() []V
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

	DFS(v1, v2 V) []Edge[V, W]
	BFS(v1, v2 V) []Edge[V, W]
	AStar(v1, v2 V, h func(e Edge[V, W]) W) []Edge[V, W]

	AllCliques() [][]V
}

func NewGraph[V Ordered]() Graph[V, int] {
	g := new(graph[V, int])
	g.vertices = make(map[V]vertex[V])
	g.originmap = make(map[V][]edge[V, int])
	g.neighbors = make(map[V][]V)
	g.directed = false
	return g
}

func NewWeightedGraph[V Ordered, W Number]() Graph[V, W] {
	g := new(graph[V, W])
	g.vertices = make(map[V]vertex[V])
	g.originmap = make(map[V][]edge[V, W])
	g.neighbors = make(map[V][]V)
	return g
}

func NewDirectedGraph[V Ordered, W Number]() Graph[V, W] {
	g := new(graph[V, W])
	g.vertices = make(map[V]vertex[V])
	g.originmap = make(map[V][]edge[V, W])
	g.directed = true
	g.neighbors = make(map[V][]V)
	return g
}

type graph[V Ordered, W Number] struct {
	vertices  map[V]vertex[V]
	originmap map[V][]edge[V, W]

	neighbors map[V][]V
	directed  bool
}

type vertex[V Ordered] struct {
	data V
}

type edge[V Ordered, W Number] struct {
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

func (g *graph[V, W]) GetVertices() []V {
	list := []V{}
	for _, v := range g.vertices {
		list = append(list, v.data)
	}
	return list
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
	g.originmap[v1] = append(g.originmap[v1], newedge)
	if !g.directed {
		indedge := edge[V, W]{weight: w, left: v2d, right: v1d}
		g.originmap[v2] = append(g.originmap[v2], indedge)
	}
	delete(g.neighbors, v1)
	delete(g.neighbors, v2)
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

	for _, e := range origins {
		if e.left.data == v1 && e.right.data == v2 {
			adj = true
			break
		}
	}
	return adj
}

func (g graph[V, W]) Neighbors(v V) []V {
	if n, ok := g.neighbors[v]; ok {
		return n
	}
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

	g.neighbors[v] = ns
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
		g.originmap[ev.right.data] = rlist
	}

	delete(g.vertices, v)
	delete(g.originmap, v)
	delete(g.neighbors, v)
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
	delete(g.neighbors, v1)
	delete(g.neighbors, v2)
}

func (g graph[V, W]) Print() {
	fmt.Println("vertices")
	i := 0
	for _, v := range g.vertices {
		fmt.Print(v)
		i++
		if i%10 == 0 {
			fmt.Println()
		}
	}
	fmt.Println()
	fmt.Println("edges")
	for _, v := range g.originmap {
		for _, e := range v {
			fmt.Println(e.left.data, "->", e.right.data)
		}
	}
	fmt.Println()
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

func (g graph[V, W]) AllCliques() [][]V {
	p := g.GetVertices()
	sort.Slice(p, func(i, j int) bool {
		return p[i] < p[j]
	})
	return g.bronKerbosch([]V{}, p, []V{})
}

func (g graph[V, W]) bronKerbosch(r, p, x []V) [][]V {
	if len(p) == 0 && len(x) == 0 {
		return [][]V{r}
	}

	var pivot V
	if len(p) > 0 {
		pivot = p[0]
	} else {
		pivot = x[0]
	}

	var cliques [][]V
	for _, v := range p {
		vn := g.Neighbors(v)
		if !contains(vn, pivot) {
			newR := append(r, v)
			sort.Slice(vn, func(i, j int) bool {
				return cmp.Less[V](vn[i], vn[j])
			})

			newP := intersection(p, vn)
			newX := intersection(x, vn)

			cliques = append(cliques, g.bronKerbosch(newR, newP, newX)...)
			p = remove(p, v)
			x = append(x, v)
		}
	}

	return cliques
}

func intersection[V comparable](a, b []V) []V {
	var result []V
	for _, x := range a {
		for _, y := range b {
			if x == y {
				result = append(result, x)
				break
			}
		}
	}
	return result
}

func contains[V Ordered](list []V, item V) bool {
	for _, v := range list {
		if v == item {
			return true
		}
	}
	return false
}

func remove[V Ordered](s []V, e V) []V {
	var result []V
	for _, v := range s {
		if v != e {
			result = append(result, v)
		}
	}
	return result
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
