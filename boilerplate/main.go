package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jasontconnell/advent/common"
)

type input = []string
type output = int

func main() {
	startTime := time.Now()

	in, err := common.ReadStrings(common.InputFilename(os.Args))
	if err != nil {
		log.Fatal(err)
	}

	p1 := part1(in)
	p2 := part2(in)

	w := common.TeeOutput(os.Stdout)
	fmt.Fprintln(w, "--{year} day {day} solution--")
	fmt.Fprintln(w, "Part 1:", p1)
	fmt.Fprintln(w, "Part 2:", p2)
	fmt.Println("Time", time.Since(startTime))
}

func part1(in input) output {
	return 0
}

func part2(in input) output {
	return 0
}

// reg := regexp.MustCompile("-?[0-9]+")
/*
if groups := reg.FindStringSubmatch(txt); groups != nil && len(groups) > 1 {
				fmt.Println(groups[1:])
			}
*/
