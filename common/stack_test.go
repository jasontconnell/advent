package common

import (
	"testing"
)

func TestStack(t *testing.T) {
	stk := NewStack[int]()
	for i := 0; i < 10; i++ {
		stk.Push(i)
	}

	for stk.Any() {
		p := stk.Pop()
		t.Log(p)
	}
}
