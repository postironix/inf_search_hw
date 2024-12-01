package pattern

import (
	"search/index"
	"strconv"
	"strings"

	"github.com/RoaringBitmap/roaring"
	lsm "github.com/krasun/lsmtree"
)

type PatternIndex struct {
	index       *index.Index
	tree        *lsm.LSMTree
	coef_n_gram int
}

func NewPatternIndex(proc string, stopWords []string, coef int) (*PatternIndex, error) {
	tree, err := lsm.Open("./data")
	if err != nil {
		return nil, err
	}
	index, err := index.NewIndex(proc, stopWords)
	if err != nil {
		return nil, err
	}
	return &PatternIndex{
		index:       index,
		tree:        tree,
		coef_n_gram: coef,
	}, nil
}

func (pi *PatternIndex) InsertPrefixDocuments(text string, indexDoc int) {
	lowerText := strings.ToLower(text)
	words := strings.Fields(lowerText)
	var allPrefixes []string
	for _, word := range words {
		prefixes := make([]string, len(word))
		for i := 1; i <= len(word); i++ {
			prefixes[i-1] = word[:i]
		}
		allPrefixes = append(allPrefixes, prefixes...)
	}
	indices := make([]int, len(allPrefixes))
	for i := 0; i < len(allPrefixes); i++ {
		indices[i] = indexDoc
	}
	pi.index.AddBatchDocument(allPrefixes, indices)
	pi.tree.Put([]byte(strconv.Itoa(indexDoc)), []byte(text))
}

func (pi *PatternIndex) InsertPatternDocuments(text string, indexDoc int) {
	lower := strings.ToLower(text)
	words := strings.Fields(lower)
	var allNgrams []string
	for _, word := range words {
		ngrams := make([]string, 0)
		for i := 0; i < len(word); i++ {
			for j := i + 1; j <= len(word) && j-i <= pi.coef_n_gram; j++ {
				ngrams = append(ngrams, word[i:j])
			}
		}
		allNgrams = append(allNgrams, ngrams...)
	}

	indices := make([]int, len(allNgrams))
	for i := 0; i < len(allNgrams); i++ {
		indices[i] = indexDoc
	}
	pi.index.AddBatchDocument(allNgrams, indices)
	pi.tree.Put([]byte(strconv.Itoa(indexDoc)), []byte(text))
}

func (pi *PatternIndex) SearchByPrefix(prefix string, limit int) ([]string, error) {
	lower := strings.ToLower(prefix)
	indices, err := pi.index.GetListDocuments(lower)
	if err != nil {
		return []string{}, err
	}
	docs := make([]string, 0)
	for _, ind := range indices {
		doc, contains, err := pi.tree.Get([]byte(strconv.Itoa(ind)))
		if err != nil {
			return nil, err
		}
		if !contains {
			continue
		}
		docs = append(docs, string(doc))
	}
	return docs, nil
}

func (pi *PatternIndex) SearchByPattern(pattern string, limit int) ([]string, error) {
	lower := strings.ToLower(pattern)

	var words []string
	parts := strings.Split(lower, "*")

	for _, part := range parts {
		if len(part) <= limit {
			words = append(words, part)
		} else {
			for i := 0; i <= len(part)-limit; i++ {
				words = append(words, part[i:i+limit])
			}
		}
	}

	bitmap := roaring.NewBitmap()
	for _, word := range words {
		cur_bitmap, err := pi.index.GetMergedBitmapDocuments(word, limit)
		if err != nil {
			return []string{}, err
		}
		bitmap.Or(cur_bitmap)
	}

	docs := make([]string, 0)
	for _, ind := range bitmap.ToArray() {
		doc, contains, err := pi.tree.Get([]byte(strconv.Itoa(int(ind))))
		if err != nil {
			return nil, err
		}
		if !contains {
			continue
		}
		if !MatchPatternToText(string(doc), pattern) {
			continue
		}
		docs = append(docs, string(doc))
	}

	return docs, nil
}

func CheckPattern(word, pattern string) bool {
	parts := strings.Split(pattern, "*")
	start := 0

	for _, part := range parts {
		if part == "" {
			continue
		}
		index := strings.Index(word[start:], part)
		if index == -1 {
			return false
		}

		start += index + len(part)
	}
	// когда суффикс не совпадает с паттерном полностью
	if !strings.HasSuffix(pattern, "*") && start != len(word) {
		return false
	}

	return true
}

func MatchPatternToText(text, pattern string) bool {
	words := strings.Fields(text)
	for _, word := range words {
		if CheckPattern(word, pattern) {
			return true
		}
	}
	return false
}
