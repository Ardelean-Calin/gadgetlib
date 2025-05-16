package gadget_test

import (
	"path/filepath"
	"testing"

	"usbgadgets/functions"
	"usbgadgets/gadget"
)

func TestCreateCompositeGadget(t *testing.T) {
	gadgetBase := filepath.Join(t.TempDir(), "test_gadget")
	cleanup := gadget.SetGadgetBaseDir(gadgetBase)
	defer cleanup()

	opts := gadget.GadgetOptions{
		Name:         "foo",
		Manufacturer: "calin",
		Serial:       "foobar123",
		Controller:   "dummy_udc.0",
		Configs: []gadget.Config{
			{
				Number: 1,
				Functions: []gadget.Function{
					functions.FunctionACM{
						InstanceName: "usb0",
					},
					functions.FunctionECM{
						InstanceName: "usb0",
						DevAddr:      "06:00:0d:ea:f7:12",
						HostAddr:     "02:00:0d:ea:f7:12",
					},
				},
			},
		},
	}

	g, err := gadget.New(opts)
	if err != nil {
		t.Error(err)
	}
	defer g.Teardown()

	// TODO 1: Check that all files and symlinks were properly created
	assertPathExists(filepath.Join(gadgetBase, "configs", "c.1"), t)
	assertPathExists(filepath.Join(gadgetBase, "functions", "acm.usb0"), t)
	assertPathExists(filepath.Join(gadgetBase, "functions", "ecm.usb0"), t)
	// TODO 2: Check that Enable/Disable writes the proper UDC controller
}
