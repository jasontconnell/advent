package main

import (
	"fmt"
	"log"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/jasontconnell/advent/common"
)

var inputFilename = "input.txt"

type input = []string
type output = int

type room struct {
	raw, sorted string
	sector      int
	checksum    string
}

type runeCount struct {
	ch    rune
	count int
}

func main() {
	startTime := time.Now()

	in, err := common.ReadStrings(inputFilename)
	if err != nil {
		log.Fatal(err)
	}

	p1 := part1(in)
	p2 := part2(in)

	fmt.Println("Part 1:", p1)
	fmt.Println("Part 2:", p2)

	fmt.Println("Time", time.Since(startTime))
}

func part1(in input) output {
	rooms := parseInput(in)
	return sumOfRooms(rooms)
}

func part2(in input) output {
	s := 0
	rooms := parseInput(in)
	for _, r := range rooms {
		if isRoom(r) {
			n := decrypt(r)
			if n == "northpole object storage" {
				s = r.sector
			}
		}
	}
	return s
}

func sumOfRooms(rooms []room) output {
	s := 0
	for _, r := range rooms {
		s += processRoom(r)
	}
	return s
}

func isRoom(r room) bool {
	cmap := make(map[rune]runeCount)
	for _, c := range r.sorted {
		if entry, ok := cmap[c]; ok {
			entry.count++
			cmap[c] = entry
		} else {
			cmap[c] = runeCount{ch: c, count: 1}
		}
	}

	list := []runeCount{}
	for _, rc := range cmap {
		list = append(list, rc)
	}

	sort.Slice(list, func(i, j int) bool {
		less := list[i].count < list[j].count
		equal := list[i].count == list[j].count
		if equal {
			less = list[i].ch > list[j].ch
		}

		return less
	})

	realRoom := true
	for _, rc := range list[len(list)-5:] {
		if !strings.ContainsRune(r.checksum, rc.ch) {
			realRoom = false
		}
	}
	return realRoom
}

func processRoom(r room) output {
	result := 0
	if isRoom(r) {
		result = r.sector
	}

	return result
}

func decrypt(r room) string {
	newstr := ""
	for _, c := range r.raw {
		if c != '-' {
			newstr += string(((int(c) - 97 + r.sector) % 26) + 97)
		} else {
			newstr += " "
		}
	}

	return newstr
}

func parseInput(in input) []room {
	rooms := []room{}
	reg := regexp.MustCompile("(.*?)-([0-9]+)\\[([a-z]+)\\]")
	for _, line := range in {
		if groups := reg.FindStringSubmatch(line); groups != nil && len(groups) > 1 {
			inputPristine := groups[1]
			input := sortRunes(strings.Replace(inputPristine, "-", "", -1))
			sector, _ := strconv.Atoi(groups[2])
			checksum := groups[3]

			r := room{raw: groups[1], sorted: input, sector: sector, checksum: checksum}
			rooms = append(rooms, r)
		}
	}
	return rooms
}

func sortRunes(s string) string {
	b := []byte(s)
	sort.Slice(b, func(i, j int) bool {
		return b[i] < b[j]
	})
	return string(b)
}
