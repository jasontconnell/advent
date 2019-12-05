package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

var input = "02.txt"

func main() {
	startTime := time.Now()

	f, err := os.Open(input)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	scanner := bufio.NewScanner(f)
	opcodes := []int{}

	if scanner.Scan() {
		var txt = scanner.Text()
		//var txt = "1,9,10,3,2,3,11,0,99,30,40,50"
		sopcodes := strings.Split(txt, ",")
		for _, s := range sopcodes {
			i, err := strconv.Atoi(s)
			if err != nil {
				fmt.Println(err)
				continue
			}

			opcodes = append(opcodes, i)
		}
	}

	p1in := make([]int, len(opcodes))
	copy(p1in, opcodes)
	p1in[1] = 12
	p1in[2] = 2
	out1 := exec(p1in)

	out2 := part2(opcodes, 19690720)

	fmt.Println("Part 1: ", out1)
	fmt.Println("Part 2: ", out2)

	fmt.Println("Time", time.Since(startTime))
}

func part2(opcodes []int, goal int) int {
	vals := []int{0, 0}
	done := false

	for i := 0; i < 100 && !done; i++ {
		for j := 0; j < 100 && !done; j++ {
			c := make([]int, len(opcodes))
			copy(c, opcodes)

			c[1] = i
			c[2] = j

			vals[0] = i
			vals[1] = j

			t := exec(c)

			if t == goal {
				done = true
				break
			}
		}
	}

	return 100*vals[0] + vals[1]
}

func exec(ops []int) int {
	done := false
	for i := 0; i < len(ops) && !done; i += 4 {
		op := ops[i]
		switch op {
		case 1:
			ops[ops[i+3]] = ops[ops[i+1]] + ops[ops[i+2]]
			break
		case 2:
			ops[ops[i+3]] = ops[ops[i+1]] * ops[ops[i+2]]
			break
		case 99:
			done = true
			break
		}
	}

	return ops[0]
}