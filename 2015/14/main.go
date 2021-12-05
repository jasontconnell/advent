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

var reg *regexp.Regexp = regexp.MustCompile(`^([\w]+) can fly ([0-9]+) km/s for ([0-9]+) seconds, but then must rest for ([0-9]+) seconds.$`)

type Deer struct {
	Name       string
	Speed      int
	Duration   int
	Rest       int
	Distance   int
	BurstStart int
	BurstStop  int
	NextBurst  int
	Points     int
}

func (deer Deer) String() string {
	return deer.Name + " traveled " + strconv.Itoa(deer.Distance) + " miles. He earned " + strconv.Itoa(deer.Points) + " points."
}

type DeerList []Deer

type DeerListDistanceSorter struct {
	Entries DeerList
}

func (p DeerListDistanceSorter) Len() int {
	return len(p.Entries)
}
func (p DeerListDistanceSorter) Less(i, j int) bool {
	return p.Entries[i].Distance < p.Entries[j].Distance
}
func (p DeerListDistanceSorter) Swap(i, j int) {
	p.Entries[i], p.Entries[j] = p.Entries[j], p.Entries[i]
}

type DeerListPointsSorter struct {
	Entries DeerList
}

func (p DeerListPointsSorter) Len() int {
	return len(p.Entries)
}
func (p DeerListPointsSorter) Less(i, j int) bool {
	return p.Entries[i].Points < p.Entries[j].Points
}
func (p DeerListPointsSorter) Swap(i, j int) {
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
	list := getInput(in)
	result := race(list, 2503, false)
	return output(result[0].Distance)
}

func part2(in input) output {
	list := getInput(in)
	result := race(list, 2503, true)
	return output(result[0].Points)
}

func race(list DeerList, max int, points bool) DeerList {
	for i := 0; i < max; i++ {
		for d := 0; d < len(list); d++ {
			deer := &list[d]

			if i >= deer.BurstStart && i < deer.BurstStop {
				deer.Distance += deer.Speed
			}

			if i == deer.BurstStop {
				deer.BurstStart = i + deer.Rest
				deer.BurstStop = deer.BurstStart + deer.Duration
			}
		}

		if points {
			sorter := DeerListDistanceSorter{Entries: list}
			sort.Sort(sort.Reverse(sorter))

			inLead := sorter.Entries[0].Distance
			for d := 0; d < len(list); d++ {
				if list[d].Distance == inLead {
					list[d].Points++
				}
			}
		}
	}

	var entries DeerList
	if !points {
		sorter := DeerListDistanceSorter{Entries: list}
		sort.Sort(sort.Reverse(sorter))
		entries = sorter.Entries
	} else {
		pointsSorter := DeerListPointsSorter{Entries: list}
		sort.Sort(sort.Reverse(pointsSorter))
		entries = pointsSorter.Entries
	}
	return entries
}

func getInput(in input) DeerList {
	list := DeerList{}

	for _, txt := range in {
		if groups := reg.FindStringSubmatch(txt); groups != nil && len(groups) > 1 {
			speed, _ := strconv.Atoi(groups[2])
			duration, _ := strconv.Atoi(groups[3])
			rest, _ := strconv.Atoi(groups[4])

			deer := Deer{Name: groups[1], Speed: speed, Duration: duration, Rest: rest, BurstStart: 0, BurstStop: duration}
			list = append(list, deer)
		}
	}

	return list
}
