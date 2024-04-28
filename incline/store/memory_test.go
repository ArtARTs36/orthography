package store_test

import (
	"context"
	"testing"

	"github.com/artarts36/orthography/incline/store"
	"github.com/artarts36/orthography/incline/word"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMemoryTest_SaveAndGet(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		st := store.NewMemory()

		err := st.Save(context.Background(), map[string]*word.Word{})

		require.NoError(t, err)

		res, err := st.Get(context.Background(), []string{
			"test",
		})

		require.NoError(t, err)

		assert.Empty(t, res.Found)
		assert.Equal(t, []string{"test"}, res.NotFound)
	})

	t.Run("success", func(t *testing.T) {
		st := store.NewMemory()

		err := st.Save(context.Background(), map[string]*word.Word{
			"w": {
				Nominative: "w",
			},
		})

		require.NoError(t, err)

		res, err := st.Get(context.Background(), []string{
			"w",
		})

		require.NoError(t, err)

		assert.Equal(t, map[string]*word.Word{
			"w": {
				Nominative: "w",
			},
		}, res.Found)
		assert.Empty(t, res.NotFound)
	})
}
