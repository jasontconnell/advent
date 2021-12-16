package main

import (
	"crypto/md5"
	"fmt"
	"log"
	"os"
	"sort"
	"time"

	"github.com/jasontconnell/advent/common"
)

type input = string
type output = string

type xy struct {
	x, y int
}

type state struct {
	pt   xy
	path string
}

func main() {
	startTime := time.Now()

	in, err := common.ReadString(common.InputFilename(os.Args))
	if err != nil {
		log.Fatal(err)
	}

	p1 := part1(in)
	p2 := part2(in)

	w := common.TeeOutput(os.Stdout)
	fmt.Fprintln(w, "--2016 day 17 solution--")
	fmt.Fprintln(w, "Part 1:", p1)
	fmt.Fprintln(w, "Part 2:", p2)
	fmt.Println("Time", time.Since(startTime))
}

func part1(in input) output {
	shortest, _ := solve(xy{0, 0}, xy{3, 3}, in)
	return shortest
}

func part2(in input) int {
	_, longest := solve(xy{0, 0}, xy{3, 3}, in)
	return len(longest)
}

func solve(start, goal xy, hash string) (string, string) {
	queue := []state{}
	queue = append(queue, state{pt: start, path: ""})
	solvePaths := []string{}

	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]

		if cur.pt == goal {
			solvePaths = append(solvePaths, cur.path)
			continue
		}

		moves := getDoors(cur.pt, hash, cur.path)

		for _, mv := range moves {
			point := cur.pt
			if mv == 'U' {
				point.y--
			} else if mv == 'D' {
				point.y++
			} else if mv == 'L' {
				point.x--
			} else if mv == 'R' {
				point.x++
			}
			path := cur.path + string(mv)

			queue = append(queue, state{pt: point, path: path})
		}
	}

	sort.Slice(solvePaths, func(i, j int) bool {
		return len(solvePaths[i]) < len(solvePaths[j])
	})

	return solvePaths[0], solvePaths[len(solvePaths)-1]
}

func getDoors(point xy, hash, path string) (moves []rune) {
	md5 := MD5s(hash + path)
	if isOpen(rune(md5[0])) && point.y > 0 {
		moves = append(moves, 'U')
	}
	if isOpen(rune(md5[1])) && point.y < 3 {
		moves = append(moves, 'D')
	}
	if isOpen(rune(md5[2])) && point.x > 0 {
		moves = append(moves, 'L')
	}
	if isOpen(rune(md5[3])) && point.x < 3 {
		moves = append(moves, 'R')
	}

	return moves
}

func isOpen(char rune) bool {
	return char == 'b' || char == 'c' || char == 'd' || char == 'e' || char == 'f'
}

func MD5(content []byte) string {
	sum := md5.Sum(content)
	return fmt.Sprintf("%x", sum)
}

func MD5s(content string) string {
	return MD5([]byte(content))
}
