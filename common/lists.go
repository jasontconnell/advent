package common

func Permutate[T any](list []T) [][]T {
	var ret [][]T

	if len(list) == 2 {
		ret = append(ret, []T{list[0], list[1]})
		ret = append(ret, []T{list[1], list[0]})
	} else {

		for i := 0; i < len(list); i++ {
			c := make([]T, len(list))
			copy(c, list)

			t := c[i]
			sh := append(c[:i], c[i+1:]...)
			perm := Permutate(sh)

			for _, p := range perm {
				p = append([]T{t}, p...)
				ret = append(ret, p)
			}
		}
	}

	return ret
}

func CartesianProduct[T any](list1, list2 []T) [][]T {
	if len(list1) == 0 {
		return [][]T{list2}
	} else if len(list2) == 0 {
		return [][]T{list1}
	}
	flen := len(list1) * len(list2)
	nlist := make([][]T, flen)

	nidx := 0
	for i := 0; i < len(list2); i++ {
		for j := 0; j < len(list1); j++ {
			nlist[nidx] = append(nlist[nidx], list1[j], list2[i])
			nidx++
		}
	}

	return nlist
}

func AllCombinations[T any](list []T, count int) [][]T {
	indices := getCombinationsIndices(len(list)-1, count)
	ret := [][]T{}
	for _, x := range indices {
		next := []T{}
		for i := 0; i < len(x); i++ {
			next = append(next, list[x[i]])
		}
		ret = append(ret, next)
	}
	return ret
}

func getCombinationsIndices(max int, num int) [][]int {
	res := [][]int{}
	cur := make([]int, num)
	for i := 0; i < num; i++ {
		cur[i] = 0
	}
	res = append(res, cur)
	for {
		next, safe := addOne(cur, max)
		if !safe {
			break
		}
		res = append(res, next)
		cur = next
	}
	return res
}

func addOne(list []int, max int) ([]int, bool) {
	c := make([]int, len(list))
	copy(c, list)
	safe := true
	for j := len(c) - 1; j >= 0; j-- {
		if c[j] < max {
			c[j]++
			break
		} else if c[j] == max {
			if j == 0 {
				// overflow
				safe = false
				break
			} else {
				if c[j-1] == max {
					continue
				}
				c[j-1]++
				for i := j; i < len(c); i++ {
					c[i] = 0
				}
				break
			}
		}
	}
	return c, safe
}
