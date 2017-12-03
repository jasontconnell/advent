package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"time"
)

var input = "01.txt"

func main() {
	startTime := time.Now()

	f, err := os.Open(input)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		var txt = scanner.Text()
		final := getSum(txt)
		fmt.Println("final", final)

		final2 := getSumPart2(txt)
		fmt.Println("final part 2", final2)
	}

	fmt.Println("Time", time.Since(startTime))
}

func getSum(line string) int {
	sum := 0
	last := -1
	for i := 0; i <= len(line); i++ {
		idx := i

		if idx == len(line) {
			idx = 0
		}
		n, err := strconv.Atoi(string(line[idx]))
		if err != nil {
			fmt.Println(err)
			return -1
		}

		if last == n {
			sum += n
		}
		last = n
	}

	return sum
}

func getSumPart2(line string) int {
	sum := 0
	step := len(line) / 2

	for i := 0; i < len(line); i++ {
		cur, err := strconv.Atoi(string(line[i]))

		if err != nil {
			fmt.Println("parsing", err)
			return -1
		}

		half := (i + step) % len(line)

		cmp, err := strconv.Atoi(string(line[half]))
		if err != nil {
			fmt.Println("parsing", err)
			return -1
		}

		if cur == cmp {
			sum += cur
		}
	}

	return sum
}
