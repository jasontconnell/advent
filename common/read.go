package common

import (
	"bufio"
	"os"
	"strconv"
)

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
