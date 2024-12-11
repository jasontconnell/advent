package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/jasontconnell/advent/common"
)

type input = string
type output = int

type stone struct {
	digits string
	next   *stone
	prev   *stone
}

func main() {
	in, err := common.ReadString(common.InputFilename(os.Args))
	if err != nil {
		log.Fatal(err)
	}

	p1, p1time := common.Time(part1, in)
	p2, p2time := common.Time(part2, in)

	w := common.TeeOutput(os.Stdout)
	fmt.Fprintln(w, "--2024 day 11 solution--")
	fmt.Fprintln(w, "Part 1:", p1)
	fmt.Fprintln(w, "Part 2:", p2)
	fmt.Printf("Time %v (%v, %v)", p1time+p2time, p1time, p2time)
}

func part1(in input) output {
	m := parse(in)
	return count(blink(m, 25))
}

func part2(in input) output {
	m := parse(in)
	return count(blink(m, 75))
}

func count(m map[string]int) int {
	c := 0
	for k, v := range m {
		if v < 0 {
			log.Println("< 0", k, v)
			break
		}
		c += v
	}
	return c
}

func blink(m map[string]int, n int) map[string]int {
	for i := 0; i < n; i++ {
		log.Println(i, len(m))
		m = change(m)
	}
	return m
}

func change(m map[string]int) map[string]int {
	list := []string{}
	counts := make(map[string]int)
	for k, v := range m {
		if v == 0 {
			continue
		}
		counts[k] = v
		list = append(list, k)
	}
	for _, k := range list {
		count := counts[k]
		x, _ := strconv.Atoi(k)
		if k == "0" {
			m["1"] += count
			m["0"] -= count
		} else if len(k)%2 == 0 {
			ldigits := k[:len(k)/2]
			rdigits := k[len(k)/2:]

			lni, _ := strconv.Atoi(ldigits)
			rni, _ := strconv.Atoi(rdigits)
			ln := fmt.Sprintf("%d", lni)
			rn := fmt.Sprintf("%d", rni)

			m[k] -= count
			m[ln] += count
			m[rn] += count
		} else {
			n := x * 2024
			nn := strconv.Itoa(n)
			m[nn] += count
			m[k] -= count
		}
	}
	return m
}

func parse(in string) map[string]int {
	sp := strings.Fields(in)
	m := make(map[string]int)
	for _, s := range sp {
		m[s]++
	}
	return m
}
