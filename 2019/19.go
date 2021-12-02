package main

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"strconv"
	"strings"

	"github.com/jasontconnell/advent/2019/intcode"
)

type xy struct {
	x, y int
}

var input = "19.txt"

func main() {
	startTime := time.Now()

	f, err := os.Open(input)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	scanner := bufio.NewScanner(f)

	opcodes := []int{}
	if scanner.Scan() {
		var txt = scanner.Text()
		sopcodes := strings.Split(txt, ",")
		for _, s := range sopcodes {
			i, err := strconv.Atoi(s)
			if err != nil {
				fmt.Println(err)
				continue
			}

			opcodes = append(opcodes, i)
		}
	}

	p1 := part1(opcodes, 50)
	p2 := part2(opcodes, 100)

	fmt.Println("Part 1:", p1)
	fmt.Println("Part 2:", p2)

	fmt.Println("Time", time.Since(startTime))
}

func part1(ops []int, size int) int {
	count := 0

	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			c := intcode.NewComputer(ops)
			c.AddInput(x, y)
			c.Exec()

			if c.Outs[0] == 1 {
				count++
			}
		}
	}

	return count
}

func part2(ops []int, size int) int {
	pos := findSquare(ops, size)

	return pos.x*10000 + pos.y
}

func findSquare(ops []int, size int) xy {
	threshold := size * 2
	x, y := 0, 0
	found := false
	var pos xy

	lastx := 0

	for !found {
		at := atPos(ops, x, y)

		if at == 0 {
			x++
			if x > lastx+threshold {
				y++
				x = lastx
				if x < 0 {
					x = 0
				}
			}
		} else {
			lastx = x
			sq, sp := checkSquare(ops, x, y, size)
			if !sq {
				y++
				x = lastx
				if x < 0 {
					x = 0
				}
			} else {
				found = true
				pos = sp
			}
		}
	}
	return pos
}

func checkSquare(ops []int, x, y, size int) (bool, xy) {
	c := size - 1
	bottom := xy{x, y + c}

	shift := 0
	for i := bottom.x; i < bottom.x+c; i++ {
		p := atPos(ops, i, bottom.y)
		if p == 1 {
			shift = i - bottom.x
			break
		}
	}

	corners := []xy{
		xy{x + shift, y},
		xy{x + shift + c, y},
		xy{x + shift, y + c},
		xy{x + shift + c, y + c},
	}

	val := true
	for _, corner := range corners {
		p := atPos(ops, corner.x, corner.y)
		if p == 0 {
			val = false
			break
		}
	}

	return val, corners[0]
}

func atPos(ops []int, x, y int) int {
	c := intcode.NewComputer(ops)
	c.AddInput(x, y)
	c.Exec()
	return c.Outs[0]
}
