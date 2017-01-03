package main

import (
    "fmt"
    "time"
    "crypto/md5"
    //"regexp"
    //"strconv"
    //"strings"
    //"math/rand"
)
var input = "njfxhljp"

type Point struct {
    X, Y int
}

type State struct {
    Point Point
    Path string
}

func main() {
    startTime := time.Now()
    
    start := Point{ X: 0, Y: 0 }
    goal := Point{ X: 3, Y: 3 }

    solves := solve(start, goal, input)
    
    fmt.Println(solves[0].Path)
    fmt.Println(len(solves[len(solves)-1].Path))

    fmt.Println("Time", time.Since(startTime))
}

func solve(point, goal Point, hash string) []State {
    queue := []State{}
    queue = append(queue, State{ Point: point, Path: "" })
    solves := []State{}

    for len(queue) > 0 {
        state := queue[0]
        queue = queue[1:]

        moves := getDoors(state.Point, hash, state.Path)

        for _, mv := range moves {
            point := state.Point
            if mv == 'U' {
                point.Y--
            } else if mv == 'D' {
                point.Y++
            } else if mv == 'L' {
                point.X--
            } else if mv == 'R' {
                point.X++
            }
            path := state.Path + string(mv)

            if point.X == goal.X && point.Y == goal.Y {
                solves = append(solves, State{ Point: point, Path: path })
                continue
            }

            queue = append(queue, State{ Point: point, Path: path })
        }
    }

    return solves
}

func getDoors(point Point, hash, path string) (moves []rune) {
    md5 := MD5s(hash+path)
    if isOpen(rune(md5[0])) && point.Y > 0 { moves = append(moves, 'U') }
    if isOpen(rune(md5[1])) && point.Y < 3 { moves = append(moves, 'D') }
    if isOpen(rune(md5[2])) && point.X > 0 { moves = append(moves, 'L') }
    if isOpen(rune(md5[3])) && point.X < 3 { moves = append(moves, 'R') }

    return moves
}

func isOpen(char rune) bool {
    return char == 'b' || char == 'c' || char == 'd' || char == 'e' || char == 'f'
}

func MD5(content []byte) string {
    sum := md5.Sum(content)
    return fmt.Sprintf("%x", sum)
}

func MD5s(content string) string {
    return MD5([]byte(content))
}