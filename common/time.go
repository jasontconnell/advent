package common

import "time"

func Time[T any, Q any](fn func(in T) Q, in T) (Q, time.Duration) {
	start := time.Now()
	v := fn(in)
	return v, time.Since(start)
}
