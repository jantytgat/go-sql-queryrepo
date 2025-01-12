package queryrepo

import (
	"database/sql"
	"fmt"
	"io/fs"
	"path/filepath"
	"sync"
)

func NewFromFs(fsys fs.FS, rootPath string) (*Repository, error) {
	repo := &Repository{
		queries: make(map[string]collection),
	}

	return repo, loadFromFs(repo, fsys, rootPath)
}

type Repository struct {
	queries map[string]collection
	mux     sync.Mutex
}

func (r *Repository) add(c collection) error {
	r.mux.Lock()
	defer r.mux.Unlock()

	if _, ok := r.queries[c.name]; ok {
		return fmt.Errorf("statement %s already exists", c.name)
	}
	r.queries[c.name] = c
	return nil
}

func (r *Repository) DbPrepare(db *sql.DB, collectionName, statementName string) (*sql.Stmt, error) {
	var err error
	var statement string

	if statement, err = r.Get(collectionName, statementName); err != nil {
		return nil, err
	}

	return db.Prepare(statement)
}

func (r *Repository) Get(collection, query string) (string, error) {
	r.mux.Lock()
	defer r.mux.Unlock()

	if s, ok := r.queries[collection]; ok {
		return s.get(query)
	}
	return "", fmt.Errorf("collection %s not found", collection)
}

func (r *Repository) TxPrepare(tx *sql.Tx, collectionName, statementName string) (*sql.Stmt, error) {
	var err error
	var statement string

	if statement, err = r.Get(collectionName, statementName); err != nil {
		return nil, err
	}

	return tx.Prepare(statement)
}

func loadFromFs(r *Repository, f fs.FS, rootPath string) error {
	var err error
	var files []fs.DirEntry
	if files, err = fs.ReadDir(f, rootPath); err != nil {
		return err
	}

	for _, file := range files {
		if file.IsDir() {
			var c collection
			if c, err = loadFilesFromDir(f, rootPath, file.Name()); err != nil {
				return err
			}

			if err = r.add(c); err != nil {
				return err
			}
		}
	}
	return nil
}

func loadFilesFromDir(f fs.FS, rootPath, dirName string) (collection, error) {
	var err error
	var c = newCollection(dirName)
	var fullPath = filepath.Join(rootPath, dirName)

	var files []fs.DirEntry
	if files, err = fs.ReadDir(f, fullPath); err != nil {
		return c, err
	}

	for _, file := range files {
		if file.IsDir() {
			return c, fmt.Errorf("nested directories are not supported, %s is a directory in %s", file.Name(), fullPath)
		}

		var contents []byte
		if contents, err = fs.ReadFile(f, filepath.Join(dirName, file.Name())); err != nil {
			return c, fmt.Errorf("failed to read file %s from directory %s: %w", file.Name(), fullPath, err)
		}

		if err = c.add(dirName, string(contents)); err != nil {
			return c, err
		}
	}
	return c, nil
}
