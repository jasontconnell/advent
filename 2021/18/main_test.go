package main

import (
	"testing"
)

func TestExplode(t *testing.T) {
	str := "[[[[0,[4,5]],[0,0]],[[[4,5],[2,6]],[9,5]]],[7,[[[3,7],[4,3]],[[6,3],[8,8]]]]]"

	list := parseInput([]string{str})

	s := list[0]
	exp := toExplode(s)
	printTree(s)
	t.Log(exp)
	t.Log(s)
}
