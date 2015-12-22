package main

import (
	"bufio"
	"fmt"
	"os"
	"time"
	"regexp"
	"io/ioutil"
	"crypto/md5"
	"sort"
	//"strconv"
	"strings"
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
	return len(p.Entries[i].Out) < len(p.Entries[j].Out)
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
				rep := Replacement{ In: groups[1], Out: groups[2] }
				list = append(list, rep)
			}
		}
	}

	sorter := ReplacementListSorter{ Entries: list }
	sort.Sort(sort.Reverse(sorter))


	inputMolecule := ""
	if tmp,err := ioutil.ReadFile(input2); err == nil {
		inputMolecule = string(tmp)
	}

	outputmap := make(map[string]Replacement)
	for _,rep := range list {
		if _,exists := outputmap[rep.Out]; !exists {
			outputmap[rep.Out] = rep
			fmt.Println(rep.Out)
		}
	}

	//inputMolecule = "ORnFArSiThRnPMgAr"

	// do replace across all molecules, add md5 hash of output to map if not exists
	results := make(map[string]int)
	
	//inputMolecule = "HOH"
	for _, rep := range sorter.Entries {
		replacements := AllReplacements(rep.In, rep.Out, inputMolecule)
		
		for _, r := range replacements {
			md5 := MD5s(r)
			if _,exists := results[md5]; !exists {
				results[md5] = 1
			} else {
				results[md5]++
			}	
		}
	}

	// c, s := Build2(list, []string {"HF"}, inputMolecule, 0)
	// fmt.Println(c,s)
	// return

	minreplacements := 1000000
	// just those that are e are electrons
	for _,electron := range getReplacementsForMolecule("e", list) {
		startString := electron.Out
		fmt.Println(startString)
		//completed, total := Build(sorter.Entries, startString, inputMolecule, 0, 0)
		//completed, total := Build2(sorter.Entries, []string { startString }, inputMolecule, 0, 0)
		completed,total := Build3(outputmap, inputMolecule, 0)
		total = total + 1
		fmt.Println("\n####### on path", startString, "took", total, "steps and completed =", completed, " #######\n")
	}

	fmt.Println("Steps to build desired molecule", minreplacements)
	fmt.Println("Unique molecules", len(results))

	fmt.Println("Time", time.Since(startTime))
}

func Build3(outmap map[string]Replacement, current string, step int) (completed bool, steps int) {
	for i := len(current)-1; i >= 0; i-- {
		if len(current) == 1 && current == "e" {
			return true, step
		} else {
			token := current[i:]

			//fmt.Println(outmap[token])


			if rep,ok := outmap[token]; ok {
				newcur := current[:i] + rep.In

				fmt.Println("found token", token, rep)
				return

				fmt.Println(token, newcur)

				// c,s := Build3(outmap, newcur, step+1)
				// if c {
				// 	return c, s
				// }
			}
		}
	}
	return
}

func Build2(list []Replacement, current []string, desired string, step, matched int) (completed bool, steps int) {
	var unique = make(map[string]string)

	for _,rep := range list {
		for _,cur := range current {
			values := AllReplacements(rep.In, rep.Out, cur)
			
			for _,val := range values {
				if _,ok := unique[val]; !ok {
					unique[val] = val
				}
			}
		}
	}

	maxmatch := matched
	newstrings := []string{}
	for k,_ := range unique {
		newstrings = append(newstrings, k)
		if m := getMatched(k, desired); m > maxmatch {
			maxmatch = m
		}
	}

	fmt.Println("new max match ", maxmatch)
	recurse := []string{}

	if maxmatch > matched {
		for _,str := range newstrings {
			if m := getMatched(str, desired); m >= maxmatch {
				recurse = append(recurse, str)
			}
		}
	}

	if len(recurse) > 0 {
		fmt.Println("testing", len(recurse), "strings", maxmatch)
		fmt.Println(recurse)
		return Build2(list, recurse, desired, step+1, maxmatch)
	} else {
		return false, -1
	}
}

func Build(list []Replacement, current, desired string, curpos, curstep int) (completed bool, steps int) {
	fmt.Println(current)
	for _,rep := range list {
		replacements := AllReplacements(rep.In, rep.Out, current)
		if len(replacements) > 0 { fmt.Println(replacements) }

		for _,str := range replacements {
			matched := getMatched(str, desired)

			if strings.Contains(current,"NRnFYFArF") {
				fmt.Println("in", rep.In, "out", rep.Out, "current", current, "str", str, "matched", matched, desired[:len(rep.In)])
			}

			if len(str) > len(desired){
				return false, -1
			} else if matched > 0 {
				fmt.Println("recursing with ", str[matched:], "des", desired[matched:])
				c,s := Build(list, str[matched:], desired[matched:], matched, curstep+1)
				fmt.Println(" - done recursing with ", str[matched:], "des", desired[matched:])
				if c {
					return c, s	
				}
			}

			fmt.Println(matched, str, desired)

			// if matched == -2 {
			// 	return true, curstep
			// } else if matched != -1 {
			// 	fmt.Println("-------------calling build------------------")
			// 	subbuildcomplete,substeps := Build(list, str[matched:], desired, matched, curstep+1)
			// 	fmt.Println("-------------returned from call------------------")
			// 	if subbuildcomplete {
			// 		fmt.Println("found match", rep, str, desired)
			// 		return true, substeps
			// 	}
			// }
		}
	}

	return
}

func getReplacementsForMolecule(molecule string, list []Replacement) []Replacement {
	ret := []Replacement{}
	for _, rep := range list {
		if rep.In == molecule {
			ret = append(ret, rep)
		}
	}
	return ret
}

func getMatched(current, desired string) int {
	matched := -1 
	//fmt.Println("start get matched", current, desired, matched)
	for i := 0; i < len(current) && i < len(desired); i++ {
		if current[i] == desired[i] {
			//fmt.Println(current[i], "==", desired[i], i)
			matched = i+1
		} else {
			break
		}
	}
	//fmt.Println("get matched", current, desired, matched, len(current))
	if matched == len(current) { matched = -2 }
	return matched
}

func AllReplacements(in, out, input string) []string{
	cp := input
	repreg := regexp.MustCompile("(" + in + ")")

	list := []string{}
	loc := repreg.FindAllStringIndex(cp, -1)

	for _,indices := range loc {
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