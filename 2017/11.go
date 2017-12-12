package main

import (
    "bufio"
    "fmt"
    "os"
    "time"
    "strings"
)
var input = "11.txt"

type Coord struct {
    X, Y int
}

func (c Coord) String() string {
    s := fmt.Sprintf("(%v, %v)", c.X, c.Y)
    return s
}

func (c *Coord) Add(d Coord) Coord {
    e := *c
    e.X += d.X
    e.Y += d.Y

    return e
}

type Hex struct {
    Coords Coord
    NE *Hex
    N  *Hex
    NW *Hex
    SE *Hex
    S  *Hex
    SW *Hex
}

func main() {
    startTime := time.Now()

    f, err := os.Open(input)

    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
    
    scanner := bufio.NewScanner(f)
    var moves []string

    for scanner.Scan() {
        var txt = scanner.Text()
        moves = strings.Split(txt, ",")

        home := makeHex(Coord{ X: 0, Y: 0 })
        
        end, farthest := traverse(moves, home)
        if end != nil {
            dist := traverseTo(end, home)

            fmt.Println("Distance from endpoint to home:  ", dist)
            fmt.Println("Farthest distance to home:       ", farthest)
        }
    }

    fmt.Println("Time", time.Since(startTime))
}

func traverse(moves []string, start *Hex) (*Hex, int) {
    ptr := start
    max := 0
    for _, mv := range moves {
        ptr = getMove(mv, ptr)
        d := traverseTo(ptr, start)
        if d > max {
            max = d
        }
    }
    return ptr, max
}

func traverseTo(start, end *Hex) int {
    ptr := start
    d := 0

    for ptr.Coords.X != end.Coords.X || ptr.Coords.Y != end.Coords.Y {
        mv := ""
        n := ptr.Coords.X % 2 != 0

        if ptr.Coords.Y > end.Coords.Y {
            mv = "s"
        } else if ptr.Coords.Y < end.Coords.Y {
            mv = "n"
        } else if ptr.Coords.Y == end.Coords.Y {
            mv = "s"
            if n {
                mv = "n"
            }
        }


        if ptr.Coords.X > end.Coords.X {
            mv += "w"
        } else if ptr.Coords.X < end.Coords.X {
            mv += "e"
        }

        ptr = getMove(mv, ptr)

        d++
    }

    return d
}

func makeHex(coords Coord) *Hex {
    hex := new(Hex)
    hex.Coords = coords
    return hex
}

func getMove(dir string, hex *Hex) *Hex {
    var dest *Hex
    var delta Coord
    n := hex.Coords.X % 2 != 0

    switch dir {
        case "n": 
            if hex.N == nil {
                delta.Y = 1
                hex.N = makeHex(hex.Coords.Add(delta))
                hex.N.S = hex
            }
            dest = hex.N
        case "s":
            if hex.S == nil {
                delta.Y = -1
                hex.S = makeHex(hex.Coords.Add(delta))
                hex.S.N = hex
            }
            dest = hex.S
        case "ne":
            if hex.NE == nil {
                delta.X = 1
                if !n {
                    delta.Y = 1
                }
                hex.NE = makeHex(hex.Coords.Add(delta))
                hex.NE.SW = hex
            }
            
            dest = hex.NE
        case "nw":
            if hex.NW == nil {
                delta.X = -1
                if !n {
                    delta.Y = 1
                }

                hex.NW = makeHex(hex.Coords.Add(delta))
                hex.NW.SE = hex
            }
            
            dest = hex.NW
        case "se":
            if hex.SE == nil {
                delta.X = 1
                if n {
                    delta.Y = -1
                }

                hex.SE = makeHex(hex.Coords.Add(delta))
                hex.SE.NW = hex
            }
            dest = hex.SE

        case "sw":
            if hex.SW == nil {
                delta.X = -1
                if n {
                    delta.Y = -1
                }
                hex.SW = makeHex(hex.Coords.Add(delta))
                hex.SW.NE = hex
            }
            dest = hex.SW
    }
    return dest
}
