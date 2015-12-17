package main

import (
	"bufio"
	"fmt"
	"os"
	"time"
	//"regexp"
	"strconv"
	//"strings"
	//"sort"
	"math/rand"
)
var input = "17.txt"

func main() {
	startTime := time.Now()
	if f, err := os.Open(input); err == nil {
		scanner := bufio.NewScanner(f)

		ary := []int{}
		max := 0

		for scanner.Scan() {
			var txt = scanner.Text()
			i,_ := strconv.Atoi(txt)
			if i > max {
				max = i
			}
			ary = append(ary, i)
		}

		
		//fmt.Println(perms)
		buckets := 0

		combos := Combinate(ary, 3) // n choose r for 20,10 should be 184756
		fmt.Println(len(combos))
		fmt.Println(combos)

		fmt.Println(buckets)
	}

	fmt.Println("Time", time.Since(startTime))
}

func Combinate(ints []int, size int) (combos [][]int) {
	max := len(ints)
	for i := 0; i < size; i++ {
		// choose 1 random then combinate size-1

		x := rand.Intn(max)

		sub := append(ints[:x], ints[x+1:]...)
		others := Combinate(sub, size-1)
		for _,other := range others {
			fmt.Println(other)
			combos = append(combos, append([]int { ints[x] }, other...))
		}
	}
	return
}



func Permutate(ints []int) [][]int {
	var ret [][]int

	if len(ints) == 2 {
		ret = append(ret, []int{ ints[0], ints[1] })
		ret = append(ret, []int{ ints[1], ints[0] })
	} else {

		for i := 0; i < len(ints); i++ {
			strc := make([]int, len(ints))
			copy(strc, ints)

			t := strc[i]
			sh := append(strc[:i], strc[i+1:]...)
			perm := Permutate(sh)
			
			for _,p := range perm {
				p = append([]int{ t }, p...)
				ret = append(ret, p)
			}
		}
	}

	return ret
}

func Perms(total, num int) [][]int {
	ret := [][]int{}

	if num == 2 {
		for i := 0; i < total/2 + 1; i++ {
			ret = append(ret, []int{ total-i, i })
			if i != total - i {
				ret = append(ret, []int{ i, total-i })
			}
		}
	} else {
		for i := 0; i <= total; i++ {
			perms := Perms(total-i, num-1)
			for _, p := range perms {
				q := append([]int{ i }, p...)
				ret = append(ret, q)
			}
		}
	}
	return ret
}


// reg := regexp.MustCompile("-?[0-9]+")
/* 			
if groups := reg.FindStringSubmatch(txt); groups != nil && len(groups) > 1 {
				fmt.Println(groups[1:])
			}
			*/
