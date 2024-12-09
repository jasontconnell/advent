package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/jasontconnell/advent/common"
)

type input = string
type output = int

type block struct {
	pos int
	id  *int
}

func (b block) String() string {
	var x string
	if b.id == nil {
		x = "nil"
	} else {
		x = strconv.Itoa(*b.id)
	}
	return fmt.Sprintf("block[%d]%s", b.pos, x)
}

type filesystem struct {
	blocks []block
}

func (fs filesystem) String() string {
	s := ""
	for _, b := range fs.blocks {
		c := "."
		if b.id != nil {
			c = fmt.Sprintf("%d", *b.id)
		}
		s += c
	}
	return s
}

func main() {
	in, err := common.ReadString(common.InputFilename(os.Args))
	if err != nil {
		log.Fatal(err)
	}

	p1, p1time := common.Time(part1, in)
	p2, p2time := common.Time(part2, in)

	w := common.TeeOutput(os.Stdout)
	fmt.Fprintln(w, "--2024 day 09 solution--")
	fmt.Fprintln(w, "Part 1:", p1)
	fmt.Fprintln(w, "Part 2:", p2)
	fmt.Printf("Time %v (%v, %v)", p1time+p2time, p1time, p2time)
}

func part1(in input) output {
	fs := parse(in)
	return defrag(fs)
}

func part2(in input) output {
	fs := parse(in)
	return defragWhole(fs)
}

func defrag(fs filesystem) int {
	ptr := 0
	for i := len(fs.blocks) - 1; i >= 0 && ptr < i-1; i-- {
		if fs.blocks[i].id != nil {
			for ptr < len(fs.blocks) && fs.blocks[ptr].id != nil {
				ptr++
			}
			if ptr < len(fs.blocks) {
				fs.blocks[ptr].id = fs.blocks[i].id
				fs.blocks[i].id = nil
			}
		}
	}
	return checksum(fs)
}

func defragWhole(fs filesystem) int {
	eptr := len(fs.blocks) - 1
	for eptr >= 0 {
		eidx, esize := nextBlock(fs, eptr)
		if esize == 0 {
			break
		}
		if eidx != -1 {
			sidx := nextFreeChunk(fs, esize)
			if sidx != -1 && sidx < eidx {
				n := 0
				for i := sidx; i < sidx+esize; i++ {
					fs.blocks[i].id = fs.blocks[eidx+n].id
					fs.blocks[eidx+n].id = nil
					n++
				}
			}
		}
		eptr = eidx - 1
	}
	return checksum(fs)
}

func nextFreeChunk(fs filesystem, size int) int {
	start, end := -1, -1
	for i := 0; i < len(fs.blocks); i++ {
		if fs.blocks[i].id == nil && start == -1 {
			start = i
			continue
		}
		if fs.blocks[i].id != nil && start != -1 {
			end = i
			if end-start >= size {
				break
			} else {
				start = -1
			}
		}
	}
	if end-start >= size {
		return start
	}
	return -1
}

func nextBlock(fs filesystem, ptr int) (int, int) {
	idx, size := -1, -1
	start, end := -1, -1
	id := -1
	i := ptr
	for start == -1 || end == -1 {
		cur := i
		i--
		if i < 0 {
			break
		}
		boundary := fs.blocks[cur].id == nil || (id != -1 && *fs.blocks[cur].id != id)
		if boundary {
			if idx != -1 {
				start = cur
				break
			}
			continue
		}
		if idx == -1 {
			id = *fs.blocks[cur].id
			idx = cur
			end = cur
		}
	}

	idx = start + 1
	size = end - start
	return idx, size
}

func checksum(fs filesystem) int {
	sum := 0
	for i := 0; i < len(fs.blocks); i++ {
		x := 0
		if fs.blocks[i].id != nil {
			x = *fs.blocks[i].id
		}
		sum += x * i
	}
	return sum
}

func parse(in string) filesystem {
	fs := filesystem{blocks: []block{}}
	id := 0
	free := false
	for i := 0; i < len(in); i++ {
		n, _ := strconv.Atoi(string(in[i]))
		for j := 0; j < n; j++ {
			b := block{pos: len(fs.blocks)}
			if !free {
				x := id
				b.id = &x
			}
			fs.blocks = append(fs.blocks, b)
		}
		free = !free
		if i%2 == 0 {
			id++
		}
	}
	return fs
}
