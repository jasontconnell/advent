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

		hex,_ := regexp.Compile(`\\x[abcdef0-9]{2}`)
		q,_ := regexp.Compile(`\\"`)
		bs,_ := regexp.Compile(`\\\\`)

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
			txt = txt[1:len(txt)-1]
			txt = bs.ReplaceAllString(txt, `|`)
			txt = q.ReplaceAllString(txt, `|`)
			txt = hex.ReplaceAllString(txt, "|")
		
			if debug {
				fmt.Println(txt)
			}
			replaced += len(txt)
		}
		fmt.Println(total - replaced)
	}
}