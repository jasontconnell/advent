package main

import (
    //"bufio"
    "fmt"
    //"os"
    "time"
    //"regexp"
    "strconv"
    "strings"
    //"math"
    "crypto/md5"
)
var input = "ojvtpuvg"

func main() {
    startTime := time.Now()

    i := 0

    done := false

    pwd := ""
    pwd2 := "        "

    for !done {
        tmp := input + strconv.Itoa(i)
        hash := MD5s(tmp)

        if strings.HasPrefix(hash, "00000") {
            if len(pwd) < len(input){
                pwd += string(hash[5])
            }

            if p2index,err := strconv.Atoi(string(hash[5])); err == nil && p2index < 8 && pwd2[p2index] == ' ' {
                
                p2char := hash[6]
                
                pwd2 = pwd2[:p2index] + string(p2char) + pwd2[p2index+1:]

                fmt.Println(pwd2)
            }
        }

        done = len(pwd) == len(input) && strings.Count(pwd2, " ") == 0
        i++
    }

    fmt.Println(pwd)
    fmt.Println(pwd2)
    fmt.Println("Time", time.Since(startTime))
}

func MD5(content []byte) string {
    sum := md5.Sum(content)
    return fmt.Sprintf("%x", sum)
}

func MD5s(content string) string {
    return MD5([]byte(content))
}

// reg := regexp.MustCompile("-?[0-9]+")
/*          
if groups := reg.FindStringSubmatch(txt); groups != nil && len(groups) > 1 {
                fmt.Println(groups[1:])
            }
            */
