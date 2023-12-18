package main

import (
	"fmt"
	"image"
	"os"
	"regexp"
	"strconv"
)

func main() {
	input, _ := os.ReadFile("input.txt")
	re := regexp.MustCompile(`(.) (.*) \(#(.*)(.)\)`)
	delta := map[string]image.Point{
		"U": {0, -1}, "R": {1, 0}, "D": {0, 1}, "L": {-1, 0},
		"0": {1, 0}, "1": {0, 1}, "2": {-1, 0}, "3": {0, -1},
	}

	run := func(dpart, lpart, base int) int {
		dig, area := image.Point{0, 0}, 0
		for _, m := range re.FindAllStringSubmatch(string(input), -1) {
			l, _ := strconv.ParseInt(m[lpart], base, strconv.IntSize)
			n := dig.Add(delta[m[dpart]].Mul(int(l)))
			area += dig.X*n.Y - dig.Y*n.X + int(l)
			dig = n
		}
		return area/2 + 1
	}

	fmt.Println(run(1, 2, 10))
	fmt.Println(run(4, 3, 16))
}
