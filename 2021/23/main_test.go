package main

import (
	"fmt"
	"testing"

	"github.com/jasontconnell/advent/common"
)

func TestDeadlock(t *testing.T) {
	lines, err := common.ReadStrings("nodeadlock.txt")
	if err != nil {
		t.Fail()
	}

	g, pods, rooms := parseInput(lines)

	result := isDeadlock(g, rooms, pods)

	if result {
		t.Fail()
	}
	printGrid(g)
}

func TestGetPodMoves(t *testing.T) {
	lines, err := common.ReadStrings("example3.txt")
	if err != nil {
		t.Fail()
	}

	g, _, rooms := parseInput(lines)
	pts := getPodMoves(g, rooms, xy{3, 2})
	fmt.Println(pts)
}

func TestGetMoves(t *testing.T) {
	lines, err := common.ReadStrings("example4.txt")
	if err != nil {
		t.Fail()
	}

	g, pods, rooms := parseInput(lines)
	printGrid(g)

	movable := allMovablePods(g, rooms, pods)
	for _, pt := range movable {
		fmt.Println("movable", pt, pods[pt])
		mvs := getMoves(g, pt, rooms)
		for _, mv := range mvs {
			fmt.Println(mv)
		}
	}
	roomable := roomablePods(g, rooms, pods)
	for len(roomable) != 0 {
		for _, pt := range roomable {
			fmt.Println("roomable", pt, pods[pt])
			rmvs := getMoves(g, pt, rooms)
			for _, mv := range rmvs {
				fmt.Println("roomable move", mv, mv.steps)
				g[mv.pos.y][mv.pos.x].pod = Empty
				g[mv.to.y][mv.to.x].pod = pods[pt]
			}
		}
		for pt, pod := range pods {
			mvs := getMoves(g, pt, rooms)
			for _, mv := range mvs {
				fmt.Println(pod, pt, mv)
			}
		}
		printGrid(g)

		roomable = roomablePods(g, rooms, pods)
		fmt.Println("new roomable")
		fmt.Println(roomable)
	}
}

func TestMovablePods(t *testing.T) {
	lines, err := common.ReadStrings("example2.txt")
	if err != nil {
		t.Fail()
	}

	g, pods, rooms := parseInput(lines)
	list := allMovablePods(g, rooms, pods)
	fmt.Println("pods", pods)
	for _, pt := range list {
		b := g[pt.y][pt.x]
		fmt.Println(pt, b.pod)
	}
}
