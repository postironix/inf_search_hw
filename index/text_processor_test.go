package index

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLemProcessing(t *testing.T) {
	stopWords := []string{"the", "is", "at", "which", "on"}
	processor := NewSimpleProcessor(stopWords)
	result, err := processor.Lem("the cringe based", true)

	assert.Equal(t, err, nil)
	assert.Equal(t, result, []string{"cringe", "based"})
}

func TestStemProcessing(t *testing.T) {
	stopWords := []string{"the", "is", "at", "which", "on"}
	processor := NewSimpleProcessor(stopWords)
	result, err := processor.Lem("the cringe based lemming", true)

	assert.Equal(t, err, nil)
	assert.Equal(t, result, []string{"cringe", "based", "lemm"})
}
