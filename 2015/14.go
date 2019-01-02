package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"time"
	//"strings"
	"sort"
)

var input = "14.txt"

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
	if f, err := os.Open(input); err == nil {
		scanner := bufio.NewScanner(f)
		reg := regexp.MustCompile(`^([\w]+) can fly ([0-9]+) km/s for ([0-9]+) seconds, but then must rest for ([0-9]+) seconds.$`)

		list := DeerList{}

		for scanner.Scan() {
			var txt = scanner.Text()
			if groups := reg.FindStringSubmatch(txt); groups != nil && len(groups) > 1 {
				speed, _ := strconv.Atoi(groups[2])
				duration, _ := strconv.Atoi(groups[3])
				rest, _ := strconv.Atoi(groups[4])

				deer := Deer{Name: groups[1], Speed: speed, Duration: duration, Rest: rest, BurstStart: 0, BurstStop: duration}
				list = append(list, deer)
			}
		}

		max := 2503

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

			sorter := DeerListDistanceSorter{Entries: list}
			sort.Sort(sort.Reverse(sorter))

			inLead := sorter.Entries[0].Distance
			for d := 0; d < len(sorter.Entries); d++ {
				if sorter.Entries[d].Distance == inLead {
					sorter.Entries[d].Points++
				}
			}
		}

		sorter := DeerListDistanceSorter{Entries: list}
		sort.Sort(sort.Reverse(sorter))
		fmt.Println("First place for furthest distance:", sorter.Entries[0])

		pointsSorter := DeerListPointsSorter{Entries: list}
		sort.Sort(sort.Reverse(pointsSorter))

		fmt.Println("First place for most points:", sorter.Entries[0])

	}

	fmt.Println("Time", time.Since(startTime))
}
