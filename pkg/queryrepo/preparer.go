package queryrepo

import (
	"database/sql"
	"io/fs"
)

// Preparer defines the interface to create a prepared statement.
type Preparer interface {
	Prepare(query string) (*sql.Stmt, error)
}

// Prepare creates a prepared statement for the supplied Preparer by looking up a query in the supplied repository.
// It returns an nil pointer and an error if either the query cannot be found in the supplied repository, or the statement preparation fails.
func Prepare[T Preparer](t T, r *Repository, collectionName, queryName string) (*sql.Stmt, error) {
	var err error
	var query string
	
	if query, err = r.Get(collectionName, queryName); err != nil {
		return nil, err
	}
	return t.Prepare(query)
}

// PrepareFromFs creates a prepared statement for the supplied Preparer by looking up a query in the supplied filesystem.
// It returns an nil pointer and an error if either the query cannot be found in the supplied filesystem, or the statement preparation fails.
func PrepareFromFs[T Preparer](t T, f fs.FS, rootPath, collectionName, queryName string) (*sql.Stmt, error) {
	var err error
	var query string
	
	if query, err = LoadQueryFromFs(f, rootPath, collectionName, queryName); err != nil {
		return nil, err
	}
	return t.Prepare(query)
}
