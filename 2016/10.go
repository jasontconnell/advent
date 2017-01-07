package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strconv"
	"time"
	//"strings"
	//"math"
)

var input = "10.txt"

func main() {
	startTime := time.Now()
	if f, err := os.Open(input); err == nil {
		scanner := bufio.NewScanner(f)

		reg := regexp.MustCompile("^(bot|output) (\\d+) gives low to (output|bot) (\\d+) and high to (output|bot) (\\d+)$")
		inreg := regexp.MustCompile("^value (\\d+) goes to (output|bot) (\\d+)$")

		bots := make(map[string]*Automaton)

		for scanner.Scan() {
			var txt = scanner.Text()
			if groups := reg.FindStringSubmatch(txt); groups != nil && len(groups) > 1 {
				id := groups[1] + groups[2]
				lowDestID := groups[3] + groups[4]
				highDestID := groups[5] + groups[6]

				if bot, ok := bots[id]; ok {
					bot.LowDestID = lowDestID
					bot.HighDestID = highDestID
				} else {
					newbot := Automaton{IsBot: groups[1] == "bot"}
					newbot.ID = id
					newbot.LowDestID = lowDestID
					newbot.HighDestID = highDestID
					bots[id] = &newbot
				}

				if _, ok := bots[lowDestID]; !ok {
					bots[lowDestID] = &Automaton{ID: lowDestID, IsBot: groups[3] == "bot"}
				}

				if _, ok := bots[highDestID]; !ok {
					bots[highDestID] = &Automaton{ID: highDestID, IsBot: groups[5] == "bot"}
				}
			}

			if groups := inreg.FindStringSubmatch(txt); groups != nil && len(groups) > 1 {
				value, _ := strconv.Atoi(groups[1])
				id := groups[2] + groups[3]

				if bot, ok := bots[id]; ok {
					bot.Chips = append(bot.Chips, value)
					sort.Ints(bot.Chips)
				} else {
					newbot := Automaton{IsBot: groups[2] == "bot"}
					newbot.Chips = append(newbot.Chips, value)
					newbot.ID = id
					newbot.LowDestID = ""
					newbot.HighDestID = ""
					bots[id] = &newbot
				}
			}
		}

		botlist := []*Automaton{}
		for _, bot := range bots {
			botlist = append(botlist, bot)
		}
		for _, bot := range botlist {
			if bot.IsBot {
				if dest, ok := bots[bot.LowDestID]; ok && bot.LowDestID != "" {
					bot.LowDest = dest
				}

				if dest, ok := bots[bot.HighDestID]; ok && bot.HighDestID != "" {
					bot.HighDest = dest
				}
			}
		}

		finished := process(botlist)
		for !finished {
			finished = process(botlist)
		}

		for _, bot := range bots {
			if !bot.IsBot {
				if bot.ID == "output0" || bot.ID == "output1" || bot.ID == "output2" {
					fmt.Println(bot)
					value := 1
					for _, c := range bot.Chips {
						value = value * c
					}
					fmt.Println(value)
				}
			}
		}
	}

	fmt.Println("Time", time.Since(startTime))
}

func process(bots []*Automaton) bool {
	moves := 0
	for _, bot := range bots {
		if bot.IsBot && len(bot.Chips) == 2 {
			if bot.Chips[0] == 17 && bot.Chips[1] == 61 {
				fmt.Println("found bot", bot)
			}
			if bot.LowDest != nil && bot.HighDest != nil {
				if len(bot.LowDest.Chips) < 2 {
					bot.LowDest.Chips = append(bot.LowDest.Chips, bot.Chips[0])
					sort.Ints(bot.LowDest.Chips)
					moves++
					bot.Chips = bot.Chips[1:]

				}

				if len(bot.HighDest.Chips) < 2 {
					c := len(bot.Chips)
					bot.HighDest.Chips = append(bot.HighDest.Chips, bot.Chips[c-1])
					sort.Ints(bot.HighDest.Chips)
					bot.Chips = bot.Chips[c-1:]
					moves++
				}
			}
		}
	}

	return moves == 0
}

type Automaton struct {
	ID string

	IsBot bool

	LowDestID string
	LowDest   *Automaton

	HighDestID string
	HighDest   *Automaton

	Chips []int
}

func (a Automaton) String() string {
	return a.ID + " LowDest: " + a.LowDestID + " HighDest: " + a.HighDestID
}
