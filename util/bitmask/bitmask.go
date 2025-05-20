// Package bitmask provides utility functions for manipulating bitmask.
package bitmask

import "golang.org/x/exp/constraints"

// Has checks if the flag is set in the mask.
func Has[T constraints.Unsigned](mask, flag T) bool {
	return (mask & flag) != 0
}

// Add a flag to the mask.
func Add[T constraints.Unsigned](mask, flag T) T {
	return mask | flag
}

// Remove a flag in the mask.
func Remove[T constraints.Unsigned](mask, flag T) T {
	return mask &^ flag
}

// Toggle a flag state in the mask.
func Toggle[T constraints.Unsigned](mask, flag T) T {
	return mask ^ flag
}
