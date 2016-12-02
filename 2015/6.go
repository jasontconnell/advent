package main

import (
	"os"
	"fmt"
	"bufio"
	"regexp"
	"strconv"
)

var input = "6.txt"

type Point struct {
	x, y int
}

func main(){
	if f,err := os.Open(input); err == nil {
		scanner := bufio.NewScanner(f)
		

		pattern := "^(turn on|toggle|turn off) ([0-9]+),([0-9]+) through ([0-9]+),([0-9]+)$"
		reg,rerr := regexp.Compile(pattern)
		if rerr != nil {
			panic(rerr)
		}

		lights := [1000][1000]bool{}
		turn(&lights, Point{x:0,y:0},Point{x:999,y:999}, false)

		lines := 0

		for scanner.Scan() {
			var txt = scanner.Text()

			if groups := reg.FindStringSubmatch(txt); groups != nil && len(groups) > 1 {
				command := groups[1]

				coord1x,_ := strconv.Atoi(groups[2])
				coord1y,_ := strconv.Atoi(groups[3])
				coord2x,_ := strconv.Atoi(groups[4])
				coord2y,_ := strconv.Atoi(groups[5])

				start := Point{ x: coord1x, y: coord1y }
				end := Point{ x: coord2x, y: coord2y }

				switch command {
					case "turn on": turn(&lights, start, end, true)
					break
					case "turn off": turn(&lights, start, end, false)
					break
					case "toggle": toggle(&lights, start, end)
					break
				}
			}
			lines++
		}

		status(lights)
		fmt.Println(lines, "lines processed")
	}
}

func status(lights [1000][1000]bool){
	on, off := 0, 0
	for i := 0; i < 1000; i++ {
		for j := 0; j < 1000; j++ {
			if lights[i][j] { 
				on++ 
			} else { 
				off++ 
			}
		}
	}

	fmt.Println("Off:", off, "On:", on)
}

func turn(lights *[1000][1000]bool, start, end Point, value bool){
	for i := start.x; i <= end.x; i++ {
		for j := start.y; j <= end.y; j++ {
			lights[i][j] = value
		}
	}
}

func toggle(lights *[1000][1000]bool, start, end Point){
	for i := start.x; i <= end.x; i++ {
		for j := start.y; j <= end.y; j++ {
			lights[i][j] = !lights[i][j]
		}
	}
}