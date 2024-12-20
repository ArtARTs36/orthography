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
			Found:    map[string]Gender{},
			NotFound: []string{},
		}, nil
	}

	preparedNames := make([]string, 0, len(names))
	origNamesMap := make(map[string][]string)

	for _, name := range names {
		preparedName := strings.ToLower(name)

		if _, ok := origNamesMap[preparedName]; !ok {
			origNamesMap[preparedName] = []string{}
			preparedNames = append(preparedNames, preparedName)
		}

		origNamesMap[preparedName] = append(origNamesMap[preparedName], name)
	}

	placeholders := make([]string, 0, len(preparedNames))
	args := make([]interface{}, 0, len(preparedNames))

	for i, name := range preparedNames {
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

	gendersMap := map[string]Gender{}

	type row struct {
		Name   string `db:"name"`
		Gender Gender `db:"gender"`
	}

	for rows.Next() {
		var r row

		err = rows.Scan(&r.Name, &r.Gender)
		if err != nil {
			return nil, fmt.Errorf("failed to scan struct: %w", err)
		}

		if origNames, ok := origNamesMap[r.Name]; ok {
			for _, origName := range origNames {
				gendersMap[origName] = r.Gender
			}
		}
	}

	if len(gendersMap) == len(names) {
		return &GetResult{
			Found:    gendersMap,
			NotFound: []string{},
		}, nil
	}

	nf := make([]string, 0, len(names)-len(gendersMap))
	for _, name := range names {
		if _, found := gendersMap[name]; !found {
			nf = append(nf, name)
		}
	}

	return &GetResult{
		Found:    gendersMap,
		NotFound: nf,
	}, nil
}
