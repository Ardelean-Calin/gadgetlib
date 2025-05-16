package gadget_test

import (
	"os"
	"testing"
)

func assertPathExists(path string, t *testing.T) {
	_, err := os.Stat(path)
	if err != nil {
		t.Error(err)
	}
}
