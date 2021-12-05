package main

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"time"

	"github.com/jasontconnell/advent/common"
)

var inputFilename = "input.txt"

type input []string

var reg *regexp.Regexp = regexp.MustCompile("^(turn on|toggle|turn off) ([0-9]+),([0-9]+) through ([0-9]+),([0-9]+)$")

type command struct {
	action string
	start  Point
	end    Point
}

type Point struct {
	x, y int
}

func main() {
	startTime := time.Now()

	lines, err := common.ReadStrings(inputFilename)
	if err != nil {
		log.Fatal(err)
	}

	p1 := part1(lines)
	p2 := part2(lines)

	fmt.Println("Part 1:", p1)
	fmt.Println("Part 2:", p2)

	fmt.Println("Time", time.Since(startTime))
}

func part1(in input) int {
	lights := [1000][1000]int{}
	cmds := getCommands(in)
	turn(&lights, Point{x: 0, y: 0}, Point{x: 999, y: 999}, 0)
	navigate(cmds, &lights)
	on, _ := status(lights)
	return on
}

func part2(in input) int {
	lights := [1000][1000]int{}
	cmds := getCommands(in)
	turn(&lights, Point{x: 0, y: 0}, Point{x: 999, y: 999}, 0)
	navigateIlluminate(cmds, &lights)
	on, _ := status(lights)
	return on
}

func getCommands(in input) []command {
	cmds := []command{}
	for _, txt := range in {
		if groups := reg.FindStringSubmatch(txt); groups != nil && len(groups) > 1 {
			action := groups[1]

			coord1x, _ := strconv.Atoi(groups[2])
			coord1y, _ := strconv.Atoi(groups[3])
			coord2x, _ := strconv.Atoi(groups[4])
			coord2y, _ := strconv.Atoi(groups[5])

			start := Point{x: coord1x, y: coord1y}
			end := Point{x: coord2x, y: coord2y}

			c := command{action: action, start: start, end: end}
			cmds = append(cmds, c)
		}
	}
	return cmds
}

func navigateIlluminate(cmds []command, lights *[1000][1000]int) {
	for _, cmd := range cmds {
		switch cmd.action {
		case "turn on":
			illuminate(lights, cmd.start, cmd.end, 1)
			break
		case "turn off":
			illuminate(lights, cmd.start, cmd.end, -1)
			break
		case "toggle":
			illuminate(lights, cmd.start, cmd.end, 2)
			break
		}
	}
}

func illuminate(lights *[1000][1000]int, start, end Point, val int) {
	for j := start.y; j <= end.y; j++ {
		for i := start.x; i <= end.x; i++ {
			lights[j][i] += val

			if lights[j][i] < 0 {
				lights[j][i] = 0
			}
		}
	}
}

func navigate(cmds []command, lights *[1000][1000]int) {
	for _, cmd := range cmds {
		switch cmd.action {
		case "turn on":
			turn(lights, cmd.start, cmd.end, 1)
			break
		case "turn off":
			turn(lights, cmd.start, cmd.end, 0)
			break
		case "toggle":
			toggle(lights, cmd.start, cmd.end)
			break
		}
	}
}

func status(lights [1000][1000]int) (on, off int) {
	for j := 0; j < 1000; j++ {
		for i := 0; i < 1000; i++ {
			if lights[j][i] == 0 {
				off++
			} else {
				on += lights[j][i]
			}
		}
	}

	return on, off
}

func turn(lights *[1000][1000]int, start, end Point, value int) {
	for j := start.y; j <= end.y; j++ {
		for i := start.x; i <= end.x; i++ {
			lights[j][i] = value
		}
	}
}

func toggle(lights *[1000][1000]int, start, end Point) {
	for j := start.y; j <= end.y; j++ {
		for i := start.x; i <= end.x; i++ {
			x := lights[j][i]

			if x == 0 {
				x = 1
			} else {
				x = 0
			}

			lights[j][i] = x
		}
	}
}
