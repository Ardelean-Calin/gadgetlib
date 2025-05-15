package gadget

func SetGadgetBase(base string) func() {
	oldBase := BasePath
	BasePath = base
	return func() {
		BasePath = oldBase
	}
}
