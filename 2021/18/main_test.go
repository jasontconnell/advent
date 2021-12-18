package main

import (
	"testing"
)

func TestExplode(t *testing.T) {
	str := "[1,1]"

	list := parseInput([]string{str})

	s := list[0]
	cp := copyTree(s, nil)
	t.Log(cp)
	t.Log(s)
}
