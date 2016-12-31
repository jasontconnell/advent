package main

import (
    "fmt"
    "time"
    "strings"
)
var input = "11101000110010100"

func main() {
    startTime := time.Now()

    disclen := 35651584 //  part 1 is 272
    cp := input

    for len(cp) < disclen {
        cp = dragonCurve(cp)
    }

    disccp := string(cp[:disclen])
    sum, even := checksum(disccp)

    for even {
        sum, even = checksum(sum)
    }

    fmt.Println(sum, even)
    fmt.Println("Time", time.Since(startTime))
}

func dragonCurve(str string) string {
    cp := reverse(str)
    cp = strings.Replace(cp, "0", "_", -1)
    cp = strings.Replace(cp, "1", "0", -1)
    cp = strings.Replace(cp, "_", "1", -1)

    return str + "0" + cp
}

func reverse(str string) string {
    n := len(str)
    runes := make([]rune, n)
    for i := 0; i < n; i++ {
        runes[i] = rune(str[i])
    }

    for i := 0; i < n / 2; i++ {
        runes[i], runes[n-i-1] = runes[n-i-1], runes[i]
    }
    return string(runes)
}

func checksum(str string) (string, bool) {
    sum := make([]rune, len(str)/2+2)
    n := 0
    for i := 0; i < len(str)-1; i+=2 {
        c := '1'
        if str[i] != str[i+1] { c = '0'}
        sum[n] = c
        n++
    }

    return string(sum[:n]), len(sum) % 2 == 0
}