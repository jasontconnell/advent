package main

import (
	"fmt"
	"strconv"
	"time"
	"sync"
	"runtime"
	"math"
)
//var input = "3113322113"

var input = []int{ 3,1,1,3,3,2,2,1,1,3 }

func main(){
	startTime := time.Now()
	str := input
	fmt.Println("start", str)
	for i := 0; i < 50; i++ {
		str = lookAndSayDelegation(str)
		fmt.Println("on iteration", i, time.Since(startTime), "string length =", len(str))
	}

	fmt.Println(len(str))
	duration := time.Since(startTime)
	fmt.Println("time", duration)
}

func lookAndSayDelegation(array []int) (output []int) {
	workers := 1
	l := len(array)
	if l > 1500000 {
		workers = runtime.NumCPU()
	} else if l > 300000 {
		workers = 4
	} else if l > 1000 {
		workers = 2
	}

	runtime.GOMAXPROCS(workers)
	fmt.Println("using", workers, "workers")

	result := make([]int, workers)
	start := 0
	length := len(str) / workers
	lastStrLen := length

	wg := sync.WaitGroup{}
	wg.Add(workers)

	for i := 0; i < workers; i++ {
		end := int(math.Min(float64(len(str)), float64(start+lastStrLen)))
		
		part := str[start:end]

		if end < len(str) && str[end-1] == str[end] {
			part = part + string(str[end])
			if end < len(str) - 1 && str[end] == str[end+1] {
				part = part + string(str[end+1])
			}
		}
		lastStrLen = len(part)
		start = start+lastStrLen

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

func lookAndSay(str string) (output string) {
	count,digit := 0, rune(str[0])
	for _,c := range str {
		if digit == c {
			count++
		} else {
			output += strconv.Itoa(count) + string(digit)
			count = 1
		}
		digit = c
	}
	output += strconv.Itoa(count) + string(digit)
	return
}