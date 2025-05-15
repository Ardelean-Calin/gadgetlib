package gadget

import (
	"path/filepath"
	"testing"
)

func TestSerialFunctionPathAndName(t *testing.T) {
	tmp := t.TempDir()
	defer SetGadgetBase(tmp)()

	g, err := CreateGadget("testg")
	if err != nil {
		t.Fatalf("CreateGadget failed: %v", err)
	}
	f := CreateSerialFunction(g, "0")
	wantName := "acm.0"
	if f.Name() != wantName {
		t.Errorf("Name = %q; want %q", f.Name(), wantName)
	}
	wantPath := filepath.Join(tmp, "testg", FunctionsDir, wantName)
	if f.Path() != wantPath {
		t.Errorf("Path = %q; want %q", f.Path(), wantPath)
	}
}

func TestECMFunctionPathAndName(t *testing.T) {
	tmp := t.TempDir()
	defer SetGadgetBase(tmp)()

	g, err := CreateGadget("teste")
	if err != nil {
		t.Fatalf("CreateGadget failed: %v", err)
	}
	f := CreateECMFunction(g, "1")
	wantName := "ecm.1"
	if f.Name() != wantName {
		t.Errorf("Name = %q; want %q", f.Name(), wantName)
	}
	wantPath := filepath.Join(tmp, "teste", FunctionsDir, wantName)
	if f.Path() != wantPath {
		t.Errorf("Path = %q; want %q", f.Path(), wantPath)
	}
}
