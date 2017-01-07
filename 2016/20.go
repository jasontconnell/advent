package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strconv"
	"time"
)

var input = "20.txt"

type IPRange struct {
	Low, High int64
}

type IPRangeList []IPRange

type IPRangeListSorter struct {
	Entries IPRangeList
}

func (p IPRangeListSorter) Len() int {
	return len(p.Entries)
}
func (p IPRangeListSorter) Less(i, j int) bool {
	return p.Entries[i].Low < p.Entries[j].Low
}
func (p IPRangeListSorter) Swap(i, j int) {
	p.Entries[i], p.Entries[j] = p.Entries[j], p.Entries[i]
}

func main() {
	startTime := time.Now()

	reg := regexp.MustCompile("^([0-9]+)-([0-9]+)$")
	ipranges := IPRangeList{}

	if f, err := os.Open(input); err == nil {
		scanner := bufio.NewScanner(f)

		for scanner.Scan() {
			var txt = scanner.Text()

			if groups := reg.FindStringSubmatch(txt); groups != nil && len(groups) > 1 {
				low, _ := strconv.ParseInt(groups[1], 10, 64)
				high, _ := strconv.ParseInt(groups[2], 10, 64)

				ipranges = append(ipranges, IPRange{Low: low, High: high})
			}
		}
	}

	sorter := IPRangeListSorter{Entries: ipranges}
	sort.Sort(sorter)

	merged := IPRangeList{}

	for _, ip := range ipranges {
		extended := false
		for i := 0; i < len(merged); i++ {
			m := merged[i]
			if withinRange(ip.Low, m.Low, m.High) && withinRange(ip.High, m.Low, m.High) { // fully contained already
				extended = true
				break
			}

			if !withinRange(ip.Low, m.Low, m.High) && withinRange(ip.High, m.Low, m.High) { // merged low needs extension down to ip low
				merged[i].Low = ip.Low
				extended = true
				break
			} else if withinRange(ip.Low, m.Low, m.High) && !withinRange(ip.High, m.Low, m.High) { // merged high needs extension to ip high
				merged[i].High = ip.High
				extended = true
				break
			}
		}

		if len(merged) == 0 || !extended {
			merged = append(merged, ip)
		}
	}

	sorter.Entries = merged
	sort.Sort(sorter)

	var count int64 = 0
	for i := 0; i < len(merged); i++ {
		count += (merged[i].High + 1 - merged[i].Low) // inclusive of high for total count
	}

	fmt.Println(merged[0].High + 1)
	fmt.Println("total valid", 4294967295+1-count)

	fmt.Println("Time", time.Since(startTime))
}

func withinRange(num, low, high int64) bool {
	return num >= low-1 && num <= high+1
}
