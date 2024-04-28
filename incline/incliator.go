package incline

import (
	"context"

	"github.com/artarts36/orthography/incline/word"
)

type Inclinator interface {
	InclineNouns(ctx context.Context, nouns []string) (map[string]*word.Word, error)
}
