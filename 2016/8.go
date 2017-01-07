package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"time"
	//"strings"
	//"math"
)

var input = "8.txt"

func main() {
	startTime := time.Now()
	display := [][]bool{}

	w, h := 50, 6

	for i := 0; i < h; i++ {
		display = append(display, []bool{})
		for j := 0; j < w; j++ {
			display[i] = append(display[i], false)
		}
	}

	rectreg := regexp.MustCompile("^rect (\\d+)x(\\d+)$")
	rotreg := regexp.MustCompile("^rotate (row|column) (x|y)=(\\d+) by (\\d+)$")

	if f, err := os.Open(input); err == nil {
		scanner := bufio.NewScanner(f)

		for scanner.Scan() {
			var txt = scanner.Text()

			if groups := rectreg.FindStringSubmatch(txt); groups != nil {
				x, _ := strconv.Atoi(groups[1])
				y, _ := strconv.Atoi(groups[2])
				rect(display, x, y)
			}

			if groups := rotreg.FindStringSubmatch(txt); groups != nil {
				ordinal, _ := strconv.Atoi(groups[3])
				count, _ := strconv.Atoi(groups[4])

				switch groups[1] {
				case "column":
					rotateColumn(display, ordinal, count)
					break
				case "row":
					rotateRow(display, ordinal, count)
					break
				}
			}
		}
	}

	oncount := 0
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			if display[y][x] {
				oncount++
			}
		}
	}

	print(display, w, h)

	fmt.Println("on count", oncount)
	fmt.Println("Time", time.Since(startTime))
}

func print(display [][]bool, w, h int) {
	for y := 0; y < h; y++ {

		for x := 0; x < w; x++ {
			ch := " "
			if display[y][x] {
				ch = "#"
			}

			fmt.Print(ch)
		}
		fmt.Println("")
	}
}

func shift(slice []bool, count int) []bool {
	cp := make([]bool, len(slice))
	copy(cp, slice)

	for i := 0; i < len(slice); i++ {
		dstIndex := (i + count) % len(slice)
		cp[dstIndex] = slice[i]
	}
	return cp
}

func rotateColumn(display [][]bool, col, count int) {
	shifted := []bool{}

	for i := 0; i < len(display); i++ {
		shifted = append(shifted, display[i][col])
	}

	shifted = shift(shifted, count)

	for i := 0; i < len(display); i++ {
		display[i][col] = shifted[i]
	}
}

func rotateRow(display [][]bool, row, count int) {
	shifted := []bool{}

	for i := 0; i < len(display[row]); i++ {
		shifted = append(shifted, display[row][i])
	}

	shifted = shift(shifted, count)

	for i := 0; i < len(display[row]); i++ {
		display[row][i] = shifted[i]
	}
}

func rect(display [][]bool, w, h int) {
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			display[y][x] = true
		}
	}
}
