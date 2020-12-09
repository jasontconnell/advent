package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"time"
)

var input = "09_test.txt"

type infile struct {
	filename string
	preamble int
}

func main() {
	startTime := time.Now()

	infile1 := infile{"09.txt", 25}
	infile2 := infile{"09_test.txt", 5}

	inmap := make(map[string]infile)

	inmap["main"] = infile1
	inmap["test"] = infile2

	file := inmap["main"]

	f, err := os.Open(file.filename)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	scanner := bufio.NewScanner(f)

	nums := []int{}
	for scanner.Scan() {
		var txt = scanner.Text()
		i, err := strconv.Atoi(txt)
		if err != nil {
			fmt.Println(err)
			continue
		}
		nums = append(nums, i)
	}

	p1 := firstInvalidNumber(nums, file.preamble)
	fmt.Println("part 1:", p1)

	a := findContiguousSum(p1, nums)
	sort.Ints(a)

	p2 := -1
	if a != nil && len(a) > 1 {
		p2 = a[0] + a[len(a)-1]
	}
	fmt.Println("Part 2:", p2)

	fmt.Println("read", len(nums), "numbers")
	fmt.Println("Time", time.Since(startTime))
}

func firstInvalidNumber(nums []int, size int) int {
	p1 := 0
	for i := size; i < len(nums); i++ {
		preamble := getPreamble(nums, i-size, size)
		if !hasMatch(nums[i], preamble) {
			p1 = nums[i]
			break
		}
	}
	return p1
}

func findContiguousSum(num int, nums []int) []int {
	q := []int{}
	for i := 0; i < len(nums); i++ {
		if nums[i] == num {
			continue
		}

		if nums[i] > num {
			break
		}

		s := sum(q)
		for s > num {
			q = q[1:]
			s = sum(q)
		}

		if s == num {
			return q
		}

		q = append(q, nums[i])

	}
	return nil
}

func sum(a []int) int {
	s := 0
	for _, n := range a {
		s += n
	}
	return s
}

func hasMatch(num int, preamble map[int]int) bool {
	for _, v := range preamble {
		x := num - v
		_, ok := preamble[x]
		if ok {
			return true
		}
	}
	return false
}

func getPreamble(nums []int, start, size int) map[int]int {
	preamble := make(map[int]int)
	for i := start; i < start+size; i++ {
		preamble[nums[i]] = nums[i]
	}
	return preamble
}
