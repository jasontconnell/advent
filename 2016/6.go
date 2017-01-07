package main

import (
	"bufio"
	"fmt"
	"os"
	"time"
	//"regexp"
	"sort"
	//"strconv"
	//"strings"
	//"math"
)

var input = "6.txt"

func main() {
	startTime := time.Now()
	if f, err := os.Open(input); err == nil {
		scanner := bufio.NewScanner(f)

		resultMap := []map[string]int{}
		for i := 0; i < 8; i++ {
			resultMap = append(resultMap, make(map[string]int))
		}

		for scanner.Scan() {
			var txt = scanner.Text()
			processLine(txt, resultMap)
		}

		var runeCounts [8]RuneCountList
		result := ""
		p2result := ""

		for col, mp := range resultMap {
			for c, v := range mp {
				runeCounts[col] = append(runeCounts[col], RuneCount{Rune: c, Count: v})
			}

			SortRuneCounts(runeCounts[col])
			last := runeCounts[col][0]
			first := runeCounts[col][len(runeCounts[col])-1]
			result += last.Rune
			p2result += first.Rune
		}

		fmt.Println(result)
		fmt.Println("part 2", p2result)
	}

	fmt.Println("Time", time.Since(startTime))
}

func processLine(line string, resultMap []map[string]int) {
	for i, c := range line {

		cs := string(c)
		if _, ok := resultMap[i][cs]; ok {
			resultMap[i][cs]++
		} else {
			resultMap[i][cs] = 1
		}
	}
}

type RuneCount struct {
	Rune  string
	Count int
}
type RuneCountList []RuneCount

func (p RuneCountList) Len() int { return len(p) }
func (p RuneCountList) Less(i, j int) bool {
	less := p[i].Count < p[j].Count
	equal := p[i].Count == p[j].Count
	if equal {
		less = p[i].Rune > p[j].Rune
	}

	return less
}
func (p RuneCountList) Swap(i, j int) { p[i], p[j] = p[j], p[i] }

func SortRuneCounts(list RuneCountList) RuneCountList {
	sort.Sort(sort.Reverse(list))

	return list
}

// reg := regexp.MustCompile("-?[0-9]+")
/*
if groups := reg.FindStringSubmatch(txt); groups != nil && len(groups) > 1 {
                fmt.Println(groups[1:])
            }
*/
