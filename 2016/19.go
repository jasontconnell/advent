package main

import (
	"fmt"
	"time"
)

var input = 3014603

type Elf struct {
	Number   int
	Presents int
	Next     *Elf
	Prev     *Elf
}

func (elf *Elf) String() string {
	return fmt.Sprintf("number: %v, presents: %v. Next num: %v, Prev num %v", elf.Number, elf.Presents, elf.Next.Number, elf.Prev.Number)
}

func main() {
	startTime := time.Now()

	elf := &Elf{Number: 1, Presents: 1, Next: nil, Prev: nil}
	first := elf
	for i := 1; i < input; i++ {
		next := Elf{Number: i + 1, Presents: 1, Next: nil, Prev: elf}
		elf.Next = &next
		elf = &next
		if i == input-1 {
			next.Next = first
			first.Prev = &next
		}
	}

	part2 := true

	elf = first

	solved := false
	var across *Elf = first
	for i := 0; i < input/2; i++ {
		across = across.Next
	}
	count := input

	for !solved {
		if elf.Presents > 0 {
			if part2 {
				elf.Presents += across.Presents
				across.Next.Prev = across.Prev
				across.Prev.Next = across.Next
			} else {
				elf.Presents += elf.Next.Presents
				elf.Next = elf.Next.Next
				elf.Next.Prev = elf
			}
		}
		count--
		elf = elf.Next
		across = across.Next
		if count%2 == 0 {
			across = across.Next
		}
		solved = elf.Next == elf
	}

	fmt.Println(elf)
	fmt.Println("Time", time.Since(startTime))
}
