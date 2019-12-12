package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"time"
)

var input = "08.txt"

type image struct {
	w, h   int
	layers []layer
}

type layer struct {
	rows [][]int
}

func main() {
	startTime := time.Now()

	f, err := os.Open(input)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	scanner := bufio.NewScanner(f)
	var pixels []int

	for scanner.Scan() {
		var txt = scanner.Text()
		pixels = parse(txt)
	}

	img := getImage(pixels, 25, 6)
	p1 := part1(img)

	fmt.Println("Part 1: ", p1)
	fmt.Println("Time", time.Since(startTime))
}

func parse(str string) []int {
	pixels := []int{}
	for _, c := range str {
		i, err := strconv.Atoi(string(c))
		if err != nil {
			panic(err)
		}
		pixels = append(pixels, i)
	}
	return pixels
}

func part1(img image) int {
	fewest := -1
	fidx := -1

	for lidx, l := range img.layers {
		n := 0
		for _, r := range l.rows {
			for _, p := range r {
				if p == 0 {
					n++
				}
			}
		}

		if fewest == -1 || n < fewest {
			fewest = n
			fidx = lidx
		}
	}

	ones, twos := 0, 0

	for _, r := range img.layers[fidx].rows {
		for _, p := range r {
			if p == 1 {
				ones++
			}
			if p == 2 {
				twos++
			}
		}
	}

	return ones * twos
}

func getImage(pixels []int, w, h int) image {
	img := image{w: w, h: h}
	max := len(pixels) / (w * h)
	img.layers = make([]layer, max)

	for i := range img.layers {
		img.layers[i].rows = make([][]int, h)
		for j := range img.layers[i].rows {
			img.layers[i].rows[j] = make([]int, w)
		}
	}

	var lidx, cidx, ridx int
	for _, p := range pixels {
		img.layers[lidx].rows[ridx][cidx] = p

		cidx = (cidx + 1) % w

		if cidx == 0 {
			ridx = (ridx + 1) % h
		}

		if ridx == 0 && cidx == 0 {
			lidx = (lidx + 1) % max
		}
	}

	return img
}
