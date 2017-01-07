package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"time"
	//"strconv"
	"strings"
	//"math"
)

var input = "7.txt"

func main() {
	startTime := time.Now()
	if f, err := os.Open(input); err == nil {
		scanner := bufio.NewScanner(f)
		bracketReg := regexp.MustCompile("\\[([a-z]+)\\]")
		nonBracketReg := regexp.MustCompile("([a-z]+)")
		count := 0
		sslcount := 0

		for scanner.Scan() {
			var txt = scanner.Text()

			midtls := false
			updated := txt

			sslChecks := []string{}

			if matches := bracketReg.FindAllStringSubmatch(txt, -1); matches != nil && len(matches) > 0 {
				for _, groups := range matches {
					for _, g := range groups[1:] {
						midtls = midtls || isTLS(g)
						sslChecks = append(sslChecks, sslCheck(g)...)
						updated = strings.Replace(updated, g, "", -1)
					}
				}
			}

			outertls := false
			ssl := false
			if matches := nonBracketReg.FindAllStringSubmatch(updated, -1); matches != nil && len(matches) > 0 {
				for _, groups := range matches {
					for _, g := range groups[1:] {
						outertls = outertls || isTLS(g)
						ssl = ssl || isSSL(g, sslChecks)
					}
				}
			}

			if !midtls && outertls {
				count++
			}

			if ssl {
				fmt.Println(txt)
				sslcount++
			}

		}

		fmt.Println("Total TLS Enabled IPs", count)
		fmt.Println("Total SSL Enabled IPs", sslcount)
	}

	fmt.Println("Time", time.Since(startTime))
}

func isTLS(str string) bool {
	tls := false
	for i := 0; i < len(str)-3 && !tls; i++ {
		tls = tls || (str[i] == str[i+3] && str[i+1] == str[i+2] && str[i] != str[i+1])
	}

	return tls
}

func sslCheck(str string) []string {
	ret := []string{}
	for i := 0; i < len(str)-2; i++ {
		ssl := (str[i] == str[i+2] && str[i] != str[i+1])
		if ssl {
			ret = append(ret, string(str[i+1])+string(str[i])+string(str[i+1]))
		}
	}

	return ret
}

func isSSL(str string, checks []string) bool {
	ssl := false
	for i := 0; i < len(str)-3 && !ssl; i++ {
		for _, chk := range checks {
			fmt.Println(chk)
			ssl = ssl || strings.Contains(str, chk)
		}
	}

	return ssl
}

// reg := regexp.MustCompile("-?[0-9]+")
/*
if groups := reg.FindStringSubmatch(txt); groups != nil && len(groups) > 1 {
                fmt.Println(groups[1:])
            }
*/
