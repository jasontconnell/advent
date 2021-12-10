package main

import (
	"fmt"
	"log"
	"regexp"
	"sort"
	"strconv"
	"time"

	"github.com/jasontconnell/advent/common"
)

var inputFilename = "input.txt"

type input = []string
type output = int

type Automaton struct {
	Key string

	IsBot    bool
	IsOutput bool

	LowDestKey string
	LowDest    *Automaton

	HighDestKey string
	HighDest    *Automaton

	Chips []int
}

func main() {
	startTime := time.Now()

	in, err := common.ReadStrings(inputFilename)
	if err != nil {
		log.Fatal(err)
	}

	p1 := part1(in)
	p2 := part2(in)

	fmt.Println("--2016 day 10 solution--")
	fmt.Println("Part 1:", p1)
	fmt.Println("Part 2:", p2)

	fmt.Println("Time", time.Since(startTime))
}

func part1(in input) string {
	bots := parseInput(in)
	id, _ := run(bots, 17, 61)
	return id
}

func part2(in input) output {
	bots := parseInput(in)
	_, outs := run(bots, 17, 61)

	return outs["output0"] * outs["output1"] * outs["output2"]
}

func parseInput(in input) []*Automaton {
	reg := regexp.MustCompile("^bot (\\d+) gives low to (output|bot) (\\d+) and high to (output|bot) (\\d+)$")
	inreg := regexp.MustCompile("^value (\\d+) goes to (output|bot) (\\d+)$")
	bots := make(map[string]*Automaton)

	for _, txt := range in {
		if groups := reg.FindStringSubmatch(txt); groups != nil && len(groups) > 1 {
			botkey := fmt.Sprintf("bot%s", groups[1])
			lowDestKey := fmt.Sprintf(groups[2]+"%s", groups[3])
			highDestKey := fmt.Sprintf(groups[4]+"%s", groups[5])

			var bot *Automaton
			var ok bool
			if bot, ok = bots[botkey]; ok {
				bot.LowDestKey = lowDestKey
				bot.HighDestKey = highDestKey
			} else {
				bot = &Automaton{IsBot: true}
				bot.Key = botkey
				bot.LowDestKey = lowDestKey
				bot.HighDestKey = highDestKey
				bots[botkey] = bot
			}

			if _, ok := bots[lowDestKey]; !ok {
				bot.LowDest = &Automaton{Key: lowDestKey, IsBot: groups[2] == "bot", IsOutput: groups[2] == "output"}
				bots[lowDestKey] = bot.LowDest
			}

			if _, ok := bots[highDestKey]; !ok {
				bot.HighDest = &Automaton{Key: highDestKey, IsBot: groups[4] == "bot", IsOutput: groups[4] == "output"}
				bots[highDestKey] = bot.HighDest
			}
		}

		if groups := inreg.FindStringSubmatch(txt); groups != nil && len(groups) > 1 {
			value, _ := strconv.Atoi(groups[1])
			botkey := fmt.Sprintf("bot%s", groups[3])

			var bot *Automaton
			var ok bool
			if bot, ok = bots[botkey]; ok {
				bot.Chips = append(bot.Chips, value)
				sort.Ints(bot.Chips)
			} else {
				bot := &Automaton{Key: botkey, IsBot: groups[2] == "bot", IsOutput: groups[2] == "output"}
				bot.Chips = append(bot.Chips, value)
				bot.LowDestKey = ""
				bot.HighDestKey = ""
				bots[botkey] = bot
			}
		}
	}

	for _, bot := range bots {
		if bot.IsBot && bot.LowDest == nil && bot.LowDestKey != "" {
			bot.LowDest = bots[bot.LowDestKey]
		}

		if bot.IsBot && bot.HighDest == nil && bot.HighDestKey != "" {
			bot.HighDest = bots[bot.HighDestKey]
		}
	}

	botlist := []*Automaton{}
	for _, bot := range bots {
		botlist = append(botlist, bot)
	}

	return botlist
}

func processBot(bot *Automaton) {
	if bot.LowDest != nil && bot.HighDest != nil {
		low, high := bot.Chips[0], bot.Chips[1]
		if len(bot.LowDest.Chips) < 2 {
			bot.LowDest.Chips = append(bot.LowDest.Chips, low)
			low = -1
			sort.Ints(bot.LowDest.Chips)
		}

		if len(bot.HighDest.Chips) < 2 {
			bot.HighDest.Chips = append(bot.HighDest.Chips, high)
			high = -1
			sort.Ints(bot.HighDest.Chips)
		}

		if low == -1 && high == -1 {
			bot.Chips = []int{}
		} else if low == -1 {
			bot.Chips = bot.Chips[1:]
		} else if high == -1 {
			bot.Chips = bot.Chips[:1]
		}
	}
}

func run(bots []*Automaton, chip1, chip2 int) (string, map[string]int) {
	botID := ""

	proc := []*Automaton{}
	for _, bot := range bots {
		if bot.IsBot && len(bot.Chips) == 2 {
			proc = append(proc, bot)
		}
	}

	for len(proc) > 0 {
		bot := proc[0]
		proc = proc[1:]
		if bot.Chips[0] == chip1 && bot.Chips[1] == chip2 {
			botID = bot.Key
		}

		processBot(bot)

		if bot.LowDest.IsBot && len(bot.LowDest.Chips) == 2 {
			proc = append(proc, bot.LowDest)
		}
		if bot.HighDest.IsBot && len(bot.HighDest.Chips) == 2 {
			proc = append(proc, bot.HighDest)
		}
	}

	outputs := make(map[string]int)
	for _, bot := range bots {
		if bot.IsOutput && len(bot.Chips) > 0 {
			outputs[bot.Key] = bot.Chips[0]
		}
	}

	return botID, outputs
}
