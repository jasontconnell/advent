package common

import "fmt"

type state[V Ordered, W Number] struct {
	vertex V
	path   []Edge[V, W]
}

type edgestate[V Ordered, W Number] struct {
	edge  Edge[V, W]
	total W
}

type statekey[V Ordered] struct {
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

func (g *graph[V, W]) AStar(v1, v2 V, h func(e Edge[V, W]) W) []Edge[V, W] {
	// queue := NewPriorityQueue[edgestate[V, W], W](func(s edgestate[V, W]) W {
	// 	return h(s.edge)
	// })

	// gscore := make(map[V]W)
	// cameFrom := make(map[V]V)
	// for _, edge := range g.originsFrom(v1) {
	// 	queue.Enqueue(edgestate[V, W]{edge: edge, total: edge.GetWeight()})
	// }

	return nil
}

func getMoves[V Ordered, W Number](g graph[V, W], from V) []Edge[V, W] {
	mvs := []Edge[V, W]{}
	for _, e := range g.GetEdgesFrom(from) {
		mvs = append(mvs, e)
	}
	return mvs
}
