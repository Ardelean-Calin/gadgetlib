package gadget

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

var GadgetBaseDir = "/sys/kernel/config/usb_gadget"

type GadgetOptions struct {
	Name         string
	Manufacturer string
	// Optional, will generate one if not given
	Serial string
	// PID, VID // <= set by default by us
	Configs []Config
	// eg. /sys/class/udc/dummy_udc.0
	// Optionally obtain via `UDCScan()`
	Controller string
}

type Config struct {
	Number    int
	Functions []Function
	Path      string
}

func (c Config) apply(gadgetPath string) error {
	cfgDir := filepath.Join(gadgetPath, "configs", c.Name())
	err := os.MkdirAll(cfgDir, os.ModePerm)
	if err != nil {
		return fmt.Errorf("cannot create config directory %q: %w", cfgDir, err)
	}

	// TODO: populate other configuration related parameters (such as strings)

	return nil
}

func (c Config) Name() string {
	return fmt.Sprintf("cfg.%d", c.Number)
}

type Gadget interface {
	Enable() error
	Disable() error
	IsEnabled() bool
	Teardown() error
	Path() string
}

type Function interface {
	Apply(configPath, gadgetPath string) error
	Name() string
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

func (g *USBGadget) Path() string {
	return g.base
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

func New(opts GadgetOptions) (Gadget, error) {
	g := &USBGadget{
		base: filepath.Join(GadgetBaseDir, opts.Name),
		udc:  opts.Controller,
	}

	for _, c := range opts.Configs {
		cfg := config{
			name:      c.Name(),
			functions: make([]*function, 0),
		}
		if err := c.apply(g.base); err != nil {
			return nil, fmt.Errorf("cannot create gadget config %q: %w", c.Name(), err)
		}
		_ = cfg
		// TODO: Create the config directory, functions, symlinks, etc.
		// cfg.Apply()
		// fnc.Apply()

		for _, f := range c.Functions {
			fn := function{
				name: f.Name(),
			}

			if err := f.Apply(c.Path, g.base); err != nil {
				return nil, err
			}
			g.functions = append(g.functions, fn)
		}
	}
	return g, nil
}
