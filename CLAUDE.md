# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a Go version metadata library (`github.com/rbaliyan/go-version`) that captures and manages application version, build, and git information. It supports multiple injection methods with automatic fallback. Designed to be imported by other Go projects to embed version info into binaries.

## Build Commands

Standard Go commands (no Makefile):
- `go build ./...` - Build the package and CLI
- `go test ./...` - Run tests (note: no tests exist yet)
- `go vet ./...` - Run static analysis
- `go run ./cmd/go-version <command>` - Run CLI tool

When consuming this package, set version info via ldflags:
```bash
go build -ldflags="-X github.com/rbaliyan/go-version.VersionInfo=1.0.0 \
  -X github.com/rbaliyan/go-version.GitCommit=$(git rev-parse HEAD) \
  -X github.com/rbaliyan/go-version.GitBranch=$(git branch --show-current) \
  -X github.com/rbaliyan/go-version.BuildTimestamp=$(date -u '+%a %b %d %H:%M:%S %Z %Y')"
```

## Architecture

**Project structure:**
- `version.go` - Library code (imported by other projects)
- `cmd/go-version/main.go` - CLI tool for generating `.version` files

Stdlib-only dependencies (`os/exec`, `bufio`, `path/filepath`).

**Version source priority (highest to lowest):**
1. ldflags - Build-time injection via `-X` flags
2. Setters - Runtime calls to `SetVersion()`, `SetGitInfo()`, etc.
3. Version file - `.version` file in predefined locations
4. Git - Auto-detected from git repository (useful for `go run`)

**Key types:**
- `Version` - Semantic version (Major.Minor.Patch) with suffix support
- `BuildInfo` - Build timestamp and embedded GitInfo
- `GitInfo` - Commit, branch, repo metadata
- `AppInfo` - Application name, description, changelog

**Design patterns:**
- Package-level singleton state via global variables
- Idempotent setters (each setter only works once to prevent overwriting)
- Build-time injection through exported package variables (`VersionInfo`, `GitCommit`, etc.)
- `init()` function auto-initializes from injected variables

**Version parsing:**
- Handles optional 'v' prefix (v1.2.3 → 1.2.3)
- Supports suffixes after dash (1.2.3-dev.100 stores "dev.100" in Prefix field)
- Note: `Version.Prefix` is misnamed—it actually stores the *suffix* (text after the dash)
- Timestamp format: `time.UnixDate` (e.g., "Mon Jan 2 15:04:05 MST 2006")

**Version file format** (Key=Value):
```
VERSION=1.2.3
GIT_COMMIT=abc123
GIT_BRANCH=main
GIT_REPO=github.com/user/repo
BUILD_TIMESTAMP=Mon Jan 2 15:04:05 UTC 2006
```

**Version file search locations** (in order):
1. `./.version` (current working directory)
2. `<executable_dir>/.version` (next to binary)
3. `~/.config/<appname>/.version` (user config, requires SetAppInfo first)
4. `/etc/<appname>/.version` (system config, requires SetAppInfo first)

**Git auto-detection:**
- Runs automatically in `init()` if ldflags/file didn't provide values
- Uses: `git rev-parse HEAD`, `git describe --tags`, `git remote get-url origin`
- Fails silently if not in a git repo

**API usage:**
- `init()` auto-populates from ldflags, then version file, then git
- Consuming apps call `SetAppInfo()` at startup for app name/description
- `LoadFromFile(path)` - Manually load from specific file
- `LoadFromGit()` - Manually trigger git detection
- Retrieve with `Get()`, `Build()`, `Git()`, `App()` or display with `Print()`
