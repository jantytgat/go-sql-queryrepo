package queryrepo

import (
	"context"
	"database/sql"
	"errors"
	"io/fs"
)

// Preparer defines the interface to create a prepared statement.
type Preparer interface {
	Prepare(query string) (*sql.Stmt, error)
	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
}

// Prepare creates a prepared statement for the supplied Preparer by looking up a query in the supplied repository.
// It returns an nil pointer and an error if either the query cannot be found in the supplied repository, or the statement preparation fails.
func Prepare[T Preparer](t T, r *Repository, collectionName, queryName string) (*sql.Stmt, error) {
	if r == nil {
		return nil, errors.New("repository is nil")
	}
	
	var err error
	var query string
	
	if query, err = r.Get(collectionName, queryName); err != nil {
		return nil, err
	}
	return t.Prepare(query)
}

// Prepare creates a prepared statement for the supplied Preparer by looking up a query in the supplied repository using a context.
// It returns an nil pointer and an error if either the query cannot be found in the supplied repository, or the statement preparation fails.
func PrepareContext[T Preparer](ctx context.Context, t T, r *Repository, collectionName, queryName string) (*sql.Stmt, error) {
	if r == nil {
		return nil, errors.New("repository is nil")
	}
	
	var err error
	var query string
	
	if query, err = r.Get(collectionName, queryName); err != nil {
		return nil, err
	}
	return t.PrepareContext(ctx, query)
}

// PrepareFromFs creates a prepared statement for the supplied Preparer by looking up a query in the supplied filesystem.
// It returns an nil pointer and an error if either the query cannot be found in the supplied filesystem, or the statement preparation fails.
func PrepareFromFs[T Preparer](t T, f fs.FS, rootPath, collectionName, queryName string) (*sql.Stmt, error) {
	if f == nil {
		return nil, errors.New("invalid filesystem")
	}
	var err error
	var query string
	
	if query, err = LoadQueryFromFs(f, rootPath, collectionName, queryName); err != nil {
		return nil, err
	}
	return t.Prepare(query)
}

// PrepareFromFs creates a prepared statement for the supplied Preparer by looking up a query in the supplied filesystem using a context.
// It returns an nil pointer and an error if either the query cannot be found in the supplied filesystem, or the statement preparation fails.
func PrepareFromFsContext[T Preparer](ctx context.Context, t T, f fs.FS, rootPath, collectionName, queryName string) (*sql.Stmt, error) {
	if f == nil {
		return nil, errors.New("invalid filesystem")
	}
	var err error
	var query string
	
	if query, err = LoadQueryFromFs(f, rootPath, collectionName, queryName); err != nil {
		return nil, err
	}
	return t.PrepareContext(ctx, query)
}
