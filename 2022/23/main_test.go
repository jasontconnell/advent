package main

import (
	"testing"
)

func TestMoveToBack(t *testing.T) {
	list := []int{4, 6, 2, 5, 0, 12}

	list = moveToBack(list, 5)

	t.Log(list)
}
