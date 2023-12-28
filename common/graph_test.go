package common

import (
	"testing"
)

func getSampleGraph() Graph[int, int] {
	g := NewGraph[int]()
	g.AddVertices(4, 2, 5, 7, 1, 12, 23, 17)
	g.AddEdge(5, 7)
	g.AddEdge(1, 12)
	g.AddEdge(1, 4)
	g.AddEdge(1, 23)
	g.AddEdge(4, 5)
	g.AddEdge(7, 23)
	g.AddEdge(4, 12)
	return g
}

func TestGraphAdjacent(t *testing.T) {
	g := getSampleGraph()
	tests := []struct {
		left, right int
		expected    bool
	}{
		{5, 8, false},
		{4, 7, false},
		{5, 7, true},
		{4, 5, true},
	}

	for _, test := range tests {
		adj := g.Adjacent(test.left, test.right)
		if adj != test.expected {
			t.Log("unexpected value ", adj, "on test", test.left, test.right)
			t.Fail()
		}
	}
}

func TestGraphNeighbors(t *testing.T) {
	g := getSampleGraph()
	tests := []struct {
		v        int
		expected []int
	}{
		{5, []int{7}},
		{1, []int{4, 12, 23}},
	}

	for _, test := range tests {
		n := g.Neighbors(test.v)
		t.Log(n)
	}
}

func TestRemoveVertex(t *testing.T) {
	g := getSampleGraph()
	g.RemoveVertex(5)
	g.Print()
	adj := g.Adjacent(5, 7) || g.Adjacent(5, 4)
	if adj {
		// if still adjacent, remove failed
		t.Fail()
	}
}

func TestRemoveEdge(t *testing.T) {
	g := getSampleGraph()
	g.AddEdge(7, 5) // add a reverse of what we're removing

	g.RemoveEdge(5, 7) // just remove the 5-7 edge, not the 7-5 one

	adj := g.Adjacent(5, 7)
	if adj {
		t.Fail()
	}
}

func TestGetPaths(t *testing.T) {
	g := getSampleGraph()
	paths := g.GetPaths(1, 23)
	t.Log(paths)
	if len(paths) < 2 {
		t.Fail()
	}

	nopath := g.GetPaths(1, 17)
	t.Log(nopath)
	if len(nopath) > 0 {
		t.Fail()
	}
}

type xy struct {
	x, y int
}

func TestGetPath(t *testing.T) {
	g := NewGraph[xy]()
	g.AddVertices(xy{0, 0}, xy{0, 1}, xy{0, 2}, xy{1, 2}, xy{3, 4}, xy{3, 5}, xy{5, 5})
	g.AddEdge(xy{0, 0}, xy{0, 1})
	g.AddEdge(xy{0, 1}, xy{1, 2})
	g.AddEdge(xy{1, 2}, xy{3, 4})
	g.AddEdge(xy{3, 4}, xy{5, 5})
	p := g.GetPaths(xy{0, 0}, xy{5, 5})
	t.Log(p)
}
