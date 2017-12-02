package main

import (
    "bufio"
    "fmt"
    "os"
    "time"
    "regexp"
    "strconv"
    "sort"
)
var input = "02.txt"

var reg = regexp.MustCompile("-?([0-9]+)")

func main() {
    startTime := time.Now()

    f, err := os.Open(input)

    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
    
    scanner := bufio.NewScanner(f)
    lines := []string{}

    for scanner.Scan() {
        var txt = scanner.Text()
        lines = append(lines, txt)
    }

    ch := getChecksum(lines)
    d := getDivisionResult(lines)

    fmt.Println("checksum", ch)
    fmt.Println("division", d)

    fmt.Println("Time", time.Since(startTime))
}

func getChecksum(lines []string) int {
    sum := 0

    for _, line := range lines {
        min, max := getMinMax(line)
        sum += (max-min)
    }

    return sum
}

func getDivisionResult(lines []string) int {
    res := 0

    for _, line := range lines {
        div := getDivision(line)
        res += div
    }

    return res
}

func getMinMax(line string) (min int, max int){
    min = 100000
    max = 0
    if groups := reg.FindAllStringSubmatch(line, -1); groups != nil && len(groups) > 1 {
        for _, g := range groups {
            val, err := strconv.Atoi(g[0])
            if err != nil {
                fmt.Println("parsing", err)
            }

            if val < min {
                min = val
            }

            if val > max {
                max = val
            }
        }
        
    }

    return min, max
}

func getDivision(line string) (div int){
    nums := []int{}
    if groups := reg.FindAllStringSubmatch(line, -1); groups != nil && len(groups) > 1 {
        for _, g := range groups {
            val, err := strconv.Atoi(g[0])
            if err != nil {
                fmt.Println("parsing", err)
            }

            nums = append(nums, val)
        }
    }

    sort.Ints(nums)

    for i := 0; i < len(nums); i++ {
        for j := 0; j < len(nums); j++ {
            if i == j {
                continue
            }
            if nums[i] % nums[j] == 0 {
                return nums[i] / nums[j]
            }
        }
    }

    return 0
}