package queryrepo

import (
	"database/sql"
	"io/fs"
)

type Preparer interface {
	Prepare(query string) (*sql.Stmt, error)
}

func Prepare[T Preparer](t T, r *Repository, collectionName, queryName string) (*sql.Stmt, error) {
	var err error
	var query string

	if query, err = r.Get(collectionName, queryName); err != nil {
		return nil, err
	}
	return t.Prepare(query)
}

func PrepareFromFs[T Preparer](t T, f fs.FS, rootPath, collectionName, queryName string) (*sql.Stmt, error) {
	var err error
	var query string

	if query, err = LoadQueryFromFs(f, rootPath, collectionName, queryName); err != nil {
		return nil, err
	}
	return t.Prepare(query)
}
