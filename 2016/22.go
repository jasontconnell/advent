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

var input = "22.txt"

type Node struct {
	X, Y                      int
	Size, Used, Avail, UsePct int
}

type Point struct {
	X, Y int
}

func main() {
	startTime := time.Now()

	reg := regexp.MustCompile("^/dev/grid/node-x([0-9]+)-y([0-9]+) *([0-9]+)T *([0-9]+)T *([0-9]+)T *([0-9]+)%$")

	// 39 x 25
	nodes := make([][]Node, 25)
	for i := 0; i < 25; i++ {
		nodes[i] = make([]Node, 39)
	}

	if f, err := os.Open(input); err == nil {
		scanner := bufio.NewScanner(f)

		for scanner.Scan() {
			var txt = scanner.Text()
			if groups := reg.FindStringSubmatch(txt); groups != nil && len(groups) > 1 {
				x, _ := strconv.Atoi(groups[1])
				y, _ := strconv.Atoi(groups[2])

				size, _ := strconv.Atoi(groups[3])
				used, _ := strconv.Atoi(groups[4])
				avail, _ := strconv.Atoi(groups[5])
				usepct, _ := strconv.Atoi(groups[6])
				nodes[y][x] = Node{X: x, Y: y, Size: size, Used: used, Avail: avail, UsePct: usepct}
			}
		}
	}

	cnt := 0
	for y := 0; y < 25; y++ {
		for x := 0; x < 39; x++ {
			if nodes[y][x].Used > 0 {
				cnt += avail(nodes, nodes[y][x].Used, x, y)
			}
		}
	}

	print(nodes)

	fmt.Println("count", cnt)
	fmt.Println("Time", time.Since(startTime))
}

func print(nodes [][]Node) {
	for y := 0; y < 25; y++ {
		for x := 0; x < 39; x++ {
			ch := '#'
			if nodes[y][x].Used > 0 && nodes[y][x].Used < 100 {
				ch = '.'
			} else if nodes[y][x].Used == 0 {
				ch = '_'
			}

			fmt.Print(string(ch))
		}
		fmt.Println("")
	}
}

func avail(nodes [][]Node, used, curx, cury int) int {
	cnt := 0
	for y := 0; y < 25; y++ {
		for x := 0; x < 39; x++ {
			if x == curx && y == cury {
				continue
			}
			// for _, pt := range []Point{ Point{ X: x, Y: y+1 }, Point{ X: x, Y: y-1 }, Point{ X: x+1, Y: y }, Point{ X: x-1, Y: y } } {
			//     if pt.X > -1 && pt.X < 39 && pt.Y > -1 && pt.Y < 25 {
			if used <= nodes[y][x].Avail {
				cnt++
			}
			//     }
			// }
		}
	}
	return cnt
}

// reg := regexp.MustCompile("-?[0-9]+")
/*
if groups := reg.FindStringSubmatch(txt); groups != nil && len(groups) > 1 {
                fmt.Println(groups[1:])
            }
*/
