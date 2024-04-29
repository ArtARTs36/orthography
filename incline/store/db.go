package store

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/artarts36/orthography/incline/word"
)

type DB struct {
	tableName string
	conn      *sql.DB
}

func NewDB(tableName string, conn *sql.DB) *DB {
	return &DB{
		tableName: tableName,
		conn:      conn,
	}
}

func (d *DB) All(ctx context.Context) (map[string]*word.Word, error) {
	q := fmt.Sprintf("SELECT * FROM %s", d.tableName)

	rows, err := d.conn.QueryContext(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	if rowsErr := rows.Err(); rowsErr != nil {
		return nil, fmt.Errorf("failed to execute query: %w", rowsErr)
	}

	words := map[string]*word.Word{}

	for rows.Next() {
		var w word.Word

		err = rows.Scan(&w.Nominative, &w.Genitive, &w.Dative, &w.Accusative, &w.Instrumental, &w.Prepositional)
		if err != nil {
			return nil, fmt.Errorf("failed to scan struct: %w", err)
		}

		words[w.Nominative] = &w
	}

	return words, nil
}

func (d *DB) Get(ctx context.Context, nouns []string) (*GetResult, error) {
	if len(nouns) == 0 {
		return &GetResult{
			Found:    map[string]*word.Word{},
			NotFound: []string{},
		}, nil
	}

	placeholders := make([]string, 0, len(nouns))
	args := make([]interface{}, 0, len(nouns))

	for i, noun := range nouns {
		placeholders = append(placeholders, fmt.Sprintf("$%d", i+1))
		args = append(args, noun)
	}

	q := fmt.Sprintf(
		"SELECT * FROM %s WHERE nominative IN (%s)",
		d.tableName,
		strings.Join(placeholders, ","),
	)

	rows, err := d.conn.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	if rowsErr := rows.Err(); rowsErr != nil {
		return nil, fmt.Errorf("failed to execute query: %w", rowsErr)
	}

	words := map[string]*word.Word{}

	for rows.Next() {
		var w word.Word

		err = rows.Scan(&w.Nominative, &w.Genitive, &w.Dative, &w.Accusative, &w.Instrumental, &w.Prepositional)
		if err != nil {
			return nil, fmt.Errorf("failed to scan struct: %w", err)
		}

		words[w.Nominative] = &w
	}

	if len(words) == len(nouns) {
		return &GetResult{
			Found:    words,
			NotFound: []string{},
		}, nil
	}

	nf := make([]string, 0, len(nouns)-len(words))
	for _, noun := range nouns {
		if _, found := words[noun]; !found {
			nf = append(nf, noun)
		}
	}

	return &GetResult{
		Found:    words,
		NotFound: nf,
	}, nil
}

func (d *DB) Save(ctx context.Context, nouns map[string]*word.Word) error {
	const columnsCount = 6

	placeholders := make([]string, 0, len(nouns))
	vals := make([]interface{}, 0, len(nouns)*columnsCount)

	var i int64 = 1
	pID := func() int64 {
		v := i
		i++

		return v
	}

	for _, row := range nouns {
		placeholders = append(placeholders, fmt.Sprintf(
			"($%d,$%d,$%d,$%d,$%d,$%d)", pID(), pID(), pID(), pID(), pID(), pID(),
		))

		vals = append(
			vals,
			row.Nominative, row.Genitive, row.Dative, row.Accusative, row.Instrumental, row.Prepositional,
		)
	}

	q := fmt.Sprintf(
		"INSERT INTO words (nominative, genitive, dative, accusative, instrumental, prepositional) VALUES %s",
		strings.Join(placeholders, ","),
	)

	_, err := d.conn.ExecContext(
		ctx,
		q,
		vals...,
	)
	if err != nil {
		return fmt.Errorf("failed to exec statement: %w", err)
	}

	return nil
}
