package main

import (
	"os"
	"fmt"
	"bufio"
	"regexp"
	//"strconv"
	"strings"
)

var input = "8.txt"

func main(){
	if f,err := os.Open(input); err == nil {
		scanner := bufio.NewScanner(f)

		bs,_ := regexp.Compile(`\\`)
		quote,_ := regexp.Compile(`"`)

		total := 0
		replaced := 0
		for scanner.Scan() {
			var txt = scanner.Text()
			debug := false
			if strings.Contains(txt, `\x`){
				debug = true
				fmt.Println(txt)
			}
			total += len(txt)
			txt = bs.ReplaceAllString(txt, `||`)
			txt = quote.ReplaceAllString(txt, `|"`)
			txt = `"` + txt + `"`
		
			if debug {
				fmt.Println(txt)
			}
			replaced += len(txt)
		}
		fmt.Println(replaced-total)
	}
}