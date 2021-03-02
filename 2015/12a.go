package main

import (
	//"bufio"
	"fmt"
	"os"
	"time"

	//"regexp"
	//"strconv"
	"encoding/json"
	//"strings"
)

var input = "12.txt"

func main() {
	startTime := time.Now()
	if bytes, err := os.ReadFile(input); err == nil {
		var data interface{}
		var sum int

		if err := json.Unmarshal(bytes, &data); err == nil {
			switch data.(type) {
			case map[string]interface{}:
				sum = getSum(data.(map[string]interface{}))
				break
			case []interface{}:
				sum = getArraySum(data.([]interface{}))
				break
			}

		} else {
			panic(err)
		}

		fmt.Println("sum", sum)
	}

	fmt.Println("Time", time.Since(startTime))

}

func getSum(m map[string]interface{}) (sum int) {
	red := false
	for _, v := range m {
		switch v.(type) {
		case string:
			red = v.(string) == "red"
			break
		}

		if red {
			return
		}
	}

	for _, v := range m {
		switch v.(type) {
		case int:
			sum += v.(int)
			break
		case []interface{}:
			sum += getArraySum(v.([]interface{}))
			break
		case float64:
			sum += int(v.(float64))
		case map[string]interface{}:
			sum += getSum(v.(map[string]interface{}))
			break
		}
	}

	return
}

func getArraySum(arr []interface{}) (sum int) {
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
			sum += getArraySum(v.([]interface{}))
			break
		case map[string]interface{}:
			sum += getSum(v.(map[string]interface{}))
			break
		}
	}
	return
}
