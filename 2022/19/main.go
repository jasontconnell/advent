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

func (s state) String() string {
	return fmt.Sprintf("minute: %d building: %s  bots: %v collected: %v", s.minute, s.building, s.bots, s.collected)
}

type cachekey struct {
	bc, bo, bg, bob int
	min             int
	building        composite

	cc, co, cg, cob int
}

func (k cachekey) String() string {
	return fmt.Sprintf("minute %d bots: [ore %d clay %d obs %d geo %d] stock: [ore %d clay %d obs %d geo %d] building: %s", k.min, k.bo, k.bc, k.bob, k.bg, k.co, k.cc, k.cob, k.cg, k.building)
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
	return 0
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
		fmt.Println("running blueprint", bp.id, "out of", len(blueprints))
		res := simulate(bp, minutes)
		fmt.Println("geodes:", res)
		sum += (res * bp.id)
	}
	return sum
}

func simulate(bp blueprint, minutes int) int {
	o, c, ob, g := bp.bots[ore], bp.bots[clay], bp.bots[obsidian], bp.bots[geode]
	orepergeode := g.costs[ore] * o.costs[ore]
	orepergeode += g.costs[ore] * g.costs[obsidian] * o.costs[ore]
	orepergeode += c.costs[ore] * ob.costs[clay] * o.costs[ore]

	claypergeode := g.costs[obsidian] * ob.costs[clay]
	obsidianpergeode := g.costs[obsidian]

	// queue := common.NewQueue[state, int]()
	queue := common.NewPriorityQueue[state, int](func(s state) int {
		return ((obsidianpergeode-s.collected[obsidian])*100 + (claypergeode-s.collected[clay])*10 +
			(orepergeode - s.collected[ore]) +
			s.bots[geode]*10000 + s.bots[obsidian]*1000 + s.bots[clay]*10 + s.bots[ore]) *
			s.minute
	})
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

		// if len(visited) > 1000 {
		// 	fmt.Println("queue len", queue.Len(), "visited len", len(visited))
		// }

		if cur.collected[geode] > maxc {
			maxc = cur.collected[geode]
		}

		key := cachekey{
			bc:  cur.bots[clay],
			bo:  cur.bots[ore],
			bg:  cur.bots[geode],
			bob: cur.bots[obsidian],
			cc:  cur.collected[clay],
			co:  cur.collected[ore],
			cg:  cur.collected[geode],
			cob: cur.collected[obsidian],
			// min:      cur.minute,
			building: cur.building,
		}

		if _, ok := visited[key]; ok {
			// fmt.Println("cache hit", key)
			continue
		}
		visited[key] = true

		if cur.collected[geode] < maxc {
			continue
		}

		// if maxc > 0 {
		// 	fmt.Println("maxc", maxc, cur.collected[geode])
		// }

		// if cur.collected[geode] <= maxc {
		// 	timeLeft := minutes - cur.minute

		// 	orepergeode := bp.bots[geode].costs[ore] * bp.bots[ore].costs[ore]
		// 	orepergeode += bp.bots[geode].costs[ore] * bp.bots[geode].costs[obsidian] * bp.bots[ore].costs[ore]
		// 	orepergeode += bp.bots[clay].costs[ore] * bp.bots[obsidian].costs[clay] * bp.bots[ore].costs[ore]

		// 	obsidianpergeode := bp.bots[geode].costs[obsidian]
		// 	claypergeode := bp.bots[geode].costs[obsidian] * bp.bots[obsidian].costs[clay]

		// 	// example: 1 ore bot, 4 cost per ore bot
		// 	// minute 4: build an ore bot, 0 ore 2 bots
		// 	// minute 6: build an ore bot, 0 ore, 3 bots
		// 	// minute 8: build an ore bot, 2 ore, 4 bots
		// 	// 24-8 = 16
		// 	// minus 3 to account for building at least 1 other type of bot
		// 	maxorebots := minutes - minutes/bp.bots[ore].costs[ore] - 3
		// 	maxtotalore := timeLeft*(maxorebots-cur.bots[ore]) + cur.collected[ore]

		// 	maxclaybots := minutes - minutes/bp.bots[clay].costs[ore] - 3
		// 	maxtotalclay := timeLeft*(maxclaybots-cur.bots[clay]) + cur.collected[clay]

		// 	maxobsbots := minutes - minutes/bp.bots[obsidian].costs[ore] - 3
		// 	maxtotalobs := timeLeft*(maxobsbots-cur.bots[obsidian]) + cur.collected[obsidian]

		// 	if orepergeode*maxc > maxtotalore || claypergeode*maxc > maxtotalclay || obsidianpergeode*maxc > maxtotalobs {
		// 		fmt.Println("per geode: ore", orepergeode, "clay", claypergeode, "obsidian", obsidianpergeode)
		// 		fmt.Println("max ore bots?", maxorebots)
		// 		fmt.Println("max total ore", maxtotalore, "with time left", timeLeft)
		// 		fmt.Println("time left", timeLeft)
		// 		fmt.Println(bp.bots)
		// 		continue
		// 	}
		// }

		cur = collect(cur)
		if cur.building != nothing {
			cur.bots[cur.building] = cur.bots[cur.building] + 1
			cur.building = nothing
		}

		mvs := getMoves(cur, bp, minutes, maxc)
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

func getMoves(st state, bp blueprint, minutes, maxgeodes int) []move {
	mvs := []move{}

	maxbots := minutes - minutes/bp.bots[ore].costs[ore] - 3 - maxgeodes // ore will be the most needed

	for _, c := range priority {
		if st.bots[c] > maxbots {
			continue
		}
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
