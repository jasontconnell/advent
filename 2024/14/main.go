package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"

	"github.com/jasontconnell/advent/common"
)

type input = []string
type output = int

type bot struct {
	vx, vy int
}
type xy struct {
	x, y int
}

func (p xy) add(p2 xy) xy {
	return xy{p.x + p2.x, p.y + p2.y}
}

func main() {
	in, err := common.ReadStrings(common.InputFilename(os.Args))
	if err != nil {
		log.Fatal(err)
	}

	p1, p1time := common.Time(part1, in)
	p2, p2time := common.Time(part2, in)

	w := common.TeeOutput(os.Stdout)
	fmt.Fprintln(w, "--2024 day 14 solution--")
	fmt.Fprintln(w, "Part 1:", p1)
	fmt.Fprintln(w, "Part 2:", p2)
	fmt.Printf("Time %v (%v, %v)", p1time+p2time, p1time, p2time)
}

func part1(in input) output {
	bots := parse(in)
	log.Println(len(bots))
	return 0
}

func part2(in input) output {
	return 0
}

func simulate(bots map[xy]bot, w, h int) {

}

func countQuadrants(bots map[xy]bot, w, h int) int {

}

func parse(in []string) map[xy]bot {
	reg := regexp.MustCompile(`^p=([0-9]+),([0-9]+) v=([0-9]+),([0-9]+)$`)

	bm := make(map[xy]bot)
	for _, line := range in {
		m := reg.FindStringSubmatch(line)
		if len(m) != 5 {
			continue
		}
		px, _ := strconv.Atoi(m[1])
		py, _ := strconv.Atoi(m[2])

		vx, _ := strconv.Atoi(m[3])
		vy, _ := strconv.Atoi(m[4])

		b := bot{vx, vy}
		bm[xy{px, py}] = b
	}
	return bm
}
