package common

type state[V comparable, W Number] struct {
	edge Edge[V, W]
	path []Edge[V, W]
}

type statekey[V comparable] struct {
	left, right V
}

func (g graph[V, W]) GetPaths(v1, v2 V) [][]Edge[V, W] {
	paths := [][]Edge[V, W]{}
	queue := []state[V, W]{}
	vis := make(map[statekey[V]]bool)
	mvs := g.originsFrom(v1)
	if mvs == nil {
		return nil
	}

	for _, mv := range mvs {
		queue = append(queue, state[V, W]{edge: mv, path: []Edge[V, W]{mv}})
	}

	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]

		if cur.edge.GetRight() == v2 {
			paths = append(paths, cur.path)
			continue
		}

		mvs := getMoves[V, W](g, cur.edge.GetRight())
		for _, mv := range mvs {
			k := statekey[V]{left: mv.GetLeft(), right: mv.GetRight()}
			if _, ok := vis[k]; ok {
				continue
			}
			vis[k] = true
			queue = append(queue, state[V, W]{path: append(cur.path, mv), edge: mv})
		}
	}
	return paths
}

func getMoves[V comparable, W Number](g graph[V, W], from V) []Edge[V, W] {
	mvs := []Edge[V, W]{}
	for _, e := range g.GetEdgesFrom(from) {
		mvs = append(mvs, e)
	}
	return mvs
}
