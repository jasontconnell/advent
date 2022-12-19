package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/jasontconnell/advent/common"
)

type input = []string
type output = int

type composite string

const (
	nothing  composite = "nothing"
	ore      composite = "ore"
	clay     composite = "clay"
	obsidian composite = "obsidian"
	geode    composite = "geode"
)

var priority []composite = []composite{
	geode,
	obsidian,
	clay,
	ore,
}

type blueprint struct {
	id   int
	bots map[composite]bot
}

type bot struct {
	costs map[composite]int
}

type state struct {
	building composite
	minute   int

	collected map[composite]int
	bots      map[composite]int
}

type cachekey struct {
	bc, bo, bg, bob int
	min             int
	building        composite
}

func (s state) String() string {
	return fmt.Sprintf("minute: %d building: %s  bots: %v collected: %v", s.minute, s.building, s.bots, s.collected)
}

type move struct {
	bot composite
}

func main() {
	startTime := time.Now()

	in, err := common.ReadStrings(common.InputFilename(os.Args))
	if err != nil {
		log.Fatal(err)
	}

	p1 := part1(in)
	p2 := part2(in)

	w := common.TeeOutput(os.Stdout)
	fmt.Fprintln(w, "--2022 day 19 solution--")
	fmt.Fprintln(w, "Part 1:", p1)
	fmt.Fprintln(w, "Part 2:", p2)
	fmt.Println("Time", time.Since(startTime))
}

func part1(in input) output {
	blueprints := parseInput(in)
	return getQualitySum(blueprints, 24)
}

func part2(in input) output {
	blueprints := parseInput(in)
	return getGeodeProduct(blueprints[:3], 32)
}

func getGeodeProduct(blueprints []blueprint, minutes int) int {
	prod := 1
	for _, bp := range blueprints {
		res := simulate(bp, minutes)
		prod *= res
	}
	return prod
}

func getQualitySum(blueprints []blueprint, minutes int) int {
	sum := 0
	for _, bp := range blueprints {
		res := simulate(bp, minutes)
		sum += (res * bp.id)
	}
	return sum
}

func simulate(bp blueprint, minutes int) int {
	queue := common.NewQueue[state, int]()
	queue.Enqueue(state{
		minute:    0,
		building:  nothing,
		bots:      map[composite]int{ore: 1, clay: 0, obsidian: 0, geode: 0},
		collected: map[composite]int{ore: 0, clay: 0, obsidian: 0, geode: 0}})

	visited := make(map[cachekey]bool)
	maxc := 0

	for queue.Any() {
		cur := queue.Dequeue()
		if cur.minute > minutes {
			continue
		}

		if cur.collected[geode] > maxc {
			maxc = cur.collected[geode]
		}

		key := cachekey{
			bc:       cur.bots[clay],
			bo:       cur.bots[ore],
			bg:       cur.bots[geode],
			bob:      cur.bots[obsidian],
			min:      cur.minute,
			building: cur.building,
		}
		if _, ok := visited[key]; ok {
			continue
		}
		visited[key] = true

		cur = collect(cur)
		if cur.building != nothing {
			cur.bots[cur.building] = cur.bots[cur.building] + 1
			cur.building = nothing
		}

		mvs := getMoves(cur, bp)
		if len(mvs) == 0 {
			cur.minute++
			queue.Enqueue(cur)
			continue
		}

		for _, mv := range mvs {
			st := state{building: cur.building, minute: cur.minute + 1, collected: copyMap(cur.collected), bots: copyMap(cur.bots)}
			if mv.bot != nothing {
				st.building = mv.bot
				for c, cost := range bp.bots[st.building].costs {
					st.collected[c] = st.collected[c] - cost
				}
			}
			queue.Enqueue(st)
		}
	}
	return maxc
}

func collect(st state) state {
	for _, c := range priority {
		st.collected[c] = st.collected[c] + st.bots[c]
	}
	return st
}

func getMoves(st state, bp blueprint) []move {
	mvs := []move{}

	for _, c := range priority {
		bot := bp.bots[c]
		build := true
		for needed, cost := range bot.costs {
			if st.collected[needed] < cost {
				build = false
			}
		}
		if build {
			mvs = append(mvs, move{bot: c})
		}
	}
	// add a wait move
	mvs = append(mvs, move{bot: nothing})

	return mvs
}

func copyMap[K comparable, V any](m map[K]V) map[K]V {
	nm := make(map[K]V)
	for k, v := range m {
		nm[k] = v
	}
	return nm
}

func parseInput(in input) []blueprint {
	reg := regexp.MustCompile("Blueprint ([0-9]+): Each ore robot costs ([0-9]+) ore. Each clay robot costs ([0-9]+) ore. Each obsidian robot costs ([0-9]+) ore and ([0-9]+) clay. Each geode robot costs ([0-9]+) ore and ([0-9]+) obsidian.")

	list := []blueprint{}
	for _, line := range in {
		m := reg.FindStringSubmatch(line)

		if len(m) == 8 {
			id, _ := strconv.Atoi(m[1])
			bp := blueprint{id: id}

			bp.bots = make(map[composite]bot)

			orebotcost, _ := strconv.Atoi(m[2])
			bp.bots[ore] = bot{costs: map[composite]int{ore: orebotcost}}

			claybotcost, _ := strconv.Atoi(m[3])
			bp.bots[clay] = bot{costs: map[composite]int{ore: claybotcost}}

			obsorecost, _ := strconv.Atoi(m[4])
			obsclaycost, _ := strconv.Atoi(m[5])
			bp.bots[obsidian] = bot{costs: map[composite]int{ore: obsorecost, clay: obsclaycost}}

			geodeorecost, _ := strconv.Atoi(m[6])
			geodeobscost, _ := strconv.Atoi(m[7])
			bp.bots[geode] = bot{costs: map[composite]int{ore: geodeorecost, obsidian: geodeobscost}}

			list = append(list, bp)
		}
	}
	return list
}
