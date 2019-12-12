package intcode

type Computer struct {
	Name     string
	Ops      []int
	Outs     []int
	Ins      []int
	InstPtr  int
	Prev     *Computer
	Next     *Computer
	Complete bool
}

type op struct {
	code   int
	params []int
}

type mode int

const (
	position  mode = 0
	immediate mode = 1
)

func NewComputer(ops []int) *Computer {
	c := &Computer{Ops: ops, InstPtr: 0, Complete: false}
	return c
}

func (c *Computer) Exec() {
	for !c.Complete {
		c.ExecOne()
	}
}

func (c *Computer) GetNextOutput() int {
	for !c.Complete && len(c.Outs) == 0 {
		c.ExecOne()
	}
	out := c.Outs[0]
	c.Outs = c.Outs[1:]
	return out
}

func (c *Computer) ExecOne() {
	if c.Complete {
		return
	}

	opcode := getOp(c.Ops[c.InstPtr])
	switch opcode.code {
	case 1:
		value := c.getValue(1) + c.getValue(2)
		c.setValue(3, value)
		c.InstPtr += 4
		break
	case 2:
		value := c.getValue(1) * c.getValue(2)
		c.setValue(3, value)
		c.InstPtr += 4
		break
	case 3:
		if len(c.Ins) == 0 && c.Prev != nil {
			out := c.Prev.GetNextOutput()
			c.Ins = append(c.Ins, out)
		}
		c.setValue(1, c.Ins[0])
		c.Ins = c.Ins[1:]
		c.InstPtr += 2
		break
	case 4:
		outval := c.getValue(1)
		c.Outs = append([]int{outval}, c.Outs...)
		c.InstPtr += 2
		break
	case 5, 6:
		value := c.getValue(1)
		if (opcode.code == 5 && value != 0) || (opcode.code == 6 && value == 0) {
			c.InstPtr = c.getValue(2)
		} else {
			c.InstPtr += 3
		}
		break
	case 7, 8:
		value1 := c.getValue(1)
		value2 := c.getValue(2)
		if (opcode.code == 7 && value1 < value2) || (opcode.code == 8 && value1 == value2) {
			c.setValue(3, 1)
		} else {
			c.setValue(3, 0)
		}
		c.InstPtr += 4
		break
	case 99:
		c.Complete = true
		break
	default:
		c.InstPtr += 4
	}
}

func (c *Computer) getValue(pos int) int {
	var val int
	opcode := getOp(c.Ops[c.InstPtr])
	mode := getMode(opcode, pos)
	index := c.InstPtr

	if mode == position {
		val = c.Ops[c.Ops[index+pos]]
	} else {
		val = c.Ops[index+pos]
	}
	return val
}

func (c *Computer) setValue(pos, value int) {
	opcode := getOp(c.Ops[c.InstPtr])
	mode := getMode(opcode, pos)
	index := c.InstPtr

	if mode == position {
		c.Ops[c.Ops[index+pos]] = value
	} else {
		c.Ops[index+pos] = value
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
