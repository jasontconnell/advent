package main

import (
	"bufio"
	"fmt"
	"os"
	"time"
	"regexp"
	"strconv"
	"sort"
)
var input = "13.txt"

type Arrangement struct {
	People []string
	Units int
}

type ArrangementList []Arrangement

type ArrangementListSorter struct {
	Entries ArrangementList
}
func (p ArrangementListSorter) Len() int {
	return len(p.Entries)
}
func (p ArrangementListSorter) Less(i, j int) bool {
	return p.Entries[i].Units < p.Entries[j].Units
}
func (p ArrangementListSorter) Swap(i, j int) {
	p.Entries[i], p.Entries[j] = p.Entries[j], p.Entries[i]
}

func main() {
	startTime := time.Now()
	if f, err := os.Open(input); err == nil {
		scanner := bufio.NewScanner(f)

		reg := regexp.MustCompile("^([a-zA-Z]+) would (gain|lose) ([0-9]+) happiness units by sitting next to ([A-Za-z]+).$")

		list := make(map[string]int) // key is like Name-NextTo
		people := []string{}

		includeMe := true
		jason := "Jason"

		for scanner.Scan() {
			var txt = scanner.Text()
			if groups := reg.FindStringSubmatch(txt); groups != nil && len(groups) > 1 {
				name := groups[1]
				nextTo := groups[4]
				key := name + "-" + nextTo
				if units,err := strconv.Atoi(groups[3]); err == nil {
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

		if includeMe {
			people = append(people, jason)
		}

		fmt.Println("getting permutations at", time.Since(startTime))
		perms := Permutate(people)
		fmt.Println("got", len(perms), " permutation, finished at", time.Since(startTime))
		arrangements := ArrangementList{}

		for _,p := range perms {
			arr := Arrangement{ People: p }
			for i := 0; i < len(arr.People); i++ {
				// last person sits next to first person
				left := i
				right := (i+1) % len(arr.People)

				key := arr.People[left] + "-" + arr.People[right]
				key2 := arr.People[right] + "-" + arr.People[left]

				if units,ok := list[key]; ok {
					arr.Units += units
				}

				if units,ok := list[key2]; ok {
					arr.Units += units
				}
			}

			arrangements = append(arrangements, arr)
		}

		sorter := ArrangementListSorter{ Entries: arrangements }
		sort.Sort(sorter)

		fmt.Println(sorter.Entries[0])
		fmt.Println(sorter.Entries[len(sorter.Entries)-1])

	}

	fmt.Println("Time", time.Since(startTime))
}

func Permutate(str []string) [][]string {
	var ret [][]string

	if len(str) == 2 {
		ret = append(ret, []string{ str[0], str[1] })
		ret = append(ret, []string{ str[1], str[0] })
	} else {

		for i := 0; i < len(str); i++ {
			strc := make([]string, len(str))
			copy(strc, str)

			t := strc[i]
			sh := append(strc[:i], strc[i+1:]...)
			perm := Permutate(sh)
			
			for _,p := range perm {
				p = append([]string{ t }, p...)
				ret = append(ret, p)
			}
		}
	}

	return ret
}

func Contains(ss []string, s string) bool {
	for _,t := range ss {
		if t == s {
			return true
		}
	}
	return false
}