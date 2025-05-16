//go:build !integration
// +build !integration

package gadget

func SetGadgetBaseDir(newDir string) func() {
	oldDir := gadgetBaseDir
	gadgetBaseDir = newDir
	return func() {
		gadgetBaseDir = oldDir
	}
}
