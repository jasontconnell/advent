package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/jasontconnell/advent/common"
)

type input = []string
type output = int

type Particle struct {
	ID int
	Vector
	Veloc   Vector
	Accel   Vector
	Removed bool
}

func (p *Particle) Distance() int {
	x := p.X
	if x < 0 {
		x = -x
	}

	y := p.Y
	if y < 0 {
		y = -y
	}

	z := p.Z
	if z < 0 {
		z = -z
	}

	return x + y + z
}

type Vector struct {
	X, Y, Z int
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
	fmt.Fprintln(w, "--2017 day 20 solution--")
	fmt.Fprintln(w, "Part 1:", p1)
	fmt.Fprintln(w, "Part 2:", p2)
	fmt.Println("Time", time.Since(startTime))
}

func part1(in input) output {
	particles := parseInput(in)
	return determine(particles)
}

func part2(in input) output {
	particles := parseInput(in)
	determine(particles)
	return countLeft(particles)
}

func countLeft(particles []*Particle) int {
	left := 0
	for _, p := range particles {
		if !p.Removed {
			left++
		}
	}
	return left
}

func removeCollisions(particles *[]*Particle) {
	vmap := make(map[Vector][]*Particle)
	for i := 0; i < len(*particles); i++ {
		p := (*particles)[i]
		if !p.Removed {
			v := Vector{X: p.X, Y: p.Y, Z: p.Z}
			vmap[v] = append(vmap[v], p)
		}
	}

	for _, plist := range vmap {
		if len(plist) > 1 {
			for _, p := range plist {
				p.Removed = true
			}
		}

	}
}

func determine(particles []*Particle) int {
	i := 0
	dmap := make(map[int]int)
	for i < 1000 {
		animateAll(particles)
		removeCollisions(&particles)

		for _, p := range particles {
			dmap[p.ID] = p.Distance()
		}
		i++
	}

	min := -1
	id := -1

	for k, p := range dmap {
		if min == -1 || p < min {
			min = p
			id = k
		}
	}

	return id
}

func animateAll(particles []*Particle) {
	for _, p := range particles {
		animate(p)
	}
}

func animate(particle *Particle) {
	particle.Veloc.X += particle.Accel.X
	particle.Veloc.Y += particle.Accel.Y
	particle.Veloc.Z += particle.Accel.Z

	particle.X += particle.Veloc.X
	particle.Y += particle.Veloc.Y
	particle.Z += particle.Veloc.Z
}

func getVector(xs, ys, zs string) Vector {
	x, err := strconv.Atoi(xs)
	if err != nil {
		panic(err)
	}

	y, err := strconv.Atoi(ys)
	if err != nil {
		panic(err)
	}

	z, err := strconv.Atoi(zs)
	if err != nil {
		panic(err)
	}

	return Vector{X: x, Y: y, Z: z}
}

func parseInput(in input) []*Particle {
	coordreg := "<(-?[0-9]+),(-?[0-9]+),(-?[0-9]+)>"
	reg := regexp.MustCompile("p=" + coordreg + ", v=" + coordreg + ", a=" + coordreg)

	particles := []*Particle{}
	for i, line := range in {
		m := reg.FindStringSubmatch(line)
		if len(m) == 10 {
			p := &Particle{ID: i, Vector: getVector(m[1], m[2], m[3]), Veloc: getVector(m[4], m[5], m[6]), Accel: getVector(m[7], m[8], m[9])}
			particles = append(particles, p)
		}
	}
	return particles
}
