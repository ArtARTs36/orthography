package store

import (
	"context"
	"strings"

	"github.com/artarts36/orthography/incline/word"
)

type Memory struct {
	words map[string]*word.Word
}

func NewMemory() *Memory {
	return &Memory{
		words: map[string]*word.Word{},
	}
}

func (m *Memory) All(_ context.Context) (map[string]*word.Word, error) {
	return m.words, nil
}

func (m *Memory) Get(_ context.Context, nouns []string) (*GetResult, error) {
	res := &GetResult{
		Found:    map[string]*word.Word{},
		NotFound: []string{},
	}

	for _, noun := range nouns {
		n := strings.ToLower(noun)

		w, exists := m.words[n]
		if exists {
			res.Found[w.Nominative] = w
		} else {
			res.NotFound = append(res.NotFound, n)
		}
	}

	return res, nil
}

func (m *Memory) Save(_ context.Context, words map[string]*word.Word) error {
	for _, w := range words {
		m.words[strings.ToLower(w.Nominative)] = w
	}

	return nil
}
