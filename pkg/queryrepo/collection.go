package queryrepo

import (
	"fmt"
)

func newCollection(name string) collection {
	return collection{
		name:    name,
		queries: make(map[string]string),
	}
}

type collection struct {
	name    string
	queries map[string]string
}

func (c *collection) add(name, query string) error {
	if _, ok := c.queries[name]; ok {
		return fmt.Errorf("query %s already exists", name)
	}
	c.queries[name] = query
	return nil
}

func (c *collection) get(name string) (string, error) {
	if _, ok := c.queries[name]; !ok {
		return "", fmt.Errorf("query %s not found in collection %s", name, c.name)
	}
	return c.queries[name], nil
}
