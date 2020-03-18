# kubectl-bindrole

Summarize RBAC roles for the specified subject (ServiceAccount, User and Group).

![screenshot](./img/screenshot.png)

## Installation

### for macOS

bindrole supports [homebrew](https://docs.brew.sh/Installation) :beer:

```
brew install Ladicle/kubectl-bindrole/kubectl-bindrole
```

### for other devices

The easiest way is to download binary from the [release page](https://github.com/Ladicle/kubectl-bindrole/releases).
You can also download this repository and install it using Makefile.

## Usage

```bash
$ kubectl bindrole -h  # or kubectl-bindrole -h
Summarize RBAC roles for the specified subject

Usage:
  kubectl bindrole [options] <SubjectName>

Examples:
  # Summarize roles tied to the "ci-bot" ServiceAccount.
  kubectl bindrole ci-bot

  # Summarize roles tied to the "developer" Group.
  kubectl bindrole -k Group developer

SubjectKinds:
  - ServiceAccount (default)
  - User
  - Group

Options:
  -h, --help                   Display this help message
  -n, --namespace string       Change the namespace scope for this CLI request
  -k, --subject-kind string    Set SubjectKind to summarize (default: ServiceAccount)
  -o, --options                List of all options for this command
      --version                Show version for this command

Use "kubectl bindrole --options" for a list of all options (applies to this command).
```

This command supports both kubectl-plugin mode and standalone mode.
