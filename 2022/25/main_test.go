package main

import "testing"

func TestSnafu(t *testing.T) {
	list := []struct {
		from string
		to   int
	}{
		// {"2=-1=0", 4890},
		// {"1", 1},
		// {"2", 2},
		{"1=", 3},
		{"1-", 4},
		{"10", 5},
		{"11", 6},
		{"2=", 8},
		{"2-", 9},
		{"20", 10},
		{"1=11-2", 2022},
		{"1121-1110-1=0", 314159265},
	}

	for _, s := range list {
		t.Log("testing", s.from, "<->", s.to)
		n := getSnafu(s.from)
		if n != s.to {
			t.Log(s.to, "!=", n)
			t.Fail()
		}
		sn := asSnafu(n)
		if s.from != sn {
			t.Log(s.from, "!=", sn)
			t.Fail()
		}
	}
}
