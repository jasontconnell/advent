package main

import (
	"fmt"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
)

var input = "2.txt"

func main() {

	if s, err := os.ReadFile(input); err == nil {
		totalArea := 0
		totalRibbon := 0

		lines := strings.Split(string(s), "\r\n")

		for _, line := range lines {
			sides := strings.Split(line, "x")
			l, _ := strconv.Atoi(sides[0])
			w, _ := strconv.Atoi(sides[1])
			h, _ := strconv.Atoi(sides[2])

			a1 := l * w
			a2 := w * h
			a3 := h * l

			list := []int{l, w, h}
			sort.Ints(list)

			ribbon := 2*list[0] + 2*list[1]
			totalRibbon += ribbon + (l * w * h)

			area := 2*a1 + 2*a2 + 2*a3
			smallest := int(math.Min(float64(a1), math.Min(float64(a2), float64(a3))))
			area += smallest

			totalArea += area
		}

		fmt.Println("Total area is ", totalArea)
		fmt.Println("Total ribbon length is ", totalRibbon)
	}
}
