package version

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// AppInfo ...
type AppInfo struct {
	Name        string
	Description string
}

// GitInfo ...
type GitInfo struct {
	Commit string
	Branch string
	Repo   string
}

// BuildInfo ...
type BuildInfo struct {
	Timestamp time.Time
	Git       GitInfo
}

// Version ...
type Version struct {
	Raw    string
	Prefix string
	Major  int
	Minor  int
	Patch  int
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
	app.Name = name
	app.Description = description
}

// SetGitInfo ...
func SetGitInfo(commit, branch, repo string) {
	build.Git.Branch = branch
	build.Git.Commit = commit
	build.Git.Repo = repo
}

// SetBuildInfo set build info
func SetBuildInfo(timestamp string) {
	build.Timestamp, _ = time.Parse(time.UnixDate, timestamp)
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
