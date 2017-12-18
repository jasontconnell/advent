package main

import (
    "bufio"
    "fmt"
    "os"
    "time"
    "strconv"
    "strings"
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
    Name rune
    Value int
}

type Program struct {
    Registers map[rune]*Register
    Insts []Inst
    ID int
    Queue []int
    Send chan int
    Receive chan int
}

func (r *Register) String() string {
    return fmt.Sprintf("%v: %v", string(r.Name), r.Value)
}

func main() {
    startTime := time.Now()

    f, err := os.Open(input)


    registers := make(map[rune]*Register)

    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
    
    scanner := bufio.NewScanner(f)

    insts := []Inst{}
    for scanner.Scan() {
        var txt = scanner.Text()
        inst := getInst(txt, registers)
        insts = append(insts, inst)
    }

    rcv := &Register{ Name:'x', Value: 0}
    rval := handleInsts(insts,  rcv)

    fmt.Println("Received value", rval)

    fmt.Println("Time", time.Since(startTime))
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
            rcv.Value =  valueOfArg1(inst)
        case "set": setArg1(&inst, valueOfArg2(inst))
        case "add": setArg1(&inst, valueOfArg1(inst) + valueOfArg2(inst))
        case "mul": setArg1(&inst, valueOfArg1(inst) * valueOfArg2(inst))
        case "mod": setArg1(&inst, valueOfArg1(inst) % valueOfArg2(inst))
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

func setArg1(inst *Inst, val int){
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

    inst := Inst{ Cmd: s[0], Line: line }

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
            r = &Register{ Name: c, Value: 0}
            registers[c] = r
        }

        return 0, r
    }

    return argInt, nil
}