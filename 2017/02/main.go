package main

import (
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
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
	fmt.Fprintln(w, "--2017 day 02 solution--")
	fmt.Fprintln(w, "Part 1:", p1)
	fmt.Fprintln(w, "Part 2:", p2)
	fmt.Println("Time", time.Since(startTime))
}

func part1(in input) output {
	rows := parseInput(in)
	return getChecksum(rows)
}

func part2(in input) output {
	rows := parseInput(in)
	return getDivisionResult(rows)
}

func parseInput(in input) [][]int {
	rows := make([][]int, len(in))

	for i := 0; i < len(in); i++ {
		flds := strings.Fields(in[i])

		for _, f := range flds {
			n, _ := strconv.Atoi(f)
			rows[i] = append(rows[i], n)
		}

		sort.Ints(rows[i])
	}
	return rows
}

func getChecksum(rows [][]int) int {
	sum := 0

	for _, row := range rows {
		min, max := getMinMax(row)
		sum += (max - min)
	}

	return sum
}

func getDivision(row []int) (div int) {
	for i := 0; i < len(row); i++ {
		for j := len(row) - 1; j >= 0; j-- {
			if i == j {
				continue
			}
			if row[j]%row[i] == 0 {
				return row[j] / row[i]
			}
		}
	}

	return 0
}

func getDivisionResult(rows [][]int) int {
	res := 0

	for _, line := range rows {
		div := getDivision(line)
		res += div
	}

	return res
}

func getMinMax(row []int) (min int, max int) {
	return row[0], row[len(row)-1]
}
