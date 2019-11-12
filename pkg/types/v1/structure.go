package v1

import (
	"github.com/jpnauta/remote-structure-test/pkg/drivers"
	types "github.com/jpnauta/remote-structure-test/pkg/types/unversioned"
)

type StructureTest struct {
	DriverImpl         func(drivers.DriverConfig) (drivers.Driver, error)
	DriverArgs         drivers.DriverConfig
	SchemaVersion      string              `yaml:"schemaVersion"`
	CommandTests       []CommandTest       `yaml:"commandTests"`
	FileExistenceTests []FileExistenceTest `yaml:"fileExistenceTests"`
	FileContentTests   []FileContentTest   `yaml:"fileContentTests"`
}

func (st *StructureTest) NewDriver() (drivers.Driver, error) {
	return st.DriverImpl(st.DriverArgs)
}

func (st *StructureTest) SetDriverImpl(f func(drivers.DriverConfig) (drivers.Driver, error), args drivers.DriverConfig) {
	st.DriverImpl = f
	st.DriverArgs = args
}

func (st *StructureTest) RunAll(channel chan interface{}, file string) {
	fileProcessed := make(chan bool, 1)
	go st.runAll(channel, fileProcessed)
	<-fileProcessed
}

func (st *StructureTest) runAll(channel chan interface{}, fileProcessed chan bool) {
	st.RunCommandTests(channel)
	st.RunFileContentTests(channel)
	st.RunFileExistenceTests(channel)
	fileProcessed <- true
}

func (st *StructureTest) RunCommandTests(channel chan interface{}) {
	for _, test := range st.CommandTests {
		if !test.Validate(channel) {
			continue
		}
		res := &types.TestResult{
			Name: test.Name,
			Pass: false,
		}
		driver, err := st.NewDriver()
		if err != nil {
			res.Errorf("error creating driver: %s", err.Error())
			channel <- res
			continue
		}
		defer driver.Destroy()
		if err = driver.Setup(); err != nil {
			res.Errorf("error in setup: %s", err.Error())
			channel <- res
			continue
		}
		channel <- test.Run(driver)
	}
}

func (st *StructureTest) RunFileExistenceTests(channel chan interface{}) {
	for _, test := range st.FileExistenceTests {
		if !test.Validate(channel) {
			continue
		}
		res := &types.TestResult{
			Name: test.Name,
			Pass: false,
		}
		driver, err := st.NewDriver()
		if err != nil {
			res.Errorf("error creating driver: %s", err.Error())
			channel <- res
			continue
		}
		channel <- test.Run(driver)
		driver.Destroy()
	}
}

func (st *StructureTest) RunFileContentTests(channel chan interface{}) {
	for _, test := range st.FileContentTests {
		if !test.Validate(channel) {
			continue
		}
		res := &types.TestResult{
			Name: test.Name,
			Pass: false,
		}
		driver, err := st.NewDriver()
		if err != nil {
			res.Errorf("error creating driver: %s", err.Error())
			channel <- res
			continue
		}
		channel <- test.Run(driver)
		driver.Destroy()
	}
}
