package main

import (
	"fmt"
	"log"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/jasontconnell/advent/common"
)

var inputFilename = "input.txt"

type input = []string
type output = int

type Dir struct {
	Turn  int
	Moves int
}

type xy struct {
	x, y int
}

func main() {
	startTime := time.Now()

	in, err := common.ReadStringsCsv(inputFilename)
	if err != nil {
		log.Fatal(err)
	}

	p1 := part1(in)
	p2 := part2(in)

	fmt.Println("Part 1:", p1)
	fmt.Println("Part 2:", p2)

	fmt.Println("Time", time.Since(startTime))
}

func part1(in input) output {
	dirs := parseInput(in)
	pt, _ := navigate(dirs)

	x := int(math.Abs(float64(pt.x)))
	y := int(math.Abs(float64(pt.y)))
	return x + y
}

func part2(in input) output {
	dirs := parseInput(in)
	_, pt := navigate(dirs)

	x := int(math.Abs(float64(pt.x)))
	y := int(math.Abs(float64(pt.y)))
	return x + y
}

func parseInput(in input) []Dir {
	dirs := []Dir{}
	for _, s := range in {
		tt := string(s[0])
		turn := -1
		if tt == "R" {
			turn = 1
		}
		moves, _ := strconv.Atoi(string(s[1:]))
		dirs = append(dirs, Dir{Turn: turn, Moves: moves})
	}
	return dirs
}

func getHeading(dir Dir, curheading string) string {
	heading := ""
	switch dir.Turn {
	case 1:
		switch curheading {
		case "N":
			heading = "E"
		case "E":
			heading = "S"
		case "S":
			heading = "W"
		case "W":
			heading = "N"
		}
		break
	case -1:
		switch curheading {
		case "N":
			heading = "W"
		case "W":
			heading = "S"
		case "S":
			heading = "E"
		case "E":
			heading = "N"
		}
		break
	}
	return heading
}

func navigate(dirs []Dir) (xy, xy) {
	x, y := 0, 0
	visits := make(map[string]int)
	visits["0,0"] = 0
	heading := "N"
	var twice *xy

	for _, dir := range dirs {
		heading = getHeading(dir, heading)
		var repeat *xy
		switch heading {
		case "N":
			repeat = track(y, dir.Moves, strconv.Itoa(x)+",%d", visits)
			y += dir.Moves
			break
		case "E":
			repeat = track(x, dir.Moves, "%d,"+strconv.Itoa(y), visits)
			x += dir.Moves
			break
		case "W":
			repeat = track(x, -dir.Moves, "%d,"+strconv.Itoa(y), visits)
			x -= dir.Moves
			break
		case "S":
			repeat = track(y, -dir.Moves, strconv.Itoa(x)+",%d", visits)
			y -= dir.Moves
		}

		if twice == nil && repeat != nil {
			twice = repeat
		}
	}

	return xy{x, y}, *twice
}

func track(changeval, moves int, format string, visits map[string]int) *xy {
	sign := 1
	if moves < 0 {
		sign = -1
	}

	var twice *xy
	for i := 0; i < int(math.Abs(float64(moves))); i++ {
		c := fmt.Sprintf(format, changeval+i*sign)
		visits[c]++

		if twice == nil && visits[c] == 2 {
			pts := strings.Split(c, ",")
			x, _ := strconv.Atoi(pts[0])
			y, _ := strconv.Atoi(pts[1])
			twice = &xy{x, y}
		}
	}
	return twice
}
