package main

import (
	"fmt"
	"log"
	"os"

	"github.com/jasontconnell/advent/common"
)

type input = []int
type output = int

type sellPoint struct {
	a, b, c, d int
}

type highPrice struct {
	sp    sellPoint
	price int
}

func main() {
	in, err := common.ReadInts(common.InputFilename(os.Args))
	if err != nil {
		log.Fatal(err)
	}

	p1, p1time := common.Time(part1, in)
	p2, p2time := common.Time(part2, in)

	w := common.TeeOutput(os.Stdout)
	fmt.Fprintln(w, "--2024 day 22 solution--")
	fmt.Fprintln(w, "Part 1:", p1)
	fmt.Fprintln(w, "Part 2:", p2)
	fmt.Printf("Time %v (%v, %v)", p1time+p2time, p1time, p2time)
}

func part1(in input) output {
	sum := 0
	for _, v := range in {
		sum += getNSecretNumber(v, 2000)
	}
	return sum
}

func part2(in input) output {
	mem := make(map[int]map[sellPoint]int)
	for _, v := range in {
		chartSellPoints(v, 2000, mem)
	}

	counts := make(map[sellPoint]int)
	for _, v := range mem {
		for k2 := range v {
			counts[k2]++
		}
	}

	maxes := make(map[sellPoint]int)
	for sp, c := range counts {
		if sp.d <= 0 {
			continue
		}
		if m, ok := maxes[sp]; !ok || c > m {
			maxes[sp] = c
		}
	}
	// log.Println(maxes)

	max := 0
	for sp := range maxes {
		total := 0
		for _, v := range in {
			lm := mem[v]
			if price, ok := lm[sp]; ok {
				total += price
			}
		}
		if total > max {
			max = total
		}
	}
	return max
}

func getSellPrice(start, n int, check sellPoint, mem map[int]map[int]highPrice) int {
	if sm, ok := mem[start]; ok {
		for _, v := range sm {
			if v.sp == check {
				return v.price
			}
		}
	}
	return 0
}

func chartSellPoints(start, n int, mem map[int]map[sellPoint]int) {
	var sm map[sellPoint]int
	if ex, ok := mem[start]; ok {
		sm = ex
	} else {
		sm = make(map[sellPoint]int)
	}
	cur := start
	lastPrice := start % 10
	last4 := []int{}

	for i := 0; i < n; i++ {
		cur = getNextSecret(cur)
		price := cur % 10
		if len(last4) == 4 {
			last4 = append(last4[1:], price-lastPrice)
			sp := sellPoint{last4[0], last4[1], last4[2], last4[3]}
			if _, ok := sm[sp]; !ok {
				sm[sp] = price
			}
		} else {
			last4 = append(last4, price-lastPrice)
		}
		lastPrice = price
	}
	mem[start] = sm
}

func getNSecretNumber(start, n int) int {
	cur := start
	for i := 0; i < n; i++ {
		cur = getNextSecret(cur)
	}
	return cur
}

func getNextSecret(n int) int {
	n = ((n * 64) ^ n) % 16777216
	n = ((n / 32) ^ n) % 16777216
	n = ((n * 2048) ^ n) % 16777216
	return n
}
