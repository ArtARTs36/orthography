package firstname

import (
	"context"
	"fmt"
)

type FindResult GetResult

type GenderFinder interface {
	Find(ctx context.Context, names []string) (*Gender, error)
}

type StorableFinder struct {
	store Store
}

func NewStorableGenderFinder(store Store) *StorableFinder {
	return &StorableFinder{store: store}
}

func (f *StorableFinder) Find(ctx context.Context, names []string) (*FindResult, error) {
	storeNames, err := f.store.Get(ctx, names)
	if err != nil {
		return nil, fmt.Errorf("failed to find names: %w", err)
	}
	return &FindResult{
		Found:    storeNames.Found,
		NotFound: storeNames.NotFound,
	}, nil
}
