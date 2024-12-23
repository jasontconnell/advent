package common

import (
	"testing"
)

func TestCombinations(t *testing.T) {
	list := []int{0, 1, 2}
	c := AllCombinations(list, 5)
	t.Log(c, len(c))
}

func TestCombinationsString(t *testing.T) {
	list := []string{"drax", "frankenstein", "connell"}
	list2 := []string{"rules", "stinks", "poop", "1", "2", "3"}

	cp := CartesianProduct(list, list2)

	t.Log(cp)
}
