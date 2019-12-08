package intcode

type op struct {
	code   int
	params []int
}

type mode int

const (
	position  mode = 0
	immediate mode = 1
)

func Exec(ops []int, inputs []int) ([]int, []int) {
	outputs := []int{}
	done := false
	step := 4
	inputIndex := 0
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
			setValue(ops, opcode, i, 1, inputs[inputIndex])
			inputIndex++
			if inputIndex > len(inputs) - 1 {
				inputIndex = len(inputs) - 1
			}
			step = 2
			break
		case 4:
			outval := getValue(ops, opcode, i, 1)
			outputs = append(outputs, outval)
			step = 2
			break
		case 5, 6:
			value := getValue(ops, opcode, i, 1)
			if (opcode.code == 5 && value != 0) || (opcode.code == 6 && value == 0) {
				i = getValue(ops, opcode, i, 2)
				step = 0
			} else {
				step = 3
			}
			break
		case 7, 8:
			value1 := getValue(ops, opcode, i, 1)
			value2 := getValue(ops, opcode, i, 2)
			if (opcode.code == 7 && value1 < value2) || (opcode.code == 8 && value1 == value2) {
				setValue(ops, opcode, i, 3, 1)
			} else {
				setValue(ops, opcode, i, 3, 0)
			}
			step = 4
			break
		case 99:
			done = true
			break
		default:
			step = 4
		}
	}

	return ops, outputs
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
	fp := pos - 1 // pos is 1 based index
	if fp >= max || fp < 0 {
		return position
	}

	return mode(opcode.params[fp])
}

func getOp(val int) op {
	opcode := val % 100

	params := Digits(val / 100)

	return op{opcode, params}
}

func Digits(val int) []int {
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
