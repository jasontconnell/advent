package main

import (
    "fmt"
    "time"
    "crypto/md5"
    "strconv"
    "strings"
)
var input = "jlmsuwbz"

type Key struct {
    Index int
    Value string
    Seed string
    Char rune
}

func (key Key) String() string {
    return fmt.Sprintf("%v) Value: %v with Seed: %v .  Repeated char: %v", key.Index, key.Value, key.Seed, string(key.Char))
}

func getMap(start, count int, part2 bool) map[string]string {
    mp := make(map[string]string)
    for hk := 0; len(mp) < count; hk++ {
        s := input + strconv.Itoa(hk)
        md5 := MD5s(s)

        if part2 {
            for i := 0; i < 2016; i++ {
                md5 = MD5s(md5)
            }
        }
        if rep,_ := hasRepeat(md5, 3); rep {
            mp[s] = md5
        }
    }
    return mp
}

func main() {
    startTime := time.Now()
    i := 0
    keys := []Key{}
    part2 := true

    hashes := getMap(0, 5000, part2)

    for len(keys) < 64 {
        s := input + strconv.Itoa(i)
        if md5, ok := hashes[s]; ok {
            if rep, char := hasRepeat(md5, 3); rep {
                k := Key{ Seed: s, Value: md5, Index: i, Char: char }

                for si := i+1; si <= i+1001; si++ {
                    s2 := input + strconv.Itoa(si)
                    if md52,ok := hashes[s2]; ok {
                        if hasCharRepeat(md52, 5, char) {
                            keys = append(keys, k)
                            break
                        }
                    }
                }
            }
        }

        i++
    }

    for i, k := range keys {
        fmt.Println(i+1, k)
    }

    fmt.Println("Time", time.Since(startTime))
}

func hasCharRepeat(input string, l int, char rune) bool {
    substr := strings.Repeat(string(char), l)
    return strings.Contains(input, substr)
}

func hasRepeat(input string, l int) (val bool, char rune ) {
    val = false
    for i := 0; i < len(input) - l + 1 && !val; i++ {
        for j := 1; j < l; j++ {
            if input[i] != input[i+j] { break }
            if j == l-1 { 
                char = rune(input[i])
                val = true //!hasCharRepeat(input, l+1, char)
            }
        }
    }
    return
}

func MD5(content []byte) string {
    sum := md5.Sum(content)
    return fmt.Sprintf("%x", sum)
}

func MD5s(content string) string {
    return MD5([]byte(content))
}