package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"

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

func (p xy) dist(p2 xy) int {
	dx := p.x - p2.x
	dy := p.y - p2.y
	return int(math.Abs(float64(dx)) + math.Abs(float64(dy)))
}

type pair struct {
	a, b xy
}

type dir struct {
	dir    xy
	bpress byte
}

type state struct {
	pt        xy
	stepIndex int
	path      string
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
	npmoves := mapAllMoves(numpad)
	dpmoves := mapAllMoves(dpad)
	return solve(in, numpad, dpad, npmoves, dpmoves, 2)
}

func part2(in input) output {
	return 0
}

func solve(keycodes []string, numpad, dpad map[xy]byte, npadmoves, dpadmoves map[pair][]string, ndpads int) int {
	numpadc := getCoords(numpad)
	dpadc := getCoords(dpad)

	na := numpadc['A']
	da := dpadc['A']

	complexities := make(map[string]int)

	total := 0
	for _, kc := range keycodes {
		seqs := getSequences(kc, na, numpad, npadmoves)
		for i := 0; i < ndpads; i++ {
			var sseqs []string
			for _, seq := range seqs {
				dseqs := getSequences(seq, da, dpad, dpadmoves)
				sseqs = append(sseqs, dseqs...)
			}
			seqs = sseqs
		}
		for _, fseq := range seqs {
			check := getComplexity(kc, fseq)
			if c, ok := complexities[kc]; !ok || check < c {
				complexities[kc] = check
			}
		}
		total += complexities[kc]
	}
	return total
}

func getComplexity(code, path string) int {
	v, _ := strconv.Atoi(code[:len(code)-1])
	return v * len(path)
}

func getSequences(str string, start xy, pad map[xy]byte, pmoves map[pair][]string) []string {
	rev := getCoords(pad)

	cur := start
	list := []string{}
	for i := 0; i < len(str); i++ {
		e := rev[str[i]]

		rev := common.CartesianProduct(list, pmoves[pair{cur, e}])
		sub := make([]string, len(rev))
		if len(rev) == 0 {
			sub = pmoves[pair{cur, e}]
		}

		for j := 0; j < len(rev); j++ {
			s := strings.Join(rev[j], "")
			sub[j] = s
		}
		list = sub
		cur = e
	}

	minlen := math.MaxInt32
	for i := range list {
		if len(list[i]) < minlen {
			minlen = len(list[i])
		}
	}

	for i := len(list) - 1; i >= 0; i-- {
		if len(list[i]) > minlen {
			list = append(list[:i], list[i+1:]...)
		}
	}

	return list
}

func mapAllMoves(pad map[xy]byte) map[pair][]string {
	dirs := []dir{
		{xy{0, -1}, '^'},
		{xy{1, 0}, '>'},
		{xy{0, 1}, 'v'},
		{xy{-1, 0}, '<'},
	}

	mps := make(map[pair][]string)
	for k := range pad {
		for k2 := range pad {
			list := []string{}
			queue := []state{{pt: k, stepIndex: 0, path: ""}}
			for len(queue) > 0 {
				cur := queue[0]
				queue = queue[1:]

				if _, ok := pad[cur.pt]; !ok {
					continue
				}

				if len(cur.path) > len(pad) {
					continue
				}

				if cur.pt == k2 {
					list = append(list, cur.path+"A")
					continue
				}

				for _, d := range dirs {
					nst := state{pt: cur.pt.add(d.dir), stepIndex: cur.stepIndex + 1, path: cur.path + string(d.bpress)}
					if nst.pt.dist(k2) < cur.pt.dist(k2) {
						queue = append(queue, nst)
					}
				}
			}

			mps[pair{k, k2}] = list
		}
	}
	return mps
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
