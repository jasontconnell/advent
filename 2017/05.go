package main

import (
    "bufio"
    "fmt"
    "os"
    "time"
    "strconv"
)
var input = "05.txt"

func main() {
    startTime := time.Now()

    f, err := os.Open(input)

    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
    
    inst := []int{}
    scanner := bufio.NewScanner(f)

    for scanner.Scan() {
        var txt = scanner.Text()
        n,err := strconv.Atoi(txt)
        if err != nil {
            fmt.Println("parsing", err)
            break
        }

        inst = append(inst, n)
    }

    p1 := make([]int, len(inst))
    p2 := make([]int, len(inst))

    copy(p1, inst)
    steps := process(p1)

    copy(p2, inst)
    steps2 := processPart2(p2)

    fmt.Println("steps", steps)
    fmt.Println("steps part 2", steps2)

    fmt.Println("Time", time.Since(startTime))
}

func process(inst []int) int {
    i := 0
    steps := 0
    max := len(inst)
    for i > -1 && i < max {
        steps++
        cur := inst[i]

        inst[i] = cur+1
        i = i + cur
    }

    return steps
}

func processPart2(inst []int) int {
    i := 0
    steps := 0
    max := len(inst)
    for i > -1 && i < max {
        steps++
        cur := inst[i]
        jmp := 1
        if cur >= 3 {
            jmp = -1
        }

        inst[i] = cur+jmp
        i = i + cur
    }

    return steps
}