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
	p2 := matchTimestamp(sched)

	fmt.Println("Part 1:", p1)

	fmt.Println("Part 2:", p2)
	fmt.Println("Time", time.Since(startTime))
}

func firstBus(sched schedule) int {
	min := math.MaxInt32
	bid := 0

	for _, b := range sched.buses {
		if b == -1 {
			continue
		}
		x := sched.timestamp / b
		y := x*b + b - sched.timestamp

		if y < min {
			min = y
			bid = b
		}
	}
	return bid * min
}

func matchesTimestamp(ts, pos, bus int64) bool {
	b := (ts+pos)%bus == 0
	return b
}

func matchTimestamp(sched schedule) int64 {
	m := make(map[int64]int64)
	indices := []int64{}

	for i, b := range sched.buses {
		if b == -1 {
			continue
		}
		m[int64(i)] = int64(b)
		indices = append(indices, int64(i))
	}

	done := false
	var factor int64 = 1
	var factorStep int64 = 1
	var x int64
	var loops int64

	lasts := make(map[int]int64)
	maxes := make(map[int]int64)

	for !done {
		allFound := true
		x = m[indices[0]] * factor
		for idx, i := range indices[1:] {
			b := m[i]
			if (x+i)%b != 0 {
				allFound = false
				break
			} else if idx > 0 {
				if cmax, ok := maxes[idx]; ok {
					last, lok := lasts[idx]
					if lok {
						if factor-last > cmax {
							maxes[idx] = factor - last

							var nmax int64
							for _, mx := range maxes {
								if mx > nmax {
									nmax = mx
								}
							}
							factorStep = nmax
						}
					}

					lasts[idx] = factor
				} else {
					maxes[idx] = factor
					lasts[idx] = factor
				}
			}
		}
		done = allFound

		factor += factorStep
		loops++
	}

	return x
}

func getSchedule(lines []string) schedule {
	ts, _ := strconv.Atoi(lines[0])

	bbs := strings.Split(lines[1], ",")

	sched := schedule{timestamp: ts}

	for _, bb := range bbs {
		b, err := strconv.Atoi(bb)
		if err != nil {
			b = -1
		}

		sched.buses = append(sched.buses, b)
	}
	return sched
}
