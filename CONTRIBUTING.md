# Contributing to goGenerateCFToken

Contributions to this project are welcomed.

## Table of Contents

- [Development Environment Setup](#development-environment-setup)
  - [Tools](#tools)
  - [Pulling the Project](#pulling-the-project)
- [Building Locally](#building-locally)
  - [Using Go](#using-go)
  - [Using GoReleaser](#using-goreleaser)
- [Style Guides](#style-guides)
  - [Semantic Branch Naming Guide](#semantic-branch-naming-guide)
  - [Commit Message Style Guide](#commit-message-style-guide)
  - [Pull Request Message Style Guide](#pull-request-message-style-guide)
  - [Release Message Format](#release-message-format)
- [GitHub Actions Workflows](#github-actions-workflows)
  - [Pull Requests](#pull-requests)
  - [Building Releases](#building-releases)
  - [Linting](#linting)
  - [Unit Testing](#go-source-code-unit-testing)
- [Developer Resources](#resources)

## Development Environment Setup

### Tools

- [Go](https://go.dev/)
- [Make](https://www.gnu.org/software/make/)
- [Cobra CLI](https://github.com/spf13/cobra-cli)
- [Golangci-Lint](https://golangci-lint.run/)
- [Mockery](https://vektra.github.io/mockery/latest/)
- [GoReleaser](https://goreleaser.com/)

### Pulling the Project

```bash
git clone https://github.com/nicholas-fedor/goGenerateCFToken.git
```

## Building Locally

### Using Go

- Windows Binary (exe):

  ```powershell
  go build -o gogeneratecftoken.exe .
  ```

- Linux Binary:

  ```bash
  go build -o gogeneratecftoken .
  ```

### Using GoReleaser

```bash
goreleaser build --single-target --snapshot --clean --config build/goreleaser/goreleaser.yaml
```

## Style Guides

### Semantic Branch Naming Guide

goGenerateCFToken uses semantic branch naming for structured branch names.

| Type       | Description                                                        | Examples                     |
|------------|--------------------------------------------------------------------|------------------------------|
| `chore`    | updating grunt tasks etc; no production code change                | `chore/update-build-script`  |
| `docs`     | changes to the documentation                                       | `docs/update-readme`         |
| `feat`     | new feature for the user, not a new feature for build script       | `feat/user-authentication`   |
| `fix`      | bug fix for the user, not a fix to a build script                  | `fix/login-error`            |
| `refactor` | refactoring production code, eg. renaming a variable               | `refactor/rename-user-class` |
| `style`    | formatting, missing semi colons, etc; no production code change    | `style/format-components`    |
| `test`     | adding missing tests, refactoring tests; no production code change | `test/add-unit-tests`        |

- If the branch is addressing an issue, then include the issue number in the branch name:

  ```text
  feat/#1-add-cli-feature
  ```

### Commit Message Style Guide

goGenerateCFToken uses the [Conventional Commits specification](https://www.conventionalcommits.org/en/v1.0.0/#summary) for structured commit messages.

Commit Message Structure:

```text
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

### Pull Request Message Style Guide

- Similar to commit messages
- Use imperative tense
- Title Format: Conventional Commit styling (i.e. type(scope): brief description)
- Body Format:

  ```text
  Brief Description
  ## Changes
  ```

### Release Message Format

It's expected that release notes will primarily address changes to the source code and largely omit changes to the project's CI/CD pipelines.

```text
Brief Overview
## Features
## Fixes
## Enhancements
## (Go) Dependency Updates
## CI Improvements/Dependency Updates
```

## GitHub Actions Workflows

### Pull Requests

- [`.github/workflows/pull-request.yaml`](.github/workflows/pull-request.yaml)
- Invoked when changes are made to the Go source code.

### Building Releases

[GoReleaser](https://github.com/goreleaser/goreleaser-action) is used to build the release files, which are uploaded to the associated release as tarball and zip archives.

To run a new build, update the tag, as follows:

```bash
git tag -a v0.1.0 -m "First release"
git push origin v0.1.0
```

Alternatively, use the GitHub website to [draft a new release](https://github.com/nicholas-fedor/goGenerateCFToken/releases/new).

### Linting

#### Go

- [.github/workflows/lint-go.yaml](.github/workflows/lint-go.yaml)

[Golangci-lint](https://golangci-lint.run/) is used lint the Go source code via the [Golangci-lint Action](https://github.com/golangci/golangci-lint-action).

This uses the configuration file `build/golangci-lint/golangci.yaml`.

If changes are detected by the workflow, i.e. linting needs to be done to the respective file, then the workflow will error and provide the relevant information.

#### GitHub Actions

- [.github/workflows/lint-gh.yaml](.github/workflows/lint-gh.yaml)

[Rhysd's actionlint](https://github.com/rhysd/actionlint) is used to lint the worflows in `.github/workflows`.

### Go Source Code Unit Testing

- [.github/workflows/test.yaml](.github/workflows/test.yaml)

Unit testing is run using Go's built-in `go test` functionality.
The results are uploaded to Codecov.

### Security Scanning

- [.github/workflows/security.yaml](.github/workflows/security.yaml)

Security scanning of the Go source code is done using the [Gosec Security Scanner](https://github.com/securego/gosec) and [Govulncheck-action](https://github.com/nicholas-fedor/govulncheck-action).
The results are uploaded to GitHub's CodeQL, which allows for integration with GitHub.

## Resources

- [Cloudflare's API Endpoints](https://developers.cloudflare.com/api-next)
- [Cloudflare's Go SDK](https://github.com/cloudflare/cloudflare-go)
- [spf13/Cobra](https://github.com/spf13/cobra)
- [spf13/Viper](https://github.com/spf13/viper)
