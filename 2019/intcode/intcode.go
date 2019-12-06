package intcode

import (
	"fmt"
)

type op struct {
	code   int
	params []int
}

type mode int

const (
	position  mode = 0
	immediate mode = 1
)


func Exec(ops []int) int {
	done := false
	step := 4
	for i := 0; i < len(ops) && !done; i += step {
		opcode := getOp(ops[i])
		switch opcode.code {
		case 1:
			value := getValue(ops, opcode, i, 1) + getValue(ops, opcode, i, 2)
			setValue(ops, opcode, i, 3, value)
			step = 4
			break
		case 2:
			value := getValue(ops, opcode, i, 1) * getValue(ops, opcode, i, 2)
			setValue(ops, opcode, i, 3, value)
			step = 4
			break
		case 3:
			value := getValue(ops, opcode, i, 1)
			setValue(ops, opcode, i, 3, value)
			step = 2
		case 4:
			outval := getValue(ops, opcode, i, 1)
			fmt.Println("output", outval)
			step = 2
		case 99:
			done = true
			break
		}
	}

	return ops[0]
}

func getValue(ops []int, opcode op, index, pos int) int {
	var val int
	mode := getMode(opcode, pos)

	if mode == position {
		val = ops[ops[index+pos]]
	} else {
		val = ops[index+pos]
	}
	return val
}

func setValue(ops []int, opcode op, index, pos, value int) {
	mode := getMode(opcode, pos)
	if mode == position {
		ops[ops[index+pos]] = value
	} else {
		ops[index+pos] = value
	}
}

func getMode(opcode op, pos int) mode {
	max := len(opcode.params)
	fp := pos-1 // pos is 1 based index
	if fp >= max || fp < 0 {
		return position
	}

	return mode(opcode.params[fp])
}

func getOp(val int) op {
	opcode := val % 100

	params := digits(val / 100)

	return op{opcode, params}
}

func digits(val int) []int {
	v := []int{}
	c := val
	div := 10
	done := false

	var x int
	for !done {
		x = c % div
		v = append(v, x)
		c = c / div
		done = c == 0
	}

	// reverse slice
	// not reversed: "read right-to-left from the opcode"
	// for i := 0; i < len(v)/2; i++ {
	// 	v[i], v[len(v)-i-1] = v[len(v)-i-1], v[i]
	// }

	return v
}
