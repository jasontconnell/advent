package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/jasontconnell/advent/common"
)

type input = []string
type output = int

const lscmd = "$ ls"

type dir struct {
	parent  *dir
	name    string
	files   []file
	subdirs []*dir
}

type file struct {
	parent *dir
	name   string
	size   int
}

func (d *dir) getSize(recurse bool) int {
	sz := 0
	if recurse {
		for _, sd := range d.subdirs {
			sz += sd.getSize(recurse)
		}
	}

	for _, f := range d.files {
		sz += f.size
	}
	return sz
}

func (d *dir) print(level int) {
	tab := strings.Repeat(" ", level)

	fmt.Println(tab, "/"+d.name)
	for _, dir := range d.subdirs {
		dir.print(level + 1)
	}

	for _, f := range d.files {
		fmt.Println(tab, "-"+f.name, f.size)
	}
}

func main() {
	startTime := time.Now()

	in, err := common.ReadStrings(common.InputFilename(os.Args))
	if err != nil {
		log.Fatal(err)
	}

	p1 := part1(in)
	p2 := part2(in)

	w := common.TeeOutput(os.Stdout)
	fmt.Fprintln(w, "--2022 day 07 solution--")
	fmt.Fprintln(w, "Part 1:", p1)
	fmt.Fprintln(w, "Part 2:", p2)
	fmt.Println("Time", time.Since(startTime))
}

func part1(in input) output {
	root := parseInput(in)
	big := getSmallDirs(root, 100000)
	sum := 0
	for _, d := range big {
		sum += d.getSize(true)
	}
	return sum
}

func part2(in input) output {
	return 0
}

func getSmallDirs(root *dir, maxsize int) []*dir {
	dirs := []*dir{}

	for _, d := range root.subdirs {
		if d.getSize(true) < maxsize {
			dirs = append(dirs, d)
		}

		subs := getSmallDirs(d, maxsize)
		dirs = append(dirs, subs...)
	}

	return dirs
}

func parseInput(in input) *dir {
	root := &dir{name: "/"}
	cur := root
	isls := false
	dirs := make(map[string]*dir)
	dirs["/"] = root

	cdreg := regexp.MustCompile(`\$ cd ([a-zA-Z\\\./]*)?`)
	dirreg := regexp.MustCompile("dir (.*)")
	freg := regexp.MustCompile("([0-9]+) (.*)")

	for i := 0; i < len(in); i++ {
		m := cdreg.FindStringSubmatch(in[i])
		if len(m) > 0 {
			isls = false
			nm := m[1]
			if nm == ".." {
				cur = cur.parent
				continue
			}
			if _, ok := dirs[nm]; !ok {
				dirs[nm] = &dir{name: nm, parent: cur}
			}
			cur = dirs[nm]
			continue
		}

		if in[i] == lscmd {
			isls = true
			continue
		}

		if isls {
			dm := dirreg.FindStringSubmatch(in[i])
			fm := freg.FindStringSubmatch(in[i])

			if len(dm) > 0 {
				sub := &dir{name: dm[1], parent: cur}
				cur.subdirs = append(cur.subdirs, sub)
				dirs[sub.name] = sub
			}

			if len(fm) > 0 {
				sz, _ := strconv.Atoi(fm[1])
				cur.files = append(cur.files, file{parent: cur, name: fm[2], size: sz})
			}
		}
	}
	return root
}
