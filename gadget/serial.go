package gadget

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

const SerialFunctionTypeName = "acm"

type SerialFunction struct {
	name     string
	path     string
	instance string

	g *Gadget
}

func (e *SerialFunction) Path() string {
	return filepath.Join(e.path, e.name)
}

func (e *SerialFunction) Name() string {
	return e.name
}

// func (e *SerialFunction) GetDev() string {
// 	return e.instance
// }

func CreateSerialFunction(g *Gadget, instance string) *SerialFunction {
	basePath := filepath.Join(g.Path(), FunctionsDir)
	name := fmt.Sprintf("%s.%s", SerialFunctionTypeName, instance)
	path := filepath.Join(basePath, name)

	function := &SerialFunction{
		name:     name,
		path:     basePath,
		instance: instance,

		g: g,
	}

	err := os.Mkdir(path, os.ModePerm)
	if err != nil && !os.IsExist(err) {
		log.Fatal(err)
	}
	return function
}
