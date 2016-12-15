package main

import (
    "bufio"
    "fmt"
    "os"
    "time"
    "regexp"
    "strconv"
    "strings"
    //"math"
)
var input = "9.txt"
var reg = regexp.MustCompile("(\\d+)x(\\d+)\\)")

func main() {
    startTime := time.Now()
    if f, err := os.Open(input); err == nil {
        scanner := bufio.NewScanner(f)

        output := 0
        output2 := 0

        for scanner.Scan() {
            var txt = scanner.Text()
            
            output = explode(txt, false)
            output2 = explode(txt, true)
        }

        fmt.Println("len output", output)
        fmt.Println("len output 2", output2)
    }

    fmt.Println("Time", time.Since(startTime))
}

func explode(txt string, dec bool) int {
    result := 0
    for i := 0; i < len(txt); i++ {
        cs := string(txt[i])

        if cs == "(" {
            m := i + 10
            if i + m > len(txt) { m = len(txt) - i }
            if groups := reg.FindStringSubmatch(string(txt[i:i+m])); groups != nil && len(groups) > 1 {
                chars, _ := strconv.Atoi(groups[1])
                repeat, _ := strconv.Atoi(groups[2])

                i += len(groups[0])
                rep := txt[i+1:i+chars+1]

                if dec && strings.Contains(rep, "(") {
                    result += repeat * explode(rep, dec)
                } else {
                    result += repeat * len(rep)
                }

                i += chars
            }
        } else {
            result ++
        }
    }
    return result
}