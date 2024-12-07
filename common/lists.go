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
