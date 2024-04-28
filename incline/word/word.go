package word

import "database/sql"

type Word struct {
	Nominative    string         `db:"nominative"`
	Genitive      sql.NullString `db:"genitive"`
	Dative        sql.NullString `db:"dative"`
	Accusative    sql.NullString `db:"accusative"`
	Instrumental  sql.NullString `db:"instrumental"`
	Prepositional sql.NullString `db:"prepositional"`
}
