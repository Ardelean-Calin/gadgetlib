package functions

import (
	"fmt"
	"os"
	"path/filepath"
)

type FunctionECM struct {
	InstanceName string
	DevAddr      string
	HostAddr     string
}

func (f FunctionECM) Apply(base string) error {
	err := os.WriteFile(filepath.Join(base, "host_addr"), []byte(f.HostAddr), os.ModePerm)
	if err != nil {
		return err
	}
	err = os.WriteFile(filepath.Join(base, "dev_addr"), []byte(f.DevAddr), os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

func (f FunctionECM) Name() string {
	return fmt.Sprintf("ecm.%s", f.InstanceName)
}
