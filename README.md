Remote Structure Tests
====================

The Remote Structure Tests provide a powerful framework to validate the structure
of a remote host. These tests can be used to check the output of commands
in a host, as well as verify contents of the filesystem.

Tests can be run either through a Docker image or a standalone binary.

Inspired by and adapted from
[Google Container Structure Tests](https://github.com/GoogleContainerTools/container-structure-test).

## Installation

### Docker

A container image for running tests can be found at `jpnauta/remote-structure-test:latest`.

## Setup
To use remote structure tests to validate your host, you'll need the following:
- The remote structure test docker image or binary
- A remote SSH host to test against
- A test `.yaml` or `.json` file with user defined structure tests to run inside of the specified container image

## Example Run

An example run using Docker Compose:

```
version: '3'
services:
  structure_test:
    image: jpnauta/remote-structure-test-docker
    volumes:
      - "./config.yaml:/root/config.yaml"
    command: test --host localhost --username user --config ./config.yaml
```

Tests within this framework are specified through a YAML or JSON config file,
which is provided to the test driver via a CLI flag. Multiple config files may
be specified in a single test run. The config file will be loaded in by the
test driver, which will execute the tests in order. Within this config file,
three types of tests can be written:

- Command Tests (testing output/error of a specific command issued)
- File Existence Tests (making sure a file is, or isn't, present in the
file system of the host)
- File Content Tests (making sure files in the file system of the host
contain, or do not contain, specific contents)

## Command Tests
Command tests ensure that certain commands run properly in the target host.
Regexes can be used to check for expected or excluded strings in both `stdout`
and `stderr`. Additionally, any number of flags can be passed to the argument
as normal.

#### Supported Fields:

**NOTE: `schemaVersion` must be specified in all remote-structure-test yamls. The current version is `2.0.0`.**

- Name (`string`, **required**): The name of the test
- Command (`string`, **required**): The command to run in the test.
- Expected Output (`[]string`, *optional*): List of regexes that should
match the stdout from running the command.
- Excluded Output (`[]string`, *optional*): List of regexes that should **not**
match the stdout from running the command.
- Expected Error (`[]string`, *optional*): List of regexes that should
match the stderr from running the command.
- Excluded Error (`[]string`, *optional*): List of regexes that should **not**
match the stderr from running the command.

Example:
```yaml
commandTests:
  - name: "gunicorn flask"
    command: "which gunicorn"
    expectedOutput: ["/env/bin/gunicorn"]
  - name:  "apt-get upgrade"
    command: "apt-get -qqs upgrade"
    excludedOutput: [".*Inst.*Security.* | .*Security.*Inst.*"]
    excludedError: [".*Inst.*Security.* | .*Security.*Inst.*"]
```

## File Existence Tests
File existence tests check to make sure a specific file (or directory) exist
within the file system of the host. No contents of the files or directories
are checked. These tests can also be used to ensure a file or directory is
**not** present in the file system.

#### Supported Fields:

- Name (`string`, **required**): The name of the test
- Path (`string`, **required**): Path to the file or directory under test
- ShouldExist (`boolean`, **required**): Whether or not the specified file or
directory should exist in the file system
- Permissions (`string`, *optional*): The expected Unix permission string (e.g.
  drwxrwxrwx) of the files or directory.
- Uid (`int`, *optional*): The expected Unix user ID of the owner of the file
  or directory.
- Gid (`int`, *optional*): The expected Unix group ID of the owner of the file or directory.
- IsExecutableBy (`string`, *optional*): Checks if file is executable by a given user.
  One of `owner`, `group`, `other` or `any`

Example:
```yaml
fileExistenceTests:
- name: 'Root'
  path: '/'
  shouldExist: true
  permissions: '-rw-r--r--'
  uid: 1000
  gid: 1000
  isExecutableBy: 'group'
```

## File Content Tests
File content tests open a file on the file system and check its contents.
These tests assume the specified file **is a file**, and that it **exists**
(if unsure about either or these criteria, see the above
**File Existence Tests** section). Regexes can again be used to check for
expected or excluded content in the specified file.

#### Supported Fields:

- Name (`string`, **required**): The name of the test
- Path (`string`, **required**): Path to the file under test
- ExpectedContents (`string[]`, *optional*): List of regexes that
should match the contents of the file
- ExcludedContents (`string[]`, *optional*): List of regexes that
should **not** match the contents of the file

Example:
```yaml
fileContentTests:
- name: 'Debian Sources'
  path: '/etc/apt/sources.list'
  expectedContents: ['.*httpredir\.debian\.org.*']
  excludedContents: ['.*gce_debian_mirror.*']
```
