package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"time"
)

var input = "25.txt"

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

	keys := getKeys(lines)
	doorenc, cardenc := decrypt(keys)
	if doorenc != cardenc {
		panic("wrong!!")
	}

	fmt.Println("Part 1:", doorenc)
	fmt.Println("Time", time.Since(startTime))
}

func decrypt(keys []int) (doorenc, cardenc int) {
	door, card := keys[1], keys[0]
	dloops := getLoopsize(door)
	cloops := getLoopsize(card)

	cardenc = getEncryptionKey(card, dloops)
	doorenc = getEncryptionKey(door, cloops)
	return
}

func getEncryptionKey(pk int, loops int) int {
	subj := pk
	modv := 20201227

	cur := 1
	for i := 0; i < loops; i++ {
		cur = (subj * cur) % modv
	}
	return cur
}

func getLoopsize(k int) int {
	subj := 7
	modv := 20201227
	loops := 0

	var cur int = 1
	for cur != k {
		cur = (subj * cur) % modv
		loops++
	}
	return loops
}

func getKeys(lines []string) []int {
	keys := []int{}
	for _, line := range lines {
		k, _ := strconv.Atoi(line)
		keys = append(keys, k)
	}
	return keys
}
