package version

import (
	"fmt"
	"runtime"
)

// The current version of remote-structure-test
// This is a private field and is set through a compilation flag from the Makefile

var version = "v0.0.0-unset"

var buildDate string
var platform = fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH)

type Info struct {
	Version    string
	GitVersion string
	BuildDate  string
	GoVersion  string
	Compiler   string
	Platform   string
}

// Get returns the version and buildtime information about the binary
func GetVersion() *Info {
	// These variables typically come from -ldflags settings to `go build`
	return &Info{
		Version:   version,
		BuildDate: buildDate,
		GoVersion: runtime.Version(),
		Compiler:  runtime.Compiler,
		Platform:  platform,
	}
}
