package gadget_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ardelean-calin/gadgetlib/functions"
	"github.com/ardelean-calin/gadgetlib/gadget"
)

func TestCreateCompositeGadget(t *testing.T) {
	gadgetBase := filepath.Join(t.TempDir(), "gadgetlib_test")
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

	// Check that all files and symlinks were properly created
	assertPathExists(gadget.GadgetBaseDir, t)
	assertPathExists(filepath.Join(gadget.GadgetBaseDir, "foo/configs/cfg.1"), t)
	assertPathExists(filepath.Join(gadget.GadgetBaseDir, "foo/configs/cfg.1/strings/0x409/configuration"), t)
	assertPathExists(filepath.Join(gadget.GadgetBaseDir, "foo/configs/cfg.1/acm.usb0"), t)
	assertPathExists(filepath.Join(gadget.GadgetBaseDir, "foo/configs/cfg.1/ecm.usb0"), t)
	assertPathExists(filepath.Join(gadget.GadgetBaseDir, "foo/functions/acm.usb0"), t)
	assertPathExists(filepath.Join(gadget.GadgetBaseDir, "foo/functions/ecm.usb0"), t)
	assertPathExists(filepath.Join(gadget.GadgetBaseDir, "foo/functions/ecm.usb0/dev_addr"), t)
	assertPathExists(filepath.Join(gadget.GadgetBaseDir, "foo/functions/ecm.usb0/host_addr"), t)

	// Check that Enable writes the proper UDC controller
	if err := g.Enable(); err != nil {
		t.Error(err)
	}
	data, err := os.ReadFile(filepath.Join(g.Path(), "UDC"))
	if err != nil {
		t.Error(err)
	}
	if string(data) != opts.Controller+"\n" {
		t.Errorf("expected UDC file to contain %q, got %q", opts.Controller, string(data))
	}
	// Check that Disable writes a newline
	if err := g.Disable(); err != nil {
		t.Error(err)
	}
	data, err = os.ReadFile(filepath.Join(g.Path(), "UDC"))
	if err != nil {
		t.Error(err)
	}
	if string(data) != "\n" {
		t.Errorf("expected UDC file to contain newline, got %q", string(data))
	}
}
