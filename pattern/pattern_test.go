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

	indexer, err := NewPatternIndex("", stopWords, -1)
	if err != nil {
		panic(err)
	}
	indexer.InsertPrefixDocuments("biba bomba aboba", 0)
	indexer.InsertPrefixDocuments("bimba bomba aboba", 1)
	res, err := indexer.SearchByPrefix("bimb", 2)
	assert.Equal(t, err, nil)
	assert.Equal(t, res, []string{"bimba bomba aboba"})

	res, err = indexer.SearchByPrefix("abo", 100)
	assert.Equal(t, err, nil)
	assert.Equal(t, res, []string{"biba bomba aboba", "bimba bomba aboba"})

}

func TestPatternIndex(t *testing.T) {
	Rebuild()
	defer Rebuild()

	stopWords := []string{}

	indexer, err := NewPatternIndex("", stopWords, 3)
	if err != nil {
		println("aboba")
		panic(err)
	}
	indexer.InsertPatternDocuments("biba bomba aboba", 0)
	indexer.InsertPatternDocuments("bimba bomba aboba", 1)
	res, err := indexer.SearchByPattern("bim*a", 2)
	assert.Equal(t, err, nil)
	assert.Equal(t, res, []string{"bimba bomba aboba"})

	res, err = indexer.SearchByPattern("*imba", 2)
	assert.Equal(t, err, nil)
	assert.Equal(t, res, []string{"bimba bomba aboba"})

	res, err = indexer.SearchByPattern("ab*ba", 2)
	assert.Equal(t, err, nil)
	assert.Equal(t, res, []string{"biba bomba aboba", "bimba bomba aboba"})
}
