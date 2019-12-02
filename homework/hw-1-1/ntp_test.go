package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTime(t *testing.T) {
	_, err := Time()
	assert.Nil(t, err)
}
