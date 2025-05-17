package gadget

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var gadgetBaseDir = "/sys/kernel/config/usb_gadget"

type GadgetOptions struct {
	Name         string
	Manufacturer string
	// Optional, will generate one if not given
	Serial  string
	Configs []Config
	// eg. /sys/class/udc/dummy_udc.0
	// Optionally obtain via `UDCScan()`
	Controller string
}

type Config struct {
	Number    int
	Functions []Function
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
	Apply(base string) error
	Name() string
}

type config struct {
	name      string
	functions []Function
	path      string
}

func (c *config) apply() error {
	err := os.MkdirAll(c.path, os.ModePerm)
	if err != nil {
		return fmt.Errorf("cannot create config directory %q: %w", c.path, err)
	}

	p := filepath.Join(c.path, "strings/0x409/configuration")
	err = os.MkdirAll(filepath.Dir(p), os.ModePerm)
	if err != nil {
		return fmt.Errorf("cannot create strings dir for config %q: %w", c.name, err)
	}
	names := []string{}
	for _, f := range c.functions {
		names = append(names, f.Name())
	}

	err = os.WriteFile(
		p,
		fmt.Appendf([]byte{}, "Configuration [ %s ]", strings.Join(names, ", ")),
		os.ModePerm)
	if err != nil {
		return fmt.Errorf("cannot write strings for configuration %q: %w", c.name, err)
	}

	// Set MaxPower
	err = os.WriteFile(filepath.Join(c.path, "MaxPower"), []byte("250"), os.ModePerm) // 500mA
	if err != nil {
		return fmt.Errorf("cannot set configuration MaxPower to 500mA: %q - %w", c.name, err)
	}

	// Apply/create each individual function
	var errs []error
	for _, f := range c.functions {
		errs = append(errs, c.applyFunction(f))
	}

	return errors.Join(errs...)
}

// applyFunction produces the needed folder layout and symlinks for a function
// to work. For example:
// A symlink is needed to the appropriate config directory:
//
//	./configs/c.1/ncm.usb0 -> ../../../../usb_gadget/g1/functions/ncm.usb0
//
// Function layout:
//
//	./functions
//	./functions/ncm.usb0
//
// Each function can have more options/parameters (that are applied in each
// individual `Function` `Apply()` method):
//
//	./functions/ncm.usb0/ifname
//	./functions/ncm.usb0/qmult
//	./functions/ncm.usb0/host_addr
//	./functions/ncm.usb0/dev_addr
func (c *config) applyFunction(f Function) error {
	gadgetbase := filepath.Join(c.path, "../..")

	funcPath := filepath.Join(gadgetbase, "functions", f.Name())
	if err := os.MkdirAll(funcPath, os.ModePerm); err != nil {
		return fmt.Errorf("cannot create function %q: %w", f.Name(), err)
	}

	// Every function can implement its own `Apply` to do changes to
	// its base dir.
	if err := f.Apply(funcPath); err != nil {
		return fmt.Errorf("cannot apply function %q: %w", f.Name(), err)
	}

	src := filepath.Join(funcPath)
	dst := filepath.Join(c.path, f.Name())

	if err := os.Symlink(src, dst); err != nil {
		return fmt.Errorf("cannot symlink %q -> %q: %w", src, dst, err)
	}

	return nil
}

type function struct {
	name string
}

type USBGadget struct {
	// Basic parameters
	name         string
	manufacturer string
	serial       string
	// base path of this gadget, easy to tear down
	base      string
	udc       string
	functions []Function
	// Configuration names
	configs []config
}

func (g *USBGadget) Path() string {
	return g.base
}

func (g *USBGadget) Enable() error {
	return os.WriteFile(filepath.Join(g.base, "UDC"), []byte(g.udc+"\n"), os.ModePerm)
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

 // Teardown removes the USB gadget by disabling it if enabled and cleaning up all
 // configurations, function directories, string directories, and the gadget base directory.
func (g *USBGadget) Teardown() error {
	var errs []error

	if g.IsEnabled() {
		g.Disable()
	}

	for _, c := range g.configs {
		cfgBase := filepath.Join(g.base, "configs", c.name)
		// 1. Remove functions from configurations (aka the symlinks)
		for _, f := range c.functions {
			linkPath := filepath.Join(cfgBase, f.Name())
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
		funcDir := filepath.Join(g.base, "functions", f.Name())
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

 // apply initializes the USB gadget by creating its base directory and writing
 // device descriptors and string files based on the gadget's settings.
func (g *USBGadget) apply() error {
	if err := os.MkdirAll(g.base, os.ModePerm); err != nil {
		return err
	}
	err := g.ensureFile("idVendor", "0x1d6b")
	if err != nil {
		return err
	}
	err = g.ensureFile("idProduct", "0x0104")
	if err != nil {
		return err
	}
	err = g.ensureFile("bcdUSB", "0x0300")
	if err != nil {
		return err
	}
	err = g.ensureFile("bcdDevice", "0x0100")
	if err != nil {
		return err
	}
	err = g.ensureFile("strings/0x409/manufacturer", g.manufacturer)
	if err != nil {
		return err
	}
	err = g.ensureFile("strings/0x409/serialnumber", g.serial)
	if err != nil {
		return err
	}
	err = g.ensureFile("strings/0x409/product", g.name)
	if err != nil {
		return err
	}

	return nil
}

func (g *USBGadget) ensureFile(relPath string, content string) error {
	fullPath := filepath.Join(g.base, relPath)
	// Step 1: Check if parent dir exists
	if _, err := os.Stat(filepath.Dir(fullPath)); os.IsNotExist(err) {
		err = os.MkdirAll(filepath.Dir(fullPath), os.ModePerm)
		if err != nil {
			return fmt.Errorf("cannot create parent dir for %q: %w", relPath, err)
		}
	}

	// Step 2: we've created the dir, write the content
	return os.WriteFile(fullPath, []byte(content), os.ModePerm)
}

func New(opts GadgetOptions) (Gadget, error) {
	g := &USBGadget{
		base: filepath.Join(gadgetBaseDir, opts.Name),
		udc:  opts.Controller,

		name:         opts.Name,
		serial:       opts.Serial,
		manufacturer: opts.Manufacturer,
	}

	// Create gadget dir
	if err := g.apply(); err != nil {
		return nil, fmt.Errorf("cannot initialize gadget dir %q: %w", g.base, err)
	}

	for _, c := range opts.Configs {
		cfg := config{
			name:      c.Name(),
			functions: []Function{},
			path:      filepath.Join(g.base, "configs", c.Name()),
		}

		for _, fn := range c.Functions {
			g.functions = append(g.functions, fn)
			cfg.functions = append(cfg.functions, fn)
		}

		if err := cfg.apply(); err != nil {
			return nil, err
		}

		g.configs = append(g.configs, cfg)
	}
	return g, nil
}
