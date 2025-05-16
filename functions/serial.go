package functions

import (
	"fmt"
	"os"
	"path/filepath"
)

type FunctionACM struct {
	InstanceName string
}

func (f FunctionACM) Apply(gadgetPath, configPath string) error {
	functionPath := filepath.Join(gadgetPath, "functions", f.Name())

	err := os.MkdirAll(functionPath, os.ModePerm)
	if err != nil {
		return fmt.Errorf("cannot create function %q: %w", f.Name(), err)
	}

	// TODO: AAAAH!!! I need to link to a config not to a gadget!
	dst := filepath.Join(configPath, f.Name())
	src := functionPath
	if err := os.Link(src, dst); err != nil {
		return fmt.Errorf("cannot symlink %q <- %q: %w", dst, src, err)
	}

	return nil
}

func (f FunctionACM) Name() string {
	return fmt.Sprintf("acm.%s", f.InstanceName)
}
