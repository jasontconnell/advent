package main

import (
    "fmt"
    "time"
    //"regexp"
    //"strconv"
    "strings"
    //"math"
)
var input = ".^..^....^....^^.^^.^.^^.^.....^.^..^...^^^^^^.^^^^.^.^^^^^^^.^^^^^..^.^^^.^^..^.^^.^....^.^...^^.^."
var rows = 40
var safe byte = '.'
var trap byte = '^'

func main() {
    startTime := time.Now()
    
    prev := input
    count := 0
    part2 := true
    if part2 {
        rows = 400000
    }
    

    for i := 0; i < rows; i++ {
        count += strings.Count(prev, string(safe))
        prev = transform(prev)
    }

    fmt.Println("count", count)

    fmt.Println("Time", time.Since(startTime))
}

func transform(str string) string {
    ns := ""
    for i := 0; i < len(str); i++ {
        nc := '^'
        if isSafe(i, str) {
            nc = '.'
        }
        ns += string(nc)
    }
    return ns
}

func isSafe(index int, prev string) bool {
    tri := getThreePrev(index, prev)
    //fmt.Println(prev, index, tri)
    istrap := tri[0] == trap && tri[1] == trap && tri[2] != trap ||
            tri[0] != trap && tri[1] == trap && tri[2] == trap ||
            tri[0] == trap && tri[1] != trap && tri[2] != trap ||
            tri[0] != trap && tri[1] != trap && tri[2] == trap

    return !istrap
}

func getThreePrev(index int, prev string) string {
    tri := ""
    if index == 0 { 
        tri = "." + string(prev[:2]) 
    } else if index == len(prev) - 1 { 
        tri = string(prev[len(prev)-2:]) + "." 
    } else {
        tri = string(prev[index-1:index+2])
    }

    return tri
}