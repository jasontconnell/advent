package main

import (
	"fmt"
)

func main(){
	var A, B, C, D, E, F int

	cm := make(map[int]bool)
	dm := make(map[int]bool)

	C = 1250634
	D = 65536

	for {
		E = D % 256
		C += E
		C = C & 16777215
		C = C * 65899
    	C = C & 16777215
		if D < 256 {
			if _, ok := cm[C]; !ok {
				fmt.Println(C)
			}
			cm[C] = true

			D = C | 65536
			if _, ok := dm[D]; ok {
				fmt.Println("D repeats at", D)
				break
			}

			dm[D] = true

			C = 1250634
			continue
		}

		D = D / 256
	}

	fmt.Println(A, B, C, D, E, F)
}