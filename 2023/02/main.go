package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/jasontconnell/advent/common"
)

type input = []string
type output = int

type game struct {
	id       int
	cubesets []cubeset
}

type cubeset struct {
	cubes []cube
}

type cube struct {
	color string
	num   int
}

func main() {
	in, err := common.ReadStrings(common.InputFilename(os.Args))
	if err != nil {
		log.Fatal(err)
	}

	p1, p1time := common.Time(part1, in)
	p2, p2time := common.Time(part2, in)

	w := common.TeeOutput(os.Stdout)
	fmt.Fprintln(w, "--2023 day 02 solution--")
	fmt.Fprintln(w, "Part 1:", p1)
	fmt.Fprintln(w, "Part 2:", p2)
	fmt.Printf("Time %v (%v, %v)", p1time+p2time, p1time, p2time)
}

func part1(in input) output {
	games := parseInput(in)
	return getPossible(games, 12, 13, 14)
}

func part2(in input) output {
	return 0
}

func getPossible(games []game, r, g, b int) int {
	possible := 0
	for _, gm := range games {
		if isPossible(gm, r, g, b) {
			possible += gm.id
		}
	}
	return possible
}

func isPossible(gm game, r, g, b int) bool {
	return getMaxCubeDraws(gm, "red") <= r && getMaxCubeDraws(gm, "green") <= g && getMaxCubeDraws(gm, "blue") <= b
}

func getMaxCubeDraws(gm game, color string) int {
	max := -1
	for _, cs := range gm.cubesets {
		for _, c := range cs.cubes {
			if c.color == color && c.num > max {
				max = c.num
			}
		}
	}
	return max
}

func parseInput(lines input) []game {
	reg := regexp.MustCompile("^Game ([0-9]+): (.*)$")
	creg := regexp.MustCompile("([0-9]+) ([a-z]+)")
	games := []game{}
	for _, line := range lines {
		m := reg.FindStringSubmatch(line)
		id, _ := strconv.Atoi(m[1])

		glist := m[2]

		cubesets := []cubeset{}
		for _, s := range strings.Split(glist, ";") {
			cs := cubeset{}
			for _, c := range strings.Split(s, ",") {
				cm := creg.FindAllStringSubmatch(strings.Trim(c, " "), -1)
				for _, mm := range cm {
					num, _ := strconv.Atoi(mm[1])

					cb := cube{color: mm[2], num: num}
					cs.cubes = append(cs.cubes, cb)
				}
			}
			cubesets = append(cubesets, cs)
		}

		games = append(games, game{id: id, cubesets: cubesets})
	}
	return games
}
