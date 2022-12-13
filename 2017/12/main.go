package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/jasontconnell/advent/common"
)

type input = []string
type output = int

var reg = regexp.MustCompile("^([0-9]+) <-> (.*)$")

type Program struct {
	ID       string
	Piped    []*Program
	PipedIDs []string
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
	fmt.Fprintln(w, "--2017 day 12 solution--")
	fmt.Fprintln(w, "Part 1:", p1)
	fmt.Fprintln(w, "Part 2:", p2)
	fmt.Println("Time", time.Since(startTime))
}

func part1(in input) output {
	p := []*Program{}
	for _, line := range in {
		p = append(p, getProgram(line))
	}
	mapPrograms(p)
	c := getConnected(p, "0")
	return len(c)
}

func part2(in input) output {
	p := []*Program{}
	for _, line := range in {
		p = append(p, getProgram(line))
	}
	mapPrograms(p)
	return getGroups(p)
}

func getGroups(programs []*Program) int {
	// find out which group each program belongs
	// the trick to making it run fast is to remove visited programs from the list of programs to check
	groups := make(map[string]string)
	vmap := make(map[string]bool)
	cp := make([]*Program, len(programs))

	copy(cp, programs)

	for _, p := range programs {
		if _, ok := groups[p.ID]; !ok {
			r := getConnected(cp, p.ID)

			for _, rp := range r {
				vmap[rp.ID] = true
				if _, ok := groups[rp.ID]; !ok {
					groups[rp.ID] = p.ID
				}
			}

			// remove visited from the list of programs to check
			for i := len(cp) - 1; i >= 0; i-- {
				sp := cp[i]
				if _, ok := vmap[sp.ID]; ok {
					cp = append(cp[:i], cp[i+1:]...)
				}
			}
		}
	}

	// determine unique groups
	m := make(map[string]string)
	for _, v := range groups {
		m[v] = v
	}

	return len(m)
}

func getConnected(programs []*Program, id string) []*Program {
	connected := []*Program{}

	for _, p := range programs {
		vmap := make(map[string]bool)

		if connects(p, id, vmap) {
			connected = append(connected, p)
		}
	}
	return connected
}

func connects(p *Program, id string, vmap map[string]bool) bool {
	if p.ID == id {
		return true
	}

	val := false
	for _, piped := range p.Piped {
		_, visited := vmap[piped.ID]
		vmap[piped.ID] = true

		if !visited {
			r := connects(piped, id, vmap)
			if r {
				val = true
				break
			}
		}
	}

	return val
}

func mapPrograms(programs []*Program) {
	pmap := make(map[string]*Program)
	for _, p := range programs {
		pmap[p.ID] = p
	}

	for _, p := range programs {
		for _, id := range p.PipedIDs {
			piped, ok := pmap[id]
			if !ok {
				fmt.Println("couldn't find", id, "program", p.ID, "piped ids", p.PipedIDs)
				break
			}

			p.Piped = append(p.Piped, piped)
		}
	}
}

func getProgram(line string) *Program {
	var p *Program
	if groups := reg.FindStringSubmatch(line); groups != nil && len(groups) > 1 {
		id := groups[1]
		ids := strings.Split(strings.Replace(groups[2], " ", "", -1), ",")

		p = &Program{ID: id, PipedIDs: ids, Piped: []*Program{}}

	}

	return p
}
