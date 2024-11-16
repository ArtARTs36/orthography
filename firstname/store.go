package firstname

import (
	"context"
)

type GetResult struct {
	Found    map[string]Gender
	NotFound []string
}

type Store interface {
	Get(ctx context.Context, names []string) (*GetResult, error)
}
