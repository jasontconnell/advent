package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

var input = "06.txt"

type planet struct {
	name           string
	orbitedBy      []*planet
	orbitedByNames []string
	orbiting       *planet
}

func main() {
	startTime := time.Now()

	f, err := os.Open(input)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	scanner := bufio.NewScanner(f)
	planets := make(map[string]*planet)

	for scanner.Scan() {
		var txt = scanner.Text()
		sp := strings.Split(txt, ")")

		if p, ok := planets[sp[0]]; ok {
			p.orbitedByNames = append(p.orbitedByNames, sp[1])
		} else {
			planets[sp[0]] = &planet{name: sp[0], orbitedByNames: []string{sp[1]}}
		}
	}

	buildTree(planets)

	p1 := getAllOrbitCounts(planets)
	fmt.Println("Part 1: ", p1)

	from, to := planets["YOU"], planets["SAN"]
	p2 := findPathLength(from, to)

	fmt.Println("Part 2: ", p2)

	fmt.Println("Time", time.Since(startTime))
}

func buildTree(planets map[string]*planet) {
	for _, p := range planets {
		for _, orb := range p.orbitedByNames {
			if _, ok := planets[orb]; !ok {
				planets[orb] = &planet{name: orb, orbiting: p}
			}
		}
	}

	for _, p := range planets {
		for _, orb := range p.orbitedByNames {
			p.orbitedBy = append(p.orbitedBy, planets[orb])
			planets[orb].orbiting = p
		}
	}
}

func findPathLength(from, to *planet) int {
	checked := make(map[string]int)

	ref := from
	d := 0
	for ref != nil {
		checked[ref.name] = d
		d++
		ref = ref.orbiting
	}

	ref = to
	d2 := 0
	var dist int
	for ref != nil {
		if v, ok := checked[ref.name]; ok {
			dist = d2 + v
			break
		} else {
			checked[ref.name] = d2
			d2++
		}
		ref = ref.orbiting
	}

	return dist - 2
}

func getAllOrbitCounts(planets map[string]*planet) int {
	var total int
	for _, p := range planets {
		total += countTree(p)
	}

	return total
}

func countTree(p *planet) int {
	if p == nil {
		return 0
	}

	ref := p
	val := 0
	for ref != nil && ref.name != "COM" {
		val++
		ref = ref.orbiting
	}

	return val
}
