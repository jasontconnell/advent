package main

import (
	"crypto/md5"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/jasontconnell/advent/common"
)

var inputFilename = "input.txt"

type input = string
type output = string

func main() {
	startTime := time.Now()

	in, err := common.ReadString(inputFilename)
	if err != nil {
		log.Fatal(err)
	}

	p1 := part1(in)
	p2 := part2(in)

	fmt.Println("Part 1:", p1)
	fmt.Println("Part 2:", p2)

	fmt.Println("Time", time.Since(startTime))
}

func part1(in input) output {
	return getPassword(in)
}

func part2(in input) output {
	return getPassword2(in)
}

func getPassword(in input) output {
	done := false
	i := 0
	pwd := ""
	for !done {
		tmp := in + strconv.Itoa(i)
		hash := MD5s(tmp)
		if strings.HasPrefix(hash, "00000") {
			if len(pwd) < len(in) {
				pwd += string(hash[5])
			}
		}

		done = len(pwd) == len(in)
		i++
	}
	return pwd
}

func getPassword2(in input) output {
	done := false
	i := 0
	pwd := strings.Repeat(" ", len(in))
	for !done {
		tmp := in + strconv.Itoa(i)
		hash := MD5s(tmp)
		if strings.HasPrefix(hash, "00000") {
			if idx, err := strconv.Atoi(string(hash[5])); err == nil && idx < len(pwd) && pwd[idx] == ' ' {
				ch := hash[6]
				pwd = pwd[:idx] + string(ch) + pwd[idx+1:]
			}
		}

		done = strings.Count(pwd, " ") == 0
		i++
	}
	return pwd
}

func MD5(content []byte) string {
	sum := md5.Sum(content)
	return fmt.Sprintf("%x", sum)
}

func MD5s(content string) string {
	return MD5([]byte(content))
}
