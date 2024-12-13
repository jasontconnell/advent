package common

import "math"

func Min[N Number](x N, y N) N {
	return N(math.Min(float64(x), float64(y)))
}

func Max[N Number](x N, y N) N {
	return N(math.Max(float64(x), float64(y)))
}
