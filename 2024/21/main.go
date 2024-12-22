package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"strconv"

	"github.com/jasontconnell/advent/common"
)

type input = []string
type output = int

type xy struct {
	x, y int
}

func (p xy) add(p2 xy) xy {
	return xy{p.x + p2.x, p.y + p2.y}
}

type dir struct {
	dir    xy
	bpress byte
}

type state struct {
	pt        xy
	stepIndex int
	path      string
	last      dir
}

type statekey struct {
	pt     xy
	bpress byte
}

var numpadkeys []string = []string{"789", "456", "123", " 0A"}
var dirpadkeys []string = []string{" ^A", "<v>"}

func main() {
	in, err := common.ReadStrings(common.InputFilename(os.Args))
	if err != nil {
		log.Fatal(err)
	}

	p1, p1time := common.Time(part1, in)
	p2, p2time := common.Time(part2, in)

	w := common.TeeOutput(os.Stdout)
	fmt.Fprintln(w, "--2024 day 21 solution--")
	fmt.Fprintln(w, "Part 1:", p1)
	fmt.Fprintln(w, "Part 2:", p2)
	fmt.Printf("Time %v (%v, %v)", p1time+p2time, p1time, p2time)
}

func part1(in input) output {
	numpad := getPad(numpadkeys)
	dpad := getPad(dirpadkeys)
	return solve(in, numpad, dpad)
}

func part2(in input) output {
	return 0
}

func solve(keycodes []string, numpad, dpad map[xy]byte) int {
	numpadc := getCoords(numpad)
	dpadc := getCoords(dpad)

	na := numpadc['A']
	da := dpadc['A']

	complexity := 0
	for _, kc := range keycodes {
		nseq := getSequence(kc, na, numpad, dpad, false, 2)
		log.Println("numeric pad", kc, nseq, len(nseq))
		d1seq := getSequence(nseq, da, numpad, dpad, true, 1)
		log.Println("first dpad", d1seq, len(d1seq))
		d2seq := getSequence(d1seq, da, numpad, dpad, true, 0)
		log.Println("second dpad", d2seq, len(d2seq))
		complexity += getComplexity(kc, d2seq)
		break
	}
	return complexity
}

func getComplexity(code, path string) int {
	s := ""
	for i := 0; i < len(code); i++ {
		if code[i] >= 48 && code[i] <= 57 {
			s += string(code[i])
		}
	}
	v, _ := strconv.Atoi(s)
	return v * len(path)
}

func getSequence(str string, start xy, numpad, dpad map[xy]byte, isdpad bool, level int) string {
	rev := getCoords(numpad)
	if isdpad {
		rev = getCoords(dpad)
	}

	var val string
	cur := start
	for i := 0; i < len(str); i++ {
		s := str[i]
		e := rev[s]
		val += getSubsequence(cur, e, numpad, dpad, isdpad, level) + "A"
		cur = e
	}
	return val
}

func getSubsequence(start, end xy, numpad, dpad map[xy]byte, isdpad bool, level int) string {
	if start == end {
		return ""
	}
	queue := common.NewQueue[state, int]()
	initial := state{pt: start, stepIndex: 0, path: ""}
	queue.Enqueue(initial)

	pad := numpad
	rev := getCoords(dpad)
	if isdpad {
		pad = dpad
	}

	best := math.MaxInt32
	visit := make(map[xy]bool)
	bests := make(map[int]string)

	for queue.Any() {
		cur := queue.Dequeue()

		if _, ok := visit[cur.pt]; ok {
			continue
		}
		visit[cur.pt] = true

		if cur.pt == end {
			log.Println(level, cur.path)
			if len(cur.path) < 3 || level == 0 {
				bests[len(cur.path)] = cur.path
				best = len(cur.path)
			} else if level > 0 {
				x := level

				bs := ""
				for x > 0 {
					s, e := rev[cur.path[0]], rev[cur.path[len(cur.path)-1]]
					bs += getSubsequence(s, e, numpad, dpad, true, x-1)
					x--
				}
				if _, ok := bests[len(bs)]; !ok {
					bests[len(bs)] = cur.path
					if len(bs) < best {
						best = len(bs)
					}
				}
			}
		}

		mvs := getMoves(pad, cur)
		for _, mv := range mvs {
			queue.Enqueue(mv)
		}
	}
	return bests[best]
}

func getMoves(pad map[xy]byte, cur state) []state {
	dirs := []dir{{xy{1, 0}, '>'}, {xy{-1, 0}, '<'}, {xy{0, 1}, 'v'}, {xy{0, -1}, '^'}}
	mvs := []state{}
	for _, d := range dirs {
		np := cur.pt.add(d.dir)
		if _, ok := pad[np]; ok {
			mv := state{pt: np, stepIndex: cur.stepIndex + 1, path: cur.path + string(d.bpress)}
			mvs = append(mvs, mv)
		}
	}
	return mvs
}

func getCoords(pad map[xy]byte) map[byte]xy {
	cmap := make(map[byte]xy)
	for k, v := range pad {
		cmap[v] = k
	}
	return cmap
}

func getPad(str []string) map[xy]byte {
	m := make(map[xy]byte)
	for y := 0; y < len(str); y++ {
		for x := 0; x < len(str[y]); x++ {
			if str[y][x] == ' ' {
				continue
			}
			pt := xy{x, y}
			m[pt] = byte(str[y][x])
		}
	}
	return m
}
