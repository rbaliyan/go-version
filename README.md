# go-version

[![CI](https://github.com/rbaliyan/go-version/actions/workflows/ci.yml/badge.svg)](https://github.com/rbaliyan/go-version/actions/workflows/ci.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/rbaliyan/go-version.svg)](https://pkg.go.dev/github.com/rbaliyan/go-version)
[![Go Report Card](https://goreportcard.com/badge/github.com/rbaliyan/go-version)](https://goreportcard.com/report/github.com/rbaliyan/go-version)
[![Release](https://img.shields.io/github/v/release/rbaliyan/go-version)](https://github.com/rbaliyan/go-version/releases/latest)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)
[![OpenSSF Scorecard](https://api.scorecard.dev/projects/github.com/rbaliyan/go-version/badge)](https://scorecard.dev/viewer/?uri=github.com/rbaliyan/go-version)

A lightweight Go package for embedding version, build, and git metadata into your application binaries. Supports multiple injection methods with automatic fallback.

## Installation

### CLI Tool

**Install script:**
```bash
curl -sSfL https://raw.githubusercontent.com/rbaliyan/go-version/main/install.sh | sh
```

**Go install:**
```bash
go install github.com/rbaliyan/go-version/cmd/go-version@latest
```

**Download binary:**

Download from [GitHub Releases](https://github.com/rbaliyan/go-version/releases) for your platform.

### As a Library

```bash
go get github.com/rbaliyan/go-version
```

## CLI Tool

The `go-version` CLI generates `.version` files and ldflags for CI pipelines.

### Commands

```bash
go-version file      # Generate a .version file
go-version ldflags   # Generate -ldflags for go build
go-version show      # Show git version info
go-version version   # Show go-version CLI version
```

### Generate a version file

```bash
# Generate from git info
go-version file

# Custom output path
go-version file -o build/.version

# Manual version override
go-version file -v 1.2.3
```

### Generate ldflags for go build

```bash
# Use in build command (shell substitutions)
go build -ldflags="$(go-version ldflags)" ./cmd/myapp

# Static values for CI pipelines
go build -ldflags="$(go-version ldflags -static)" ./cmd/myapp

# Custom package path
go build -ldflags="$(go-version ldflags -p mycompany/myapp)" ./cmd/myapp
```

### Show current git info

```bash
go-version show
```

Output:
```
Version:  v1.0.2
Commit:   f663cfdfb69bfd922a55e56e29a7784aab73e8c3
Branch:   master
Repo:     git@github.com:user/repo.git
```

### Shell Completions

Shell completions are included in the release archives:
- Bash: `completions/go-version.bash`
- Zsh: `completions/go-version.zsh`
- Fish: `completions/go-version.fish`

## Library Usage

### Basic Setup

```go
package main

import (
    "fmt"
    version "github.com/rbaliyan/go-version"
)

func main() {
    // Set application info (optional)
    version.SetAppInfo("myapp", "My awesome application")

    // Print all version info
    version.Print()

    // Or access individual components
    fmt.Printf("Version: %d.%d.%d\n", version.Get().Major, version.Get().Minor, version.Get().Patch)
    fmt.Printf("Commit: %s\n", version.Git().Commit)
}
```

### Build with Version Info

Inject version metadata at build time using `-ldflags`:

```bash
go build -ldflags="\
  -X github.com/rbaliyan/go-version.VersionInfo=1.2.3 \
  -X github.com/rbaliyan/go-version.GitCommit=$(git rev-parse HEAD) \
  -X github.com/rbaliyan/go-version.GitBranch=$(git branch --show-current) \
  -X github.com/rbaliyan/go-version.GitRepo=$(git remote get-url origin) \
  -X 'github.com/rbaliyan/go-version.BuildTimestamp=$(date -u "+%a %b %d %H:%M:%S %Z %Y")'"
```

Or use the CLI:

```bash
go build -ldflags="$(go-version ldflags -static)" ./cmd/myapp
```

## Version Sources

Version info can be loaded from (in priority order):

1. **ldflags** - Build-time injection via `-X` flags (loaded automatically)
2. **Setters** - Runtime calls to `SetVersion()`, `SetGitInfo()`, etc.
3. **Version file** - Call `LoadFromFile()` to load from a `.version` file
4. **Git** - Call `LoadFromGit()` to detect from git repository

### Version File Format

Create a `.version` file (Key=Value format):

```
VERSION=1.2.3
GIT_COMMIT=abc123def456
GIT_BRANCH=main
GIT_REPO=github.com/user/repo
BUILD_TIMESTAMP=Mon Jan 2 15:04:05 UTC 2006
```

The package searches for `.version` in these locations (in order):
1. Current working directory (`./.version`)
2. Executable directory (`<exe_dir>/.version`)
3. User config (`~/.config/<appname>/.version`) - requires `SetAppInfo()` first
4. System config (`/etc/<appname>/.version`) - requires `SetAppInfo()` first

### Git Detection

Call `LoadFromGit()` to read version info from git commands:

```go
version.LoadFromGit()
```

This reads:
- Commit hash from `git rev-parse HEAD`
- Branch from `git rev-parse --abbrev-ref HEAD`
- Version from `git describe --tags --always`
- Remote URL from `git remote get-url origin`

## API

### Setters

| Function | Description |
|----------|-------------|
| `SetAppInfo(name, description)` | Set application name and description |
| `SetVersion(ver)` | Parse and set semantic version (supports `v` prefix and suffixes like `1.2.3-dev`) |
| `SetBuildInfo(timestamp)` | Set build timestamp (format: `time.UnixDate`) |
| `SetGitInfo(commit, branch, repo)` | Set git metadata |
| `SetChangelog(changelog)` | Set changelog text |
| `SetChangelogFromFile(path)` | Load changelog from file |
| `LoadFromFile(path)` | Load version info from a specific `.version` file |
| `LoadFromGit()` | Manually trigger git auto-detection |

All setters are idempotent—they only set values once and ignore subsequent calls.

### Getters

| Function | Returns |
|----------|---------|
| `Get()` | `Version` struct with Major, Minor, Patch, Raw, Prefix fields |
| `Build()` | `BuildInfo` struct with Timestamp and Git info |
| `Git()` | `GitInfo` struct with Commit, Branch, Repo |
| `App()` | `AppInfo` struct with Name, Description, Changelog |
| `Print()` | Outputs all version info to stdout |

### Injected Variables

These package-level variables can be set via `-ldflags -X`:

- `VersionInfo` - Version string (e.g., "1.2.3" or "v1.2.3-dev")
- `GitCommit` - Git commit hash
- `GitBranch` - Git branch name
- `GitRepo` - Git repository URL
- `BuildTimestamp` - Build time in `time.UnixDate` format

## License

MIT License - see [LICENSE](LICENSE) file.
