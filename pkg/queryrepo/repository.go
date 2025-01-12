// Package queryrepo enables the use of centralized storage for all SQL queries used in an application.
package queryrepo

import (
	"database/sql"
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"
	"sync"
)

// NewFromFs creates a new repository using a filesystem.
// It takes a filesystem and a root path to start loading files from and returns an error if files cannot be loaded.
func NewFromFs(fsys fs.FS, rootPath string) (*Repository, error) {
	repo := &Repository{
		queries: make(map[string]collection),
	}
	
	return repo, loadFromFs(repo, fsys, rootPath)
}

// A Repository stores multiple collections of queries in a map for later use.
// Queries can either be retrieved by their name, or be used to create a prepared statement.
type Repository struct {
	queries map[string]collection
	mux     sync.Mutex
}

// add adds the supplied collection to the repository.
// It returns an error if the collection already exists.
func (r *Repository) add(c collection) error {
	r.mux.Lock()
	defer r.mux.Unlock()
	
	if _, ok := r.queries[c.name]; ok {
		return fmt.Errorf("collection %s already exists", c.name)
	}
	r.queries[c.name] = c
	return nil
}

// DbPrepare creates a prepared statement for the supplied database handle.
// It takes a collection name and query name to look up the query to create the prepared statement.
func (r *Repository) DbPrepare(db *sql.DB, collectionName, queryName string) (*sql.Stmt, error) {
	var err error
	var query string
	
	if query, err = r.Get(collectionName, queryName); err != nil {
		return nil, err
	}
	
	return db.Prepare(query)
}

// Get retrieves the supplied query from the repository.
// It takes a collection name and a query name to perform the lookup and returns an empty string and an error if the query cannot be found
// in the collection.
func (r *Repository) Get(collectionName, queryName string) (string, error) {
	r.mux.Lock()
	defer r.mux.Unlock()
	
	if s, ok := r.queries[collectionName]; ok {
		return s.get(queryName)
	}
	return "", fmt.Errorf("collection %s not found", collectionName)
}

// TxPrepare creates a prepared statement for the supplied in-progress database transaction.
// It takes a collection name and query name to look up the query to create the prepared statement.
func (r *Repository) TxPrepare(tx *sql.Tx, collectionName, queryName string) (*sql.Stmt, error) {
	var err error
	var statement string
	
	if statement, err = r.Get(collectionName, queryName); err != nil {
		return nil, err
	}
	
	return tx.Prepare(statement)
}

// loadFromFs looks for directories in the root path to create collections for.
// If a directory is found, it loads all the files in the subdirectory and adds the returned collection to the repository.
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

// loadFilesFromDir loads all the files in the directory and returns a collection of queries.
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
		
		var contents string
		if contents, err = LoadQueryFromFs(f, rootPath, dirName, strings.TrimSuffix(file.Name(), filepath.Ext(file.Name()))); err != nil {
			return c, err
		}
		
		if err = c.add(strings.TrimSuffix(file.Name(), filepath.Ext(file.Name())), contents); err != nil {
			return c, err
		}
	}
	return c, nil
}
