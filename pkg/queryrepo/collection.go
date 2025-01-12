package queryrepo

import (
	"fmt"
)

// newCollection creates a new collection with the supplied name and returns it to the caller.
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

// add adds a query to the collection.
func (c *collection) add(name, query string) error {
	if _, ok := c.queries[name]; ok {
		return fmt.Errorf("query %s already exists", name)
	}
	c.queries[name] = query
	return nil
}

// get retrieves a query from the collection by name.
// If the query name cannot be found, get() returns an empty string and an error.
func (c *collection) get(name string) (string, error) {
	if _, ok := c.queries[name]; !ok {
		return "", fmt.Errorf("query %s not found in collection %s", name, c.name)
	}
	return c.queries[name], nil
}
