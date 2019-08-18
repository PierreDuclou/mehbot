package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChunk(t *testing.T) {
	assert := assert.New(t)
	buf := []string{"first", "second", "third", "fourth"}
	expected := [][]string{
		[]string{"first", "second"},
		[]string{"third", "fourth"},
	}

	assert.Equal(expected, Chunk(buf, 2), "they should be equal")
}
