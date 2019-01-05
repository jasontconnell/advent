package main

import (
	"fmt"
	"time"
)

const (
	left int = iota
	right
)

type action struct {
	move      int
	next      *state
	nextstate rune
	write     int
}

func (a action) String() string {
	return fmt.Sprintf("Next: %c Move: %d  Write: %d", a.nextstate, a.move, a.write)
}

type state struct {
	code      rune
	onaction  action
	offaction action
}

func (s *state) String() string {
	return fmt.Sprintf("%c - On: %v  Off: %v", s.code, s.onaction, s.offaction)
}

type tape struct {
	left  *tape
	right *tape
	val   int
}

func main() {
	startTime := time.Now()
	doChecksum := 12208951
	start := 'A'

	p1 := run(start, doChecksum, getStates())

	fmt.Println("Part 1:", p1)

	fmt.Println("Time", time.Since(startTime))
}

func run(start rune, doChecksum int, states map[rune]*state) int {
	cursor := &tape{left: nil, right: nil, val: 0}
	curstate, _ := states[start]
	i := 0
	for i < doChecksum {
		a := curstate.onaction
		if cursor.val == 0 {
			a = curstate.offaction
		}

		cursor.val = a.write

		if a.move == left {
			r := cursor
			cursor = cursor.left
			if cursor == nil {
				cursor = &tape{val: 0, right: r}
				r.left = cursor
			}
		} else {
			l := cursor
			cursor = cursor.right
			if cursor == nil {
				cursor = &tape{val: 0, left: l}
				l.right = cursor
			}
		}

		curstate = a.next

		i++
	}

	return sum(cursor)
}

func print(t *tape) {
	c := t
	for c.left != nil {
		c = c.left
	}

	s := ""
	for {
		f := "%d "
		if c == t {
			f = "[%d] "
		}
		s += fmt.Sprintf(f, c.val)
		c = c.right
		if c == nil {
			break
		}
	}
}

func sum(c *tape) int {
	cur := c
	for cur.left != nil {
		cur = cur.left
	}

	s := 0
	for {
		s += cur.val
		cur = cur.right
		if cur == nil {
			break
		}
	}
	return s
}

func getTestStates() map[rune]*state {
	m := make(map[rune]*state)
	a := &state{
		code:      'A',
		offaction: getAction(right, 1, 'B'),
		onaction:  getAction(left, 0, 'B'),
	}
	m[a.code] = a

	b := &state{
		code:      'B',
		offaction: getAction(left, 1, 'A'),
		onaction:  getAction(right, 1, 'A'),
	}
	m[b.code] = b

	for _, s := range m {
		offn, _ := m[s.offaction.nextstate]
		s.offaction.next = offn

		onn, _ := m[s.onaction.nextstate]
		s.onaction.next = onn
	}

	return m
}

func getStates() map[rune]*state {
	m := make(map[rune]*state)
	a := &state{
		code:      'A',
		offaction: getAction(right, 1, 'B'),
		onaction:  getAction(left, 0, 'E'),
	}
	m[a.code] = a

	b := &state{
		code:      'B',
		offaction: getAction(left, 1, 'C'),
		onaction:  getAction(right, 0, 'A'),
	}
	m[b.code] = b

	c := &state{
		code:      'C',
		offaction: getAction(left, 1, 'D'),
		onaction:  getAction(right, 0, 'C'),
	}
	m[c.code] = c

	d := &state{
		code:      'D',
		offaction: getAction(left, 1, 'E'),
		onaction:  getAction(left, 0, 'F'),
	}
	m[d.code] = d

	e := &state{
		code:      'E',
		offaction: getAction(left, 1, 'A'),
		onaction:  getAction(left, 1, 'C'),
	}
	m[e.code] = e

	f := &state{
		code:      'F',
		offaction: getAction(left, 1, 'E'),
		onaction:  getAction(right, 1, 'A'),
	}
	m[f.code] = f

	for _, s := range m {
		offn, _ := m[s.offaction.nextstate]
		s.offaction.next = offn

		onn, _ := m[s.onaction.nextstate]
		s.onaction.next = onn
	}

	return m
}

func getAction(move, write int, next rune) action {
	return action{move: move, write: write, nextstate: next}
}
