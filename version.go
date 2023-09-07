package version

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// AppInfo application name and dscription
type AppInfo struct {
	// appl name
	Name        string
	// app description
	Description string
	// changelog for this build
	Changelog   string
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
	verparts := strings.Split(ver, ".")
	if len(verparts) == 3 {
		version.Major, _ = strconv.Atoi(verparts[0])
		version.Minor, _ = strconv.Atoi(verparts[1])
		parts := strings.Split(verparts[2], "-")
		version.Patch, _ = strconv.Atoi(parts[0])
		if len(parts) > 1 {
			version.Prefix = strings.Join(parts[1:], "-")
		}
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
