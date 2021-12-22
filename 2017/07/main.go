package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/jasontconnell/advent/common"
)

type input = []string
type output = int

type Program struct {
	Name     string
	Weight   int
	Children []*Program
	Above    []string
	All      int
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
	fmt.Fprintln(w, "--2017 day 07 solution--")
	fmt.Fprintln(w, "Part 1:", p1)
	fmt.Fprintln(w, "Part 2:", p2)
	fmt.Println("Time", time.Since(startTime))
}

func part1(in input) string {
	programs := parseInput(in)
	bottom := getBottom(programs)
	n := ""
	if bottom != nil {
		n = bottom.Name
	}
	return n
}

func part2(in input) output {
	programs := parseInput(in)
	bottom := getBottom(programs)
	buildTree(programs)
	setWeights(bottom)
	diff, value := findUnbalance(bottom)
	return diff + value
}

func parseInput(in input) []*Program {
	progs := []*Program{}
	for _, line := range in {
		p := getProgram(line)
		if p == nil {
			continue
		}
		progs = append(progs, p)
	}
	return progs
}

func setWeights(p *Program) {
	p.All = sumAll(p)

	for _, c := range p.Children {
		setWeights(c)
	}
}

func findUnbalance(p *Program) (d, val int) {
	wmap := make(map[int]int)
	for _, c := range p.Children {
		w := c.All
		wmap[w]++
	}

	n := 0
	for k, v := range wmap {
		if v > 1 {
			n = k
		}
	}

	for k, v := range wmap {
		if v == 1 {
			for _, c := range p.Children {
				if c.All == k {
					// this is the mismatched one
					diff := n - c.All

					sub, w := findUnbalance(c)
					if sub == 0 {
						return diff, c.Weight
					} else {
						return sub, w
					}
				}
			}
		}
	}

	return 0, p.Weight
}

func sumAll(p *Program) int {
	w := p.Weight
	for _, c := range p.Children {
		w += sumAll(c)
	}
	return w
}

func buildTree(programs []*Program) {
	pmap := make(map[string]*Program)
	for _, p := range programs {
		pmap[p.Name] = p
	}

	for _, p := range pmap {
		for _, a := range p.Above {
			c, ok := pmap[a]

			if ok {
				p.Children = append(p.Children, c)
			}
		}
	}
}

func getBottom(programs []*Program) *Program {
	nmap := make(map[string]string)
	for _, p := range programs {
		for _, a := range p.Above {
			nmap[a] = a
		}
	}

	var bottom *Program
	for _, p := range programs {
		if _, ok := nmap[p.Name]; !ok {
			bottom = p

			break
		}
	}
	return bottom
}

func getProgram(line string) *Program {
	var reg = regexp.MustCompile("^(.*?) \\(([0-9]+)\\)(.*)$")
	if groups := reg.FindStringSubmatch(line); groups != nil && len(groups) > 1 {
		name := groups[1]
		weight, err := strconv.Atoi(groups[2])
		if err != nil {
			fmt.Println("parsing", err)
			return nil
		}

		children := groups[3]
		children = strings.Replace(children, "-> ", "", -1)
		children = strings.Replace(children, ",", "", -1)

		above := strings.Split(children, " ")

		p := &Program{Name: name, Weight: weight, Above: above}
		return p
	}
	return nil
}
