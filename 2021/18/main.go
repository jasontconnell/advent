package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"time"

	"github.com/jasontconnell/advent/common"
)

type input = []string
type output = int64

type snailfish struct {
	id     *int
	parent *snailfish
	right  *snailfish
	left   *snailfish
}

var StopTraverse = fmt.Errorf("Stop traversing, fool")

type visitfunc func(s *snailfish, d int) error

func (s *snailfish) String() string {
	str := ""
	if s.left != nil {
		str += "["
	}
	if s.id != nil {
		str += fmt.Sprintf("%d", *s.id)
	}
	if s.left != nil {
		str += s.left.String()
	}
	if s.right != nil {
		if s.left != nil {
			str += ","
		}
		str += s.right.String()
		str += "]"
	}
	return str
}

func copyTree(s *snailfish, p *snailfish) *snailfish {
	if s == nil {
		return nil
	}
	r := &snailfish{id: s.id, parent: p}
	r.left = copyTree(s.left, r)
	r.right = copyTree(s.right, r)
	return r
}

func main() {
	startTime := time.Now()

	in, err := common.ReadStrings(common.InputFilename(os.Args))
	if err != nil {
		log.Fatal(err)
	}

	p1 := part1(in)
	p2 := part2(in)

	w := common.TeeOutput(os.Stdout)
	fmt.Fprintln(w, "--2021 day 18 solution--")
	fmt.Fprintln(w, "Part 1:", p1)
	fmt.Fprintln(w, "Part 2:", p2)
	fmt.Println("Time", time.Since(startTime))
}

func part1(in input) output {
	list := parseInput(in)
	var result *snailfish = list[0]
	for i := 1; i < len(list); i++ {
		result = addSnailfish(result, list[i])
	}
	return magnitude(result)
}

func part2(in input) output {
	list := parseInput(in)
	return largestMagnitude(list)
}

func largestMagnitude(list []*snailfish) int64 {
	var largest int64 = math.MinInt64

	for i := 0; i < len(list); i++ {
		for j := 0; j < len(list); j++ {
			if i == j {
				continue
			}
			prev, cur := copyTree(list[i], nil), copyTree(list[j], nil)
			ns := addSnailfish(prev, cur)
			m := magnitude(ns)
			if m > largest {
				largest = m
			}
		}
	}

	return largest
}

func magnitude(s *snailfish) int64 {
	m := int64(0)
	if s.left != nil {
		if s.left.id != nil {
			m += 3 * int64(*s.left.id)
		} else {
			m += 3 * magnitude(s.left)
		}
	}
	if s.right != nil {
		if s.right.id != nil {
			m += 2 * int64(*s.right.id)
		} else {
			m += 2 * magnitude(s.right)
		}
	}

	return m
}

func addSnailfish(a, b *snailfish) *snailfish {
	sn := &snailfish{left: a, right: b}
	a.parent = sn
	b.parent = sn
	sn = reduce(sn)
	return sn
}

func reduce(s *snailfish) *snailfish {
	for {
		reduced := false
		exp := toExplode(s)
		for exp != nil {
			reduced = true
			explode(exp)
			exp = toExplode(s)
		}

		sp := toSplit(s)
		if sp != nil {
			reduced = true
			split(sp)
		}

		if !reduced {
			break
		}
	}

	return s
}

func explode(s *snailfish) {
	nextLeft := getNextLeft(s)
	if nextLeft != nil {
		if nextLeft.id == nil {
			nextLeft = rightLeaf(nextLeft)
		}
		lid := addIds(s.left, nextLeft)
		nextLeft.id = lid
	}

	nextRight := getNextRight(s)
	if nextRight != nil {
		if nextRight.id == nil {
			nextRight = leftLeaf(nextRight)
		}
		rid := addIds(s.right, nextRight)
		nextRight.id = rid
	}

	newid := 0
	ns := &snailfish{id: &newid, parent: s.parent}
	if s == s.parent.right {
		ns.parent.right = ns
	} else if s == s.parent.left {
		ns.parent.left = ns
	}

	removeNode(s.left)
	removeNode(s.right)
	removeNode(s)
}

func removeNode(s *snailfish) {
	if s == nil {
		return
	}
	if s == s.parent.left {
		s.parent.left = nil
	}
	if s == s.parent.right {
		s.parent.right = nil
	}
	s.parent = nil
	s.left = nil
	s.right = nil
}

func split(s *snailfish) {
	val := *s.id

	fval := float64(val) / 2
	v1, v2 := int(math.Floor(fval)), int(math.Ceil(fval))

	left := &snailfish{id: &v1, parent: s}
	right := &snailfish{id: &v2, parent: s}

	s.id = nil
	s.left = left
	s.right = right
}

func leftLeaf(p *snailfish) *snailfish {
	for p.left != nil {
		p = p.left
		if p.left == nil {
			return p
		}
	}
	return p
}

func rightLeaf(p *snailfish) *snailfish {
	for p.right != nil {
		p = p.right
		if p.right == nil {
			return p
		}
	}
	return p
}

func getNextLeft(p *snailfish) *snailfish {
	orig := p
	for p != nil {
		p = p.parent
		if p != nil && p.left != nil && p.left != orig {
			return p.left
		}
		orig = p
	}
	return nil
}

func getNextRight(p *snailfish) *snailfish {
	orig := p
	for p != nil {
		p = p.parent
		if p != nil && p.right != nil && p.right != orig {
			return p.right
		}
		orig = p
	}
	return nil
}

func addIds(a, b *snailfish) *int {
	if a == nil || b == nil || a.id == nil || b.id == nil {
		return nil //new(int)
	}
	s := *a.id + *b.id
	return &s
}

func toSplit(s *snailfish) *snailfish {
	var sp *snailfish
	traverse(s, func(c *snailfish, d int) error {
		if c.id != nil && *c.id > 9 && sp == nil {
			sp = c
			return StopTraverse
		}
		return nil
	}, 0)
	return sp
}

func toExplode(s *snailfish) *snailfish {
	var exp *snailfish
	traverse(s, func(c *snailfish, d int) error {
		if d >= 4 && c.left != nil && c.left.id != nil && c.right != nil && c.right.id != nil {
			exp = c
			return StopTraverse
		}
		return nil
	}, 0)

	return exp
}

func traverse(s *snailfish, v visitfunc, depth int) error {
	err := v(s, depth)
	if err == StopTraverse {
		return err
	}
	if s.left != nil {
		err = traverse(s.left, v, depth+1)
		if err == StopTraverse {
			return err
		}
	}
	if s.right != nil {
		err = traverse(s.right, v, depth+1)
		if err == StopTraverse {
			return err
		}
	}
	return nil
}

func parseInput(in input) []*snailfish {
	var list []*snailfish
	for _, line := range in {
		s := getSnailfish(line)
		list = append(list, s)
	}
	return list
}

func getSnailfish(r string) *snailfish {
	root := &snailfish{}
	cur := root
	for _, c := range r {
		switch c {
		case '[':
			n := &snailfish{parent: cur}
			cur.left = n
			cur = n
		case ']':
			cur = cur.parent
		case ',':
			n := &snailfish{parent: cur}
			cur.right = n
			cur = n
		default:
			id, _ := strconv.Atoi(string(c))
			cur.id = &id
			cur = cur.parent
		}
	}
	return root
}
