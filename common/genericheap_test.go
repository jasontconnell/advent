package common

import (
	"math/rand"
	"sort"
	"testing"
)

func TestHeap(t *testing.T) {
	q := NewPriorityQueue(func(i int) int {
		return i
	})

	for i := 0; i < 20; i++ {
		q.Enqueue(rand.Int() % 100)
	}

	// priority queue orders by high priority
	// so it will be sorted descending
	cp := make([]Item[int, int], len(q.items))
	for i := 0; i < len(q.items); i++ {
		cp[i] = q.items[len(q.items)-i-1]
	}

	ints := []int{}
	for _, x := range cp {
		ints = append(ints, x.item)
	}

	if !sort.IntsAreSorted(ints) {
		t.Log(q.items)
		t.Fail()
	}
}
