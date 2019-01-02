package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"strconv"
	"strings"
)

var input = "iwrupvqb"

func main() {
	found := false
	lastTest := ""
	for i := 0; !found; i++ {
		test := input + strconv.Itoa(i)

		hex := getMD5Hex(test)

		if found = strings.HasPrefix(hex, "000000"); found {
			lastTest = test
		}
	}

	fmt.Println(lastTest)
}

func getMD5Hex(str string) string {
	h := md5.New()
	io.WriteString(h, str)
	hex := fmt.Sprintf("%x", h.Sum(nil))
	return hex
}
