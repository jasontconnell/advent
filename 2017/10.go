package main

import (
	"fmt"
	"time"
)

var lengths = []int{206, 63, 255, 131, 65, 80, 238, 157, 254, 24, 133, 2, 16, 0, 1, 3}
var p2lengths = "206,63,255,131,65,80,238,157,254,24,133,2,16,0,1,3"

func main() {
	startTime := time.Now()

	vals := make([]int, 256)
	vals2 := make([]int, 256)
	for i := 0; i < 256; i++ {
		vals[i] = i
		vals2[i] = i
	}

	lengths2 := []int{}
	for i := 0; i < len(p2lengths); i++ {
		lengths2 = append(lengths2, int(p2lengths[i]))
	}
	lengths2 = append(lengths2, []int{17, 31, 73, 47, 23}...)

	_, _, res := knot(vals, lengths, 0, 0)

	h := hash(vals2, lengths2)

	fmt.Println("product of first two numbers:", res[0]*res[1])
	fmt.Println("hash value                  :", h)
	fmt.Println("Time", time.Since(startTime))
}

func hash(val, lens []int) string {
	pos := 0
	skip := 0
	for i := 0; i < 64; i++ {
		pos, skip, val = knot(val, lens, pos, skip)
	}

	xors := []int{}
	for i := 0; i < len(val); i += 16 {
		xors = append(xors, xor(val[i:i+16]))
	}

	result := ""
	for i := 0; i < len(xors); i++ {
		result += fmt.Sprintf("%x", xors[i])
	}

	return result
}

func xor(val []int) int {
	result := val[0]
	for i := 1; i < len(val); i++ {
		result ^= val[i]
	}

	return result
}

func knot(val, lens []int, pos, skip int) (p int, s int, result []int) {
	for i := 0; i < len(lens); i++ {
		val = reverse(val, pos, pos+lens[i]-1)
		pos = (pos + lens[i] + skip) % len(val)
		skip = skip + 1
	}

	return pos, skip, val
}

func reverse(val []int, start, end int) []int {
	cp := make([]int, len(val))

	copy(cp, val)
	dst := start % len(val)

	for i := end; i >= start; i-- {
		cp[dst%len(val)] = val[i%len(val)]
		dst++
	}
	return cp
}
