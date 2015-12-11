package main

import (
	"fmt"
	"time"
	"sync"
	"runtime"
	"math"
)
//var input = "3113322113"

var input = []int{ 3,1,1,3,3,2,2,1,1,3 }

func main(){
	startTime := time.Now()
	ints := input
	fmt.Println("start", ints)
	for i := 0; i < 50; i++ {
		ints = lookAndSayDelegation(ints)
		fmt.Println("on iteration", i, time.Since(startTime), "string length =", len(ints))
	}

	fmt.Println(len(ints))
	duration := time.Since(startTime)
	fmt.Println("time", duration)
}

func lookAndSayDelegation(ints []int) (output []int) {
	workers := 1
	//l := len(ints)
	// if l > 1500000 {
	// 	workers = runtime.NumCPU()
	// } else if l > 300000 {
	// 	workers = 4
	// } else if l > 1000 {
	// 	workers = 2
	// }

	runtime.GOMAXPROCS(workers)
	fmt.Println("using", workers, "workers")

	result := make([][]int, workers)
	start := 0
	length := len(ints) / workers
	lastLen := length

	wg := sync.WaitGroup{}
	wg.Add(workers)

	cp := make([]int, len(ints))
	copy(cp, ints)

	for i := 0; i < workers; i++ {
		end := int(math.Min(float64(len(cp)), float64(start+lastLen)))
		
		part := cp[start:end]

		if end < len(ints) && cp[end-1] == cp[end] {
			part = append(part, cp[end])
			if end < len(cp) - 1 && cp[end] == cp[end+1] {
				part = append(part, cp[end+1])
			}
		}
		lastLen = len(part)
		start = start+lastLen

		go func(index int){
			output := lookAndSay(part)
			result[index] = output
			wg.Done()
		}(i)
	}
	wg.Wait()

	for _,r := range result {
		output = append(output, r...)
	}
	return
}

func lookAndSay(ints []int) (output []int) {
	count,digit := 0, ints[0]
	for _,c := range ints {
		if digit == c {
			count++
		} else {
			output = append(output, count, digit)
			count = 1
		}
		digit = c
	}
	output = append(output, count, digit)
	return
}