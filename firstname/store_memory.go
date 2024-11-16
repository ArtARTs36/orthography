package firstname

import (
	"context"
	"strings"
)

type MemoryStore struct {
	names map[string]*Name
}

func NewMemoryStore(names map[string]*Name) *MemoryStore {
	return &MemoryStore{names: names}
}

func (s *MemoryStore) Get(_ context.Context, names []string) (*GetResult, error) {
	result := &GetResult{
		Found:    map[string]*Name{},
		NotFound: make([]string, 0),
	}

	for _, n := range names {
		if name, ok := s.names[strings.ToLower(n)]; ok {
			result.Found[n] = name
		} else {
			result.NotFound = append(result.NotFound, n)
		}
	}

	return result, nil
}
