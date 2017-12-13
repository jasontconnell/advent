package main

import (
    "bufio"
    "fmt"
    "os"
    "time"
    "regexp"
    //"strconv"
    "strings"
    //"math"
)
var input = "12.txt"

var reg = regexp.MustCompile("^([0-9]+) <-> (.*)$")


type Program struct {
    ID string
    Piped []*Program
    PipedIDs []string
}

func main() {
    startTime := time.Now()

    f, err := os.Open(input)

    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
    
    scanner := bufio.NewScanner(f)


    programs := []*Program{}

    for scanner.Scan() {
        var txt = scanner.Text()
        p := getProgram(txt)
        if p != nil {
            programs = append(programs, p)
        }
    }

    mapPrograms(programs)

    c := getConnected(programs, "0")
    g := getGroups(programs)

    fmt.Println("programs that connect to 0    ", len(c))
    fmt.Println("total number of groups        ", g)

    fmt.Println("Time", time.Since(startTime))
}

func getGroups(programs []*Program) int {
    startTime := time.Now()

    // find out which group each program belongs
    groups := make(map[string]string)
    vmap := make(map[string]bool)
    cp := make([]*Program, len(programs))

    copy(cp, programs)

    for _, p := range programs {
        if _, ok := groups[p.ID]; !ok {
            r := getConnected(cp, p.ID)

            for _, rp := range r {
                vmap[rp.ID] = true
                if _, ok := groups[rp.ID]; !ok {
                    groups[rp.ID] = p.ID
                }
            }

            for i := len(cp)-1; i >= 0; i-- {
                sp := cp[i]
                if _, ok := vmap[sp.ID]; ok {
                    cp = append(cp[:i], cp[i+1:]...)
                }
            }
        }
    }

    fmt.Println("Finished determining groups after", time.Since(startTime))


    // determine unique groups
    m := make(map[string]string)
    for _, v := range groups {
        m[v] = v
    }

    
    return len(m)
}

func getConnected(programs []*Program, id string) []*Program {
    connected := []*Program{}

    for _, p := range programs {
        vmap := make(map[string]bool)

        if connects(p, id, vmap) {
            connected = append(connected, p)
        }
    }
    return connected
}

func connects(p *Program, id string, vmap map[string]bool) bool {
    if p.ID == id {
        return true
    }

    val := false
    for _, piped := range p.Piped {
        _, visited := vmap[piped.ID]
        vmap[piped.ID] = true

        if !visited {
            r := connects(piped, id, vmap)
            if r {
                val = true
                break
            }
        }
    }

    return val
}

func mapPrograms(programs []*Program) {
    pmap := make(map[string]*Program)
    for _, p := range programs {
        pmap[p.ID] = p
    }

    for _, p := range programs {
        for _, id := range p.PipedIDs {
            piped, ok := pmap[id]
            if !ok {
                fmt.Println("couldn't find", id, "program", p.ID, "piped ids", p.PipedIDs)
                break
            }

            p.Piped = append(p.Piped, piped)
        }
    }
}

func getProgram(line string) *Program {
    var p *Program
    if groups := reg.FindStringSubmatch(line); groups != nil && len(groups) > 1 {
        id := groups[1]
        ids := strings.Split(strings.Replace(groups[2], " ", "", -1), ",")

        p = &Program{ ID: id, PipedIDs: ids, Piped: []*Program{} }

    }

    return p
}


// reg := regexp.MustCompile("-?[0-9]+")
/*          
if groups := reg.FindStringSubmatch(txt); groups != nil && len(groups) > 1 {
                fmt.Println(groups[1:])
            }
            */
