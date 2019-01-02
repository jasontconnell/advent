package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"time"
)

var input = "24.txt"

func main() {
	startTime := time.Now()
	if f, err := os.Open(input); err == nil {
		scanner := bufio.NewScanner(f)

		var list []int
		sum := 0
		for scanner.Scan() {
			var txt = scanner.Text()
			weight, _ := strconv.Atoi(txt)
			list = append(list, weight)
			sum += weight
		}

		// for part 1 use 3, part 2 use 4
		each := sum / 4

		sort.Ints(list)
		fmt.Println(len(list), sum, each)
		fmt.Println(list)

		for i := 0; i < len(list); i++ {
			cp := make([]int, len(list))
			copy(cp, list)
			cp = append(cp[:i], cp[i+1:]...)
			filled := Fill(cp, 0, each)
			product := 1
			for _, num := range filled {
				product = product * num
			}
			fmt.Println("product", product)
		}
	}

	fmt.Println("Time", time.Since(startTime))
}

func Fill(list []int, bucket, max int) []int {
	nums := []int{}

	for i := len(list) - 1; i >= 0; i-- {
		if bucket+list[i] <= max {
			bucket += list[i]
			nums = append(nums, list[i])
		}
	}

	return nums
}
