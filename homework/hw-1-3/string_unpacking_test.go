package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnpackString(t *testing.T) {
	result, _ := unpackString("a2b3c4d5")
	assert.Equal(t, "aabbbccccddddd", result, "Strings should be equal")

	result, _ = unpackString(`abc\r3`)
	assert.Equal(t, "abcrrr", result, "Strings should be equal")

	result, _ = unpackString(`abc\\3`)
	assert.Equal(t, `abc\\\`, result, "Strings should be equal")

	result, _ = unpackString(`abc\3\4`)
	assert.Equal(t, "abc34", result, "Strings should be equal")

	_, err := unpackString("789")
	assert.NotNil(t, err)
}
