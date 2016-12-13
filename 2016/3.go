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

var input = "3.txt"

func main() {
	startTime := time.Now()
	if f, err := os.Open(input); err == nil {
		scanner := bufio.NewScanner(f)
		triangles := 0

		allsides := [][]int{[]int{}, []int{}, []int{}}
		reg := regexp.MustCompile(" *([0-9]+) *([0-9]+) *([0-9]+)")

		for scanner.Scan() {
			var txt = scanner.Text()
			if groups := reg.FindStringSubmatch(txt); groups != nil && len(groups) > 1 {
				sides := []int{}
				col := 0
				for _, n := range groups[1:] {
					if num, err := strconv.Atoi(n); err == nil {
						sides = append(sides, num)

						allsides[col] = append(allsides[col], num)
						col = (col + 1) % 3
					}
				}

				if isTriangle(sides...) {
					triangles++
				}
			}
		}
		verttriangles := 0
		for _, column := range allsides {
			for i := 0; i < len(column); i += 3 {
				if isTriangle(column[i:(i + 3)]...) {
					verttriangles++
				}
			}
		}

		fmt.Println("number of triangles", triangles, verttriangles)
	}

	fmt.Println("Time", time.Since(startTime))
}

func isTriangle(x ...int) bool {
	a, b, c := x[0], x[1], x[2]

	fmt.Println(a, b, c)

	return a+b > c && a+c > b && b+c > a
}

// reg := regexp.MustCompile("-?[0-9]+")
/*

if groups := reg.FindStringSubmatch(txt); groups != nil && len(groups) > 1 {
                fmt.Println(groups[1:])
            }
*/
