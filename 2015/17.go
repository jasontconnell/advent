package main

import (
	"bufio"
	"fmt"
	"os"
	"time"
	//"regexp"
	"strconv"
	//"strings"
	//"sort"
	//"math/rand"
)

var input = "17.txt"

func main() {
	startTime := time.Now()
	if f, err := os.Open(input); err == nil {
		scanner := bufio.NewScanner(f)

		ary := []int{}
		max := 0

		for scanner.Scan() {
			var txt = scanner.Text()
			i, _ := strconv.Atoi(txt)
			if i > max {
				max = i
			}
			ary = append(ary, i)
		}

		fmt.Println(ary)

		buckets := count(ary, 150, len(ary), 0)

		result := 0
		i := 0
		for result == 0 {
			result = count(ary, 150, i, 0)
			i++
		}

		fmt.Println("buckets", buckets)
		fmt.Println("result", result)
	}

	fmt.Println("Time", time.Since(startTime))
}

func count(list []int, total, n, i int) int {
	if i < 0 {
		i = 0
	}

	if n < 0 {
		return 0
	} else if total == 0 {
		return 1
	} else if i == len(list) || total < 0 {
		return 0
	} else {
		return count(list, total, n, i+1) + count(list, total-list[i], n-1, i+1)
	}
}
