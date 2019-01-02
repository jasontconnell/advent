package main

import (
	"bufio"
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"sort"
	"time"
	//"strconv"
	//"strings"
)

var input = "19.txt"
var input2 = "19.2.txt"

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
	list := ReplacementList{}
	if f, err := os.Open(input); err == nil {
		scanner := bufio.NewScanner(f)
		reg := regexp.MustCompile(`^(\w+) => (\w+)$`)

		for scanner.Scan() {
			var txt = scanner.Text()
			if groups := reg.FindStringSubmatch(txt); groups != nil && len(groups) > 1 {
				rep := Replacement{In: groups[1], Out: groups[2]}
				list = append(list, rep)
			}
		}
	}

	sorter := ReplacementListSorter{Entries: list}
	sort.Sort(sort.Reverse(sorter))

	inputMolecule := ""
	if tmp, err := ioutil.ReadFile(input2); err == nil {
		inputMolecule = string(tmp)
	}

	results := make(map[string]int)

	for _, rep := range sorter.Entries {
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

	maxoutput := 0

	for _, rep := range list {
		if len(rep.Out) > maxoutput {
			maxoutput = len(rep.Out)
		}
	}

	completed, total := Fabricate(sorter.Entries, inputMolecule, maxoutput, 0)

	fmt.Println("Steps to build desired molecule", total, completed)

	// fmt.Println("Steps to build desired molecule", minreplacements)
	fmt.Println("Unique molecules", len(results))
	fmt.Println("Time", time.Since(startTime))
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
