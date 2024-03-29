package v1

import (
	"archive/tar"
	"fmt"
	"os"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/jpnauta/remote-structure-test/pkg/drivers"
	types "github.com/jpnauta/remote-structure-test/pkg/types/unversioned"
)

var defaultOwnership = -1

type FileExistenceTest struct {
	Name           string `yaml:"name"`           // name of test
	Path           string `yaml:"path"`           // file to check existence of
	ShouldExist    bool   `yaml:"shouldExist"`    // whether or not the file should exist
	Permissions    string `yaml:"permissions"`    // expected Unix permission string of the file, e.g. drwxrwxrwx
	Uid            int    `yaml:"uid"`            // ID of the owner of the file
	Gid            int    `yaml:"gid"`            // ID of the group of the file
	IsExecutableBy string `yaml:"isExecutableBy"` // name of group that file should be executable by
}

func (fe FileExistenceTest) MarshalYAML() (interface{}, error) {
	return FileExistenceTest{ShouldExist: true}, nil
}

func (fe *FileExistenceTest) UnmarshalYAML(unmarshal func(interface{}) error) error {
	// Create a type alias and call unmarshal on this type to unmarshal the yaml text into
	// struct, since calling unmarshal on FileExistenceTest will result in an infinite loop.
	type FileExistenceTestHolder FileExistenceTest
	holder := FileExistenceTestHolder{
		ShouldExist: true,
		Uid:         defaultOwnership,
		Gid:         defaultOwnership,
	}
	err := unmarshal(&holder)
	if err != nil {
		return err
	}
	*fe = FileExistenceTest(holder)
	return nil
}

func (ft FileExistenceTest) Validate(channel chan interface{}) bool {
	res := &types.TestResult{}
	if ft.Name == "" {
		res.Errorf("Please provide a valid name for every test")
	}
	res.Name = ft.Name
	if ft.Path == "" {
		res.Errorf("Please provide a valid file path for test %s", ft.Name)
	}
	if len(res.Errors) > 0 {
		channel <- res
		return false
	}
	return true
}

func (ft FileExistenceTest) LogName() string {
	return fmt.Sprintf("File Existence Test: %s", ft.Name)
}

func (ft FileExistenceTest) Run(driver drivers.Driver) *types.TestResult {
	result := &types.TestResult{
		Name:   ft.LogName(),
		Pass:   true,
		Errors: make([]string, 0),
	}
	logrus.Info(ft.LogName())
	var info os.FileInfo
	info, err := driver.StatFile(ft.Path)
	if info == nil && ft.ShouldExist {
		result.Errorf(errors.Wrap(err, "Error examining file in host").Error())
		result.Fail()
		return result
	}
	if ft.ShouldExist && err != nil {
		result.Errorf("File %s should exist but does not, got error: %s", ft.Path, err)
		result.Fail()
	} else if !ft.ShouldExist && err == nil {
		result.Errorf("File %s should not exist but does", ft.Path)
		result.Fail()
	}
	if ft.Permissions != "" && info != nil {
		perms := info.Mode()
		if perms.String() != ft.Permissions {
			result.Errorf("%s has incorrect permissions. Expected: %s, Actual: %s", ft.Path, ft.Permissions, perms.String())
			result.Fail()
		}
	}
	if ft.IsExecutableBy != "" {
		perms := info.Mode()
		switch ft.IsExecutableBy {
		case "any":
			if perms&0111 == 0 {
				result.Errorf("%s has incorrect executable bit. Expected to be executable by any, Actual: %s", ft.Path, perms.String())
				result.Fail()
			}
		case "owner":
			if perms&0100 == 0 {
				result.Errorf("%s has incorrect executable bit. Expected to be executable by owner, Actual: %s", ft.Path, perms.String())
				result.Fail()
			}
		case "group":
			if perms&0010 == 0 {
				result.Errorf("%s has incorrect executable bit. Expected to be executable by group, Actual: %s", ft.Path, perms.String())
				result.Fail()
			}
		case "other":
			if perms&0001 == 0 {
				result.Errorf("%s has incorrect executable bit. Expected to be executable by other, Actual: %s", ft.Path, perms.String())
				result.Fail()
			}
		default:
			result.Errorf("%s not recognised as a valid option", ft.IsExecutableBy)
			result.Fail()
		}
	}
	if ft.Uid != defaultOwnership || ft.Gid != defaultOwnership {
		header, ok := info.Sys().(*tar.Header)
		if ok {
			if ft.Uid != defaultOwnership && header.Uid != ft.Uid {
				result.Errorf("%s has incorrect user ownership. Expected: %d, Actual: %d", ft.Path, ft.Uid, header.Uid)
				result.Fail()
			}
			if ft.Gid != defaultOwnership && header.Gid != ft.Gid {
				result.Errorf("%s has incorrect group ownership. Expected: %d, Actual: %d", ft.Path, ft.Gid, header.Gid)
				result.Fail()
			}
		} else {
			result.Errorf("Error checking ownership of file %s", ft.Path)
			result.Fail()
		}
	}
	return result
}
