package common

import "math"

func GetMinMax(list []int) (int, int) {
	min, max := math.MaxInt32, math.MinInt32

	for i := 0; i < len(list); i++ {
		if list[i] < min {
			min = list[i]
		}
		if list[i] > max {
			max = list[i]
		}
	}
	return min, max
}
