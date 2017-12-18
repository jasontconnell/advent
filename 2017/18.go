package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

var input = "18.txt"

type Inst struct {
	Cmd string

	Arg1Reg *Register
	Arg1Val int

	Arg2Reg *Register
	Arg2Val int

	Line string
}

type Register struct {
	Name  rune
	Value int
}

type Program struct {
	Name       string
	Registers  map[rune]*Register
	Insts      []Inst
	Queue      []int
	Parallel   *Program
	Deadlocked bool
	Sends      int
	InstPtr    int
}

func (r *Register) String() string {
	return fmt.Sprintf("%v: %v", string(r.Name), r.Value)
}

func main() {
	startTime := time.Now()

	f, err := os.Open(input)

	registers := make(map[rune]*Register)
	p1reg := make(map[rune]*Register)
	p2reg := make(map[rune]*Register)

	p1, p2 := initPrograms()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	scanner := bufio.NewScanner(f)

	insts := []Inst{}

	for scanner.Scan() {
		var txt = scanner.Text()
		inst := getInst(txt, registers)
		inst1 := getInst(txt, p1reg)
		inst2 := getInst(txt, p2reg)
		insts = append(insts, inst)

		p1.Insts = append(p1.Insts, inst1)
		p2.Insts = append(p2.Insts, inst2)
	}

	p1.Registers = p1reg
	p2.Registers = p2reg

	preg, _ := p2.Registers['p']
	preg.Value = 1

	rcv := &Register{Name: 'x', Value: 0}
	rval := handleInsts(insts, rcv)

	runPrograms(p1, p2)

	fmt.Println("Received value", rval)
	fmt.Println("Program 0 sent", p1.Sends, "times")
	fmt.Println("Program 1 sent", p2.Sends, "times")

	fmt.Println("Time", time.Since(startTime))
}

func initPrograms() (*Program, *Program) {
	p1 := &Program{Name: "Program 1", Deadlocked: false}
	p2 := &Program{Name: "Program 2", Deadlocked: false}

	p1.Parallel = p2
	p2.Parallel = p1

	return p1, p2
}

func runPrograms(p1, p2 *Program) {
	i := 0
	for !p1.Deadlocked || !p2.Deadlocked {
		inst1 := p1.Insts[p1.InstPtr]
		inst2 := p2.Insts[p2.InstPtr]

		jmp1 := handleInstPart2(inst1, p1)
		p1.InstPtr += jmp1
		p1.Deadlocked = jmp1 == 0

		jmp2 := handleInstPart2(inst2, p2)
		p2.InstPtr += jmp2
		p2.Deadlocked = jmp2 == 0

		i++
	}
}

func handleInstPart2(inst Inst, p *Program) int {
	jmp := 1

	switch inst.Cmd {
	case "snd":
		sndval := valueOfArg1(inst)
		p.Parallel.Queue = append(p.Parallel.Queue, sndval)
		p.Sends++
	case "set":
		setArg1(&inst, valueOfArg2(inst))
	case "add":
		setArg1(&inst, valueOfArg1(inst)+valueOfArg2(inst))
	case "mul":
		setArg1(&inst, valueOfArg1(inst)*valueOfArg2(inst))
	case "mod":
		setArg1(&inst, valueOfArg1(inst)%valueOfArg2(inst))
	case "rcv":
		if len(p.Queue) == 0 {
			jmp = 0
		} else {
			rcvval := p.Queue[0]
			p.Queue = p.Queue[1:]
			setArg1(&inst, rcvval)
		}
	case "jgz":
		t := valueOfArg1(inst)
		if t > 0 {
			jmp = valueOfArg2(inst)
		}

	}

	return jmp
}

func handleInsts(insts []Inst, rcv *Register) int {
	cur := 0
	val := 0
	for cur < len(insts) {
		inst := insts[cur]

		jmp, t := handleInst(inst, rcv)

		if t > 0 {
			val = t
			break
		}

		cur += jmp
	}

	return val
}

func handleInst(inst Inst, rcv *Register) (int, int) {
	jmp := 1
	ans := 0
	switch inst.Cmd {
	case "snd":
		rcv.Value = valueOfArg1(inst)
	case "set":
		setArg1(&inst, valueOfArg2(inst))
	case "add":
		setArg1(&inst, valueOfArg1(inst)+valueOfArg2(inst))
	case "mul":
		setArg1(&inst, valueOfArg1(inst)*valueOfArg2(inst))
	case "mod":
		setArg1(&inst, valueOfArg1(inst)%valueOfArg2(inst))
	case "rcv":
		ans = rcv.Value
	case "jgz":
		t := valueOfArg1(inst)
		if t > 0 {
			jmp = valueOfArg2(inst)
		}

	}
	return jmp, ans
}

func valueOfArg1(inst Inst) int {
	if inst.Arg1Reg != nil {
		return inst.Arg1Reg.Value
	} else {
		return inst.Arg1Val
	}
}

func setArg1(inst *Inst, val int) {
	if inst.Arg1Reg != nil {
		inst.Arg1Reg.Value = val
	} else {
		panic("Arg 1 is a static value - " + inst.Line)
	}

}

func valueOfArg2(inst Inst) int {
	if inst.Arg2Reg != nil {
		return inst.Arg2Reg.Value
	} else {
		return inst.Arg2Val
	}
}

func getInst(line string, registers map[rune]*Register) Inst {
	s := strings.Split(line, " ")

	inst := Inst{Cmd: s[0], Line: line}

	v1, r1 := getRegOrVal(s[1], registers)

	if r1 != nil {
		inst.Arg1Reg = r1
	} else {
		inst.Arg1Val = v1
	}

	if len(s) == 3 {
		v2, r2 := getRegOrVal(s[2], registers)

		if r2 != nil {
			inst.Arg2Reg = r2
		} else {
			inst.Arg2Val = v2
		}
	}

	return inst
}

func getRegOrVal(v string, registers map[rune]*Register) (int, *Register) {
	argInt, err := strconv.Atoi(v)

	if err != nil {
		c := rune(v[0])
		r, ok := registers[c]

		if !ok {
			r = &Register{Name: c, Value: 0}
			registers[c] = r
		}

		return 0, r
	}

	return argInt, nil
}
