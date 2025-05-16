package gadget_test

import (
	"path/filepath"
	"testing"

	"gadgetlib/functions"
	"gadgetlib/gadget"
)

func TestCreateSerialGadget(t *testing.T) {
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
	assertPathExists(gadget.GadgetBaseDir, t)
	assertPathExists(filepath.Join(gadget.GadgetBaseDir, "foo/configs/cfg.1"), t)
	assertPathExists(filepath.Join(gadget.GadgetBaseDir, "foo/configs/cfg.1/strings/0x409/configuration"), t)
	assertPathExists(filepath.Join(gadget.GadgetBaseDir, "foo/configs/cfg.1/acm.usb0"), t)
	assertPathExists(filepath.Join(gadget.GadgetBaseDir, "foo/configs/cfg.1/ecm.usb0"), t)
	assertPathExists(filepath.Join(gadget.GadgetBaseDir, "foo/functions/acm.usb0"), t)
	assertPathExists(filepath.Join(gadget.GadgetBaseDir, "foo/functions/ecm.usb0"), t)
	assertPathExists(filepath.Join(gadget.GadgetBaseDir, "foo/functions/ecm.usb0/dev_addr"), t)
	assertPathExists(filepath.Join(gadget.GadgetBaseDir, "foo/functions/ecm.usb0/host_addr"), t)
	// TODO 2: Check that Enable/Disable writes the proper UDC controller
}
