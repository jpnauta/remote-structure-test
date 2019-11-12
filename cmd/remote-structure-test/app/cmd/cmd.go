package cmd

import (
	"io"

	"github.com/jpnauta/remote-structure-test/pkg/version"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	v string
)

func NewRootCommand(out, err io.Writer) *cobra.Command {
	var rootCmd = &cobra.Command{
		Use:   "remote-structure-test",
		Short: "remote-structure-test provides a framework to test the structure of a remote host",
		Long: `remote-structure-test provides a powerful framework to validate
the structure of a remote host.
These tests can be used to check the output of commands,
as well as verify contents of the filesystem.`,
	}

	rootCmd.PersistentPreRunE = func(cmd *cobra.Command, _ []string) error {
		if err := SetUpLogs(err, v); err != nil {
			return err
		}

		rootCmd.SilenceUsage = true
		logrus.Infof("remote-structure-test %+v", version.GetVersion())
		return nil
	}

	rootCmd.SilenceErrors = true
	rootCmd.AddCommand(NewCmdVersion(rootCmd, out))
	rootCmd.AddCommand(NewCmdTest(rootCmd, out))

	rootCmd.PersistentFlags().StringVarP(&v, "verbosity", "v", logrus.WarnLevel.String(), "Log level (debug, info, warn, error, fatal, panic)")

	return rootCmd
}

func SetUpLogs(out io.Writer, level string) error {
	logrus.SetOutput(out)
	lvl, err := logrus.ParseLevel(v)
	if err != nil {
		return errors.Wrap(err, "parsing log level")
	}
	logrus.SetLevel(lvl)
	return nil
}
