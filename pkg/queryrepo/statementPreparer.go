package queryrepo

import "database/sql"

type StatementPreparer interface {
	Prepare(query string) (*sql.Stmt, error)
}

func Prepare[T StatementPreparer](t T, r *Repository, collectionName, queryName string) (*sql.Stmt, error) {
	var err error
	var statement string

	if statement, err = r.Get(collectionName, queryName); err != nil {
		return nil, err
	}
	return t.Prepare(statement)
}
