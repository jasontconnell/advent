package intcode

type Computer struct {
	Name         string
	Ops          []int
	Outs         []int
	Ins          []int
	InstPtr      int
	Prev         *Computer
	Next         *Computer
	Complete     bool
	RelativeBase int
	OnOutput     func(int)
	RequestInput func() int

	memory   map[int]int
	original []int
}

type op struct {
	code   int
	params []int
}

type mode int

const (
	position  mode = 0
	immediate mode = 1
	relative  mode = 2
)

func NewComputer(ops []int) *Computer {
	cp := make([]int, len(ops))

	copy(cp, ops)
	c := &Computer{Ops: cp, InstPtr: 0, Complete: false, memory: make(map[int]int), original: ops}
	return c
}

func (c *Computer) Reset() {
	c.Ops = c.original
	c.InstPtr = 0
	c.Complete = false
	c.memory = make(map[int]int)
	c.Ins = []int{}
	c.Outs = []int{}
}

func (c *Computer) AddInput(ins ...int) {
	c.Ins = append(c.Ins, ins...)
}

func (c *Computer) AddInputs(ins ...[]int) {
	for _, a := range ins {
		c.AddInput(a...)
	}
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
			c.AddInput(out)
		}

		if len(c.Ins) == 0 && c.RequestInput != nil {
			in := c.RequestInput()
			c.AddInput(in)
		}

		if len(c.Ins) > 0 {
			c.setValue(1, c.Ins[0])
			c.Ins = c.Ins[1:]
		}
		c.InstPtr += 2
		break
	case 4:
		outval := c.getValue(1)
		c.Outs = append([]int{outval}, c.Outs...)
		if c.OnOutput != nil {
			c.OnOutput(outval)
		}
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
	case 9:
		value1 := c.getValue(1)
		c.RelativeBase += value1
		c.InstPtr += 2
		break
	case 99:
		c.Complete = true
		break
	default:
		panic("unknown op " + string(opcode.code))
	}
}

func (c *Computer) getValue(pos int) int {
	if m, ok := c.memory[pos]; ok {
		return m
	}

	if pos > len(c.Ops) {
		return 0
	}

	var val int
	opcode := getOp(c.Ops[c.InstPtr])
	mode := getMode(opcode, pos)
	index := c.InstPtr

	switch mode {
	case position:
		//val = c.Ops[c.Ops[index+pos]]
		v1 := c.access(index + pos)
		val = c.access(v1)
	case immediate:
		// val = c.Ops[index+pos]
		val = c.access(index + pos)
	case relative:
		ridx := c.RelativeBase + c.access(index+pos)

		val = c.access(ridx)
	}
	return val
}

func (c *Computer) access(index int) int {
	if index >= len(c.Ops) {
		if m, ok := c.memory[index]; ok {
			return m
		} else {
			return 0
		}
	}
	return c.Ops[index]
}

func (c *Computer) write(index int, value int) {
	if index >= len(c.Ops) {
		c.memory[index] = value
		return
	}
	c.Ops[index] = value
}

func (c *Computer) setValue(pos, value int) {
	if pos > len(c.Ops) {
		c.memory[pos] = value
	}
	opcode := getOp(c.Ops[c.InstPtr])
	mode := getMode(opcode, pos)
	index := c.InstPtr

	switch mode {
	case position:
		//c.Ops[c.Ops[index+pos]] = value
		v1 := c.access(index + pos)
		c.write(v1, value)
	case immediate:
		//c.Ops[index+pos] = value
		c.write(index+pos, value)
	case relative:
		//c.Ops[c.Ops[index+c.RelativeBase]] = value
		ridx := c.RelativeBase + c.access(index+pos)
		c.write(ridx, value)
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
	return v
}
