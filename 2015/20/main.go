package main

import (
	"fmt"
	"log"
	"math"
	"time"

	"github.com/jasontconnell/advent/common"
)

var inputFilename = "input.txt"

type input int
type output int

func main() {
	startTime := time.Now()

	in, err := common.ReadInt(inputFilename)
	if err != nil {
		log.Fatal(err)
	}

	p1 := part1(input(in))
	p2 := part2(input(in))

	fmt.Println("Part 1:", p1)
	fmt.Println("Part 2:", p2)

	fmt.Println("Time", time.Since(startTime))
}

func part1(in input) output {
	house := 100000
	presents := 0

	for presents < int(in) {
		presents = GetPresents(house)

		if presents > int(in) {
			break
		}
		house++
	}
	return output(house)
}

func part2(in input) output {
	house := 100000
	presents2 := 0
	for presents2 < int(in) {
		presents2 = GetPresents2(house)

		if presents2 > int(in) {
			fmt.Println("got presents2", presents2, "house", house)
			break
		}
		house++
	}

	return output(house)
}

func GetPresents(max int) (presents int) {
	sqrt := int(math.Sqrt(float64(max))) + 1
	for i := 1; i <= sqrt; i++ {
		if max%i == 0 {
			presents += i
			presents += max / i
		}
	}
	return presents * 10
}

func GetPresents2(max int) (presents int) {

	for i := 1; i <= 50; i++ {
		if max%i == 0 {
			presents += i
			presents += max / i
		}
	}
	return presents * 11
}
