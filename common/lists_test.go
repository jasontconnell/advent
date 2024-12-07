package common

import (
	"testing"
)

func TestCombinations(t *testing.T) {
	list := []int{0, 1, 2}
	c := AllCombinations(list, 5)
	t.Log(c, len(c))
}
