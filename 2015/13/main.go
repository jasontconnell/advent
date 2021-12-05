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

var reg *regexp.Regexp = regexp.MustCompile("^([a-zA-Z]+) would (gain|lose) ([0-9]+) happiness units by sitting next to ([A-Za-z]+).$")

type Arrangement struct {
	People []string
	Units  int
}

type ArrangementList []Arrangement

type ArrangementListSorter struct {
	Entries ArrangementList
}

func (p ArrangementListSorter) Len() int {
	return len(p.Entries)
}
func (p ArrangementListSorter) Less(i, j int) bool {
	return p.Entries[i].Units > p.Entries[j].Units
}
func (p ArrangementListSorter) Swap(i, j int) {
	p.Entries[i], p.Entries[j] = p.Entries[j], p.Entries[i]
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
	data, people := getInput(in)
	entries := getOptimalArrangement(data, people)
	return output(entries[0].Units)
}

func part2(in input) output {
	data, people := getInput(in)
	people = append(people, "Jason")
	entries := getOptimalArrangement(data, people)
	return output(entries[0].Units)
}

func getOptimalArrangement(data map[string]int, people []string) ArrangementList {
	perms := Permutate(people)
	arrangements := ArrangementList{}

	for _, p := range perms {
		arr := Arrangement{People: p}
		for i := 0; i < len(arr.People); i++ {
			// last person sits next to first person
			left := i
			right := (i + 1) % len(arr.People)

			key := arr.People[left] + "-" + arr.People[right]
			key2 := arr.People[right] + "-" + arr.People[left]

			if units, ok := data[key]; ok {
				arr.Units += units
			}

			if units, ok := data[key2]; ok {
				arr.Units += units
			}
		}

		arrangements = append(arrangements, arr)
	}

	sorter := ArrangementListSorter{Entries: arrangements}
	sort.Sort(sorter)

	return sorter.Entries
}

func getInput(in input) (map[string]int, []string) {
	list := make(map[string]int) // key is like Name-NextTo
	people := []string{}

	for _, txt := range in {
		if groups := reg.FindStringSubmatch(txt); groups != nil && len(groups) > 1 {
			name := groups[1]
			nextTo := groups[4]
			key := name + "-" + nextTo
			if units, err := strconv.Atoi(groups[3]); err == nil {
				if groups[2] == "lose" {
					units = -units
				}
				list[key] = units
			} else {
				panic(err)
			}

			if !Contains(people, name) {
				people = append(people, name)
			}
			if !Contains(people, nextTo) {
				people = append(people, nextTo)
			}
		}
	}

	return list, people
}

func Permutate(str []string) [][]string {
	var ret [][]string

	if len(str) == 2 {
		ret = append(ret, []string{str[0], str[1]})
		ret = append(ret, []string{str[1], str[0]})
	} else {

		for i := 0; i < len(str); i++ {
			strc := make([]string, len(str))
			copy(strc, str)

			t := strc[i]
			sh := append(strc[:i], strc[i+1:]...)
			perm := Permutate(sh)

			for _, p := range perm {
				p = append([]string{t}, p...)
				ret = append(ret, p)
			}
		}
	}

	return ret
}

func Contains(ss []string, s string) bool {
	for _, t := range ss {
		if t == s {
			return true
		}
	}
	return false
}
