package main

import (
	"fmt"
	"time"
	//"math"
)

var input = 165061
//var input = 59414

type elf struct {
	currentRecipe int
}

func main() {
	startTime := time.Now()

	next10 := process([]elf{
		elf{currentRecipe: 0},
		elf{currentRecipe: 1},
	}, []int{3, 7}, input)
	countBefore := process2([]elf{
		elf{currentRecipe: 0},
		elf{currentRecipe: 1},
	}, []int{3, 7}, input)
	fmt.Println("Part 1: ", str(next10))
	fmt.Println("Part 2: ", countBefore)

	fmt.Println("Time", time.Since(startTime))
}

func str(a []int) string {
	s := ""
	for _, i := range a {
		s += fmt.Sprintf("%d", i)
	}
	return s
}

func process(elves []elf, recipes []int, count int) []int {
	done := false
	for !done {
		elves, recipes = processOne(elves, recipes)
		done = len(recipes) > count+11
	}

	return recipes[count : count+10]
}

func process2(elves []elf, recipes []int, count int) int {
	val := digits(count)
	fmt.Println(count, val)
	curindex := 0
	found := false
	for !found {
		elves, recipes = processOne(elves, recipes)
		found, curindex = find(recipes, val, curindex)
	}

	return curindex-1
}

func find(recipes []int, val []int, start int) (bool, int) {
	if start+len(val) >= len(recipes){
		return false, start
	}

	found := true
	for j := 0; j < len(val); j++ {
		found = found && (recipes[start+j] == val[j])
	}
	return found, start+1
}

func processOne(elves []elf, recipes []int) ([]elf, []int) {
	newscore := 0
	for i := 0; i < len(elves); i++ {
		newscore += recipes[elves[i].currentRecipe]
	}

	d := digits(newscore)
	for _, di := range d {
		recipes = append(recipes, di)
	}

	for i := 0; i < len(elves); i++ {
		elves[i].currentRecipe = nextRecipeIndex(elves[i], recipes)
	}

	return elves, recipes
}

func nextRecipeIndex(e elf, recipes []int) int {
	c := recipes[e.currentRecipe]
	return (1 + e.currentRecipe + c) % len(recipes)
}

func digits(val int) []int {
	done := false
	c := 10
	digits := []int{}
	for !done {
		done = c > val
		s := val % c
		digits = append(digits, s)
		val = val / c
	}
	rev := []int{}
	for i := len(digits) - 1; i >= 0; i-- {
		rev = append(rev, digits[i])
	}
	return rev
}
