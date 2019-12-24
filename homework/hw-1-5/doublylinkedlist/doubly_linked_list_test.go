package doublylinkedlist

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testValues struct {
	values []interface{}
}

// Factory for list creation and filling
func createList(tests []testValues, back bool) List {
	list := List{}
	for _, test := range tests {
		for _, v := range test.values {
			switch back {
			case true:
				list.PushBack(v)
			case false:
				list.PushFront(v)
			}
		}
	}
	return list
}

func TestPushBack(t *testing.T) {
	tests := []testValues{
		{[]interface{}{2, 12, 85, 06}},
	}

	list := createList(tests, true)

	assert.Equal(t, len(tests[0].values), list.Len(), "Lists lengths should be equal")
	assert.Equal(t, tests[0].values[0], list.First().Value(), "First items should be equal")
	assert.Equal(t, tests[0].values[len(tests[0].values)-1], list.Last().Value(), "Last items should be equal")
}

func TestPushFront(t *testing.T) {
	tests := []testValues{
		{[]interface{}{"это", "твой", "номер", "номер", "номер"}},
	}

	list := createList(tests, false)

	assert.Equal(t, len(tests[0].values), list.Len(), "Lists lengths should be equal")
	assert.Equal(t, tests[0].values[len(tests[0].values)-1], list.First().Value(), "First items should be equal")
	assert.Equal(t, tests[0].values[0], list.Last().Value(), "Last items should be equal")
}

func TestRemove(t *testing.T) {
	tests := []struct {
		values []float64
	}{
		{[]float64{math.Pi, math.Phi, math.E, math.MaxFloat64, math.SmallestNonzeroFloat64}},
	}

	// Create and fill list with values
	list := List{}
	for _, test := range tests {
		for _, v := range test.values {
			list.PushBack(v)
		}
	}

	i := list.First() // get first item
	list.Remove(*i)   // remove first item from list

	i = list.Last() // get last item
	list.Remove(*i) // remove last item from list

	assert.Equalf(t, len(tests[0].values)-2, list.Len(), "List length should be %d", len(tests[0].values)-2)
	assert.Equalf(t, tests[0].values[1], list.First().Value(), "First item should be %f", tests[0].values[1])
	assert.Equalf(t, tests[0].values[3], list.Last().Value(), "Last item should be %f", tests[0].values[3])
}
