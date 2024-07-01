package slices

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFilter(t *testing.T) {
	s := []int{1, 2, 3, 4, 5}
	f := func(i int) bool {
		return i%2 == 0
	}

	r := Filter(s, f)
	assert.Equal(t, 2, len(r))
	assert.Contains(t, r, 2)
	assert.Contains(t, r, 4)
}

func TestFilterStruct(t *testing.T) {
	type user struct {
		Name string
		Age  int
	}

	s := []user{
		{Name: "Alice", Age: 20},
		{Name: "Bob", Age: 22},
		{Name: "Cathy", Age: 24},
	}
	f := func(u user) bool {
		return u.Age > 21
	}

	r := Filter(s, f)
	assert.Equal(t, 2, len(r))
	assert.Contains(t, r, user{Name: "Bob", Age: 22})
	assert.Contains(t, r, user{Name: "Cathy", Age: 24})
}

func TestMap(t *testing.T) {
	s := []int{1, 2, 3, 4, 5}
	f := func(i int) string {
		return "value: " + strconv.Itoa(i)
	}

	r := Map(s, f)
	assert.Equal(t, 5, len(r))
	assert.Contains(t, r, "value: 1")
	assert.Contains(t, r, "value: 2")
	assert.Contains(t, r, "value: 3")
	assert.Contains(t, r, "value: 4")
	assert.Contains(t, r, "value: 5")
}

func TestMapStruct(t *testing.T) {
	type user struct {
		Name string
		Age  int
	}

	s := []user{
		{Name: "Alice", Age: 20},
		{Name: "Bob", Age: 22},
		{Name: "Cathy", Age: 24},
	}
	f := func(u user) string {
		return "name: " + u.Name + ", age: " + strconv.Itoa(u.Age)
	}

	r := Map(s, f)
	assert.Equal(t, 3, len(r))
	assert.Contains(t, r, "name: Alice, age: 20")
	assert.Contains(t, r, "name: Bob, age: 22")
	assert.Contains(t, r, "name: Cathy, age: 24")
}
