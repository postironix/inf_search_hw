package index

import (
	"github.com/RoaringBitmap/roaring"
	lsm "github.com/krasun/lsmtree"
)

type Index struct {
	tree      *lsm.LSMTree
	processor Processor
	proc      string
}

func NewIndex(proc string, stopWords []string) (*Index, error) {
	tree, err := lsm.Open("./data")
	if err != nil {
		return nil, err
	}
	return &Index{
		tree:      tree,
		processor: NewSimpleProcessor(stopWords),
		proc:      proc,
	}, nil
}

func (index *Index) AddDocument(text string, index_doc int) error {
	var processed_text []string
	var err error

	if index.proc == "stem" {
		processed_text, err = index.processor.Stem(text, true)
		if err != nil {
			return err
		}
	} else {
		processed_text, err = index.processor.Lem(text, true)
		if err != nil {
			return err
		}
	}

	for _, word := range processed_text {
		wordBytes := []byte(word)
		data, contains, err := index.tree.Get(wordBytes)
		if err != nil {
			return err
		}
		var bitmap *roaring.Bitmap
		if !contains {
			bitmap = roaring.NewBitmap()
		} else {
			bitmap = roaring.NewBitmap()
			err = bitmap.UnmarshalBinary(data)
			if err != nil {
				return err
			}
		}

		bitmap.Add(uint32(index_doc))
		data, err = bitmap.MarshalBinary()
		if err != nil {
			return err
		}

		index.tree.Put(wordBytes, data)
	}
	return nil
}

func (index *Index) AddBatchDocument(texts []string, index_docs []int) error {
	wordBitmaps := make(map[string]*roaring.Bitmap)

	for i, text := range texts {
		var processed_text []string
		var err error

		if index.proc == "stem" {
			processed_text, err = index.processor.Stem(text, true)
			if err != nil {
				return err
			}
		} else {
			processed_text, err = index.processor.Lem(text, true)
			if err != nil {
				return err
			}
		}

		for _, word := range processed_text {
			if _, exists := wordBitmaps[word]; !exists {
				wordBitmaps[word] = roaring.NewBitmap()
			}
			wordBitmaps[word].Add(uint32(index_docs[i]))
		}
	}

	for word, newBitmap := range wordBitmaps {
		wordBytes := []byte(word)
		data, contains, err := index.tree.Get(wordBytes)
		if err != nil {
			return err
		}
		var existingBitmap *roaring.Bitmap
		if contains {
			existingBitmap = roaring.NewBitmap()
			if err := existingBitmap.UnmarshalBinary(data); err != nil {
				return err
			}
			existingBitmap.Or(newBitmap)
		} else {
			existingBitmap = newBitmap
		}

		data, err = existingBitmap.MarshalBinary()
		if err != nil {
			return err
		}
		index.tree.Put(wordBytes, data)
	}
	return nil
}

func (index *Index) GetListDocuments(word string) ([]int, error) {
	bitmap, err := index.GetBitmapDocuments(word)
	if err != nil {
		return []int{}, err
	}
	uint32Array := bitmap.ToArray()
	intArray := make([]int, len(uint32Array))
	for i, v := range uint32Array {
		intArray[i] = int(v)
	}
	return intArray, nil
}

func (index *Index) GetBitmapDocuments(word string) (*roaring.Bitmap, error) {
	wordBytes := []byte(word)
	val, contains, err := index.tree.Get(wordBytes)
	bitmap := roaring.NewBitmap()
	if !contains {
		return bitmap, nil
	}
	if err != nil {
		return nil, err
	}
	err = bitmap.UnmarshalBinary(val)
	if err != nil {
		return nil, err
	}

	return bitmap, nil
}

func (index *Index) GetMergedBitmapDocuments(text string, limit int) (*roaring.Bitmap, error) {
	var processed_text []string
	var err error

	if index.proc == "stem" {
		processed_text, err = index.processor.Stem(text, true)
		if err != nil {
			return nil, err
		}
	} else {
		processed_text, err = index.processor.Lem(text, true)
		if err != nil {
			return nil, err
		}
	}

	merged := roaring.NewBitmap()
	total_size := 0
	for _, word := range processed_text {
		if total_size >= limit {
			break
		}
		word_result, err := index.GetBitmapDocuments(word)
		if err != nil {
			return nil, err
		}
		merged.Or(word_result)
		total_size += int(word_result.DenseSize())
	}

	return merged, nil
}

func (index *Index) GetMergedListsDocuments(text string, limit int) ([]int, error) {
	merged, err := index.GetMergedBitmapDocuments(text, limit)
	if err != nil {
		return []int{}, err
	}

	uint32Array := merged.ToArray()
	intArray := make([]int, len(uint32Array))
	for i, v := range uint32Array {
		intArray[i] = int(v)
	}

	return intArray, nil
}
