package types

import (
	"github.com/jpnauta/remote-structure-test/pkg/drivers"
	"github.com/jpnauta/remote-structure-test/pkg/types/v1"
)

type StructureTest interface {
	SetDriverImpl(func(drivers.DriverConfig) (drivers.Driver, error), drivers.DriverConfig)
	NewDriver() (drivers.Driver, error)
	RunAll(chan interface{}, string)
}

var SchemaVersions map[string]func() StructureTest = map[string]func() StructureTest{
	"1.0.0": func() StructureTest { return new(v1.StructureTest) },
}

type SchemaVersion struct {
	SchemaVersion string `yaml:"schemaVersion"`
}

type Unmarshaller func([]byte, interface{}) error
