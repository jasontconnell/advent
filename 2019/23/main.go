package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/jasontconnell/advent/2019/intcode"
	"github.com/jasontconnell/advent/common"
)

var inputFilename = "input.txt"

type input = []int
type output = int

type packet struct {
	x, y   *int
	isgoal bool
}

func (p *packet) String() string {
	return fmt.Sprintf("(%d, %d)", p.x, p.y)
}

func main() {
	startTime := time.Now()

	in, err := common.ReadIntCsv(inputFilename)
	if err != nil {
		log.Fatal(err)
	}

	p1 := part1(in)
	p2 := part2(in)

	fmt.Println("Part 1:", p1)
	fmt.Println("Part 2:", p2)

	fmt.Println("Time", time.Since(startTime))
}

func part1(in input) output {
	computers := getComputers(in, 50)
	return run(computers)
}

func part2(in input) output {
	return 0
}

func run(computers map[int]*intcode.Computer) int {
	outqueue := sync.Map{}
	inqueue := sync.Map{}
	done := make(chan bool)
	result := 0
	for i, c := range computers {
		c.AddInput(i)
		c.OnOutputCtx = func(ctx *intcode.Computer, out int) {
			stored, _ := outqueue.LoadOrStore(ctx.ID, []int{})
			outs := stored.([]int)
			outs = append(outs, out)
			if len(outs) > 2 {
				to, x, y := outs[0], outs[1], outs[2]
				newouts := outs[3:]
				fmt.Println("process queue for", ctx.ID, outs)
				// fmt.Println(ctx.ID, "to, x, y", to, x, y)
				inqueue.Store(to, &packet{nil, &y, to == 255})
				if to == 255 {
					result = y
					done <- true
				}
				if to > 49 {
					fmt.Println("to > 49", to)
				} else {
					dest := computers[to]
					dest.AddInput(x)
				}
				outqueue.Store(ctx.ID, newouts)
				// ctx.Outs = ctx.Outs[3:]
			} else {
				outqueue.Store(ctx.ID, outs)
			}
		}

		c.RequestInputCtx = func(ctx *intcode.Computer) int {
			stored, ok := inqueue.Load(ctx.ID)
			if !ok {
				return -1
			}

			p := stored.(*packet)
			if p.y == nil {
				return -1
			}

			val := 0
			if p.x != nil {
				val = *p.x
				p.x = nil
				inqueue.Store(ctx.ID, p)
			} else if p.y != nil {
				val = *p.y
				p.y = nil
				inqueue.Delete(ctx.ID)

				if p.isgoal {
					fmt.Println("goal", p.y)
					result = *p.y
					done <- true
				}
			}
			return val
		}
	}

	for _, c := range computers {
		go c.Exec()
	}

	<-done

	return result
}

func getComputers(code []int, num int) map[int]*intcode.Computer {
	m := make(map[int]*intcode.Computer)

	for i := 0; i < num; i++ {
		c := intcode.NewComputer(code)
		c.ID = i
		m[i] = c
	}
	return m
}

// reg := regexp.MustCompile("-?[0-9]+")
/*
if groups := reg.FindStringSubmatch(txt); groups != nil && len(groups) > 1 {
				fmt.Println(groups[1:])
			}
*/
