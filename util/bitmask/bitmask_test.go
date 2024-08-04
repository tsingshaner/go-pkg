package bitmask

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type bitmaskCase struct {
	mask uint
	flag uint
}

func TestHas(t *testing.T) {
	cases := []struct {
		bitmaskCase
		want bool
	}{
		{bitmaskCase{1, 1}, true},
		{bitmaskCase{1, 2}, false},
		{bitmaskCase{3, 1}, true},
		{bitmaskCase{3, 2}, true},
	}

	for _, c := range cases {
		assert.Equal(t, c.want, Has(c.mask, c.flag), "mask: %d, flag: %d", c.mask, c.flag)
	}
}

func TestToggle(t *testing.T) {
	cases := []struct {
		bitmaskCase
		want uint
	}{
		{bitmaskCase{1, 1}, 0},
		{bitmaskCase{1, 2}, 3},
		{bitmaskCase{3, 1}, 2},
		{bitmaskCase{3, 2}, 1},
	}

	for _, c := range cases {
		assert.Equal(t, c.want, Toggle(c.mask, c.flag), "mask: %d, flag: %d", c.mask, c.flag)
	}
}
