package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"regexp"
	"strconv"
	"time"
)

var input = "12.txt"

type moon struct {
	x, y, z    int
	vx, vy, vz int
}

func main() {
	startTime := time.Now()

	f, err := os.Open(input)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	scanner := bufio.NewScanner(f)
	reg := regexp.MustCompile("([x-z])=(-?[0-9]*)")
	moons := []moon{}
	for scanner.Scan() {
		var txt = scanner.Text()

		m := moon{}
		if groups := reg.FindAllStringSubmatch(txt, -1); groups != nil && len(groups) > 1 {
			for _, g := range groups {
				val, err := strconv.Atoi(g[2])
				if err != nil {
					panic("parsing " + g[2])
				}
				switch g[1] {
				case "x":
					m.x = val
				case "y":
					m.y = val
				case "z":
					m.z = val
				}
			}
		}
		moons = append(moons, m)
	}

	mcp := make([]moon, len(moons))
	copy(mcp, moons)

	mcp = simulate(mcp, 1000)
	p1 := getEnergy(mcp)
	fmt.Println("Part 1: ", p1)

	fmt.Println("Time", time.Since(startTime))
}

func getEnergy(moons []moon) int {
	energy := 0
	for _, m := range moons {

		pot := math.Abs(float64(m.x)) + math.Abs(float64(m.y)) + math.Abs(float64(m.z))
		kin := math.Abs(float64(m.vx)) + math.Abs(float64(m.vy)) + math.Abs(float64(m.vz))

		energy += int(pot * kin)
	}
	return energy
}

func simulate(moons []moon, steps int) []moon {
	i := 0
	for i < steps {
		for mi := 0; mi < len(moons); mi++ {
			m := moons[mi]
			for mj := range moons {
				m2 := moons[mj]
				if mi == mj {
					continue
				}

				m.vx += comp(m.x, m2.x)
				m.vy += comp(m.y, m2.y)
				m.vz += comp(m.z, m2.z)
			}

			moons[mi] = m
		}

		for mi := 0; mi < len(moons); mi++ {
			m := moons[mi]
			m.x += m.vx
			m.y += m.vy
			m.z += m.vz
			moons[mi] = m
		}
		i++
	}
	return moons
}

func comp(a, b int) int {
	if a > b {
		return -1
	} else if a < b {
		return 1
	}
	return 0
}
