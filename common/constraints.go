package common

type Number interface {
	int | float64 | int64 | int32 | uint
}

type Item[T any, N Number] struct {
	item  T
	value N
}
