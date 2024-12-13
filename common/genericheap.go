package common

import (
	"fmt"
)

func (s Item[T, N]) String() string {
	return fmt.Sprintf("%v", s.value)
}

type Queue[T any, N Number] struct {
	items []Item[T, N]
	cost  func(item T) N
}

func NewPriorityQueue[T any, N Number](cost func(item T) N) *Queue[T, N] {
	h := new(Queue[T, N])
	h.cost = cost
	return h
}

func NewQueue[T any, N Number]() *Queue[T, N] {
	h := new(Queue[T, N])
	return h
}

func (h *Queue[T, N]) Enqueue(item T) {
	if h.cost == nil {
		h.items = append(h.items, Item[T, N]{item: item})
	} else {
		v := h.cost(item)
		idx := -1
		if len(h.items) > 0 && h.items[0].value >= v {
			idx = 0
		} else if len(h.items) > 0 && v < h.items[len(h.items)-1].value {
			lbound, ubound := 0, len(h.items)
			found := false
			for !found {
				idx = lbound + (ubound-lbound)/2
				midval := h.items[idx].value
				if v == midval {
					break
				} else if v > midval {
					lbound = idx
				} else {
					ubound = idx
				}

				if ubound-lbound <= 1 {
					found = true
					if v > h.items[lbound].value {
						idx = lbound + 1
					} else {
						idx = ubound - 1
					}
				}
			}
		}

		i := Item[T, N]{item: item, value: v}
		if idx == -1 || idx >= len(h.items) {
			h.items = append(h.items, i)
		} else {
			h.items = append(h.items[:idx], append([]Item[T, N]{i}, h.items[idx:]...)...)
		}
	}
}

func (h *Queue[T, N]) Dequeue() T {
	var item Item[T, N]
	if len(h.items) > 0 {
		item = h.items[0]
		h.items = h.items[1:]
	}
	return item.item
}

func (h *Queue[T, N]) Any() bool {
	return h.Len() > 0
}

func (h *Queue[T, N]) Len() int {
	return len(h.items)
}
