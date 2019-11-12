package v1

import (
	"fmt"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/jpnauta/remote-structure-test/pkg/drivers"
	types "github.com/jpnauta/remote-structure-test/pkg/types/unversioned"
	"github.com/jpnauta/remote-structure-test/pkg/utils"
)

type CommandTest struct {
	Name           string   `yaml:"name"`
	Command        string   `yaml:"command"`
	ExpectedOutput []string `yaml:"expectedOutput"`
	ExcludedOutput []string `yaml:"excludedOutput"`
	ExpectedError  []string `yaml:"expectedError"`
	ExcludedError  []string `yaml:"excludedError"` // excluded error from running command
}

func (ct *CommandTest) Validate(channel chan interface{}) bool {
	res := &types.TestResult{}
	if ct.Name == "" {
		res.Error("Please provide a valid name for every test")
	}
	res.Name = ct.Name
	if ct.Command == "" {
		res.Errorf("Please provide a valid command to run for test %s", ct.Name)
	}
	if len(res.Errors) > 0 {
		channel <- res
		return false
	}
	return true
}

func (ct *CommandTest) LogName() string {
	return fmt.Sprintf("Command Test: %s", ct.Name)
}

func (ct *CommandTest) Run(driver drivers.Driver) *types.TestResult {
	logrus.Debug(ct.LogName())
	start := time.Now()
	stdout, stderr, err := driver.RunCommand(ct.Command)
	end := time.Now()
	duration := end.Sub(start)
	result := &types.TestResult{
		Name:     ct.LogName(),
		Pass:     true,
		Errors:   make([]string, 0),
		Stderr:   stderr,
		Stdout:   stdout,
		Duration: duration,
	}
	if err != nil {
		result.Fail()
		result.Error(err.Error())
		return result
	}

	ct.CheckOutput(result, stdout, stderr)
	return result
}

func (ct *CommandTest) CheckOutput(result *types.TestResult, stdout string, stderr string) {
	for _, errStr := range ct.ExpectedError {
		if !utils.CompileAndRunRegex(errStr, stderr, true) {
			result.Errorf("Expected string '%s' not found in error '%s'", errStr, stderr)
			result.Fail()
		}
	}
	for _, errStr := range ct.ExcludedError {
		if !utils.CompileAndRunRegex(errStr, stderr, false) {
			result.Errorf("Excluded string '%s' found in error '%s'", errStr, stderr)
			result.Fail()
		}
	}
	for _, outStr := range ct.ExpectedOutput {
		if !utils.CompileAndRunRegex(outStr, stdout, true) {
			result.Errorf("Expected string '%s' not found in output '%s'", outStr, stdout)
			result.Fail()
		}
	}
	for _, outStr := range ct.ExcludedOutput {
		if !utils.CompileAndRunRegex(outStr, stdout, false) {
			result.Errorf("Excluded string '%s' found in output '%s'", outStr, stdout)
			result.Fail()
		}
	}
}
