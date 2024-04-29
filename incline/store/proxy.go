package store

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/artarts36/orthography/incline/word"
)

type Proxy struct {
	hot  Store
	cold Store
}

func NewProxy(hot Store, cold Store) *Proxy {
	return &Proxy{
		hot:  hot,
		cold: cold,
	}
}

func (p *Proxy) All(ctx context.Context) (map[string]*word.Word, error) {
	return p.cold.All(ctx)
}

func (p *Proxy) Get(ctx context.Context, nouns []string) (*GetResult, error) {
	if len(nouns) == 0 {
		return &GetResult{
			Found:    map[string]*word.Word{},
			NotFound: []string{},
		}, nil
	}

	hotRes, err := p.hot.Get(ctx, nouns)
	if err != nil {
		return nil, fmt.Errorf("failed to get words in hot store: %w", err)
	}

	if len(hotRes.Found) == len(nouns) {
		slog.DebugContext(ctx, "[orthography][proxy-store] hot store contains all words")

		return hotRes, nil
	}

	if len(hotRes.Found) == 0 {
		slog.DebugContext(ctx, "[orthography][proxy-store] hot store does not contain a single word")
	} else {
		slog.DebugContext(
			ctx,
			fmt.Sprintf("[orthography][proxy-store] hot store contains %d words", len(hotRes.Found)),
		)
	}

	coldRes, err := p.cold.Get(ctx, hotRes.NotFound)
	if err != nil {
		return hotRes, fmt.Errorf("failed to get words in cold store: %w", err)
	}

	slog.DebugContext(
		ctx,
		fmt.Sprintf("[orthography][proxy-store] saving %d words in hot store", len(coldRes.Found)),
	)

	err = p.hot.Save(ctx, coldRes.Found)
	if err != nil {
		return hotRes, fmt.Errorf("failed to save words in hot store: %w", err)
	}

	for _, w := range coldRes.Found {
		hotRes.Found[w.Nominative] = w
	}

	hotRes.NotFound = coldRes.NotFound

	return hotRes, nil
}

func (p *Proxy) Save(ctx context.Context, words map[string]*word.Word) error {
	err := p.hot.Save(ctx, words)
	if err != nil {
		return fmt.Errorf("failed to save words in hot store: %w", err)
	}

	err = p.cold.Save(ctx, words)
	if err != nil {
		return fmt.Errorf("failed to save words in cold store: %w", err)
	}

	return nil
}

func (p *Proxy) WarmUp(ctx context.Context) (int, error) {
	words, err := p.cold.All(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to get all words from cold storage: %w", err)
	}

	err = p.hot.Save(ctx, words)
	if err != nil {
		return 0, fmt.Errorf("failed to save all words to hot storage: %w", err)
	}

	return len(words), nil
}

func (p *Proxy) WarmUpPeriodically(ctx context.Context, interval time.Duration) {
	warmup := func() {
		slog.DebugContext(ctx, "[orthography][proxy-store] try warmup")

		c, err := p.WarmUp(ctx)
		if err != nil {
			slog.
				With(slog.String("err", err.Error())).
				ErrorContext(ctx, "[orthography][proxy-store] failed to warmup")
			return
		}

		slog.DebugContext(ctx, fmt.Sprintf("[orthography][proxy-store] updated %d words in hot storage", c))
	}

	warmup()

	tick := time.NewTicker(interval).C

	for {
		select {
		case <-ctx.Done():
			slog.InfoContext(ctx, "[orthography][proxy-store] warmup stopped")
			return
		case <-tick:
			warmup()
		}
	}
}
