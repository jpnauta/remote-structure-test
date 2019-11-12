package app

import (
	"os"

	"github.com/jpnauta/remote-structure-test/cmd/remote-structure-test/app/cmd"
)

func Run() error {
	c := cmd.NewRootCommand(os.Stdout, os.Stderr)
	return c.Execute()
}
