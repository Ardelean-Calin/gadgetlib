//go:build integration
// +build integration

package gadget

// SetGadgetBaseDir is a no-op under the "integration" tag.
func SetGadgetBaseDir(newDir string) func() {
    return func() {}
}
