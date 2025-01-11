package queryrepo

import (
	"database/sql"
	"fmt"
	"io/fs"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

func NewFromFs(fsys fs.FS) (*Repository, error) {
	repo := New()

	var err error
	var files []fs.DirEntry
	if files, err = fs.ReadDir(fsys, "."); err != nil {
		return nil, err
	}
	return repo, loadStatements(repo, ".", fsys, files)
}

func loadStatements(r *Repository, rootPath string, filesystem fs.FS, files []fs.DirEntry) error {
	var err error

	for _, file := range files {
		if file.IsDir() {
			var subFiles []fs.DirEntry
			if subFiles, err = fs.ReadDir(filesystem, filepath.Join(rootPath, file.Name())); err != nil {
				return err
			}
			if err = loadStatements(r, filepath.Join(rootPath, file.Name()), filesystem, subFiles); err != nil {
				return err
			}
		} else {
			var contents []byte
			if contents, err = fs.ReadFile(filesystem, filepath.Join(rootPath, file.Name())); err != nil {
				fmt.Println("Error reading file", file.Name())
				return err
			}

			var fileStatements = Statements{}
			if err = yaml.Unmarshal(contents, &fileStatements); err != nil {
				fmt.Println("Error parsing file", file.Name())
				return err
			}

			if err = r.Add(fileStatements); err != nil {
				return err
			}
		}
	}
	return nil
}

func New() *Repository {
	return &Repository{
		statements: make(map[string]Statements),
	}
}

type Repository struct {
	statements map[string]Statements
}

func (c *Repository) Add(statements Statements) error {
	if _, ok := c.statements[statements.Name]; ok {
		return fmt.Errorf("statement %s already exists", statements.Name)
	}
	c.statements[statements.Name] = statements
	return nil
}

func (c *Repository) Get(collectionName, statementName string) (string, error) {
	if s, ok := c.statements[collectionName]; ok {
		return s.Get(statementName)
	}
	return "", fmt.Errorf("collection %s not found", collectionName)
}

func (c *Repository) DbPrepare(db *sql.DB, collectionName, statementName string) (*sql.Stmt, error) {
	var err error
	var statement string

	if statement, err = c.Get(collectionName, statementName); err != nil {
		return nil, err
	}

	return db.Prepare(statement)
}

func (c *Repository) TxPrepare(tx *sql.Tx, collectionName, statementName string) (*sql.Stmt, error) {
	var err error
	var statement string

	if statement, err = c.Get(collectionName, statementName); err != nil {
		return nil, err
	}

	return tx.Prepare(statement)
}
