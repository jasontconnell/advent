package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"time"
	//"regexp"
	//"strings"
	//"math"
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

	masses := []int{}
	for scanner.Scan() {
		var txt = scanner.Text()
		m, err := strconv.Atoi(txt)
		if err != nil {
			fmt.Println(err)
			continue
		}

		masses = append(masses, m)
	}

	fuel := getSimpleFuel(masses)
	allFuel := getAllFuel(masses)

	fmt.Println("Part 1: ", fuel)
	fmt.Println("Part 2: ", allFuel)

	fmt.Println("Time", time.Since(startTime))
}

func getFuel(m int) int {
	fp := (m / 3) - 2

	return fp
}

func getSimpleFuel(masses []int) int {
	fuel := 0
	for _, m := range masses {
		fuel += getFuel(m)
	}
	return fuel
}

func getAllFuel(masses []int) int {
	total := 0
	for _, m := range masses {
		f := getFuel(m)
		fp := f
		for fp > 0 {
			fp = getFuel(fp)
			if fp > 0 {
				f += fp
			}
		}

		total += f
	}

	return total
}

// reg := regexp.MustCompile("-?[0-9]+")
/*
if groups := reg.FindStringSubmatch(txt); groups != nil && len(groups) > 1 {
				fmt.Println(groups[1:])
			}
*/
