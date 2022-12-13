package common

type Number interface {
	int | float64
}

type Queue[T any, N Number] struct {
	items []T
	value func(item T) N
}

func NewPriorityQueue[T any, N Number](val func(item T) N) *Queue[T, N] {
	h := new(Queue[T, N])
	h.value = val
	return h
}

func NewQueue[T any, N Number]() *Queue[T, N] {
	h := new(Queue[T, N])
	return h
}

func (h *Queue[T, N]) Enqueue(item T) {
	if h.value == nil {
		h.items = append(h.items, item)
	} else {
		v := h.value(item)
		idx := -1
		for i := 0; i < len(h.items); i++ {
			r := h.value(h.items[i])

			if v > r {
				idx = i - 1
				break
			}
		}

		if idx == -1 {
			h.items = append(h.items, item)
		} else {
			h.items = append(h.items[:idx], append([]T{item}, h.items[idx:]...)...)
		}
	}
}

func (h *Queue[T, N]) Dequeue() T {
	var item T
	if len(h.items) > 0 {
		item = h.items[0]
		h.items = h.items[1:]
	}
	return item
}

func (h *Queue[T, N]) Any() bool {
	return h.Len() > 0
}

func (h *Queue[T, N]) Len() int {
	return len(h.items)
}