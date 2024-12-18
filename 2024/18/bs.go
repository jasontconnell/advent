// You can edit this code!
// Click here and start typing.
package main

import (
	"fmt"
	"log"
	"math/rand"
	"sort"
)

func main() {
	list := []int{130}
	for i := 0; i < 98; i++ {
		p := rand.Int() % 1000
		list = append(list, p)
	}
	sort.Ints(list)
	idx := binsearch(list, 130)
	if idx >= 0 {
		fmt.Println(idx, list[idx])
	} else {
		fmt.Println("not found", list)
	}
}

func binsearch(list []int, target int) int {
	min, max := 0, len(list)
	idx := -1
	loops := 0
	for {
		i := (max + min) / 2
		if list[i] == target {
			idx = i
			break
		} else if list[i] > target {
			max = i - 1
		} else if list[i] < target {
			min = i + 1
		}

		if loops > len(list) {
			log.Println(loops, i, min, max)
			break
		}
		loops++
	}
	return idx
}
