package main

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"time"

	"github.com/jasontconnell/advent/common"
)

var inputFilename = "input.txt"

type input []string
type output int

var reg *regexp.Regexp = regexp.MustCompile("-?[0-9]+")

func main() {
	startTime := time.Now()

	in, err := common.ReadStrings(inputFilename)
	if err != nil {
		log.Fatal(err)
	}

	p1 := part1(input(in))
	p2 := part2(input(in))

	fmt.Println("Part 1:", p1)
	fmt.Println("Part 2:", p2)

	fmt.Println("Time", time.Since(startTime))
}

func part1(in input) output {
	return getSum(in)
}

func part2(in input) output {
	return getSumJson(in, "red")
}

func getSum(in input) output {
	sum := 0
	for _, txt := range in {
		if str := reg.FindAllString(txt, -1); len(str) > 0 {
			for _, r := range str {
				if i, err := strconv.Atoi(string(r)); err == nil {
					sum += i
				} else {
					fmt.Print(string(r), "not a number")
				}
			}
		}
	}
	return output(sum)
}

func getSumJson(in input, exclude string) output {
	var data interface{}
	var sum int

	for _, txt := range in {
		bytes := []byte(txt)
		if err := json.Unmarshal(bytes, &data); err == nil {
			switch data.(type) {
			case map[string]interface{}:
				sum = getMapSum(data.(map[string]interface{}), exclude)
				break
			case []interface{}:
				sum = getArraySum(data.([]interface{}), exclude)
				break
			}

		}
	}

	return output(sum)
}

func getMapSum(m map[string]interface{}, exclude string) int {
	skip := false
	sum := 0
	for _, v := range m {
		switch v.(type) {
		case string:
			skip = v.(string) == exclude
			break
		}

		if skip {
			return sum
		}
	}

	for _, v := range m {
		switch v.(type) {
		case int:
			sum += v.(int)
			break
		case []interface{}:
			sum += getArraySum(v.([]interface{}), exclude)
			break
		case float64:
			sum += int(v.(float64))
		case map[string]interface{}:
			sum += getMapSum(v.(map[string]interface{}), exclude)
			break
		}
	}

	return sum
}

func getArraySum(arr []interface{}, exclude string) int {
	sum := 0
	for _, v := range arr {
		switch v.(type) {
		case int:
			sum += v.(int)
			break
		case float64:
			sum += int(v.(float64))
			break
		case string:
			break
		case []interface{}:
			sum += getArraySum(v.([]interface{}), exclude)
			break
		case map[string]interface{}:
			sum += getMapSum(v.(map[string]interface{}), exclude)
			break
		}
	}
	return sum
}
