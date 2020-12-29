package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"time"
)

var input = "22_test.txt"

type deck struct {
	cards []int
}

func main() {
	startTime := time.Now()

	f, err := os.Open(input)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	scanner := bufio.NewScanner(f)

	lines := []string{}
	for scanner.Scan() {
		var txt = scanner.Text()
		lines = append(lines, txt)
	}

	player1, player2 := readDecks(lines)
	winner := play(player1, player2)
	p1 := getDeckValue(winner)

	fmt.Println("Part 1:", p1)

	p2winner := playRecursiveCombat(player1, player2)
	p2 := getDeckValue(p2winner)

	fmt.Println("Part 2:", p2)

	fmt.Println("Time", time.Since(startTime))
}

func getDeckValue(d deck) int {
	c := 1
	v := 0
	for i := len(d.cards) - 1; i >= 0; i-- {
		v += c * d.cards[i]
		c++
	}
	return v
}

func play(player1, player2 deck) deck {
	for len(player1.cards) > 0 && len(player2.cards) > 0 {
		c1, c2 := player1.cards[0], player2.cards[0]
		player1.cards = player1.cards[1:]
		player2.cards = player2.cards[1:]

		if c1 > c2 {
			player1.cards = append(player1.cards, []int{c1, c2}...)
		} else {
			player2.cards = append(player2.cards, []int{c2, c1}...)
		}
	}

	winner := player1
	if len(winner.cards) == 0 {
		winner = player2
	}
	return winner
}

func playRecursiveCombat(player1, player2 deck) deck {
	p1m := make(map[int]int)
	p2m := make(map[int]int)

	p1m[getDeckValue(player1)] = getDeckValue(player1)
	p2m[getDeckValue(player2)] = getDeckValue(player2)

	return recurse(player1, player2, p1m, p2m)
}

func recurse(player1, player2 deck, p1m map[int]int, p2m map[int]int) deck {
	p1v := getDeckValue(player1)
	p2v := getDeckValue(player2)

	_, p1ok := p1m[p1v]
	_, p2ok := p2m[p2v]

	if p1ok || p2ok {
		return player1
	}

	p1m[p1v] = p1v
	p2m[p2v] = p2v

	if len(player1.cards) == 0 {
		return player2
	}

	if len(player2.cards) == 0 {
		return player1
	}

	p1c := player1.cards[0]
	p2c := player2.cards[0]

	player1.cards = player1.cards[1:]
	player2.cards = player2.cards[1:]

	if p1c < len(player1.cards) {
		return player2
	}

	if p2c < len(player2.cards) {
		return player1
	}

	return recurse(player1, player2, p1m, p2m)
}

func readDecks(lines []string) (deck, deck) {
	player1 := deck{}
	player2 := deck{}

	cur := &player1
	for _, line := range lines {
		if line == "" {
			cur = &player2
			continue
		}

		c, err := strconv.Atoi(line)
		if err != nil {
			continue
		}

		cur.cards = append(cur.cards, c)
	}

	return player1, player2
}

// reg := regexp.MustCompile("-?[0-9]+")
/*
if groups := reg.FindStringSubmatch(txt); groups != nil && len(groups) > 1 {
				fmt.Println(groups[1:])
			}
*/
