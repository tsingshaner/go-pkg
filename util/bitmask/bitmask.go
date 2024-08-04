package bitmask

import "golang.org/x/exp/constraints"

func Has[T constraints.Unsigned](mask, flag T) bool {
	return (mask & flag) != 0
}

func Add[T constraints.Unsigned](mask, flag T) T {
	return mask | flag
}

func Remove[T constraints.Unsigned](mask, flag T) T {
	return mask &^ flag
}

func Toggle[T constraints.Unsigned](mask, flag T) T {
	return mask ^ flag
}
