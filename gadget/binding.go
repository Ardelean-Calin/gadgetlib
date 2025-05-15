package gadget

import (
	"fmt"
	"os"
	"path/filepath"
)

type Binding struct {
	name     string
	basePath string

	config   *Config
	function Function
}

// Bind a function to a given config. Multiple functions can be bound to the
// same config.
func CreateBinding(c *Config, f Function, name string) (*Binding, error) {
	linkPath := filepath.Join(c.Path(), name)

	binding := &Binding{
		name:     name,
		basePath: c.Path(),
		config:   c,
		function: f,
	}

	err := os.Symlink(f.Path(), linkPath)
	if err != nil {
		return nil, fmt.Errorf("cannot create binding: %w", err)
	}

	c.bindings = append(c.bindings, binding)

	return binding, nil
}

func (b *Binding) Path() string {
	return filepath.Join(b.basePath, b.name)
}
