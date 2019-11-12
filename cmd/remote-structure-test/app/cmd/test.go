package cmd

import (
	"fmt"
	"io"
	"io/ioutil"

	"github.com/jpnauta/remote-structure-test/cmd/remote-structure-test/app/cmd/test"
	"os"

	"github.com/jpnauta/remote-structure-test/pkg/color"
	"github.com/jpnauta/remote-structure-test/pkg/config"
	"github.com/jpnauta/remote-structure-test/pkg/drivers"
	"github.com/jpnauta/remote-structure-test/pkg/output"
	"github.com/jpnauta/remote-structure-test/pkg/types/unversioned"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	opts = &config.StructureTestOptions{}

	args       *drivers.DriverConfig
	driverImpl func(drivers.DriverConfig) (drivers.Driver, error)
)

func NewCmdTest(rootCmd *cobra.Command, out io.Writer) *cobra.Command {
	var testCmd = &cobra.Command{
		Use:   "test",
		Short: "Runs the tests",
		Long:  `Runs the tests`,
		Args: func(cmd *cobra.Command, _ []string) error {
			return test.ValidateArgs(opts)
		},
		RunE: func(cmd *cobra.Command, _ []string) error {
			if opts.TestReport != "" {
				// Force JsonOutput
				opts.JSON = true
				testReportFile, err := os.Create(opts.TestReport)
				if err != nil {
					return err
				}
				rootCmd.SetOutput(testReportFile)
				out = testReportFile // override writer
			}

			if opts.Quiet {
				out = ioutil.Discard
			}

			color.NoColor = opts.NoColor

			return run(out)
		},
	}

	AddTestFlags(testCmd)
	return testCmd
}

func run(out io.Writer) error {
	args = &drivers.DriverConfig{
		Host:     opts.Host,
		Username: opts.Username,
		Password: opts.Password,
	}

	driverImpl = drivers.InitDriverImpl(opts.Driver)
	if driverImpl == nil {
		logrus.Fatalf("unsupported driver type: %s", opts.Driver)
	}
	channel := make(chan interface{}, 1)
	go runTests(out, channel, args, driverImpl)
	return test.ProcessResults(out, opts.JSON, channel)
}

func runTests(out io.Writer, channel chan interface{}, args *drivers.DriverConfig, driverImpl func(drivers.DriverConfig) (drivers.Driver, error)) {
	for _, file := range opts.ConfigFiles {
		if !opts.JSON {
			output.Banner(out, file)
		}
		tests, err := test.Parse(file, args, driverImpl)
		if err != nil {
			channel <- &unversioned.TestResult{
				Errors: []string{
					fmt.Sprintf("error parsing config file: %s", err),
				},
			}
			continue // Continue with other config files
		}
		tests.RunAll(channel, file)
	}
	close(channel)
}

func AddTestFlags(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&opts.Host, "host", "H", "", "host and port to connect to")
	cmd.Flags().StringVarP(&opts.Username, "username", "u", "", "username to connect to host")
	cmd.Flags().StringVarP(&opts.Password, "password", "p", "", "password to connect to host")
	cmd.Flags().StringVarP(&opts.Driver, "driver", "d", "ssh", "driver to use when running tests")

	cmd.Flags().BoolVarP(&opts.Quiet, "quiet", "q", false, "flag to suppress output")
	cmd.Flags().BoolVarP(&opts.JSON, "json", "j", false, "output test results in json format")
	cmd.Flags().BoolVar(&opts.NoColor, "no-color", false, "no color in the output")

	cmd.Flags().StringArrayVarP(&opts.ConfigFiles, "config", "c", []string{}, "test config files")
	cmd.Flags().StringVar(&opts.TestReport, "test-report", "", "generate JSON test report and write it to specified file.")
}
