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
		santa := Point{x: 0, y: 0}
		robo := Point{x: 0, y: 0}

		houses := make(map[Point]int)

		for i, ch := range s {
			dir := string(ch)

			var current *Point

			if i%2 == 0 {
				current = &santa
			} else {
				current = &robo
			}

			switch dir {
			case "<":
				current.x--
				break
			case ">":
				current.x++
				break
			case "v":
				current.y--
				break
			case "^":
				current.y++
				break
			}

			if _, exists := houses[*current]; exists {
				houses[*current]++
			} else {
				houses[*current] = 1
			}
		}

		fmt.Println(len(houses))
	}
}
