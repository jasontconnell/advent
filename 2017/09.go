package main

import (
    "bufio"
    "fmt"
    "os"
    "time"
)
var input = "09.txt"

type Group struct {
    Inner []*Group
    Score int
    Garbage int
}

func main() {
    startTime := time.Now()

    f, err := os.Open(input)

    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
    
    scanner := bufio.NewScanner(f)

    var txt string

    if scanner.Scan() {
        txt = scanner.Text()
    }

    groups := getGroups(txt, 0)
    garbage := getGarbage(txt)
    score := getScore(groups)

    fmt.Println("Score", score, "Garbage", garbage)

    fmt.Println("Time", time.Since(startTime))
}

func getGarbage(line string) int {
    g := 0
    for i := 0; i < len(line); i++ {
        s, e, skipped := nextGarbage(line, i)

        if s != -1 && e != -1 {
            g += e-s-skipped - 1
            i = e
        }
    }
    return g
}

func getScore(groups []*Group) int {
    if len(groups) == 0 { return 0 }
    c := 0
    for _, grp := range groups {
        c += grp.Score + getScore(grp.Inner)
    }
    return c
}

func getGroups(line string, level int) []*Group {
    groups := []*Group{}
    done := false
    for i := 0; i < len(line) && !done; i++ {

        op, cl := nextGroup(line, i)
        var g *Group

        if op != -1 && cl != -1 && cl-op > 1 {
            innertext := line[op+1:cl]
            innergroups := getGroups(innertext, level+1)

            g = &Group {
                Inner: innergroups,
                Score: level+1,
            }
        } else if cl-op == 1 {
            g = &Group {
                Inner: nil,
                Score: level+1,
            }
        }

        if g != nil {
            groups = append(groups, g)
        }

        if op == -1 && cl == -1 {
            done = true
        } else {
            i = cl
        }
    }

    return groups
}

func nextGarbage(line string, start int) (int, int, int) {
    s,sk1 := seekNext(line, '<', start)
    e,sk2 := seekNext(line, '>', start)
    return s, e, sk2 - sk1
}

func nextGroup(line string, start int) (int, int) {
    s,_ := seekNext(line, '{', start)
    e,_ := seekNext(line, '}', start)
    return s, e
}

func seekNext(line string, seek rune, start int) (int, int) {
    c := 0
    val := -1
    grbg := false
    cncl := false
    found := false
    skipped := 0

    for i := start; i < len(line) && !found; i++ {
        switch line[i] {
            case '!': 
                skipped++
                cncl = !cncl
            case '<': 
                if !cncl {
                    grbg = true

                    if seek == '<' {
                        found = true
                        val = i
                    }
                } else {
                    skipped++
                    cncl = false
                }
                
            case '>': 
                if !cncl {
                    grbg = false
                    if seek == '>' {
                        found = true
                        val = i
                    }
                } else {
                    skipped++
                    cncl = false
                }
            case '{': 
                if !grbg && !cncl {
                    c++
                    if seek == '{' {
                        found = true
                        val = i
                    }
                } else if cncl {
                    skipped++
                    cncl = false
                }


            case '}':
                if !grbg && !cncl {
                    c--
                    
                    if c == 0 && seek == '}' {
                        val = i
                        found = true
                    }
                } else if cncl {
                    skipped++
                    cncl = false
                }
            default:
                if cncl {
                    skipped++
                }
                cncl = false
        }
    }
    return val, skipped
}