package index

import (
	"strings"

	"github.com/jdkato/prose/v2"
	"github.com/kljensen/snowball"
)

type Processor interface {
	Lem(text string, remove bool) ([]string, error)
	Stem(text string, remove bool) ([]string, error)
}

type SimpleProcessor struct {
	StopWords map[string]bool
}

func (proc *SimpleProcessor) Lem(text string, remove bool) ([]string, error) {
	doc, err := prose.NewDocument(text)
	if err != nil {
		return nil, err
	}

	lemmatizedWords := make([]string, len(doc.Tokens()))
	for i, token := range doc.Tokens() {
		word := strings.ToLower(token.Text)
		lemmatizedWords[i] = proc.lemmatizeSimple(word)
	}

	var result []string
	for _, word := range lemmatizedWords {
		if !proc.StopWords[word] {
			result = append(result, word)
		}
	}

	return result, nil
}

func (proc *SimpleProcessor) lemmatizeSimple(word string) string {
	if strings.HasSuffix(word, "ing") {
		return strings.TrimSuffix(word, "ing")
	}
	return word
}

func (proc *SimpleProcessor) Stem(text string, remove bool) ([]string, error) {
	words := strings.Fields(text)
	var stemmedWords []string

	for _, word := range words {
		stemmedWord, err := snowball.Stem(word, "english", true)
		if err != nil {
			return nil, err
		}
		stemmedWords = append(stemmedWords, stemmedWord)
	}

	var result []string
	for _, word := range stemmedWords {
		if !proc.StopWords[word] {
			result = append(result, word)
		}
	}

	return result, nil
}

func NewSimpleProcessor(stopWords []string) *SimpleProcessor {
	proc := SimpleProcessor{}
	proc.StopWords = make(map[string]bool, 0)
	for _, word := range stopWords {
		proc.StopWords[word] = true
	}
	return &proc
}
