package main

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

var input = "18.txt"

const (
	Open int = iota
	Trees
	Lumberyard
)

type block struct {
	contents int
}

type xy struct {
	x, y int
}

type rpt struct {
	start int
	end   int
	count int
}

func main() {
	startTime := time.Now()

	f, err := os.Open(input)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	scanner := bufio.NewScanner(f)

	var blocks [][]block
	for scanner.Scan() {
		var txt = scanner.Text()
		blocks = append(blocks, readLine(txt))
	}

	blocks2 := clone(blocks)

	blocks = sim(blocks, 10)
	part1 := sumAll(blocks, Trees) * sumAll(blocks, Lumberyard)
	blocks2 = simRepeats(blocks2, 1000000000)
	part2 := sumAll(blocks2, Trees) * sumAll(blocks2, Lumberyard)

	fmt.Println("Part 1:", part1)
	fmt.Println("Part 2:", part2)

	fmt.Println("Time", time.Since(startTime))
}

func simRepeats(blocks [][]block, minutes int) [][]block {
	blocks2 := clone(blocks)
	rcmap := make(map[int]rpt)
	minute := 0
	done := false
	for minute < 3000 && !done {
		blocks2 = simOne(blocks2)
		val := sumAll(blocks2, Trees) * sumAll(blocks2, Lumberyard)
		if _, ok := rcmap[val]; !ok {
			rcmap[val] = rpt{start: minute, count: 1}
		} else {
			r := rcmap[val]
			r.start = r.end
			r.end = minute
			r.count++
			rcmap[val] = r

			done = r.count > 60
		}
		minute++
	}

	max := 0
	rc := rpt{}
	for _, r := range rcmap {
		if r.count > max {
			max = r.count
			rc = r
		}
	}

	minute = rc.start
	diff := rc.end - rc.start
	for minute+diff < minutes {
		minute += diff
	}

	return sim(blocks2, minutes-minute-1)
}

func clone(blocks [][]block) [][]block {
	cln := [][]block{}
	for y := 0; y < len(blocks); y++ {
		cln = append(cln, blocks[y])
	}
	return cln
}

func print(blocks [][]block) {
	for y := 0; y < len(blocks); y++ {
		line := ""
		for x := 0; x < len(blocks[y]); x++ {
			str := "."
			switch blocks[y][x].contents {
			case Trees:
				str = "|"
			case Lumberyard:
				str = "#"
			}
			line += str
		}
		fmt.Println(line)
	}
	fmt.Println("---------------------------------")
}

func sim(blocks [][]block, minutes int) [][]block {
	for i := 0; i < minutes; i++ {
		blocks = simOne(blocks)
	}
	return blocks
}

func simOne(blocks [][]block) [][]block {
	updated := make([][]block, len(blocks))
	maxx := len(blocks[0]) - 1
	maxy := len(blocks) - 1

	for i, _ := range blocks {
		updated[i] = make([]block, len(blocks[0]))
	}

	for y := 0; y < len(blocks); y++ {
		for x := 0; x < len(blocks[y]); x++ {
			sur := surrounding(x, y, maxx, maxy)
			c := blocks[y][x].contents
			switch c {
			case Open:
				s := sumContents(blocks, sur, Trees)
				if s >= 3 {
					c = Trees
				}
			case Trees:
				s := sumContents(blocks, sur, Lumberyard)
				if s >= 3 {
					c = Lumberyard
				}
			case Lumberyard:
				s1 := sumContents(blocks, sur, Lumberyard)
				s2 := sumContents(blocks, sur, Trees)

				if s1 == 0 || s2 == 0 {
					c = Open
				}
			}

			updated[y][x].contents = c
		}
	}
	return updated
}

func sumAll(blocks [][]block, val int) int {
	count := 0
	for y := 0; y < len(blocks); y++ {
		for x := 0; x < len(blocks[y]); x++ {
			if blocks[y][x].contents == val {
				count++
			}
		}
	}
	return count
}

func sumContents(blocks [][]block, pts []xy, val int) int {
	count := 0
	for _, pt := range pts {
		if blocks[pt.y][pt.x].contents == val {
			count++
		}
	}
	return count
}

func surrounding(x0, y0, maxx, maxy int) []xy {
	ps := []xy{}
	for y := -1; y <= 1; y++ {
		if y0+y < 0 || y0+y > maxy {
			continue
		}
		for x := -1; x <= 1; x++ {
			if x0+x < 0 || x0+x > maxx {
				continue
			}

			if x == 0 && y == 0 {
				continue
			}
			ps = append(ps, xy{x: x0 + x, y: y0 + y})
		}
	}
	return ps
}

func readLine(line string) []block {
	blocks := []block{}
	for _, c := range line {
		b := block{contents: Open}
		switch c {
		case '|':
			b.contents = Trees
		case '#':
			b.contents = Lumberyard
		}
		blocks = append(blocks, b)
	}
	return blocks
}
