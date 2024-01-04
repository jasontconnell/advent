package common

import "fmt"

type state[V comparable, W Number] struct {
	vertex V
	path   []Edge[V, W]
}

type edgestate[V comparable, W Number] struct {
	edge  Edge[V, W]
	path  []Edge[V, W]
	total W
}

type statekey[V comparable] struct {
	left, right V
}

func (g graph[V, W]) DFS(v1, v2 V) []Edge[V, W] {
	stk := NewStack[state[V, W]]()

	vis := make(map[V]bool)

	var path []Edge[V, W]
	stk.Push(state[V, W]{vertex: v1, path: []Edge[V, W]{}})
	vis[v1] = true

	for stk.Any() {
		cur := stk.Pop()

		if cur.vertex == v2 {
			fmt.Println("found v2", cur)
			path = cur.path
			break
		}

		for _, edge := range g.GetEdgesFrom(cur.vertex) {
			if _, ok := vis[edge.GetRight()]; ok {
				continue
			}
			vis[edge.GetRight()] = true
			stk.Push(state[V, W]{vertex: edge.GetRight(), path: append(cur.path, edge)})
		}
	}

	return path
}

func (g graph[V, W]) BFS(v1, v2 V) []Edge[V, W] {
	queue := NewQueue[state[V, W], W]()
	vis := make(map[V]bool)

	queue.Enqueue(state[V, W]{vertex: v1, path: []Edge[V, W]{}})
	var path []Edge[V, W]

	for queue.Any() {
		cur := queue.Dequeue()

		if cur.vertex == v2 {
			path = cur.path
			break
		}

		for _, e := range g.GetEdgesFrom(cur.vertex) {
			if _, ok := vis[e.GetRight()]; ok {
				continue
			}
			vis[e.GetRight()] = true
			queue.Enqueue(state[V, W]{vertex: e.GetRight(), path: append(cur.path, e)})
		}
	}
	return path
}

// func (g graph[V, W]) GetPaths(v1, v2 V) [][]Edge[V, W] {
// 	paths := [][]Edge[V, W]{}
// 	queue := NewPriorityQueue[state[V, W], W](func(s state[V, W]) W {
// 		return s.edge.GetWeight()
// 	})

// 	mvs := g.originsFrom(v1)
// 	if mvs == nil {
// 		return nil
// 	}

// 	for _, mv := range mvs {
// 		queue.Enqueue(state[V, W]{edge: mv, path: []Edge[V, W]{mv}})
// 	}

// 	v := make(map[statekey[V]]bool)

// 	for queue.Any() {
// 		cur := queue.Dequeue()

// 		if cur.edge.GetRight() == v2 {
// 			paths = append(paths, cur.path)
// 			continue
// 		}

// 		mvs := getMoves[V, W](g, cur.edge.GetRight())
// 		for _, mv := range mvs {
// 			k := statekey[V]{left: mv.GetLeft(), right: mv.GetRight()}
// 			fmt.Println("k", k)
// 			if _, ok := v[k]; ok {
// 				continue
// 			}
// 			v[k] = true
// 			queue.Enqueue(state[V, W]{edge: mv, path: append(cur.path, mv)})
// 		}
// 	}
// 	return paths
// }

func getMoves[V comparable, W Number](g graph[V, W], from V) []Edge[V, W] {
	mvs := []Edge[V, W]{}
	for _, e := range g.GetEdgesFrom(from) {
		mvs = append(mvs, e)
	}
	return mvs
}
