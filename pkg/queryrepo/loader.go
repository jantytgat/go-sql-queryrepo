package queryrepo

import (
	"fmt"
	"io/fs"
	"path/filepath"
)

func LoadQueryFromFs(f fs.FS, rootPath, collectionName, queryName string) (string, error) {
	var err error
	var contents []byte
	if contents, err = fs.ReadFile(f, filepath.Join(rootPath, collectionName, queryName, ".sql")); err != nil {
		return "", fmt.Errorf("failed to read file %s: %w", filepath.Join(rootPath, collectionName, queryName, ".sql"), err)
	}
	return string(contents), nil
}
