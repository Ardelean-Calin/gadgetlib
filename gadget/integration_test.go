package gadget

import (
	"testing"

	o "github.com/flotter/kernos-workshop/services/otg/pkg/gadget/option"
)

func TestSerialGadgetIntegration(t *testing.T) {
	g, err := CreateGadget("test_serial")
	if err != nil {
		t.Fatalf("CreateGadget failed: %v", err)
	}
	defer func() {
		if err := g.CleanUp(); err != nil {
			t.Fatalf("CleanUp failed: %v", err)
		}
	}()

	g.SetAttrs(&GadgetAttrs{
		IdVendor:  o.Some[uint16](0x1d6b),
		IdProduct: o.Some[uint16](0x0104),
	})
	if err := g.SetStrs(&GadgetStrs{
		SerialNumber: "12345",
		Manufacturer: "Test",
		Product:      "SerialTest",
	}, LangUsEng); err != nil {
		t.Fatalf("SetStrs failed: %v", err)
	}

	// Create and bind a serial function
	sf := CreateSerialFunction(g, "0")
	cfg, err := CreateConfig(g, DefaultConfigLabel, 1)
	if err != nil {
		t.Fatalf("CreateConfig failed: %v", err)
	}
	cfg.SetAttrs(&ConfigAttrs{})
	if err := cfg.SetStrs(&ConfigStrs{Configuration: "SerialConfig"}, LangUsEng); err != nil {
		t.Fatalf("Config.SetStrs failed: %v", err)
	}
	_, err = CreateBinding(cfg, sf, sf.Name())
	if err != nil {
		t.Fatalf("CreateBinding failed: %v", err)
	}

	// Enable the gadget
	udcs := GetUdcs()
	if len(udcs) == 0 {
		t.Skip("no UDC available")
	}
	g.Enable(udcs[0])
	if !g.IsEnabled() {
		t.Fatal("gadget not enabled")
	}

	// Verify the serial device path
	dev, err := sf.GetDev()
	if err != nil {
		t.Fatalf("GetDev failed: %v", err)
	}
	if dev == "" {
		t.Fatalf("GetDev returned empty device")
	}
	t.Logf("Serial device: %s", dev)
}
