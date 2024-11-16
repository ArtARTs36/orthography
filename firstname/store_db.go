package firstname

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
)

type DBStore struct {
	tableName string
	conn      *sql.DB
}

func NewDBStore(tableName string, conn *sql.DB) *DBStore {
	return &DBStore{
		tableName: tableName,
		conn:      conn,
	}
}

func (d *DBStore) Get(ctx context.Context, names []string) (*GetResult, error) {
	if len(names) == 0 {
		return &GetResult{
			Found:    map[string]*Name{},
			NotFound: []string{},
		}, nil
	}

	placeholders := make([]string, 0, len(names))
	args := make([]interface{}, 0, len(names))

	for i, name := range names {
		placeholders = append(placeholders, fmt.Sprintf("$%d", i+1))
		args = append(args, name)
	}

	q := fmt.Sprintf(
		"SELECT name, gender FROM %s WHERE name IN (%s)",
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

	namesMap := map[string]*Name{}

	for rows.Next() {
		var n Name

		err = rows.Scan(&n.Name, &n.Gender)
		if err != nil {
			return nil, fmt.Errorf("failed to scan struct: %n", err)
		}

		namesMap[n.Name] = &n
	}

	if len(namesMap) == len(names) {
		return &GetResult{
			Found:    namesMap,
			NotFound: []string{},
		}, nil
	}

	nf := make([]string, 0, len(names)-len(namesMap))
	for _, name := range names {
		if _, found := namesMap[name]; !found {
			nf = append(nf, name)
		}
	}

	return &GetResult{
		Found:    namesMap,
		NotFound: nf,
	}, nil
}
