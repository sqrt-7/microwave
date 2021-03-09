package tools

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSameElements(t *testing.T) {
	test := assert.New(t)

	test.True(SameElements(
		[]string{},
		[]string{},
	))

	test.True(SameElements(
		[]string{"A", "B", "C", "D", "E"},
		[]string{"A", "B", "C", "D", "E"},
	))

	test.True(SameElements(
		[]string{"A", "B", "C", "D", "E"},
		[]string{"D", "C", "E", "A", "B"},
	))

	test.True(SameElements(
		[]string{"A", "A", "A", "B", "C", "D", "E"},
		[]string{"A", "A", "A", "B", "C", "D", "E"},
	))

	test.True(SameElements(
		[]string{"AA", "AA", "AA", "BB", "CC", "DD", "EE"},
		[]string{"AA", "AA", "AA", "BB", "CC", "DD", "EE"},
	))

	test.True(SameElements(
		[]string{"A", "A", "A", "A", "A"},
		[]string{"A", "A", "A", "A", "A"},
	))

	test.False(SameElements(
		[]string{"A", "B", "C", "D", "E"},
		[]string{"A", "C", "E", "B", "D", "F"},
	))

	test.False(SameElements(
		[]string{"A", "B", "C", "D", "E", "F"},
		[]string{"A", "C", "E", "B", "D"},
	))

	test.False(SameElements(
		[]string{"A", "B", "C", "D", "E"},
		[]string{"A", "A", "A", "B", "C", "D", "E"},
	))

	test.False(SameElements(
		[]string{"A", "A", "A", "B", "C", "D", "E"},
		[]string{"A", "B", "C", "D", "E"},
	))
}

func TestUniqueStringSlice(t *testing.T) {
	test := assert.New(t)

	{
		expected := []string{"A", "B", "C", "D", "E"}
		actual := UniqueStringSlice([]string{"A", "B", "C", "D", "E"})
		sort.Strings(expected)
		sort.Strings(actual)
		test.Equal(expected, actual)
	}

	{
		expected := []string{"A", "B", "C", "D", "E"}
		actual := UniqueStringSlice([]string{"A", "A", "A", "A", "B", "C", "D", "E"})
		sort.Strings(expected)
		sort.Strings(actual)
		test.Equal(expected, actual)
	}

	{
		expected := []string{"A"}
		actual := UniqueStringSlice([]string{"A", "A", "A", "A", "A", "A", "A", "A"})
		sort.Strings(expected)
		sort.Strings(actual)
		test.Equal(expected, actual)
	}

	{
		expected := []string{"A", "B", "C", "D", "E"}
		actual := UniqueStringSlice([]string{"A", "A", "A", "A", "B", "B", "C", "C", "C", "C", "D", "E", "E", "E"})
		sort.Strings(expected)
		sort.Strings(actual)
		test.Equal(expected, actual)
	}

	{
		expected := []string{"A", "B", "C", "D", "E"}
		actual := UniqueStringSlice([]string{"A", "B", "C", "D", "E", "X"})
		sort.Strings(expected)
		sort.Strings(actual)
		test.NotEqual(expected, actual)
	}

	{
		expected := []string{"A", "B", "C", "D", "E", "X"}
		actual := UniqueStringSlice([]string{"A", "B", "C", "D", "E"})
		sort.Strings(expected)
		sort.Strings(actual)
		test.NotEqual(expected, actual)
	}
}
