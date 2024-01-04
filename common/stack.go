package common

type Stack[V any] interface {
	Push(v V)
	Pop() V
	Any() bool
}

type stack[V any] struct {
	elements []V
}

func NewStack[V any]() Stack[V] {
	s := new(stack[V])
	s.elements = []V{}
	return s
}

func (s *stack[V]) Push(v V) {
	s.elements = append(s.elements, v)
}

func (s *stack[V]) Pop() V {
	v := s.elements[len(s.elements)-1]
	s.elements = s.elements[:len(s.elements)-1]
	return v
}

func (s *stack[V]) Any() bool {
	return len(s.elements) > 0
}
