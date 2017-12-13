package main

import (
	"bufio"
	"fmt"
	"os"
	"time"
	//"regexp"
	"strconv"
	"strings"
	//"math"
)

var input = "13.txt"

type Scanner struct {
	Depth   int
	Range   int
	Current int
	Dir     int
	Nil     bool
}

func main() {
	startTime := time.Now()

	f, err := os.Open(input)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	scanner := bufio.NewScanner(f)
	scanners := []*Scanner{}

	for scanner.Scan() {
		var txt = scanner.Text()
		s := getScanner(txt)
		if s == nil {
			continue
		}
		scanners = append(scanners, s)
	}

	scanners = fill(scanners)

	sev, caught := negotiateWithSeverity(scanners)
	delay := negotiateWithoutCapture(scanners)

	fmt.Println("Firewall Severity                     : ", sev, caught)
	fmt.Println("Minimum Delay to get through unscathed: ", delay)
	fmt.Println("Time", time.Since(startTime))
}

func fill(scanners []*Scanner) []*Scanner {
	max := scanners[len(scanners)-1].Depth
	for i := 0; i < max; i++ {
		if scanners[i].Depth != i {
			s := &Scanner{Depth: i, Range: 0, Nil: true}
			scanners = append(scanners[:i], append([]*Scanner{s}, scanners[i:]...)...)
		}
	}
	return scanners
}

func reset(scanners []*Scanner) {
	for i := 0; i < len(scanners); i++ {
		scanners[i].Current = 0
	}
}

func negotiateWithoutCapture(scanners []*Scanner) int {
	delay := 0

	for {
		reset(scanners)

		p := possible(scanners, delay)

		if p {
			for j := 0; j < delay; j++ {
				tick(scanners)
			}

			sev, caught := negotiateWithSeverity(scanners)
			if sev == 0 && !caught {
				break
			}
		}

		delay++
	}
	return delay
}

func possible(scanners []*Scanner, delay int) bool {
	p := true
	for _, s := range scanners {

		zero := (s.Range*2 - 2)
		if s.Nil {
			zero = 0
		}

		p = p && (s.Nil || ((delay+s.Depth)%zero != 0))
	}
	return p
}

func negotiateWithSeverity(scanners []*Scanner) (int, bool) {
	sev := 0
	caught := false
	for _, s := range scanners {
		if s.Current == 0 && !s.Nil {
			sev += (s.Depth * s.Range)
			caught = caught || true
		}

		tick(scanners)
	}
	return sev, caught
}

func list(scanners []*Scanner) string {
	f := ""
	for _, s := range scanners {
		f += fmt.Sprintf("scanner (%v) with range (%v) at cycle %v\n", s.Depth, s.Range, s.Current)
	}
	return f
}

func tick(scanners []*Scanner) {
	for _, s := range scanners {
		if s.Nil {
			continue
		}

		if s.Current == s.Range-1 {
			s.Dir = -1
		} else if s.Current == 0 {
			s.Dir = 1
		}

		s.Current += s.Dir
	}
}

func getScanner(line string) *Scanner {
	s := strings.Split(line, ":")
	var l *Scanner
	if len(s) == 2 {
		d, err := strconv.Atoi(s[0])
		if err != nil {
			fmt.Println("parse", err, s[0])
			return nil
		}

		r, err := strconv.Atoi(strings.Trim(s[1], " "))
		if err != nil {
			fmt.Println("parse", err, s[1])
			return nil
		}

		l = &Scanner{Depth: d, Range: r, Current: 0, Dir: 1, Nil: false}
	}
	return l
}
