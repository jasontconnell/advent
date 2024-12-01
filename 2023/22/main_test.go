package main

import (
	"testing"
)

func TestOverlap(t *testing.T) {
	tests := []struct {
		a1, a2, b1, b2 int
		expect         bool
	}{
		{1, 3, 2, 3, true},
		{1, 3, 3, 4, false},
		{1, 1, 0, 2, true},
		{1, 3, 2, 4, false},
	}

	for _, u := range tests {
		res := overlaps(u.a1, u.a2, u.b1, u.b2)
		if res != u.expect {
			t.Fail()
			t.Log("expected", u.expect, "got", res, u)
		}
	}
}

func TestCollidesPoint(t *testing.T) {
	tests := []struct {
		a1, a2, b1, b2 int
		expect         bool
	}{
		{1, 3, 2, 3, true},
		{1, 3, 3, 4, false},
		{1, 1, 0, 2, true},
		{1, 3, 2, 4, false},
		{1, 5, 4, 4, true},
		{1, 5, 5, 6, false},
	}

	for _, u := range tests {
		res := pointsCollide(u.a1, u.a2, u.b1, u.b2)
		if res != u.expect {
			t.Fail()
			t.Log("expected", u.expect, "got", res, u)
		}
	}
}
