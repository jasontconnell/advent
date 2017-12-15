package knot

import (
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
)

func PrimitiveKnotHash(lengths []int) []int {
	vals := make([]int, 256)
	for i := 0; i < 256; i++ {
		vals[i] = i
	}

	_, _, res := knot(vals, lengths, 0, 0)

	return res
}

func KnotHash(lengths string) string {
	lens := getLengths(lengths)

	vals := make([]int, 256)
	for i := 0; i < 256; i++ {
		vals[i] = i
	}

	return hash(vals, lens)
}

func KnotHashBinary(lengths string) string {
	h := KnotHash(lengths)

	bin := ""
	for i := 0; i < len(h); i += 4 {
		nib := string(h[i : i+4])

		f, err := hex.DecodeString(nib)
		if err != nil {
			fmt.Println("decoding hex error", err)
			continue
		}

		v := int64(f[0])
		v = v << 8
		v |= int64(f[1])

		b := strconv.FormatInt(int64(v), 2)
		if len(b) < 16 {
			b = strings.Repeat("0", 16-len(b)) + b
		}

		bin += b
	}

	return bin
}

func getLengths(lstr string) []int {
	lens := []int{}
	for i := 0; i < len(lstr); i++ {
		lens = append(lens, int(lstr[i]))
	}
	lens = append(lens, []int{17, 31, 73, 47, 23}...)

	return lens
}

func hash(val, lens []int) string {
	pos := 0
	skip := 0
	for i := 0; i < 64; i++ {
		pos, skip, val = knot(val, lens, pos, skip)
	}

	xors := []int{}
	for i := 0; i < len(val); i += 16 {
		xors = append(xors, xor(val[i:i+16]))
	}

	result := ""
	for i := 0; i < len(xors); i++ {
		x := fmt.Sprintf("%x", xors[i])
		if len(x) == 1 {
			x = "0" + x
		}
		result += x

	}

	return result
}

func xor(val []int) int {
	result := val[0]
	for i := 1; i < len(val); i++ {
		result ^= val[i]
	}

	return result
}

func knot(val, lens []int, pos, skip int) (p int, s int, result []int) {
	for i := 0; i < len(lens); i++ {
		val = reverse(val, pos, pos+lens[i]-1)
		pos = (pos + lens[i] + skip) % len(val)
		skip = skip + 1
	}

	return pos, skip, val
}

func reverse(val []int, start, end int) []int {
	cp := make([]int, len(val))

	copy(cp, val)
	dst := start % len(val)

	for i := end; i >= start; i-- {
		cp[dst%len(val)] = val[i%len(val)]
		dst++
	}
	return cp
}
