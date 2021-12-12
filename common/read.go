package common

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func InputFilename(args []string) string {
	filename := "input.txt"
	if len(args) > 1 {
		filename = args[1]
	}
	return filename
}

func ReadStrings(filename string) ([]string, error) {
	f, err := os.Open(filename)

	if err != nil {
		return nil, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	lines := []string{}
	for scanner.Scan() {
		var txt = scanner.Text()
		lines = append(lines, txt)
	}

	return lines, nil
}

func ReadString(filename string) (string, error) {
	lines, err := ReadStrings(filename)
	if err != nil {
		return "", err
	}

	if len(lines) == 0 {
		return "", fmt.Errorf("no lines to read")
	}

	return lines[0], nil
}

func ReadInts(filename string) ([]int, error) {
	lines, err := ReadStrings(filename)
	if err != nil {
		return nil, err
	}

	vals := []int{}
	for _, s := range lines {
		i, _ := strconv.Atoi(s)
		vals = append(vals, i)
	}
	return vals, nil
}

func ReadInt(filename string) (int, error) {
	ints, err := ReadInts(filename)
	if err != nil {
		return -1, err
	}

	return ints[0], nil
}

func ReadIntCsv(filename string) ([]int, error) {
	s, err := ReadString(filename)
	if err != nil {
		return nil, err
	}

	ret := []int{}
	sp := strings.Split(s, ",")
	for _, x := range sp {
		i, _ := strconv.Atoi(x)
		ret = append(ret, i)
	}
	return ret, nil
}

func ReadStringsCsv(filename string) ([]string, error) {
	s, err := ReadString(filename)
	if err != nil {
		return nil, err
	}
	ret := []string{}
	sp := strings.Split(s, ",")
	for _, s := range sp {
		ret = append(ret, strings.Trim(s, " "))
	}
	return ret, nil
}
