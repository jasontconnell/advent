package main

import "fmt"

func main() {

	h := 0
	for b := 100000 + 9900; b < 100000+9900+17000+1; b += 17 {
		for x := 2; x < b; x++ {
			if b%x == 0 {
				h++
				break
			}
		}
	}
	fmt.Println(h)
	return

}
