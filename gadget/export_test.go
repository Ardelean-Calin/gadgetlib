//go:build !integration
// +build !integration

package gadget

func SetGadgetBaseDir(newDir string) func() {
	oldDir := GadgetBaseDir
	GadgetBaseDir = newDir
	return func() {
		GadgetBaseDir = oldDir
	}
}
