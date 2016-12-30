package main

import (
    "fmt"
    "time"
    //"regexp"
    "strconv"
    "strings"
    //"math"
)
var input = 1364
var end Move = Move{ X: 31, Y: 39 }

type Move struct {
    X, Y int
}

func hasVisited(mv Move, moves []Move) bool {
    visited := false
    for _, m := range moves {
        if m.X == mv.X && m.Y == mv.Y { 
            visited = true
            break
        }
    }
    return visited || len(moves) > 120
}

func main() {
    startTime := time.Now()

    all := []Move{}
    navigatePt2(Move{ X: 1, Y: 1 }, []Move{}, &all, 50, input)
    //navigate(Move{ X: 1, Y: 1 }, end, []Move{}, input)

    mp := make(map[string]Move)

    for _, mv := range all {
        key := fmt.Sprintf("%v-%v", mv.X, mv.Y)
        if _,ok := mp[key]; !ok {
            mp[key] = mv
        }
    }

    fmt.Println("all within 50", len(mp))
    
    fmt.Println("Time", time.Since(startTime))
}

func drawMap(moves []Move, fav int){
    for y := 0; y < 42; y++ {
        for x := 0; x < 42; x++ {
            c := '.'
            move := Move{ X: x, Y: y }
            if isWall(move, fav) {
                c = '|'
            } else if hasVisited(move, moves) {
                c = 'O'
            }

            fmt.Print(string(c))
        }

        fmt.Print("\n")
    }
}

func navigatePt2(mv Move, moves []Move, allmoves *[]Move, max, fav int) {
    pmoves := getMoves(mv, fav)

    for _, m := range pmoves {
        if !hasVisited(m, moves) && len(moves) < max {
            newmoves := append(copymoves(moves), m)
            *allmoves = append(*allmoves, m)

            if len(moves)+1 == max {
                fmt.Println("we're there dude!", len(moves)+2) // add two to take the last move and starting location
                break
            }

            navigatePt2(m, newmoves, allmoves, max, fav)
        }
    }
}

func navigate(mv, goal Move, moves []Move, fav int) {
    pmoves := getMoves(mv, fav)

    for _, m := range pmoves {
        if m.X == goal.X && m.Y == goal.Y {
            fmt.Println("we're there dude!", len(moves)+1) // add one to take the last move
            break
        }
        if !hasVisited(m, moves) {
            newmoves := append(copymoves(moves), m)
            navigate(m, goal, newmoves, fav)
        }
    }
}

func copymoves(moves []Move) []Move {
    cp := make([]Move, len(moves))
    copy(cp, moves)
    return cp
}

func getMoves(mv Move, fav int) (moves []Move) {
    for _, m := range []Move{ Move{ X: mv.X+1, Y: mv.Y }, Move{ X: mv.X, Y: mv.Y+1 }, Move{ X: mv.X, Y: mv.Y-1 },  Move{ X: mv.X-1, Y: mv.Y } } {
        if m.X > -1 && m.Y > -1 && !isWall(m, fav){
            moves = append(moves, m)
        }
    }
    return
}

func isWall(mv Move, fav int) bool {
    x, y := mv.X, mv.Y
    t := x*x + 3*x + 2*x*y + y + y*y
    t += fav
    binary := strconv.FormatInt(int64(t), 2)
    return len(strings.Replace(binary, "0", "", -1)) % 2 != 0
}