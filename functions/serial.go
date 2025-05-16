package functions

import (
	"fmt"
)

type FunctionACM struct {
	InstanceName string
}

// Serial gadgets don't need additional configuration
func (f FunctionACM) Apply(base string) error {
	return nil
}

func (f FunctionACM) Name() string {
	return fmt.Sprintf("acm.%s", f.InstanceName)
}
