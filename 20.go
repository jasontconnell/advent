package main

import (
	"fmt"
	"time"
	"math"
)

var input = 36000000

func main() {
	startTime := time.Now()
	
	house := 100000
	presents := 0

	for presents < input {
		presents = GetPresents(house)

		if presents > input {
			fmt.Println("got presents", presents, "house", house)
			break
		}
		house++
	}

	presents2 := 0
	for presents2 < input {
		presents2 = GetPresents2(house)

		if presents2 > input {
			fmt.Println("got presents2", presents2, "house", house)
			break
		}
		house++
	}
	fmt.Println("Time", time.Since(startTime))
}

func GetPresents(max int) (presents int) {
	sqrt := int(math.Sqrt(float64(max)))+1	
	for i := 1; i <= sqrt; i++ {
		if max % i == 0 {
			presents += i
			presents += max/i
		}
	}	
	return presents * 10
}

func GetPresents2(max int) (presents int) {
	
	for i := 1; i <= 50; i++ {
		if max % i == 0 {
			presents += i
			presents += max/i
		}
	}	
	return presents * 11
}