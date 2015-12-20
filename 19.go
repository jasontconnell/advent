package main

import (
	"bufio"
	"fmt"
	"os"
	"time"
	"regexp"
	"io/ioutil"
	"crypto/md5"
	//"strconv"
	//"strings"
)
var input = "19.txt"
var input2 = "19.2.txt"

type Replacement struct {
	In, Out string
}

func main() {
	startTime := time.Now()
	list := []Replacement{}
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

	inputMolecule := ""
	if tmp,err := ioutil.ReadFile(input2); err == nil {
		inputMolecule = string(tmp)
	}

	
	// do replace across all molecules, add md5 hash of output to map if not exists
	results := make(map[string]int)
	
	//inputMolecule = "HOH"
	for _, rep := range list {
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

	outputMap := make(map[string]Replacement)
	for _, v := range list {
		if _,exists := outputMap[v.Out]; !exists {
			outputMap[v.Out] = v
		}
	}

	tokensMap := make(map[string]string)
	for _, v := range list {
		if _,exists := outputMap[v.In]; !exists {
			tokensMap[v.In] = v.In
		}
	}
	minreplacements := 1000000
	for _,electron := range list {
		// just those that are e are electrons
		if electron.In == "e" {
			startString := electron.Out
			desired := inputMolecule
			pos := getMatched(startString, desired)
			fmt.Println(startString)
			completed, total := Build(list, outputMap, tokensMap, startString, desired, pos, 0)

			if completed && total < minreplacements { total = minreplacements+1 }
		}
	}

	fmt.Println("Steps to build desired molecule", minreplacements)
	fmt.Println("unique molecules", len(results))

	fmt.Println("Time", time.Since(startTime))
}

func Build(list []Replacement, repmap map[string]Replacement, tokens map[string]string, current, desired string, pos, curstep int) (completed bool, steps int) {
	fmt.Println("------------------- called build")
	var cp string = current

	for i := pos+1; i <= len(current); i++ {
		// find a replacement and test
		token := current[pos:i]
		if _,ok := tokens[token]; ok {
			for j := 0; j < len(list); j++ {				
				if list[j].In == token {
					fmt.Println("checking replacement", list[j].Out, "for", token, i, j)
					cp2 := cp[:i-len(token)] + list[j].Out + cp[i:]
					nextpos := getMatched(cp2, desired)
					if nextpos > pos {
						return Build(list, repmap, tokens, cp2, desired, nextpos, curstep + 1)
					} else if nextpos == -1 {
						return true, curstep+1
					}
				}
			}
		}
	}
	completed = getMatched(current, desired) == -1
	steps = curstep
	return
}

func getMatched(current, desired string) int {
	matched := -1 // -1 indicates entire string matches

	for i := 0; i < len(current); i++ {
		if current[:i] != desired[:i] {
			matched = i-1
			fmt.Println("getmatched", current, matched)
			break
		}
	}
	return matched
}

func AllReplacements(in, out, input string) []string{
	cp := input
	repreg := regexp.MustCompile("(" + in + ")")

	list := []string{}

	// done := false
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

// reg := regexp.MustCompile("-?[0-9]+")
/* 			
if groups := reg.FindStringSubmatch(txt); groups != nil && len(groups) > 1 {
				fmt.Println(groups[1:])
			}
			*/
