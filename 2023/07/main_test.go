package main

import (
	"fmt"
	"testing"
)

func TestHandStrength(t *testing.T) {
	// 32T3K
	h := hand{cards: []card{{'3', 3}, {'2', 2}, {'T', 10}, {'3', 3}, {'K', 14}}}

	s := handStrength(h)
	fmt.Println(s)
}
