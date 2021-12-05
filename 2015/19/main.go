package main

import (
	"crypto/md5"
	"fmt"
	"log"
	"regexp"
	"sort"
	"time"

	"github.com/jasontconnell/advent/common"
)

var inputFilename = "input.txt"

type input []string
type output int

var reg *regexp.Regexp = regexp.MustCompile(`^(\w+) => (\w+)$`)

type Replacement struct {
	In, Out string
}

type ReplacementList []Replacement

type ReplacementListSorter struct {
	Entries ReplacementList
}

func (p ReplacementListSorter) Len() int {
	return len(p.Entries)
}
func (p ReplacementListSorter) Less(i, j int) bool {
	return len(p.Entries[i].Out) > len(p.Entries[j].Out)
}
func (p ReplacementListSorter) Swap(i, j int) {
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
	list, inputMolecule := getInput(in)
	m := getUniqueMolecules(list, inputMolecule)
	return output(len(m))
}

func part2(in input) output {
	list, inputMolecule := getInput(in)
	maxoutput := 0

	for _, rep := range list {
		if len(rep.Out) > maxoutput {
			maxoutput = len(rep.Out)
		}
	}

	_, total := Fabricate(list, inputMolecule, maxoutput, 0)
	return output(total)
}

func getUniqueMolecules(list ReplacementList, inputMolecule string) map[string]int {
	results := make(map[string]int)

	for _, rep := range list {
		replacements := AllReplacements(rep.In, rep.Out, inputMolecule)

		for _, r := range replacements {
			md5 := MD5s(r)
			if _, exists := results[md5]; !exists {
				results[md5] = 1
			} else {
				results[md5]++
			}
		}
	}

	return results
}

func getInput(in input) (ReplacementList, string) {
	list := ReplacementList{}
	var inputMolecule string
	for _, txt := range in {
		if groups := reg.FindStringSubmatch(txt); groups != nil && len(groups) > 1 {
			rep := Replacement{In: groups[1], Out: groups[2]}
			list = append(list, rep)
		} else if txt != "" {
			inputMolecule = txt
		}
	}

	sorter := ReplacementListSorter{Entries: list}
	sort.Sort(sort.Reverse(sorter))

	return sorter.Entries, inputMolecule
}

func Fabricate(list []Replacement, current string, maxlen, step int) (completed bool, steps int) {
	for i := len(current) - 1; i >= 0; i-- {
		if len(current) == 1 && current == "e" {
			return true, step
		} else {
			for j := 2; j <= maxlen; j++ {
				if i+j > len(current) {
					continue
				}
				token := current[i : i+j]
				for _, rep := range list {
					var newcur string
					build := false

					if token == rep.Out {
						newcur = current[:i] + rep.In + current[i+j:]
						build = true
					}

					if build {
						c, s := Fabricate(list, newcur, maxlen, step+1)
						if c {
							return c, s
						} else {
							break
						}
					}
				}
			}
		}
	}
	return
}

func AllReplacements(in, out, input string) []string {
	cp := input
	repreg := regexp.MustCompile("(" + in + ")")

	list := []string{}
	loc := repreg.FindAllStringIndex(cp, -1)

	for _, indices := range loc {
		cp2 := string(cp[:indices[0]]) + out + string(cp[indices[1]:])
		list = append(list, cp2)
	}

	return list
}

func MD5(content []byte) string {
	sum := md5.Sum(content)
	return fmt.Sprintf("%x", sum)
}

func MD5s(content string) string {
	return MD5([]byte(content))
}
