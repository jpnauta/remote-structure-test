package drivers

import (
	"os"
)

const (
	Ssh = "ssh"
)

type DriverConfig struct {
	Host     string // used by SSH driver
	Username string // used by SSH driver
	Password string // used by SSH driver
}

type Driver interface {
	Setup() error

	// run command on the host
	RunCommand(command string) (string, string, error)

	StatFile(path string) (os.FileInfo, error)

	ReadFile(path string) ([]byte, error)

	ReadDir(path string) ([]os.FileInfo, error)

	Destroy()
}

func InitDriverImpl(driver string) func(DriverConfig) (Driver, error) {
	switch driver {
	// future drivers will be added here
	case Ssh:
		return NewSshDriver
	default:
		return nil
	}
}
