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

const (
	lscmd      = "$ ls"
	totalsize  = 70000000
	updatesize = 30000000
)

type dir struct {
	parent  *dir
	name    string
	files   []file
	subdirs []*dir
	size    int
}

type file struct {
	parent *dir
	name   string
	size   int
}

func (d *dir) getSize() int {
	if d.size > 0 {
		return d.size
	}
	for _, sd := range d.subdirs {
		d.size += sd.getSize()
	}

	for _, f := range d.files {
		d.size += f.size
	}
	return d.size
}

func (d *dir) print(level int) {
	tab := strings.Repeat(" ", level)

	fmt.Println(tab, "/"+d.name, tab+tab, "( total:", d.getSize(), ")")
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
		sum += d.getSize()
	}
	return sum
}

func part2(in input) output {
	root := parseInput(in)
	del := getDelCandidate(root, totalsize, updatesize)
	return del.getSize()
}

func getDelCandidate(root *dir, total, update int) *dir {
	used := root.getSize()
	unused := total - used
	tofree := update - unused
	cands := getBigDirs(root, tofree)
	sm := root
	for _, d := range cands {
		if d.getSize() < sm.getSize() {
			sm = d
		}
	}

	return sm
}

func getBigDirs(root *dir, minsize int) []*dir {
	dirs := []*dir{}

	for _, d := range root.subdirs {
		if d.getSize() >= minsize {
			dirs = append(dirs, d)
		}

		subs := getBigDirs(d, minsize)
		dirs = append(dirs, subs...)
	}

	return dirs
}

func getSmallDirs(root *dir, maxsize int) []*dir {
	dirs := []*dir{}

	for _, d := range root.subdirs {
		if d.getSize() < maxsize {
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
			} else {
				for _, d := range cur.subdirs {
					if d.name == nm {
						cur = d
					}
				}
			}

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
			}

			if len(fm) > 0 {
				sz, _ := strconv.Atoi(fm[1])
				cur.files = append(cur.files, file{parent: cur, name: fm[2], size: sz})
			}
		}
	}
	return root
}
