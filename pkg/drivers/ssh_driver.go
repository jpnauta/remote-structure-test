package drivers

import (
	"bytes"
	"fmt"
	"github.com/pkg/errors"
	"github.com/pkg/sftp"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
	"net"
	"os"
)

type SshDriver struct {
	client ssh.Client
}

func NewSshDriver(args DriverConfig) (Driver, error) {
	sshConfig := &ssh.ClientConfig{
		User: args.Username,
		Auth: []ssh.AuthMethod{
			ssh.Password(args.Password),
		},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}

	newClient, err := ssh.Dial("tcp", args.Host, sshConfig)
	if err != nil {
		return nil, err
	}
	return &SshDriver{
		client: *newClient,
	}, nil
}

func (d *SshDriver) Destroy() {
	if err := d.client.Close(); err != nil {
		logrus.Warnf("error closing connection: %s", err)
	}
}

func (d *SshDriver) Setup() error {
	return nil
}

func (d *SshDriver) RunCommand(command string) (string, string, error) {
	stdout, stderr, err := d.exec(command)
	if err != nil {
		return "", "", err
	}

	if stdout != "" {
		logrus.Infof("stdout: %s", stdout)
	}
	if stderr != "" {
		logrus.Infof("stderr: %s", stderr)
	}
	return stdout, stderr, nil
}

func (d *SshDriver) StatFile(target string) (os.FileInfo, error) {
	file, err := d.OpenSftpFile(target)
	if err != nil {
		return nil, err
	}

	return file.Stat()
}

func (d *SshDriver) ReadFile(target string) ([]byte, error) {
	file, err := d.OpenSftpFile(target)
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	_, err = file.WriteTo(buf)
	return buf.Bytes(), err
}

func (d *SshDriver) ReadDir(target string) ([]os.FileInfo, error) {
	ftpClient, err := sftp.NewClient(&d.client)
	if err != nil {
		return nil, err
	}

	return ftpClient.ReadDir(target)
}

func (d *SshDriver) exec(command string) (string, string, error) {
	session, err := d.client.NewSession()
	if err != nil {
		return "", "", errors.Wrap(err, "Error connecting to host")
	}
	defer session.Close()

	stdout, err := session.Output(command)
	if err != nil {
		return string(stdout), err.Error(), errors.Wrap(err, fmt.Sprintf("Error running command: %s", command))
	}
	return string(stdout), "", nil
}

func (d *SshDriver) OpenSftpFile(target string) (*sftp.File, error) {
	ftpClient, err := sftp.NewClient(&d.client)
	if err != nil {
		return nil, err
	}

	return ftpClient.Open(target)
}
