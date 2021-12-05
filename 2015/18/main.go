package main

import (
	"fmt"
	"log"
	"time"

	"github.com/jasontconnell/advent/common"
)

var inputFilename = "input.txt"

type input []string
type output int

type Light struct {
	On bool
}

func (light Light) String() string {
	c := "."
	if light.On {
		c = "#"
	}
	return c
}

func main() {
	startTime := time.Now()

	in, err := common.ReadStrings(inputFilename)
	if err != nil {
		log.Fatal(err)
	}

	p1 := part1(input(in))
	p2 := part2(input(in))

	fmt.Println("Part 1:", p1)
	fmt.Println("Part 2:", p2)

	fmt.Println("Time", time.Since(startTime))
}

func part1(in input) output {
	lights := getInput(in)
	animated := runAnimation(lights, 100, false)
	return output(Count(animated, true))
}

func part2(in input) output {
	lights := getInput(in)
	animated := runAnimation(lights, 100, true)
	return output(Count(animated, true))
}

func runAnimation(lights [][]Light, seconds int, stickCorners bool) [][]Light {
	animatedLights := make([][]Light, len(lights))
	for i := 0; i < len(lights); i++ {
		animatedLights[i] = make([]Light, len(lights[i]))
		copy(animatedLights[i], lights[i])
	}

	for i := 0; i < seconds; i++ {
		animatedLights = Animate(animatedLights, stickCorners)
	}

	return animatedLights
}

func getInput(in input) [][]Light {
	lights := [][]Light{}
	for y, txt := range in {
		lights = append(lights, []Light{})
		x := 0

		for _, c := range txt {
			r := rune(c)

			light := Light{On: r == '#'}
			lights[y] = append(lights[y], light)
			x++
		}
	}
	return lights
}

// animate 1 step
func Animate(lights [][]Light, stickCorners bool) [][]Light {
	animated := make([][]Light, len(lights))
	for i := 0; i < len(lights); i++ {
		animated[i] = make([]Light, len(lights[i]))
		copy(animated[i], lights[i])
	}

	for y := 0; y < len(animated); y++ {
		row := animated[y]
		for x := 0; x < len(row); x++ {
			isCorner := (y == 0 || y == 99) && (x == 0 || x == 99)
			if isCorner && stickCorners {
				continue
			}

			c := NeighborCount(lights, y, x)
			if animated[y][x].On && (c != 2 && c != 3) {
				animated[y][x].On = false
			} else if !animated[y][x].On && c == 3 {
				animated[y][x].On = true
			}
		}
	}
	//PrintGrid(animated)
	return animated
}

func PrintGrid(lights [][]Light) {
	for y := 0; y < len(lights); y++ {
		row := lights[y]
		for x := 0; x < len(row); x++ {
			fmt.Print(lights[y][x])
		}
		fmt.Print("\n")
	}
}

func NeighborCount(lights [][]Light, y, x int) int {
	count := 0
	if x > 0 && lights[y][x-1].On {
		count++
	}
	if x < len(lights[y])-1 && lights[y][x+1].On {
		count++
	}
	if x > 0 && y > 0 && lights[y-1][x-1].On {
		count++
	}
	if x < len(lights[y])-1 && y < len(lights)-1 && lights[y+1][x+1].On {
		count++
	}
	if y > 0 && lights[y-1][x].On {
		count++
	}
	if y < len(lights)-1 && lights[y+1][x].On {
		count++
	}
	if y < len(lights)-1 && x > 0 && lights[y+1][x-1].On {
		count++
	}
	if y > 0 && x < len(lights[y])-1 && lights[y-1][x+1].On {
		count++
	}
	return count
}

func Count(lights [][]Light, state bool) int {
	count := 0
	for y := 0; y < len(lights); y++ {
		row := lights[y]
		for x := 0; x < len(row); x++ {
			if lights[y][x].On == state {
				count++
			}
		}
	}
	return count
}
