package main

import (
	"bufio"
	"fmt"
	"os"
	"time"
	//"regexp"
	"math"
	"strconv"
	"strings"
)

var input = "1.txt"

type Dir struct {
	Turn  int
	Moves int
}

func main() {
	startTime := time.Now()
	if f, err := os.Open(input); err == nil {
		scanner := bufio.NewScanner(f)

		for scanner.Scan() {
			var txt = scanner.Text()
			dirs := getDirections(txt)
			x, y := navigate(dirs)
			fmt.Println("coords", x, y, math.Abs(float64(x))+math.Abs(float64(y)))
		}
	}

	fmt.Println("Time", time.Since(startTime))
}

func navigate(dirs []Dir) (x, y int) {
	x, y = 0, 0
	heading := "N"
	visits := make(map[string]int)
	visits["0,0"] = 1

	for _, dir := range dirs {
		switch dir.Turn {
		case 1:
			if heading == "N" {
				heading = "E"
			} else if heading == "E" {
				heading = "S"
			} else if heading == "S" {
				heading = "W"
			} else if heading == "W" {
				heading = "N"
			}
			break
		case -1:
			if heading == "N" {
				heading = "W"
			} else if heading == "W" {
				heading = "S"
			} else if heading == "S" {
				heading = "E"
			} else if heading == "E" {
				heading = "N"
			}
			break
		}

		switch heading {
		case "N":
			track(y, dir.Moves, strconv.Itoa(x)+",%v", visits)
			y += dir.Moves
			break
		case "E":
			track(x, dir.Moves, "%v,"+strconv.Itoa(y), visits)
			x += dir.Moves
			break
		case "W":
			track(x, -dir.Moves, "%v,"+strconv.Itoa(y), visits)
			x -= dir.Moves
			break
		case "S":
			track(y, -dir.Moves, strconv.Itoa(x)+",%v", visits)
			y -= dir.Moves
		}

	}

	return
}

func track(changeval, moves int, format string, visits map[string]int) {
	sign := 1
	if moves < 0 {
		sign = -1
	}

	for i := 0; i < int(math.Abs(float64(moves))); i++ {
		c := fmt.Sprintf(format, changeval+i*sign)
		if _, ok := visits[c]; !ok {
			visits[c] = 1
		} else {
			fmt.Println(c)
		}
	}
}

func getDirections(txt string) []Dir {
	dirs := []Dir{}
	t := strings.Split(strings.Replace(txt, " ", "", -1), ",")
	for _, s := range t {
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

// reg := regexp.MustCompile("-?[0-9]+")
/*
if groups := reg.FindStringSubmatch(txt); groups != nil && len(groups) > 1 {
                fmt.Println(groups[1:])
            }
*/
