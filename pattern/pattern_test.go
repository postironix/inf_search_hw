package pattern

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Rebuild() {
	os.RemoveAll("./data")
	os.MkdirAll("./data", 0755)
}

func TestPrefixIndex(t *testing.T) {
	Rebuild()
	defer Rebuild()
	stopWords := []string{}

	index, err := NewPatternIndex("", stopWords, 3)
	if err != nil {
		panic(err)
	}
	index.InsertPrefixDocuments("biba bomba aboba", 0)
	index.InsertPrefixDocuments("bimba bomba aboba", 1)
	res, err := index.SearchByPrefix("bimb", 2)
	assert.Equal(t, err, nil)
	assert.Equal(t, res, []string{"bimba bomba aboba"})

	res, err = index.SearchByPrefix("abo", 100)
	assert.Equal(t, err, nil)
	assert.Equal(t, res, []string{"biba bomba aboba", "bimba bomba aboba"})

}

func TestPatternIndex(t *testing.T) {
	Rebuild()
	defer Rebuild()

	stopWords := []string{}

	index, err := NewPatternIndex("", stopWords, 3)
	if err != nil {
		println("aboba")
		panic(err)
	}
	index.InsertPatternDocuments("biba bomba aboba", 0)
	index.InsertPatternDocuments("bimba bomba aboba", 1)
	res, err := index.SearchByPattern("bim*a", 2)
	assert.Equal(t, err, nil)
	assert.Equal(t, res, []string{"bimba bomba aboba"})

	res, err = index.SearchByPattern("*imba", 2)
	assert.Equal(t, err, nil)
	assert.Equal(t, res, []string{"bimba bomba aboba"})

	res, err = index.SearchByPattern("ab*ba", 2)
	assert.Equal(t, err, nil)
	assert.Equal(t, res, []string{"biba bomba aboba", "bimba bomba aboba"})
}
