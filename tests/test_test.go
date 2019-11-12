package tests

import (
	"github.com/jpnauta/remote-structure-test/cmd/remote-structure-test/app/cmd"
	"github.com/kami-zh/go-capturer"
	"github.com/stretchr/testify/assert"
	"os"
	"regexp"
	"strconv"
	"testing"
)

var (
	sshHost     = "localhost:22"
	sshUsername = "root"
	sshPassword = "root"
)

func StandardArgs() []string {
	return []string{"test", "-H", sshHost, "-u", sshUsername, "-p", sshPassword, "--config"}
}

func RunCommand(args []string) (error, string) {
	var err error
	output := capturer.CaptureOutput(func() {
		c := cmd.NewRootCommand(os.Stdout, os.Stderr)
		c.SetArgs(args)
		err = c.Execute()
	})

	return err, output
}

func ParseTestVar(t *testing.T, output string, expr string) int {
	varMatch := regexp.MustCompile(expr).FindStringSubmatch(output)[1]
	match, err := strconv.Atoi(varMatch)
	if err != nil {
		assert.Fail(t, err.Error())
	}
	return match
}

func ParseTestResults(t *testing.T, output string) (int, int) {
	return ParseTestVar(t, output, "Passes:\\W+(\\d+)"), ParseTestVar(t, output, "Failures:\\W+(\\d+)")
}

func TestTestHelpMessage(t *testing.T) {
	err, output := RunCommand([]string{"test", "-h"})

	assert.Nil(t, err)
	contains := []string{
		"-c, --config stringArray", "-h, --help",
		"-H, --host string", "-j, --json", "--no-color",
		"-p, --password string", "-q, --quiet",
		"--test-report string",
	}
	for _, substr := range contains {
		assert.Contains(t, output, substr)
	}
}

func TestTestHostRequired(t *testing.T) {
	err, _ := RunCommand([]string{"test"})

	assert.Equal(t, err.Error(), "Please supply host to run tests on")
}

func TestTestConfigRequired(t *testing.T) {
	err, _ := RunCommand([]string{"test", "-H", sshHost})

	assert.Equal(t, err.Error(), "Please provide at least one test config file")
}

func TestTestPassingConfig(t *testing.T) {
	err, output := RunCommand(append(StandardArgs(), "./fixtures/passing.yaml"))

	passes, failures := ParseTestResults(t, output)
	assert.Equal(t, failures, 0)
	assert.Equal(t, passes, 7)
	assert.Nil(t, err)
}

func TestTestFailingConfig(t *testing.T) {
	err, output := RunCommand(append(StandardArgs(), "./fixtures/failing.yaml"))

	passes, failures := ParseTestResults(t, output)
	assert.Equal(t, failures, 7)
	assert.Equal(t, passes, 0)
	assert.Equal(t, err.Error(), "FAIL")
}

func TestTestInvalidConfig(t *testing.T) {
	err, output := RunCommand(append(StandardArgs(), "./fixtures/invalid.yaml"))

	assert.Contains(t, output, "line 4: field foo not found in type v1.CommandTest")
	assert.Equal(t, err.Error(), "FAIL")
}

func TestTestUnsupportedConfig(t *testing.T) {
	err, output := RunCommand(append(StandardArgs(), "./fixtures/unsupported.yaml"))

	assert.Contains(t, output, "Unsupported schema version: 2.0.0")
	assert.Equal(t, err.Error(), "FAIL")
}
