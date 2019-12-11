package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnpackString(t *testing.T) {
	result, _ := unpackString("a2b3c4d5")
	assert.Equal(t, result, "aabbbccccddddd", "Strings should be equal")

	result, _ = unpackString(`abc\r3`)
	assert.Equal(t, result, "abcrrr", "Strings should be equal")

	_, err := unpackString("789")
	assert.NotNil(t, err)
}
