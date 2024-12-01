package index

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Rebuild() {
	os.RemoveAll("./data")
	os.MkdirAll("./data", 0755)
}

func TestIndex(t *testing.T) {
	Rebuild()
	defer Rebuild()
	stopWords := []string{"the", "is", "at", "which", "on"}

	index, _ := NewIndex("stem", stopWords)
	index.AddDocument("Sets are fundamental data structures in computer science. Basically set of integers can be stored as array containing elements", 0)
	index.AddDocument("There are many RLE based compression algorithms developed like WAH, EWAH, COMPAX. But all of them lack fast AND, OR, XOR and other operations required to implement fast index operations. This is where roaring bitmaps shines and differ. Roaring bitmaps dosent perform optimal compression, however in most cases it does fairly good job, however in exchange it performs index operations blazingly fast. Also it only decompresses only parts which are required.", 1)

	result, err := index.GetListDocuments("set")
	assert.Equal(t, nil, err)
	assert.Equal(t, result, []int{0})

	result, err = index.GetListDocuments("are")
	assert.Equal(t, err, nil)
	assert.Equal(t, result, []int{0, 1})
}

func TestBatchIndex(t *testing.T) {
	Rebuild()
	defer Rebuild()

	Rebuild()
	defer Rebuild()
	stopWords := []string{"the", "is", "at", "which", "on"}

	index, _ := NewIndex("stem", stopWords)
	index.AddBatchDocument([]string{
		"Sets are fundamental data structures in computer science. Basically set of integers can be stored as array containing elements",
		"There are many RLE based compression algorithms developed like WAH, EWAH, COMPAX. But all of them lack fast AND, OR, XOR and other operations required to implement fast index operations. This is where roaring bitmaps shines and differ. Roaring bitmaps dosent perform optimal compression, however in most cases it does fairly good job, however in exchange it performs index operations blazingly fast. Also it only decompresses only parts which are required.",
	}, []int{0, 1})

	result, err := index.GetListDocuments("set")
	assert.Equal(t, err, nil)
	assert.Equal(t, result, []int{0})

	result, err = index.GetListDocuments("are")
	assert.Equal(t, err, nil)
	assert.Equal(t, result, []int{0, 1})
}

func TestMergeIndex(t *testing.T) {
	Rebuild()
	defer Rebuild()

	Rebuild()
	defer Rebuild()
	stopWords := []string{"the", "is", "at", "which", "on"}

	index, _ := NewIndex("stem", stopWords)
	index.AddBatchDocument([]string{
		"aboba biba boba",
		"biba boba",
	}, []int{0, 1})

	result, err := index.GetMergedListsDocuments("aboba", 2)
	assert.Equal(t, err, nil)
	assert.Equal(t, result, []int{0})

	result, err = index.GetMergedListsDocuments("biba boba", 2)
	assert.Equal(t, err, nil)
	assert.Equal(t, result, []int{0, 1})
}
