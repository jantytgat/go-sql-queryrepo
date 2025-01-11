package queryrepo

import (
	"fmt"
)

type Statements struct {
	Name       string      `yaml:"name"`
	Statements []Statement `yaml:"statements"`
}

func (s Statements) Get(name string) (string, error) {
	for _, stmt := range s.Statements {
		if stmt.Name == name {
			return stmt.Statement, nil
		}
	}
	return "", fmt.Errorf("statement %s not found", name)
}
