package incline

import (
	"context"
	"fmt"
	"github.com/artarts36/orthography/incline/store"
	"log/slog"

	"github.com/artarts36/orthography/incline/word"
)

type Persistent struct {
	inclinator Inclinator
	store      store.Store
}

func NewPersistent(inclinator Inclinator, store store.Store) *Persistent {
	return &Persistent{
		inclinator: inclinator,
		store:      store,
	}
}

func (p *Persistent) InclineNouns(ctx context.Context, nouns []string) (map[string]*word.Word, error) {
	storeRes, err := p.store.Get(ctx, nouns)
	if err != nil {
		return storeRes.Found, fmt.Errorf("failed to get from store: %w", err)
	}

	if len(storeRes.NotFound) == 0 {
		return storeRes.Found, nil
	}

	newWords, err := p.inclinator.InclineNouns(ctx, storeRes.NotFound)
	if err != nil {
		return nil, err
	}

	if len(newWords) > 0 {
		slog.
			With(slog.Any("nouns", newWords)).
			DebugContext(ctx, "[orthography][inclinator] save new words in store")

		err = p.store.Save(ctx, newWords)
		if err != nil {
			return nil, fmt.Errorf("failed to save new words in store: %w", err)
		}
	}

	for _, w := range storeRes.Found {
		newWords[w.Nominative] = w
	}

	return newWords, nil
}
