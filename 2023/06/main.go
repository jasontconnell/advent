package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/jasontconnell/advent/common"
)

type input = []string
type output = int

func main() {
	in, err := common.ReadStrings(common.InputFilename(os.Args))
	if err != nil {
		log.Fatal(err)
	}

	p1, p1time := common.Time(part1, in)
	p2, p2time := common.Time(part2, in)

	w := common.TeeOutput(os.Stdout)
	fmt.Fprintln(w, "--2023 day 06 solution--")
	fmt.Fprintln(w, "Part 1:", p1)
	fmt.Fprintln(w, "Part 2:", p2)
	fmt.Printf("Time %v (%v, %v)", p1time+p2time, p1time, p2time)
}

func part1(in input) output {
	ts, ds := parseInput(in)
	wins := 1
	for i := 0; i < len(ts); i++ {
		rwins := race(ts[i], ds[i], 1)
		wins *= rwins
	}
	fmt.Println(ts, ds)
	return wins
}

func part2(in input) output {
	return 0
}

func race(t int, d int, acc int) int {
	queue := []int{}
	queue = append(queue, 1)

	wins := 0

	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]

		speed := cur * acc

		dist := (t - cur) * speed

		if dist > d {
			wins++
		}

		if cur+1 < t {
			queue = append(queue, cur+1)
		}

		fmt.Println(cur, speed, dist)
	}

	return wins
}

func parseInput(in input) ([]int, []int) {
	timestr := strings.Fields(strings.Split(in[0], ":")[1])
	diststr := strings.Fields(strings.Split(in[1], ":")[1])

	ts := []int{}
	ds := []int{}

	for _, t := range timestr {
		v, _ := strconv.Atoi(t)
		ts = append(ts, v)
	}

	for _, d := range diststr {
		v, _ := strconv.Atoi(d)
		ds = append(ds, v)
	}

	return ts, ds
}
