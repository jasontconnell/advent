package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/jasontconnell/advent/common"
)

var inputFilename = "input.txt"

func main() {
	startTime := time.Now()

	line, err := common.ReadString(inputFilename)
	if err != nil {
		log.Fatal(err)
	}

	p1 := part1(line)
	p2 := part2(line)

	fmt.Println("Part 1:", p1)
	fmt.Println("Part 2:", p2)

	fmt.Println("Time", time.Since(startTime))
}

func part1(in string) int {
	return findHash(string(in), "00000")
}

func part2(in string) int {
	return findHash(string(in), "000000")
}

func findHash(seed string, prefix string) int {
	found := false
	lastTest := 0
	for i := 0; !found; i++ {
		test := seed + strconv.Itoa(i)

		hex := getMD5Hex(test)

		if found = strings.HasPrefix(hex, prefix); found {
			lastTest = i
		}
	}
	return lastTest
}

func getMD5Hex(str string) string {
	h := md5.New()
	io.WriteString(h, str)
	hex := fmt.Sprintf("%x", h.Sum(nil))
	return hex
}
