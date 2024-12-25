package main

import (
	"fmt"
	"log"
	"os"

	"github.com/jasontconnell/advent/common"
)

type input = []string
type output = int

type lock struct {
	heights []int
}

type key struct {
	heights []int
}

func main() {
	in, err := common.ReadStrings(common.InputFilename(os.Args))
	if err != nil {
		log.Fatal(err)
	}

	p1, p1time := common.Time(part1, in)

	w := common.TeeOutput(os.Stdout)
	fmt.Fprintln(w, "--2024 day 25 solution--")
	fmt.Fprintln(w, "Part 1:", p1)
	fmt.Printf("Time %v", p1time)
}

func part1(in input) output {
	locks, keys := parse(in)
	return lockpick(locks, keys)
}

func lockpick(locks []lock, keys []key) int {
	total := 0
	km := make(map[int][]int)
	for i := 0; i < len(locks); i++ {
		for j := 0; j < len(keys); j++ {
			l, k := locks[i], keys[j]

			fits := true
			for h := 0; h < 5; h++ {
				if l.heights[h]+k.heights[h] > 5 {
					fits = false
				}
			}
			if fits {
				km[j] = append(km[j], i)
				total++
			}
		}
	}
	return total
}

func parse(in []string) ([]lock, []key) {
	locks := []lock{}
	keys := []key{}
	for i := 0; i < len(in); i++ {
		line := in[i]
		if line == "" {
			continue
		}

		if line == "#####" {
			// in a lock
			lockLines := in[i+1 : i+6]
			h := getHeights(lockLines)
			locks = append(locks, lock{h})
			i += 6
		} else if line == "....." {
			keyLines := in[i+1 : i+6]
			h := getHeights(keyLines)
			keys = append(keys, key{h})
			i += 6
		}
	}

	return locks, keys
}

func getHeights(lines []string) []int {
	heights := []int{0, 0, 0, 0, 0}
	for _, line := range lines {
		for i := 0; i < 5; i++ {
			if line[i] == '#' {
				heights[i]++
			}
		}
	}
	return heights
}
