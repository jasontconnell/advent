package main

import (
	"fmt"
	"testing"
)

func TestSplit(t *testing.T) {
	part := parsePart("##.#../#.##../###.../....../.#..#./..#..#")

	sp := split(part)
	for _, s := range sp {
		fmt.Println()
		print(s)
	}

	t.Log(sp)
}

// func TestJoin(t *testing.T) {
// 	strparts := []string{"##./#../...", "##./#../...", "##./#../...", "##./#../..."}
// 	parts := [][][]bool{}
// 	for _, str := range strparts {
// 		p := parsePart(str)
// 		parts = append(parts, p)
// 	}

// 	joined := join(parts)
// 	print(joined)
// }
