package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"time"
	//"regexp"
	//"strconv"
	//"strings"
	//"math"
)

var input = "01.txt"

func main() {
	startTime := time.Now()

	f, err := os.Open(input)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	scanner := bufio.NewScanner(f)
	list := []int{}
	for scanner.Scan() {
		var txt = scanner.Text()
		i := getInt(txt)
		list = append(list, i)
	}

	p1 := sum(list)
	p2 := findRpt(list)

	fmt.Println("Part 1:", p1)
	fmt.Println("Part 2:", p2)

	fmt.Println("Time", time.Since(startTime))
}

func sum(list []int) int {
	s := 0
	for _, i := range list {
		s += i
	}
	return s
}

func findRpt(list []int) int {
	found := false
	m := make(map[int]bool)
	rpt := 0
	s := 0
	for !found {
		for _, i := range list {
			s += i
			if _, ok := m[s]; ok && rpt == 0 {
				rpt = s
				found = true
			}
			m[s] = true
		}
	}
	return rpt
}

func getInt(line string) int {
	i, _ := strconv.Atoi(line)
	return i
}
