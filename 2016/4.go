package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
	//"math"
	"sort"
)

var input = "4.txt"

func main() {
	startTime := time.Now()
	if f, err := os.Open(input); err == nil {
		scanner := bufio.NewScanner(f)
		reg := regexp.MustCompile("(.*?)-([0-9]+)\\[([a-z]+)\\]")
		sectorSum := 0

		for scanner.Scan() {
			var txt = scanner.Text()
			if groups := reg.FindStringSubmatch(txt); groups != nil && len(groups) > 1 {
				inputPristine := groups[1]
				input := SortRunes(strings.Replace(inputPristine, "-", "", -1))
				sector, _ := strconv.Atoi(groups[2])
				checksum := groups[3]
				lineResult := processLine(input, checksum, sector)
				sectorSum += lineResult

				decrypted := decrypt(inputPristine, sector)
				if lineResult != 0 {
					fmt.Println(decrypted, sector)
				}
			}
		}

		fmt.Println("Sector sum", sectorSum)
	}

	fmt.Println("Time", time.Since(startTime))
}

func decrypt(input string, sector int) string {
	newstr := ""
	for _, c := range input {
		if string(c) != "-" {
			newstr += string(((int(c) - 97 + sector) % 26) + 97)
		} else {
			newstr += " "
		}
	}

	return newstr
}

func processLine(input, checksum string, sector int) (result int) {
	if len(checksum) != 5 { return }
	cmap := make(map[string]RuneCount)
	for _, c := range input {
		cs := string(c)
		if entry, ok := cmap[cs]; ok {
			entry.Count++
			cmap[cs] = entry
		} else {
			cmap[cs] = RuneCount{ Rune: cs, Count: 1 }
		}
	}

	list := RuneCountList{}
	for _, rc := range cmap {
		list = append(list, rc)
	}

	SortRuneCounts(list)
	firstFive := list[:5]
	realRoom := true

	for _, rc := range firstFive {
		if !strings.ContainsRune(checksum, rune(rc.Rune[0])) {
			realRoom = false
		}
	}

	if realRoom {
		result = sector
	}

	return
}

type RuneCount struct {
	Rune string
	Count int
}
type RuneCountList []RuneCount

func (p RuneCountList) Len() int           { return len(p) }
func (p RuneCountList) Less(i, j int) bool { 
	less := p[i].Count < p[j].Count
	equal := p[i].Count == p[j].Count
	if equal {
		less = p[i].Rune > p[j].Rune
	}

	return less
}
func (p RuneCountList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

func SortRuneCounts(list RuneCountList) RuneCountList {
	sort.Sort(sort.Reverse(list))
	
	return list
}


type RuneSlice []rune

func (p RuneSlice) Len() int           { return len(p) }
func (p RuneSlice) Less(i, j int) bool { return p[i] < p[j] }
func (p RuneSlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

func SortRunes(s string) string {
	runes := []rune(s)
	sort.Sort(RuneSlice(runes))
	return string(runes)
}

// reg := regexp.MustCompile("-?[0-9]+")
/*
if groups := reg.FindStringSubmatch(txt); groups != nil && len(groups) > 1 {
                fmt.Println(groups[1:])
            }
*/
