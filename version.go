// Package version provides version, build, and git metadata for Go applications.
//
// Version info can be loaded from:
//   - ldflags: Build-time injection via -X flags (loaded automatically in init)
//   - Version file: Call LoadFromFile() to load from a .version file
//   - Git: Call LoadFromGit() to detect from git repository
//
// # Build with ldflags
//
// Use the go-version CLI tool to generate ldflags:
//
//	go build -ldflags="$(go-version ldflags)" ./cmd/myapp
//	go build -ldflags="$(go-version ldflags -static)" ./cmd/myapp  # for CI
//
// # Load from file or git
//
//	version.LoadFromFile(".version")  // load from specific file
//	version.LoadFromGit()             // detect from git repo
//
// # Basic usage
//
//	import version "github.com/rbaliyan/go-version"
//
//	func main() {
//	    version.SetAppInfo("myapp", "My application")
//	    version.Print()
//	}
package version

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// AppInfo application name and dscription
type AppInfo struct {
	// appl name
	Name string
	// app description
	Description string
	// changelog for this build
	Changelog string
}

// GitInfo git details for the this build od the app
type GitInfo struct {
	// git commit
	Commit string
	// git brnach
	Branch string
	// git repo
	Repo string
}

// BuildInfo build timestamp and git information for the repo
type BuildInfo struct {
	Timestamp time.Time
	Git       GitInfo
}

// Version application version details
type Version struct {
	// single string message
	Raw string
	// version prefix
	Prefix string
	// major version part
	Major int
	// minor version part
	Minor int
	// path version path
	Patch int
}

var (
	build   = BuildInfo{}
	version = Version{}
	app     = AppInfo{}

	// BuildTimestamp ...
	BuildTimestamp = ""
	// GitCommit ...
	GitCommit = ""
	// GitBranch ...
	GitBranch = ""
	// GitRepo ...
	GitRepo = ""
	// VersionInfo ...
	VersionInfo = ""
)

func init() {
	// Load from ldflags if provided at build time
	SetBuildInfo(BuildTimestamp)
	SetGitInfo(GitCommit, GitBranch, GitRepo)
	SetVersion(VersionInfo)
}

func (ver Version) String() string {
	return fmt.Sprintf("%s %d.%d.%d", ver.Prefix, ver.Major, ver.Minor, ver.Patch)
}

func (app AppInfo) String() string {
	return fmt.Sprintf("%s: \n\t %s", app.Name, app.Description)
}

func (build BuildInfo) String() string {
	return fmt.Sprintf("Timestamp : %v, Git: %v", build.Timestamp, build.Git)
}

func (git GitInfo) String() string {
	return fmt.Sprintf("Repo: %s, Branch: %s, Commit: %s", git.Repo, git.Branch, git.Commit)
}

// SetAppInfo ...
func SetAppInfo(name, description string) {
	if app.Name == "" {
		app.Name = name
		app.Description = description
	}
}

// SetGitInfo set git details for app
func SetGitInfo(commit, branch, repo string) {
	if build.Git.Commit == "" {
		build.Git.Branch = branch
		build.Git.Commit = commit
		build.Git.Repo = repo
	}
}

// SetBuildInfo set build info
func SetBuildInfo(timestamp string) {
	if build.Timestamp.IsZero() {
		build.Timestamp, _ = time.Parse(time.UnixDate, timestamp)
	}
}

// SetChangelog set application changelog
func SetChangelog(changelog string) {
	if app.Changelog == "" {
		app.Changelog = changelog
	}
}

// SetChangelogFromFile read changelog from a file
func SetChangelogFromFile(path string) error {
	if app.Changelog == "" {
		b, err := os.ReadFile(path) // just pass the file name
		if err != nil {
			return err
		}
		app.Changelog = string(b)
	}
	return nil
}

// SetVersion ...
func SetVersion(ver string) {
	version.Raw = ver

	// Strip 'v' prefix if present for parsing
	versionStr := strings.TrimPrefix(ver, "v")

	// Split into base version and suffix (if any)
	// Handle formats like: 1.2.3, 1.2.3-dev, 1.2.3-dev.100
	parts := strings.SplitN(versionStr, "-", 2)
	baseVersion := parts[0]
	suffix := ""
	if len(parts) > 1 {
		suffix = parts[1]
	}

	// Parse the base version (X.Y.Z)
	verparts := strings.Split(baseVersion, ".")
	if len(verparts) >= 3 {
		version.Major, _ = strconv.Atoi(verparts[0])
		version.Minor, _ = strconv.Atoi(verparts[1])
		version.Patch, _ = strconv.Atoi(verparts[2])
		version.Prefix = suffix
	}
}

// Get ...
func Get() Version {
	return version
}

// Build ...
func Build() BuildInfo {
	return build
}

// Git ...
func Git() GitInfo {
	return build.Git
}

// App ...
func App() AppInfo {
	return app
}

// Print ...
func Print() {
	fmt.Println("Running:", App())
	fmt.Println("Version:", Get())
	fmt.Println("Build:", Build())
}

// LoadFromFile loads version information from a key=value file.
// Keys: VERSION, GIT_COMMIT, GIT_BRANCH, GIT_REPO, BUILD_TIMESTAMP
func LoadFromFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		switch key {
		case "VERSION":
			if version.Raw == "" {
				SetVersion(value)
			}
		case "GIT_COMMIT":
			if build.Git.Commit == "" {
				build.Git.Commit = value
			}
		case "GIT_BRANCH":
			if build.Git.Branch == "" {
				build.Git.Branch = value
			}
		case "GIT_REPO":
			if build.Git.Repo == "" {
				build.Git.Repo = value
			}
		case "BUILD_TIMESTAMP":
			if build.Timestamp.IsZero() {
				build.Timestamp, _ = time.Parse(time.UnixDate, value)
			}
		}
	}
	return scanner.Err()
}

// LoadFromGit reads version information directly from git commands.
// This is useful during development with 'go run'.
func LoadFromGit() error {
	// Check if we're in a git repository
	if err := exec.Command("git", "rev-parse", "--git-dir").Run(); err != nil {
		return err
	}

	// Get commit hash
	if build.Git.Commit == "" {
		if out, err := exec.Command("git", "rev-parse", "HEAD").Output(); err == nil {
			build.Git.Commit = strings.TrimSpace(string(out))
		}
	}

	// Get branch name
	if build.Git.Branch == "" {
		if out, err := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD").Output(); err == nil {
			build.Git.Branch = strings.TrimSpace(string(out))
		}
	}

	// Get remote URL
	if build.Git.Repo == "" {
		if out, err := exec.Command("git", "remote", "get-url", "origin").Output(); err == nil {
			build.Git.Repo = strings.TrimSpace(string(out))
		}
	}

	// Get version from git describe (tags)
	if version.Raw == "" {
		if out, err := exec.Command("git", "describe", "--tags", "--always").Output(); err == nil {
			SetVersion(strings.TrimSpace(string(out)))
		}
	}

	return nil
}
