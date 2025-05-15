package gadget

import (
	"os"
	"strings"
	"testing"

	o "usbgadgets/gadget/option"
)

func inDocker() bool {
	if _, err := os.Stat("/.dockerenv"); err == nil {
		return true
	}
	if data, err := os.ReadFile("/proc/1/cgroup"); err == nil {
		return strings.Contains(string(data), "docker")
	}
	return false
}

func TestSerialGadgetIntegration(t *testing.T) {
	if !inDocker() {
		t.Skip("skipping integration test: not inside Docker")
	}
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
		BcdUSB:          o.Some[uint16](0x0300), // USB3.0
		BDeviceClass:    o.None[uint8](),
		BDeviceSubClass: o.None[uint8](),
		BDeviceProtocol: o.None[uint8](),
		BMaxPacketSize0: o.None[uint8](),
		IdVendor:        o.Some[uint16](0x1d6b), // Linux Foundation
		IdProduct:       o.Some[uint16](0x0104), // Multifunction Composite Gadget
		BcdDevice:       o.Some[uint16](0x0100), // v1.0.0
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
	cfg.SetAttrs(&ConfigAttrs{
		BmAttributes: o.None[uint8](),
		BMaxPower:    o.None[uint8](),
	})
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
