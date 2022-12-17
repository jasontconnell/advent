package main

import (
	"testing"
)

func TestReleased(t *testing.T) {
	list := []release{
		{minute: 3, pressure: 20},
		{minute: 6, pressure: 13},
		{minute: 10, pressure: 21},
		{minute: 18, pressure: 22},
		{minute: 22, pressure: 3},
		{minute: 25, pressure: 2},
	}

	r := getReleased(list, 30)
	t.Log(r)
}

func TestCustom(t *testing.T) {
	rel := []release{
		{5, 14, 0},
		{9, 25, 0},
		{12, 19, 0},
		{16, 10, 0},
		{20, 20, 0},
		{23, 24, 0},
		{26, 18, 0},
		{29, 8, 0},
		{30, 4, 0},
	}

	min := 30
	total := 0
	for _, r := range rel {
		total += r.pressure * (min - r.minute + 1)
	}
	t.Log(total)
}
