package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"time"
)

var input = "22.txt"

type deck struct {
	id    int
	cards []int
}

func cardstr(cards []int) string {
	s := ""
	for _, i := range cards {
		s += strconv.Itoa(i) + "_"
	}
	return s
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

	fmt.Println("Part 1: Player", winner.id, "wins! Score:", p1)

	p2winner := playRecursiveCombat(player1, player2)
	p2 := getDeckValue(p2winner)

	fmt.Println("Part 2: Player", p2winner.id, "wins! Score:", p2)

	fmt.Println("Time", time.Since(startTime))
}

func getDeckValue(d deck) int {
	if len(d.cards) == 0 {
		return 0
	}
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
	return recurse(player1, player2)
}

func recurse(player1, player2 deck) deck {

	if len(player1.cards) == 0 {
		return player2
	}

	if len(player2.cards) == 0 {
		return player1
	}

	determined := false
	p1m := make(map[int]int)
	p2m := make(map[int]int)

	var winner deck
	for !determined {

		p1c := player1.cards[0]
		p2c := player2.cards[0]

		// fmt.Println("at start", player1, player2)

		player1.cards = player1.cards[1:]
		player2.cards = player2.cards[1:]

		p1v := getDeckValue(player1)
		p2v := getDeckValue(player2)

		_, p1ok := p1m[p1v]
		_, p2ok := p2m[p2v]

		playRound := true
		if p1ok && p2ok {
			player1.cards = append(player1.cards, []int{p1c, p2c}...)
			playRound = false
		}

		if p1v != 0 {
			p1m[p1v] = p1v
		}
		if p2v != 0 {
			p2m[p2v] = p1v
		}

		if playRound {
			if p1c <= len(player1.cards) && p2c <= len(player2.cards) {
				p1copy := copyDeck(player1, p1c)
				p2copy := copyDeck(player2, p2c)
				recurseWinner := recurse(p1copy, p2copy)
				if recurseWinner.id == 1 {
					player1.cards = append(player1.cards, []int{p1c, p2c}...)
				} else if recurseWinner.id == 2 {
					player2.cards = append(player2.cards, []int{p2c, p1c}...)
				}
			} else {
				if p1c > p2c {
					player1.cards = append(player1.cards, []int{p1c, p2c}...)
				} else {
					player2.cards = append(player2.cards, []int{p2c, p1c}...)
				}
			}
		}

		determined = len(player1.cards) == 0 || len(player2.cards) == 0
	}

	if winner.id == 0 {
		winner = player1
		if len(winner.cards) == 0 {
			winner = player2
		}
	}

	return winner
}

func copyDeck(d deck, length int) deck {
	dc := deck{id: d.id}
	dc.cards = make([]int, length)
	copy(dc.cards, d.cards[:length])
	return dc
}

func readDecks(lines []string) (deck, deck) {
	player1 := deck{id: 1}
	player2 := deck{id: 2}

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
