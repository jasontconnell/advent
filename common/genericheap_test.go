package common

import (
	"fmt"
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

func TestBinarySearch(t *testing.T) {
	vals := []int{1, 3, 5, 6, 9, 12, 15, 16, 18, 20, 24, 29, 31, 33, 43}

	ins := 11
	sort.Ints(vals)

	idx := -1
	if len(vals) > 0 && vals[0] > ins {
		idx = 0
	} else if len(vals) > 0 && vals[len(vals)-1] > ins {
		lbound, ubound := 0, len(vals)
		found := false
		for !found {
			idx = (lbound + ubound) / 2
			midval := vals[idx]
			if ins > midval {
				lbound = idx
			} else {
				ubound = idx
			}

			fmt.Println(lbound, ubound)

			if ubound-lbound <= 1 {
				if ins > vals[lbound] {
					idx = ubound
				}
				fmt.Println(lbound, vals[lbound], ubound, vals[ubound])
				found = true
			}
		}
	}

	fmt.Println(idx)
	fmt.Println(ins, vals)

	if idx == -1 {
		vals = append(vals, ins)
	} else {
		vals = append(vals[:idx], append([]int{ins}, vals[idx:]...)...)
	}
	fmt.Println("inserted", vals)
}

func TestHeap(t *testing.T) {
	q := NewPriorityQueue(func(i int) int {
		return i
	})

	var count int = 400
	for i := 0; i < count; i++ {
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

	if q.Len() != count || len(ints) != count {
		t.Fail()
	}

	if !sort.IntsAreSorted(ints) {
		t.Log(q.items)
		t.Fail()
	}
}
