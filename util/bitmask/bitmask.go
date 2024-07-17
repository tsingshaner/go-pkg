package bitmask

import "golang.org/x/exp/constraints"

func Has[T constraints.Integer](mask, flag T) bool {
	return (mask & flag) == flag
}

func Add[T constraints.Integer](mask, flag T) T {
	return mask | flag
}

func Remove[T constraints.Integer](mask, flag T) T {
	return mask &^ flag
}

func Toggle[T constraints.Integer](mask, flag T) T {
	return mask ^ flag
}
