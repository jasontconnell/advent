package main

import (
	"fmt"
	"github.com/jasontconnell/advent/2017/knot"
	"time"
)

var lengths = []int{206, 63, 255, 131, 65, 80, 238, 157, 254, 24, 133, 2, 16, 0, 1, 3}
var p2lengths = "206,63,255,131,65,80,238,157,254,24,133,2,16,0,1,3"

func main() {
	startTime := time.Now()

	res := knot.PrimitiveKnotHash(lengths)

	h := knot.KnotHash(p2lengths)

	fmt.Println("product of first two numbers:", res[0]*res[1])
	fmt.Println("hash value                  :", h)
	fmt.Println("Time", time.Since(startTime))
}
