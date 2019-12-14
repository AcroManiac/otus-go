package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNormalizeText(t *testing.T) {
	result, _ := Sort([]int{1, 43, 5, 34, 2, 5, 6, 6, 10, 21}, 3)
	assert.Equal(t, result, []int{1, 2, 5, 5, 6, 6, 10, 21, 34, 43}, "Slices should be equal")
}
