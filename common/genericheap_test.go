package common

import (
	"math/rand"
	"sort"
	"testing"
)

func TestDequeue(t *testing.T) {
	q := NewPriorityQueue(func(i int) int {
		return i
	})

	var count int = 12
	for i := 0; i < count; i++ {
		q.Enqueue(rand.Int() % 100)
	}

	// visually observe that items are in order biggest to smallest
	for q.Any() {
		item := q.Dequeue()
		t.Log(item)
	}
}

func TestHeap(t *testing.T) {
	q := NewPriorityQueue(func(i int) int {
		return i
	})

	var count int = 600
	for i := 0; i < count; i++ {
		q.Enqueue(rand.Int() % 100)
	}

	// priority queue orders by priority
	// so it will be sorted ascending
	cp := make([]Item[int, int], len(q.items))
	for i := 0; i < len(q.items); i++ {
		cp[i] = q.items[i]
	}

	ints := []int{}
	for _, x := range cp {
		ints = append(ints, x.item)
	}

	if q.Len() != count || len(ints) != count {
		t.Log("lens off", q.Len(), count, len(ints))
		t.Fail()
	}

	if !sort.IntsAreSorted(ints) {
		t.Log("not sorted", q.items)
		t.Fail()
	}
}
