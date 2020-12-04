package main

import (
	"bufio"
	"fmt"
	"log"
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
		n, err := strconv.Atoi(txt)
		if err != nil {
			log.Fatal(err)
		}
		list = append(list, n)
	}

	x, y := findTwo(list, 2020)
	fmt.Println(x, y, x*y)

	x, y, z := findThree(list, 2020)
	fmt.Println(x, y, z, x*y*z)

	fmt.Println("Time", time.Since(startTime))
}

func findTwo(list []int, sum int) (int, int) {
	m := make(map[int]int)
	for _, n := range list {
		m[n] = n
	}

	for _, n := range list {
		x := sum - n
		if n2, ok := m[x]; ok {
			return n, n2
		}
	}
	return -1, -1
}

func findThree(list []int, sum int) (int, int, int) {
	m := make(map[int]int)
	for _, x := range list {
		m[x] = x
	}

	for _, x := range list {
		for _, y := range list {
			if z, ok := m[sum-x-y]; x != y && ok {
				return x, y, z
			}
		}
	}
	return -1, -1, -1
}
