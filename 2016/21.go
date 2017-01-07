package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"time"
	//"strings"
	//"math"
)

var input = "21.txt"
var pwd = "abcdefgh"
var pwd2 = "fbgdceah"

//var pwd = "abcde"

type Instruction struct {
	Cmd    string
	Param  string
	X, Y   int
	CX, CY rune
}

func main() {
	startTime := time.Now()

	swapreg := regexp.MustCompile("^swap (position|letter) ([a-z0-9]+) with (position|letter) ([a-z0-9]+)$")
	rotnreg := regexp.MustCompile("^rotate (left|right) ([0-9]) steps?$")
	rotpreg := regexp.MustCompile("^rotate based on position of letter ([a-z])$")
	revreg := regexp.MustCompile("^reverse positions ([0-9]) through ([0-9])$")
	mvreg := regexp.MustCompile("^move position ([0-9]) to position ([0-9])$")

	instructions := []Instruction{}

	if f, err := os.Open(input); err == nil {
		scanner := bufio.NewScanner(f)

		for scanner.Scan() {
			var txt = scanner.Text()
			if groups := swapreg.FindStringSubmatch(txt); groups != nil && len(groups) > 1 {
				x, y := -1, -1
				cx, cy := ' ', ' '
				if groups[1] == "position" {
					x, _ = strconv.Atoi(groups[2])
					y, _ = strconv.Atoi(groups[4])
				} else {
					cx = rune(groups[2][0])
					cy = rune(groups[4][0])
				}
				instructions = append(instructions, Instruction{Cmd: "swap", Param: groups[1], X: x, Y: y, CX: cx, CY: cy})

			} else if groups := rotnreg.FindStringSubmatch(txt); groups != nil && len(groups) > 1 {
				x, _ := strconv.Atoi(groups[2])
				if groups[1] == "right" {
					x = -x
				}
				instructions = append(instructions, Instruction{Cmd: "rotateN", X: x, Y: -1})

			} else if groups := rotpreg.FindStringSubmatch(txt); groups != nil && len(groups) > 1 {
				cx := rune(groups[1][0])

				instructions = append(instructions, Instruction{Cmd: "rotateC", CX: cx})
			} else if groups := revreg.FindStringSubmatch(txt); groups != nil && len(groups) > 1 {
				x, _ := strconv.Atoi(groups[1])
				y, _ := strconv.Atoi(groups[2])
				instructions = append(instructions, Instruction{Cmd: "reverse", X: x, Y: y})

			} else if groups := mvreg.FindStringSubmatch(txt); groups != nil && len(groups) > 1 {
				x, _ := strconv.Atoi(groups[1])
				y, _ := strconv.Atoi(groups[2])
				instructions = append(instructions, Instruction{Cmd: "move", X: x, Y: y})
			}
		}
	}

	fmt.Println(process(pwd, instructions))

	combos := Permutate([]rune(pwd2))

	for _, c := range combos {
		s := process(string(c), instructions)
		if s == pwd2 {
			fmt.Println("unscrambled", string(c))
		}
	}

	fmt.Println("Time", time.Since(startTime))
}

func process(password string, instructions []Instruction) string {
	str := make([]rune, len(password))
	copy(str, []rune(password))

	for _, instr := range instructions {
		switch instr.Cmd {
		case "swap":
			if instr.Param == "position" {
				swapN(str, instr.X, instr.Y)
			} else {
				swapC(str, instr.CX, instr.CY)
			}
			break
		case "rotateN":
			rotateN(str, instr.X)
			break
		case "rotateC":
			rotateC(str, instr.CX)
			break
		case "move":
			move(str, instr.X, instr.Y)
			break
		case "reverse":
			reverse(str, instr.X, instr.Y)
			break
		}
	}

	return string(str)
}

func move(str []rune, x, y int) {
	cp := []rune{}
	for i := 0; i < len(str); i++ {
		if i != x {
			cp = append(cp, str[i])
		}
	}

	cp = append(cp, str[x])
	for i := len(cp) - 1; i > y; i-- {
		cp[i], cp[i-1] = cp[i-1], cp[i]
	}

	copy(str, cp)
}

func rotateN(str []rune, x int) {
	if x == 0 {
		return
	}
	if x < 0 {
		x = len(str) + x
	}
	cp := make([]rune, len(str))
	for i := 0; i < len(str); i++ {
		idx := (i + x) % len(str)
		cp[i] = str[idx]
	}

	copy(str, cp)
}

func rotateC(str []rune, cx rune) {
	idx := index(str, cx)
	rots := -idx - 1
	if idx > 3 {
		rots--
	}
	rots = rots % len(str)
	rotateN(str, rots)
}

func swapN(str []rune, x, y int) {
	str[x], str[y] = str[y], str[x]
}

func swapC(str []rune, cx, cy rune) {
	x, y := index(str, cx), index(str, cy)
	swapN(str, x, y)
}

func index(str []rune, c rune) int {
	for i, cs := range str {
		if cs == c {
			return i
		}
	}
	return -1
}

func reverse(str []rune, x, y int) {
	runes := make([]rune, y+1-x)
	for i := x; i < y+1; i++ {
		runes[i-x] = str[i]
	}

	n := len(runes)
	for i := 0; i < n/2; i++ {
		runes[i], runes[n-i-1] = runes[n-i-1], runes[i]
	}

	str = append(str[:x], append(runes, str[y+1:]...)...)
}

func Permutate(str []rune) [][]rune {
	var ret [][]rune

	if len(str) == 2 {
		ret = append(ret, []rune{str[0], str[1]})
		ret = append(ret, []rune{str[1], str[0]})
	} else {

		for i := 0; i < len(str); i++ {
			strc := make([]rune, len(str))
			copy(strc, str)

			t := strc[i]
			sh := append(strc[:i], strc[i+1:]...)
			perm := Permutate(sh)

			for _, p := range perm {
				p = append([]rune{t}, p...)
				ret = append(ret, p)
			}
		}
	}

	return ret
}
