package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
	//"regexp"
	//"strconv"
	//"strings"
	//"math"
)

type passport struct {
	Fields map[string]string
}

func newPassport() *passport {
	pp := new(passport)
	pp.Fields = make(map[string]string)
	return pp
}

var input = "04.txt"

func main() {
	startTime := time.Now()

	f, err := os.Open(input)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	scanner := bufio.NewScanner(f)
	lines := []string{}

	for scanner.Scan() {
		var txt = scanner.Text()
		lines = append(lines, txt)
	}

	passports := getPassports(lines)
	valid := getValid(passports)

	fullValid := fullValidate(valid)

	fmt.Println("Valid: ", len(valid))
	fmt.Println("Full Valid: ", len(fullValid))

	fmt.Println("Time", time.Since(startTime))
}

var propreg *regexp.Regexp = regexp.MustCompile("([a-z]+):([^ ]+)")

func getValid(passports []*passport) []*passport {
	mightHave := "cid"
	expFields := 8
	vpps := []*passport{}

	for _, pp := range passports {
		valid := len(pp.Fields) == expFields
		if _, ok := pp.Fields[mightHave]; !ok {
			valid = len(pp.Fields) == 7
		}

		if valid {
			vpps = append(vpps, pp)
		}
	}
	return vpps
}

func fullValidate(passports []*passport) []*passport {
	funcs := map[string]func(string) bool{
		"byr": validbyr,
		"iyr": validiyr,
		"eyr": valideyr,
		"hgt": validhgt,
		"hcl": validhcl,
		"ecl": validecl,
		"pid": validpid,
		"cid": validcid,
	}

	fullValid := []*passport{}
	for _, pp := range passports {
		valid := true
		for prop, value := range pp.Fields {
			valid = valid && funcs[prop](value)
		}
		if valid {
			fullValid = append(fullValid, pp)
		}
	}
	return fullValid
}

func validbyr(value string) bool {
	i, _ := strconv.Atoi(value)
	return len(value) == 4 && i >= 1920 && i <= 2002
}

func validiyr(value string) bool {
	i, _ := strconv.Atoi(value)
	return len(value) == 4 && i >= 2010 && i <= 2020
}

func valideyr(value string) bool {
	i, _ := strconv.Atoi(value)
	return len(value) == 4 && i >= 2020 && i <= 2030
}

func validhgt(value string) bool {
	if strings.HasSuffix(value, "cm") {
		i, _ := strconv.Atoi(strings.TrimSuffix(value, "cm"))
		return i >= 150 && i <= 193
	} else if strings.HasSuffix(value, "in") {
		i, _ := strconv.Atoi(strings.TrimSuffix(value, "in"))
		return i >= 59 && i <= 76
	}
	return false
}

var hclreg *regexp.Regexp = regexp.MustCompile("^#([0-9a-f]){6}$")

func validhcl(value string) bool {
	return hclreg.MatchString(value)
}

var cmap map[string]bool = map[string]bool{
	"amb": true,
	"blu": true,
	"brn": true,
	"gry": true,
	"grn": true,
	"hzl": true,
	"oth": true,
}

func validecl(value string) bool {
	_, ok := cmap[value]
	return ok
}

func validpid(value string) bool {
	_, err := strconv.Atoi(value)
	return len(value) == 9 && err == nil
}

func validcid(value string) bool {
	return true
}

func getPassports(lines []string) []*passport {
	pp := newPassport()
	pps := []*passport{}

	for i, line := range lines {
		if line != "" {
			matches := propreg.FindAllStringSubmatch(line, -1)
			for _, m := range matches {
				pp.Fields[m[1]] = m[2]
			}
		}

		if line == "" || i == len(lines)-1 {
			pps = append(pps, pp)
			pp = newPassport()
		}
	}

	return pps
}
