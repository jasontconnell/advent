package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"time"
	//"strings"
	//"math"
)

var input = "20.txt"

var coordreg string = "<(-?[0-9]+),(-?[0-9]+),(-?[0-9]+)>"
var reg *regexp.Regexp = regexp.MustCompile("p=" + coordreg + ", v=" + coordreg + ", a=" + coordreg)

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

	f, err := os.Open(input)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	scanner := bufio.NewScanner(f)
	particles := []*Particle{}

	for scanner.Scan() {
		var txt = scanner.Text()
		p := getParticle(txt)
		if p != nil {
			p.ID = len(particles)
			particles = append(particles, p)
		}
	}

	id := determine(particles)

	left := countLeft(particles)

	fmt.Println("Closest to 0            :", id)
	fmt.Println("Left after collisions   :", left)

	fmt.Println("Time", time.Since(startTime))
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

func getParticle(line string) *Particle {
	if groups := reg.FindStringSubmatch(line); groups != nil && len(groups) > 1 {
		p := &Particle{Vector: getVector(groups[1], groups[2], groups[3]), Veloc: getVector(groups[4], groups[5], groups[6]), Accel: getVector(groups[7], groups[8], groups[9])}
		return p
	}
	return nil
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

// reg := regexp.MustCompile("-?[0-9]+")
/*
if groups := reg.FindStringSubmatch(txt); groups != nil && len(groups) > 1 {
                fmt.Println(groups[1:])
            }
*/
