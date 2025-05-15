package gadgetv2

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

const (
	ACMFunction string = "acm"
	ECMFunction        = "ecm"
)

var GadgetBaseDir = "/sys/kernel/config/usb_gadget"

type GadgetConfig struct {
	Name         string
	Manufacturer string
	// Optional, will generate one if not given
	Serial string
	// PID, VID // <= set by default by us
	Configs []struct {
		Number    int // eg. 1
		Functions []Function
	}
	// eg. /sys/class/udc/dummy_udc.0
	// Optionally obtain via `UDCScan()`
	Controller string
}

type Gadget interface {
	Enable() error
	Disable() error
	IsEnabled() bool
	Teardown() error
}

type Function interface {
	Apply(gadgetPath string) error
	Name() string
}

type FunctionACM struct {
	InstanceName string
}

func (f FunctionACM) Apply(gadgetPath string) error {
	functionPath := filepath.Join(gadgetPath, "functions", f.Name())

	err := os.MkdirAll(functionPath, os.ModePerm)
	if err != nil {
		return fmt.Errorf("cannot create function %q: %w", f.Name(), err)
	}

	configPath := filepath.Join(gadgetPath, "configs")
	dst := filepath.Join(configPath, f.Name())
	src := functionPath
	if err := os.Link(src, dst); err != nil {
		return fmt.Errorf("cannot symlink %q <- %q: %w", dst, src, err)
	}

	return nil
}

func (f FunctionACM) Name() string {
	return fmt.Sprintf("acm.%s", f.InstanceName)
}

type FunctionECM struct {
	InstanceName string
	DevAddr      string
	HostAddr     string
}

type config struct {
	name      string
	functions []*function
}

type function struct {
	name string
}

type USBGadget struct {
	// base path of this gadget, easy to tear down
	base      string
	udc       string
	functions []function
	// Configuration names
	configs []struct {
		name      string
		functions []*function
	}
}

func (g *USBGadget) Enable() error {
	return os.WriteFile(filepath.Join(g.base, "UDC"), []byte(g.udc), os.ModePerm)
}

func (g *USBGadget) Disable() error {
	return os.WriteFile(filepath.Join(g.base, "UDC"), []byte("\n"), os.ModePerm)
}

func (g *USBGadget) IsEnabled() bool {
	data, err := os.ReadFile(filepath.Join(g.base, "UDC"))
	if err != nil {
		return false
	}
	return string(data) == g.udc
}

// Add a descriptive method description AI!
func (g *USBGadget) Teardown() error {
	var errs []error

	if g.IsEnabled() {
		g.Disable()
	}

	for _, c := range g.configs {
		cfgBase := filepath.Join(g.base, "configs", c.name)
		// 1. Remove functions from configurations (aka the symlinks)
		for _, f := range c.functions {
			linkPath := filepath.Join(cfgBase, f.name)
			err := os.Remove(linkPath)
			if err != nil {
				errs = append(errs, err)
			}
		}
		// 2. Remove strings directories in configurations
		err := os.RemoveAll(filepath.Join(cfgBase, "strings", "0x409"))
		if err != nil {
			errs = append(errs, err)
		}
		// 3. Remove the configurations
		err = os.RemoveAll(cfgBase)
		if err != nil {
			errs = append(errs, err)
		}
	}
	// 4. Remove the functions
	for _, f := range g.functions {
		funcDir := filepath.Join(g.base, "functions", f.name)
		err := os.RemoveAll(funcDir)
		if err != nil {
			errs = append(errs, err)
		}
	}
	// 5. Remove strings directories in the gadget
	err := os.RemoveAll(filepath.Join(g.base, "strings", "0x409"))
	if err != nil {
		errs = append(errs, err)
	}
	// 6. Remove the gadget
	err = os.RemoveAll(g.base)
	if err != nil {
		errs = append(errs, err)
	}

	return errors.Join(errs...)
}

func New(cfg GadgetConfig) (Gadget, error) {
	g := &USBGadget{
		base: filepath.Join(GadgetBaseDir, cfg.Name),
		udc:  cfg.Controller,
	}

	for _, c := range cfg.Configs {
		cfg := config{
			name:      fmt.Sprintf("cfg.%d", c.Number),
			functions: make([]*function, 0),
		}
		_ = cfg
		// TODO: Create the config directory, functions, symlinks, etc.
		// cfg.Apply()
		// fnc.Apply()

		for _, f := range c.Functions {
			fn := function{
				name: f.Name(),
			}

			if err := f.Apply(g.base); err != nil {
				return nil, err
			}
			g.functions = append(g.functions, fn)
		}
	}
	return g, nil
}
