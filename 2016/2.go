package main

import (
	"bufio"
	"fmt"
	"os"
	"time"
	//"regexp"
	//"strconv"
	//"strings"
	//"math"
)

var input = "2.txt"
var numpad [][]int = [][]int{[]int{1, 2, 3}, []int{4, 5, 6}, []int{7, 8, 9}}
var numpad2 [][]int = [][]int{[]int{0, 0, 1, 0, 0}, []int{0, 2, 3, 4, 0}, []int{5, 6, 7, 8, 9}, []int{0, 10, 11, 12, 0}, []int{0, 0, 13, 0, 0}}

func main() {
	startTime := time.Now()
	if f, err := os.Open(input); err == nil {
		scanner := bufio.NewScanner(f)

		x, y := 1, 1

		x2, y2 := 0, 2

		chars := make(map[int]string)
		chars[1] = "1"
		chars[2] = "2"
		chars[3] = "3"
		chars[4] = "4"
		chars[5] = "5"
		chars[6] = "6"
		chars[7] = "7"
		chars[8] = "8"
		chars[9] = "9"
		chars[10] = "A"
		chars[11] = "B"
		chars[12] = "C"
		chars[13] = "D"

		for scanner.Scan() {
			var txt = scanner.Text()
			x, y = processLine(txt, x, y)

			x2, y2 = processLineV2(txt, x2, y2)

			fmt.Println("v1", x, y, numpad[y][x])
			fmt.Println("v2", x2, y2, chars[numpad2[y2][x2]])

		}
	}

	fmt.Println("Time", time.Since(startTime))
}

func processLine(line string, startx, starty int) (endx, endy int) {
	endx, endy = startx, starty
	for _, k := range line {
		sk := string(k)
		switch sk {
		case "U":
			if endy > 0 {
				endy--
			}
		case "D":
			if endy < 2 {
				endy++
			}
		case "L":
			if endx > 0 {
				endx--
			}
		case "R":
			if endx < 2 {
				endx++
			}
		}
	}

	return endx, endy
}

func processLineV2(line string, startx, starty int) (endx, endy int) {
	endx, endy = startx, starty
	for _, k := range line {
		sk := string(k)
		switch sk {
		case "U":
			if endy > 0 && numpad2[endy-1][endx] != 0 {
				endy--
			}
		case "D":
			if endy < 4 && numpad2[endy+1][endx] != 0 {
				endy++
			}
		case "L":
			if endx > 0 && numpad2[endy][endx-1] != 0 {
				endx--
			}
		case "R":
			if endx < 4 && numpad2[endy][endx+1] != 0 {
				endx++
			}
		}
	}

	return endx, endy
}

// reg := regexp.MustCompile("-?[0-9]+")
/*
if groups := reg.FindStringSubmatch(txt); groups != nil && len(groups) > 1 {
                fmt.Println(groups[1:])
            }
*/
