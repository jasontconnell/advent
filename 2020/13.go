package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
	"time"
)

var input = "13.txt"

type schedule struct {
	timestamp int
	buses     []int
}

func main() {
	startTime := time.Now()

	f, err := os.Open(input)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	scanner := bufio.NewScanner(f)

	lines := []string{}
	for scanner.Scan() {
		var txt = scanner.Text()
		lines = append(lines, txt)
	}

	sched := getSchedule(lines)

	p1 := firstBus(sched)

	fmt.Println("Part 1:", p1)
	fmt.Println("Time", time.Since(startTime))
}

func firstBus(sched schedule) int {
	min := math.MaxInt32
	bid := 0

	for _, b := range sched.buses {
		x := sched.timestamp / b
		y := x*b + b - sched.timestamp

		if y < min {
			min = y
			bid = b
		}
	}
	return bid * min
}

func getSchedule(lines []string) schedule {
	ts, _ := strconv.Atoi(lines[0])

	bbs := strings.Split(lines[1], ",")

	sched := schedule{timestamp: ts}

	for _, bb := range bbs {
		if bb == "x" {
			continue
		}

		b, _ := strconv.Atoi(bb)

		sched.buses = append(sched.buses, b)
	}
	return sched
}
