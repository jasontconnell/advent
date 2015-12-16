package main

import (
	"bufio"
	"fmt"
	"os"
	"time"
	"regexp"
	"strconv"
	"sort"
	//"strings"
)
var input = "16.txt"

type Property struct {
	Name string
	Value int
}

type Sue struct {
	Number string
	Properties map[string]Property
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

func NewSue(num string) Sue {
	sue := Sue{Number: num, Properties: make(map[string]Property)}
	return sue
}



func main() {
	startTime := time.Now()
	if f, err := os.Open(input); err == nil {
		scanner := bufio.NewScanner(f)
		reg := regexp.MustCompile(`^Sue ([0-9]+): ([\w]+): ([0-9]+), ([\w]+): ([0-9]+), ([\w]+): ([0-9]+)$`)

		list := SueList{}

		for scanner.Scan() {
			var txt = scanner.Text()
			if groups := reg.FindStringSubmatch(txt); groups != nil && len(groups) > 1 {
				prop1 := groups[2]
				val1,_ := strconv.Atoi(groups[3])
				prop2 := groups[4]
				val2,_ := strconv.Atoi(groups[5])
				prop3 := groups[6]
				val3,_ := strconv.Atoi(groups[7])

				sue := NewSue(groups[1])
				sue.Properties[prop1] = Property{ Name: prop1, Value: val1 }
				sue.Properties[prop2] = Property{ Name: prop2, Value: val2 }
				sue.Properties[prop3] = Property{ Name: prop3, Value: val3 }
				list = append(list, sue)
			}
		}
		
		answer := []Property{}
		answer = append(answer, Property{ Name: "children", Value: 3})
		answer = append(answer, Property{ Name: "cats", Value: 7})
		answer = append(answer, Property{ Name: "samoyeds", Value: 2})
		answer = append(answer, Property{ Name: "pomeranians", Value: 3})
		answer = append(answer, Property{ Name: "akitas", Value: 0})
		answer = append(answer, Property{ Name: "vizslas", Value: 0})
		answer = append(answer, Property{ Name: "goldfish", Value: 5})
		answer = append(answer, Property{ Name: "trees", Value: 3})
		answer = append(answer, Property{ Name: "cars", Value: 2})
		answer = append(answer, Property{ Name: "perfumes", Value: 1})

		for i := 0; i < len(list); i++ {
			sue := &list[i]
			var total float64
			for _,prop := range answer {
				if p,ok := sue.Properties[prop.Name]; ok {
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
							pct = 100
						} else if greater && diff > 0 {
							pct = 100
						}
					} else if diff != 0 {
						pct = 1.0/float64(diff)
						
					} else {
						pct = 100
					}

					total += pct
				}
			}

			sue.MatchPercent = total / 3
		}

		sorter := SueListSorter{ Entries: list }
		sort.Sort(sort.Reverse(sorter))
		fmt.Println("First:", sorter.Entries[0])
		fmt.Println("Last:", sorter.Entries[len(sorter.Entries)-1])

	}

	fmt.Println("Time", time.Since(startTime))
}
/* 			
if groups := reg.FindStringSubmatch(txt); groups != nil && len(groups) > 1 {
				fmt.Println(groups[1:])
			}
			*/
