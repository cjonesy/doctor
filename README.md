# Doctor

[![Go Report Card](https://goreportcard.com/badge/github.com/cjonesy/doctor)](https://goreportcard.com/report/github.com/cjonesy/doctor) [![Release](https://img.shields.io/github/release/cjonesy/doctor.svg)](https://github.com/cjonesy/doctor/releases/latest)

## Introduction

The goal of Doctor is to be an easy use framework for troubleshooting a user's
environment. It is inspired heavily by the [`brew doctor`](https://docs.brew.sh/Manpage#doctor-dr---list-checks---audit-debug-diagnostic_check-)
command found in [Homebrew](https://brew.sh/).

Checks are defined in a `.doctor.yml` file, typically found at the root of a
repository. Once Doctor is installed, users can simply run `doctor` from the
root of their repo checkout to check for problems with their environment setup.

## Installation

The latest version of Doctor can be found on the
[Releases](https://github.com/cjonesy/doctor/releases) tab.

## Configuration

**Example Config:**

```yaml
checks:
  - description: Ensure that go is installed
    fix: Run `brew install go`
    type: command-in-path
    command: go
  - description: Ensure ssh key exists
    fix: See https://docs.github.com/en/authentication/connecting-to-github-with-ssh/generating-a-new-ssh-key-and-adding-it-to-the-ssh-agent
    type: file-exists
    path: ~/.ssh/id_rsa
  - description: Ensure bashrc has eval statement
    fix: Run `echo 'eval "$(some_command init bash)"' > ~/.bashrc`
    type: file-contains
    path: ~/.bashrc
    contents: eval "$(some_command init bash)"
  - description: Ensure that terraform is the correct version
    fix: See https://www.terraform.io/downloads
    type: output-contains
    command: terraform --version
    contents: 1.2.2
```

### `checks`

This is where you define the checks you wish Doctor to perform. Each check must
have a [`description`](#description), [`fix`](#fix), and [`type`](#type).
Additional attributes may need to be set depending on the type of check.

#### `description`

A description of this check. This text will be displayed when `doctor` runs.

#### `fix`

Instructions the user can follow to correct an issue discovered by a check. This
will be displayed whenever a check fails.

#### `type`

The type of check to perform. See below for more details on the types of checks
that are supported.

##### `command-in-path`

Checks that a command is in the user's path.

The following attributes must also be set:

- `command` - the command to check for existence in user's path

##### `file-exists`

Checks that a file exists on the user's system.

The following attributes must also be set:

- `path` - the path in which the file is expected to exist

##### `file-contains`

Checks that a file contains specific text.

The following attributes must also be set:

- `path` - the path in which the file exists
- `content` - the text that is expected to exist within the file

##### `output-contains`

Checks that a command's output contains specific text.

The following attributes must also be set:

- `command` - the command to run
- `content` - the text that is expected to exist within the command's output

## How to contribute

This project has some clear Contribution Guidelines and expectations that you
can read here ([CONTRIBUTING](CONTRIBUTING.md)).

The contribution guidelines outline the process that you'll need to follow to
get a patch merged.

And you don't just have to write code. You can help out by writing
documentation, tests, or even by giving feedback about this work.

Thank you for contributing!
