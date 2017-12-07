package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var input = "07.txt"

var reg = regexp.MustCompile("^(.*?) \\(([0-9]+)\\)(.*)$")

type Program struct {
	Name     string
	Weight   int
	Children []*Program
	Above    []string
	All      int
	Level    int
}

func main() {
	startTime := time.Now()

	f, err := os.Open(input)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	scanner := bufio.NewScanner(f)

	programs := []*Program{}

	for scanner.Scan() {
		var txt = scanner.Text()
		p := getProgram(txt)
		if p == nil {
			continue
		}
		programs = append(programs, p)
	}

	bottom := getBottom(programs)

	buildTree(programs)
	setLevels(bottom)
	setWeights(bottom)
	diff, value := findUnbalance(bottom)

	fmt.Println("part 1, bottom is           ", bottom.Name)
	fmt.Println("part 2, change weight to    ", diff+value)

	fmt.Println("Time", time.Since(startTime))
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

func setLevels(p *Program) {
	for _, c := range p.Children {
		c.Level = p.Level + 1
		setLevels(c)
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
			bottom.Level = 0

			break
		}
	}
	return bottom
}

func getProgram(line string) *Program {
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
