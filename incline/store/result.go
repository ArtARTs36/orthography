package store

import "github.com/artarts36/orthography/incline/word"

type GetResult struct {
	Found    map[string]*word.Word
	NotFound []string
}
