package firstname

import (
	"context"
)

type GetResult struct {
	Found    map[string]*Name
	NotFound []string
}

type Store interface {
	Get(ctx context.Context, names []string) (*GetResult, error)
}
