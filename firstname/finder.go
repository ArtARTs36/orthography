package firstname

import (
	"context"
	"fmt"
)

var ErrNameNotFound = fmt.Errorf("name not found")

type Finder interface {
	Find(ctx context.Context, name string) (*Name, error)
}

type StorableFinder struct {
	store Store
}

func NewStorableFinder(store Store) *StorableFinder {
	return &StorableFinder{store: store}
}

func (f *StorableFinder) Find(ctx context.Context, name string) (*Name, error) {
	res, err := f.store.Get(ctx, []string{name})
	if err != nil {
		return nil, fmt.Errorf("failed to find name: %w", err)
	}

	foundName, ok := res.Found[name]
	if !ok {
		return nil, ErrNameNotFound
	}

	return foundName, nil
}
