package main

import (
    "bufio"
    "fmt"
    "os"
    "time"
    //"regexp"
    //"strconv"
    //"strings"
    //"math"
)
var input = "2.txt"

func main() {
    startTime := time.Now()
    if f, err := os.Open(input); err == nil {
        scanner := bufio.NewScanner(f)

        for scanner.Scan() {
            var txt = scanner.Text()
            fmt.Println(txt)
        }
    }

    numpad := [][]int { []int { 1, 2, 3 }, []int { 4, 5, 6}, []int{ 7, 8, 9 }}

    fmt.Println(numpad)

    fmt.Println("Time", time.Since(startTime))
}


// reg := regexp.MustCompile("-?[0-9]+")
/*          
if groups := reg.FindStringSubmatch(txt); groups != nil && len(groups) > 1 {
                fmt.Println(groups[1:])
            }
            */
