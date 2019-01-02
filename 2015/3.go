package main

import (
	"fmt"
	"io/ioutil"
)

var input = "3.txt"

type Point struct {
	x, y int
}

func main() {
	if s, err := ioutil.ReadFile(input); err == nil {
		x := 0
		y := 0

		houses := make(map[Point]int)

		for _, ch := range s {
			dir := string(ch)

			switch dir {
			case "<":
				x--
				break
			case ">":
				x++
				break
			case "v":
				y--
				break
			case "^":
				y++
				break
			}

			p := Point{x: x, y: y}

			if _, exists := houses[p]; exists {
				houses[p]++
			} else {
				houses[p] = 1
			}
		}

		fmt.Println(len(houses))
	}
}
