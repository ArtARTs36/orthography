package store

import (
	"context"
	"github.com/artarts36/orthography/incline/word"
)

type Store interface {
	All(ctx context.Context) (map[string]*word.Word, error)
	Get(ctx context.Context, nouns []string) (*GetResult, error)
	Save(ctx context.Context, words map[string]*word.Word) error
}
