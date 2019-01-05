package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

var input = "24.txt"

type component struct {
	id           int
	port1, port2 int
}

type connection struct {
	id          int
	left, right int
}

type state struct {
	conn       connection
	connectors int
	sum        int
	conns      []int
}

func main() {
	startTime := time.Now()

	f, err := os.Open(input)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	scanner := bufio.NewScanner(f)
	components := []component{}
	id := 1
	for scanner.Scan() {
		var txt = scanner.Text()
		c := readComponent(txt, id)
		components = append(components, c)
		id++
	}

	p1, p2 := solve(components)

	fmt.Println("Part 1:", p1)
	fmt.Println("Part 2:", p2)

	fmt.Println("Time", time.Since(startTime))
}

func solve(components []component) (int, int) {
	start := connection{left: 0, right: 0}
	strsolves := getPath(start, components, true)
	maxstr := 0

	for _, s := range strsolves {
		if s.sum > maxstr {
			maxstr = s.sum
		}
	}

	lensolves := getPath(start, components, false)
	maxlen := 0
	maxlenstr := 0
	for _, s := range lensolves {
		if s.connectors > maxlen {
			maxlen = s.connectors
			maxlenstr = s.sum
		}
	}
	return maxstr, maxlenstr
}

func getPath(start connection, components []component, strength bool) []state {
	maxsum := 0
	maxlen := 0
	queue := []state{state{conn: start, connectors: 0, sum: maxsum, conns: []int{}}}
	solves := []state{}

	for len(queue) > 0 {
		st := queue[0]
		queue = queue[1:]

		mvs := findConnections(st.conn, components, st.conns)
		if strength && len(mvs) == 0 && st.sum > maxsum {
			maxsum = st.sum
			solves = append(solves, st)
		} else if !strength && len(mvs) == 0 && st.connectors > maxlen {
			maxlen = st.connectors
			solves = append(solves, st)
		}

		for _, mv := range mvs {
			cp := make([]int, len(st.conns)+1)
			copy(cp, st.conns)
			cp[len(cp)-1] = mv.id
			mvstate := state{conn: mv, connectors: st.connectors + 1, sum: st.sum + mv.left + mv.right, conns: cp}
			if strength {
				queue = append(queue, mvstate)
			} else {
				if st.connectors+1 > maxlen {
					queue = append(queue, mvstate)
				}
			}
		}
	}

	return solves
}

func contains(id int, ids []int) bool {
	c := false
	for _, ii := range ids {
		if id == ii {
			c = true
			break
		}
	}
	return c
}

func findConnections(to connection, components []component, ids []int) []connection {
	conns := []connection{}
	for _, c := range components {
		if contains(c.id, ids) {
			continue
		}
		add := false
		cn := connection{}
		if c.port1 == to.right || c.port2 == to.right {
			cn.id = c.id
			add = true
			if c.port1 == to.right {
				cn.left = c.port1
				cn.right = c.port2
			} else if c.port2 == to.right {
				cn.left = c.port2
				cn.right = c.port1
			}
		}

		if add {
			conns = append(conns, cn)
		}
	}
	return conns
}

func readComponent(line string, id int) component {
	ps := strings.Split(line, "/")
	p1, _ := strconv.Atoi(ps[0])
	p2, _ := strconv.Atoi(ps[1])

	return component{id: id, port1: p1, port2: p2}
}
