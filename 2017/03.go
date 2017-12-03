package main

import (
    "fmt"
    "time"
    "math"
)

var input = 368078

func main() {
    startTime := time.Now()

    d := dist(input)
    fmt.Println("distance", d)

    fmt.Println("Time", time.Since(startTime))
}

func ring(i int) int {
    if i == 1 {
        return 1
    }
    s := 1
    itr := 1

    for (s*s) < i {
        s += 2
        itr++
    }

    return s-2
}

func coords(i int) (int, int, int) {
    sq := ring(i)
    start := sq*sq + 1
    max := sq+1
    startx, starty := max, max-1

    diff := i - start
    turns := 0

    extra := diff
    if diff > max - 1 {
        for j := 0; j < 4; j++ {
            m := max
            if j == 0 { m = m - 1}

            if extra - m >= 0 {
                extra = extra - m
                turns ++
            }
        }
    }

    x, y := 0, 0
    switch turns {
        case 0: 
            x = max
            y = starty - extra
        case 1: 
            x = startx - extra
            y = 0
        case 2: 
            x = 0
            y = extra
        case 3: 
            x = extra
            y = max
    }

    return x, y, max
}

func dist(i int) float64 {
    x, y, max := coords(i)
    one := max / 2 // position of 1

    return math.Abs(float64(x - one)) + math.Abs(float64(y - one))
}