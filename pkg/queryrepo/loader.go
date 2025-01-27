package queryrepo

import (
	"embed"
	"fmt"
	"io/fs"
	"path"
	"path/filepath"
)

// LoadQueryFromFs retrieves a query from a filesystem.
// It needs the root path to start the search from, as well as a collection name and a query name.
// The collection name equals to a direct directory name in the root path.
// The query name is the file name (without extension) to load the contents from.
// It returns and empty string and an error if the file cannot be found.
func LoadQueryFromFs(f fs.FS, rootPath, collectionName, queryName string) (string, error) {
	var err error
	var contents []byte
	switch f.(type) {
	case embed.FS:
		if contents, err = fs.ReadFile(f, path.Join(rootPath, collectionName, queryName)+".sql"); err != nil {
			return "", fmt.Errorf("failed to read file %s: %w", path.Join(rootPath, collectionName, queryName)+".sql", err)
		}
	default:
		if contents, err = fs.ReadFile(f, filepath.Join(rootPath, collectionName, queryName)+".sql"); err != nil {
			return "", fmt.Errorf("failed to read file %s: %w", filepath.Join(rootPath, collectionName, queryName)+".sql", err)
		}
	}
	return string(contents), nil
}
