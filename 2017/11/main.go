package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/jasontconnell/advent/common"
)

type input = string
type output = int

type Coord struct {
	X, Y int
}

func (c *Coord) Add(d Coord) Coord {
	e := *c
	e.X += d.X
	e.Y += d.Y

	return e
}

type Hex struct {
	Coords Coord
	NE     *Hex
	N      *Hex
	NW     *Hex
	SE     *Hex
	S      *Hex
	SW     *Hex
}

func main() {
	startTime := time.Now()

	in, err := common.ReadString(common.InputFilename(os.Args))
	if err != nil {
		log.Fatal(err)
	}

	p1 := part1(in)
	p2 := part2(in)

	w := common.TeeOutput(os.Stdout)
	fmt.Fprintln(w, "--2017 day 11 solution--")
	fmt.Fprintln(w, "Part 1:", p1)
	fmt.Fprintln(w, "Part 2:", p2)
	fmt.Println("Time", time.Since(startTime))
}

func part1(in input) output {
	mvs := strings.Split(in, ",")
	home := makeHex(Coord{X: 0, Y: 0})
	end, _ := traverse(mvs, home)
	return traverseTo(end, home)
}

func part2(in input) output {
	mvs := strings.Split(in, ",")
	home := makeHex(Coord{X: 0, Y: 0})
	_, farthest := traverse(mvs, home)
	return farthest
}

func traverse(moves []string, start *Hex) (*Hex, int) {
	ptr := start
	max := 0
	for _, mv := range moves {
		ptr = getMove(mv, ptr)
		d := traverseTo(ptr, start)
		if d > max {
			max = d
		}
	}
	return ptr, max
}

func traverseTo(start, end *Hex) int {
	ptr := start
	d := 0

	for ptr.Coords.X != end.Coords.X || ptr.Coords.Y != end.Coords.Y {
		mv := ""
		n := ptr.Coords.X%2 != 0

		if ptr.Coords.Y > end.Coords.Y {
			mv = "s"
		} else if ptr.Coords.Y < end.Coords.Y {
			mv = "n"
		} else if ptr.Coords.Y == end.Coords.Y {
			mv = "s"
			if n {
				mv = "n"
			}
		}

		if ptr.Coords.X > end.Coords.X {
			mv += "w"
		} else if ptr.Coords.X < end.Coords.X {
			mv += "e"
		}

		ptr = getMove(mv, ptr)

		d++
	}

	return d
}

func makeHex(coords Coord) *Hex {
	hex := new(Hex)
	hex.Coords = coords
	return hex
}

func getMove(dir string, hex *Hex) *Hex {
	var dest *Hex
	var delta Coord
	n := hex.Coords.X%2 != 0

	switch dir {
	case "n":
		if hex.N == nil {
			delta.Y = 1
			hex.N = makeHex(hex.Coords.Add(delta))
			hex.N.S = hex
		}
		dest = hex.N
	case "s":
		if hex.S == nil {
			delta.Y = -1
			hex.S = makeHex(hex.Coords.Add(delta))
			hex.S.N = hex
		}
		dest = hex.S
	case "ne":
		if hex.NE == nil {
			delta.X = 1
			if !n {
				delta.Y = 1
			}
			hex.NE = makeHex(hex.Coords.Add(delta))
			hex.NE.SW = hex
		}

		dest = hex.NE
	case "nw":
		if hex.NW == nil {
			delta.X = -1
			if !n {
				delta.Y = 1
			}

			hex.NW = makeHex(hex.Coords.Add(delta))
			hex.NW.SE = hex
		}

		dest = hex.NW
	case "se":
		if hex.SE == nil {
			delta.X = 1
			if n {
				delta.Y = -1
			}

			hex.SE = makeHex(hex.Coords.Add(delta))
			hex.SE.NW = hex
		}
		dest = hex.SE

	case "sw":
		if hex.SW == nil {
			delta.X = -1
			if n {
				delta.Y = -1
			}
			hex.SW = makeHex(hex.Coords.Add(delta))
			hex.SW.NE = hex
		}
		dest = hex.SW
	}
	return dest
}
