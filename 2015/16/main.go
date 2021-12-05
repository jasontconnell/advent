package main

import (
	"fmt"
	"log"
	"regexp"
	"sort"
	"strconv"
	"time"

	"github.com/jasontconnell/advent/common"
)

var inputFilename = "input.txt"

type input []string
type output int

var reg *regexp.Regexp = regexp.MustCompile(`^Sue ([0-9]+): ([\w]+): ([0-9]+), ([\w]+): ([0-9]+), ([\w]+): ([0-9]+)$`)

type Property struct {
	Name  string
	Value int
}

type Sue struct {
	Number       int
	Properties   map[string]Property
	MatchPercent float64
}

type SueList []Sue

type SueListSorter struct {
	Entries SueList
}

func (p SueListSorter) Len() int {
	return len(p.Entries)
}
func (p SueListSorter) Less(i, j int) bool {
	return p.Entries[i].MatchPercent < p.Entries[j].MatchPercent
}
func (p SueListSorter) Swap(i, j int) {
	p.Entries[i], p.Entries[j] = p.Entries[j], p.Entries[i]
}

func NewSue(num int) Sue {
	sue := Sue{Number: num, Properties: make(map[string]Property)}
	return sue
}

func main() {
	startTime := time.Now()

	in, err := common.ReadStrings(inputFilename)
	if err != nil {
		log.Fatal(err)
	}

	p1 := part1(input(in))
	p2 := part2(input(in))

	fmt.Println("Part 1:", p1)
	fmt.Println("Part 2:", p2)

	fmt.Println("Time", time.Since(startTime))
}

func part1(in input) output {
	answers := []Property{}
	answers = append(answers, Property{Name: "children", Value: 3})
	answers = append(answers, Property{Name: "cats", Value: 7})
	answers = append(answers, Property{Name: "samoyeds", Value: 2})
	answers = append(answers, Property{Name: "pomeranians", Value: 3})
	answers = append(answers, Property{Name: "akitas", Value: 0})
	answers = append(answers, Property{Name: "vizslas", Value: 0})
	answers = append(answers, Property{Name: "goldfish", Value: 5})
	answers = append(answers, Property{Name: "trees", Value: 3})
	answers = append(answers, Property{Name: "cars", Value: 2})
	answers = append(answers, Property{Name: "perfumes", Value: 1})

	sueList := getInput(in)

	list := process(sueList, answers, false)

	return output(list[0].Number)
}

func part2(in input) output {
	answers := []Property{}
	answers = append(answers, Property{Name: "children", Value: 3})
	answers = append(answers, Property{Name: "cats", Value: 7})
	answers = append(answers, Property{Name: "samoyeds", Value: 2})
	answers = append(answers, Property{Name: "pomeranians", Value: 3})
	answers = append(answers, Property{Name: "akitas", Value: 0})
	answers = append(answers, Property{Name: "vizslas", Value: 0})
	answers = append(answers, Property{Name: "goldfish", Value: 5})
	answers = append(answers, Property{Name: "trees", Value: 3})
	answers = append(answers, Property{Name: "cars", Value: 2})
	answers = append(answers, Property{Name: "perfumes", Value: 1})

	sueList := getInput(in)

	list := process(sueList, answers, true)

	return output(list[0].Number)
}

func process(list SueList, answers []Property, fuzzy bool) SueList {
	for i := 0; i < len(list); i++ {
		sue := &list[i]

		if !fuzzy && match(*sue, answers) {
			return SueList{*sue}
		} else {
			var total float64
			for _, prop := range answers {
				if p, ok := sue.Properties[prop.Name]; ok {
					diff := p.Value - prop.Value

					rng := false
					greater := false
					less := false

					if prop.Name == "cats" || prop.Name == "trees" {
						rng, greater = true, true
					}

					if prop.Name == "pomeranians" || prop.Name == "goldfish" {
						rng, less = true, true
					}

					pct := 0.0
					if rng {
						if less && diff < 0 {
							pct = 1
						} else if greater && diff > 0 {
							pct = 1
						}
					} else if diff != 0 {
						pct = 1.0 / float64(diff)
					} else {
						pct = 1
					}

					total += pct
				}
			}
			sue.MatchPercent = total / 3
		}
	}

	sorter := SueListSorter{Entries: list}
	sort.Sort(sort.Reverse(sorter))

	return sorter.Entries
}

func match(sue Sue, props []Property) bool {
	m := true

	for _, p := range props {
		if sp, ok := sue.Properties[p.Name]; ok {
			if p.Value != sp.Value {
				m = false
				break
			}
		}
	}
	return m
}

func getInput(in input) SueList {
	list := SueList{}

	for _, txt := range in {
		if groups := reg.FindStringSubmatch(txt); groups != nil && len(groups) > 1 {
			prop1 := groups[2]
			val1, _ := strconv.Atoi(groups[3])
			prop2 := groups[4]
			val2, _ := strconv.Atoi(groups[5])
			prop3 := groups[6]
			val3, _ := strconv.Atoi(groups[7])

			n, _ := strconv.Atoi(groups[1])

			sue := NewSue(n)
			sue.Properties[prop1] = Property{Name: prop1, Value: val1}
			sue.Properties[prop2] = Property{Name: prop2, Value: val2}
			sue.Properties[prop3] = Property{Name: prop3, Value: val3}
			list = append(list, sue)
		}
	}
	return list
}
