package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"time"
)

var input = "12.txt"

type op string

const (
	N op = "N"
	S op = "S"
	E op = "E"
	W op = "W"
	L op = "L"
	R op = "R"
	F op = "F"
)

type command struct {
	instr op
	value int
}

type xy struct {
	x, y int
}

func (p xy) String() string {
	return fmt.Sprintf("(%d,%d)", p.x, p.y)
}

type ship struct {
	bearing  xy
	position xy
	waypoint xy
}

func main() {
	startTime := time.Now()

	f, err := os.Open(input)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	scanner := bufio.NewScanner(f)

	lines := []string{}
	for scanner.Scan() {
		var txt = scanner.Text()
		lines = append(lines, txt)
	}

	commands := getCommands(lines)

	sh := &ship{bearing: xy{1, 0}, position: xy{0, 0}}
	pos := sail(commands, sh)
	fmt.Println("Part 1:", int(math.Abs(float64(pos.x))+math.Abs(float64(pos.y))), pos)

	sh.waypoint = xy{10, -1}
	sh.position = xy{0, 0}

	wpos := sailWaypoint(commands, sh)
	fmt.Println("Part 2:", int(math.Abs(float64(wpos.x))+math.Abs(float64(wpos.y))), wpos)
	fmt.Println("Time", time.Since(startTime))
}

func turnWaypoint(sh *ship, val int) {
	lp := val / 90

	for i := 0; i < lp; i++ {
		x := sh.waypoint.x - sh.position.x
		y := sh.waypoint.y - sh.position.y

		if x >= 0 && y >= 0 {
			sh.waypoint.x = sh.position.x + y
			sh.waypoint.y = sh.position.y + -x
		} else if x >= 0 && y <= 0 {
			sh.waypoint.x = sh.position.x + y
			sh.waypoint.y = sh.position.y + -x
		} else if x <= 0 && y <= 0 {
			sh.waypoint.x = sh.position.x + y
			sh.waypoint.y = sh.position.y + -x
		} else if x <= 0 && y >= 0 {
			sh.waypoint.x = sh.position.x + y
			sh.waypoint.y = sh.position.y + -x
		}
	}
}

func turnShip(sh *ship, val int) {
	lp := val / 90

	for i := 0; i < lp; i++ {
		if sh.bearing.y == -1 {
			sh.bearing.x = -1
			sh.bearing.y = 0
		} else if sh.bearing.y == 1 {
			sh.bearing.x = 1
			sh.bearing.y = 0
		} else if sh.bearing.x == -1 {
			sh.bearing.x = 0
			sh.bearing.y = 1
		} else if sh.bearing.x == 1 {
			sh.bearing.x = 0
			sh.bearing.y = -1
		}
	}
}

func sail(commands []command, sh *ship) xy {
	for _, cmd := range commands {
		switch cmd.instr {
		case F:
			if sh.bearing.x != 0 {
				sh.position.x += (sh.bearing.x * cmd.value)
			} else {
				sh.position.y += (sh.bearing.y * cmd.value)
			}
		case L: // only L
			turnShip(sh, cmd.value)
		case N:
			sh.position.y += -cmd.value
		case E:
			sh.position.x += cmd.value
		case W:
			sh.position.x += -cmd.value
		case S:
			sh.position.y += cmd.value
		}
	}

	return sh.position
}

func sailWaypoint(commands []command, sh *ship) xy {
	for _, cmd := range commands {
		switch cmd.instr {
		case F:
			x := (sh.waypoint.x - sh.position.x) * cmd.value
			y := (sh.waypoint.y - sh.position.y) * cmd.value

			sh.position.x += x
			sh.position.y += y

			sh.waypoint.x += x
			sh.waypoint.y += y
		case L:
			turnWaypoint(sh, cmd.value)
		case N:
			sh.waypoint.y += -cmd.value
		case E:
			sh.waypoint.x += cmd.value
		case W:
			sh.waypoint.x += -cmd.value
		case S:
			sh.waypoint.y += cmd.value
		}
	}

	return sh.position
}

func getCommands(lines []string) []command {
	cmds := []command{}
	for _, line := range lines {
		ch := line[0]

		cmd := command{}
		val, err := strconv.Atoi(string(line[1:]))
		if err != nil {
			panic(err.Error())
		}
		cmd.value = val
		cmd.instr = op(string(ch))

		// just handle all left hand turns
		if cmd.instr == R {
			cmd.instr = L
			cmd.value = 360 - cmd.value
		}

		cmds = append(cmds, cmd)
	}
	return cmds
}
