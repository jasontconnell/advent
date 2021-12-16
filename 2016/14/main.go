package main

import (
	"crypto/md5"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/jasontconnell/advent/common"
)

type input = string
type output = int

type Key struct {
	Index int
	Value string
	Seed  string
	Char  rune
}

func main() {
	startTime := time.Now()

	in, err := common.ReadString(common.InputFilename(os.Args))
	if err != nil {
		log.Fatal(err)
	}

	p1 := part1(in)
	p2 := part2(in)

	w := common.TeeOutput(os.Stdout)
	fmt.Fprintln(w, "--2016 day 14 solution--")
	fmt.Fprintln(w, "Part 1:", p1)
	fmt.Fprintln(w, "Part 2:", p2)
	fmt.Println("Time", time.Since(startTime))
}

func part1(in input) output {
	val := process(in, false)
	return val
}

func part2(in input) output {
	val := process(in, true)
	return val
}

func process(in string, stretch bool) int {
	i := 0
	keys := []Key{}

	hashes := getMap(in, 0, 5000, stretch)

	for len(keys) < 64 {
		i++
		s := in + strconv.Itoa(i)
		if md5, ok := hashes[s]; ok {
			if rep, char := hasRepeat(md5, 3); rep {
				k := Key{Seed: s, Value: md5, Index: i, Char: char}

				for si := i + 1; si <= i+1001; si++ {
					s2 := in + strconv.Itoa(si)
					if md52, ok := hashes[s2]; ok {
						if hasCharRepeat(md52, 5, char) {
							keys = append(keys, k)
							break
						}
					}
				}
			}
		}

	}

	return i
}

func getMap(in string, start, count int, part2 bool) map[string]string {
	mp := make(map[string]string)
	for hk := 0; len(mp) < count; hk++ {
		s := in + strconv.Itoa(hk)
		md5 := MD5s(s)

		if part2 {
			for i := 0; i < 2016; i++ {
				md5 = MD5s(md5)
			}
		}
		if rep, _ := hasRepeat(md5, 3); rep {
			mp[s] = md5
		}
	}
	return mp
}

func hasCharRepeat(input string, l int, char rune) bool {
	substr := strings.Repeat(string(char), l)
	return strings.Contains(input, substr)
}

func hasRepeat(input string, l int) (val bool, char rune) {
	val = false
	for i := 0; i < len(input)-l+1 && !val; i++ {
		for j := 1; j < l; j++ {
			if input[i] != input[i+j] {
				break
			}
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
