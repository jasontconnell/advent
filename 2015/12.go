package main

import (
	"bufio"
	"fmt"
	"os"
	"time"
	"regexp"
	"strconv"
	//"strings"
)
var input = "12.txt"

func main() {
	startTime := time.Now()
	if f, err := os.Open(input); err == nil {
		scanner := bufio.NewScanner(f)

		reg := regexp.MustCompile("-?[0-9]+")
		sum := 0
		for scanner.Scan() {
			var txt = scanner.Text()
			if str := reg.FindAllString(txt,-1); len(str) > 0 { 
				for _,r := range str {
					if i,err := strconv.Atoi(string(r)); err == nil {
						sum += i
					} else {
						fmt.Print(string(r), "not a number")
					}
				}
			}
		}
		fmt.Println("sum", sum)
	}

	fmt.Println("Time", time.Since(startTime))

}
