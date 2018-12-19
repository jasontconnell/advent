package main

import (
	"fmt"
	"time"
)

type node struct {
	left  *node
	right *node
	val   int
}

type player struct {
	num   int
	score int
}

func main() {
	startTime := time.Now()

	playerCount := 458
	lastMarble := 72019

	highScore := play(playerCount, lastMarble)
	highScore2 := play(playerCount, lastMarble*100)

	fmt.Println("High Score:", highScore)
	fmt.Println("High Score 2:", highScore2)

	fmt.Println("Time", time.Since(startTime))
}

func print(n *node) {
	tmp := n.right

	fmt.Print("(", n.val, ") ")
	for tmp != n {
		fmt.Print(tmp.val, " ")
		tmp = tmp.right
	}

	fmt.Println("")
}

func play(playerCount, lastMarble int) int {
	cur := &node{val: 0}
	cur.left = cur
	cur.right = cur

	players := []*player{}
	for i := 0; i < playerCount; i++ {
		p := &player{num: i, score: 0}
		players = append(players, p)
	}

	for i := 0; i < lastMarble; i++ {
		pid := i % playerCount
		p := players[pid]

		var scoreDelta int
		cur, scoreDelta = insert(i+1, cur)
		p.score += scoreDelta
	}

	return highScore(players)
}

func highScore(players []*player) int {
	high := 0
	for _, p := range players {
		if p.score > high {
			high = p.score
		}
	}
	return high
}

func move(n *node, steps, dir int) *node {
	for i := 0; i < steps; i++ {
		switch dir {
		case -1:
			n = n.left
		case 1:
			n = n.right
		}
	}
	return n
}

func moveCW(n *node, steps int) *node {
	return move(n, steps, 1)
}

func moveCCW(n *node, steps int) *node {
	return move(n, steps, -1)
}

func insert(val int, cur *node) (*node, int) {
	if val%23 != 0 {
		ins := &node{val: val}
		cur = moveCW(cur, 1)
		tmp := cur.right
		cur.right = ins
		ins.left = cur
		ins.right = tmp
		tmp.left = ins
		return ins, 0
	}

	score := val
	cur = moveCCW(cur, 7)
	tmp := moveCW(cur, 1)
	score += remove(cur)

	return tmp, score
}

func remove(cur *node) int {
	val := cur.val
	tmp := cur.left
	tmp.right = cur.right
	cur.right.left = tmp
	cur = nil
	return val
}

func sum(ints []int) int {
	s := 0
	for _, i := range ints {
		s += i
	}
	return s
}
