package version

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// resetState clears all package-level globals so tests don't interfere with each other.
func resetState() {
	version = Version{}
	build = BuildInfo{}
	app = AppInfo{}
}

// --- SetVersion tests ---

func TestSetVersion_SemVer(t *testing.T) {
	resetState()
	SetVersion("1.2.3")

	v := Get()
	if v.Raw != "1.2.3" {
		t.Errorf("Raw = %q, want %q", v.Raw, "1.2.3")
	}
	if v.Major != 1 || v.Minor != 2 || v.Patch != 3 {
		t.Errorf("version = %d.%d.%d, want 1.2.3", v.Major, v.Minor, v.Patch)
	}
	if v.Prefix != "" {
		t.Errorf("Prefix = %q, want empty", v.Prefix)
	}
}

func TestSetVersion_VPrefix(t *testing.T) {
	resetState()
	SetVersion("v2.5.10")

	v := Get()
	if v.Raw != "v2.5.10" {
		t.Errorf("Raw = %q, want %q", v.Raw, "v2.5.10")
	}
	if v.Major != 2 || v.Minor != 5 || v.Patch != 10 {
		t.Errorf("version = %d.%d.%d, want 2.5.10", v.Major, v.Minor, v.Patch)
	}
}

func TestSetVersion_WithSuffix(t *testing.T) {
	resetState()
	SetVersion("v1.0.0-beta.1")

	v := Get()
	if v.Raw != "v1.0.0-beta.1" {
		t.Errorf("Raw = %q, want %q", v.Raw, "v1.0.0-beta.1")
	}
	if v.Major != 1 || v.Minor != 0 || v.Patch != 0 {
		t.Errorf("version = %d.%d.%d, want 1.0.0", v.Major, v.Minor, v.Patch)
	}
	if v.Prefix != "beta.1" {
		t.Errorf("Prefix = %q, want %q", v.Prefix, "beta.1")
	}
}

func TestSetVersion_DevSuffix(t *testing.T) {
	resetState()
	SetVersion("3.1.4-dev")

	v := Get()
	if v.Prefix != "dev" {
		t.Errorf("Prefix = %q, want %q", v.Prefix, "dev")
	}
}

func TestSetVersion_MultipleDashSuffix(t *testing.T) {
	resetState()
	SetVersion("1.2.3-rc-1")

	v := Get()
	if v.Major != 1 || v.Minor != 2 || v.Patch != 3 {
		t.Errorf("version = %d.%d.%d, want 1.2.3", v.Major, v.Minor, v.Patch)
	}
	if v.Prefix != "rc-1" {
		t.Errorf("Prefix = %q, want %q", v.Prefix, "rc-1")
	}
}

func TestSetVersion_Empty(t *testing.T) {
	resetState()
	SetVersion("")

	v := Get()
	if v.Raw != "" {
		t.Errorf("Raw = %q, want empty", v.Raw)
	}
	if v.Major != 0 || v.Minor != 0 || v.Patch != 0 {
		t.Errorf("version = %d.%d.%d, want 0.0.0", v.Major, v.Minor, v.Patch)
	}
}

func TestSetVersion_TooFewParts(t *testing.T) {
	resetState()
	SetVersion("1.2")

	v := Get()
	if v.Raw != "1.2" {
		t.Errorf("Raw = %q, want %q", v.Raw, "1.2")
	}
	// With fewer than 3 parts, Major/Minor/Patch should remain zero
	if v.Major != 0 || v.Minor != 0 || v.Patch != 0 {
		t.Errorf("version = %d.%d.%d, want 0.0.0 (not enough parts)", v.Major, v.Minor, v.Patch)
	}
}

func TestSetVersion_NonNumeric(t *testing.T) {
	resetState()
	SetVersion("abc.def.ghi")

	v := Get()
	// strconv.Atoi returns 0 on failure
	if v.Major != 0 || v.Minor != 0 || v.Patch != 0 {
		t.Errorf("version = %d.%d.%d, want 0.0.0 for non-numeric", v.Major, v.Minor, v.Patch)
	}
}

func TestSetVersion_Overwrites(t *testing.T) {
	resetState()
	SetVersion("1.0.0")
	SetVersion("2.0.0")

	v := Get()
	// SetVersion is NOT idempotent - it always overwrites
	if v.Raw != "2.0.0" {
		t.Errorf("Raw = %q, want %q (SetVersion should overwrite)", v.Raw, "2.0.0")
	}
	if v.Major != 2 {
		t.Errorf("Major = %d, want 2", v.Major)
	}
}

// --- SetAppInfo tests ---

func TestSetAppInfo(t *testing.T) {
	resetState()
	SetAppInfo("myapp", "My application")

	a := App()
	if a.Name != "myapp" {
		t.Errorf("Name = %q, want %q", a.Name, "myapp")
	}
	if a.Description != "My application" {
		t.Errorf("Description = %q, want %q", a.Description, "My application")
	}
}

func TestSetAppInfo_Idempotent(t *testing.T) {
	resetState()
	SetAppInfo("first", "First app")
	SetAppInfo("second", "Second app")

	a := App()
	if a.Name != "first" {
		t.Errorf("Name = %q, want %q (should be idempotent)", a.Name, "first")
	}
	if a.Description != "First app" {
		t.Errorf("Description = %q, want %q", a.Description, "First app")
	}
}

// --- SetGitInfo tests ---

func TestSetGitInfo(t *testing.T) {
	resetState()
	SetGitInfo("abc123", "main", "github.com/user/repo")

	g := Git()
	if g.Commit != "abc123" {
		t.Errorf("Commit = %q, want %q", g.Commit, "abc123")
	}
	if g.Branch != "main" {
		t.Errorf("Branch = %q, want %q", g.Branch, "main")
	}
	if g.Repo != "github.com/user/repo" {
		t.Errorf("Repo = %q, want %q", g.Repo, "github.com/user/repo")
	}
}

func TestSetGitInfo_Idempotent(t *testing.T) {
	resetState()
	SetGitInfo("abc123", "main", "repo1")
	SetGitInfo("def456", "dev", "repo2")

	g := Git()
	if g.Commit != "abc123" {
		t.Errorf("Commit = %q, want %q (should be idempotent)", g.Commit, "abc123")
	}
	if g.Branch != "main" {
		t.Errorf("Branch = %q, want %q", g.Branch, "main")
	}
	if g.Repo != "repo1" {
		t.Errorf("Repo = %q, want %q", g.Repo, "repo1")
	}
}

func TestSetGitInfo_AlsoInBuild(t *testing.T) {
	resetState()
	SetGitInfo("abc123", "main", "repo")

	b := Build()
	if b.Git.Commit != "abc123" {
		t.Errorf("Build().Git.Commit = %q, want %q", b.Git.Commit, "abc123")
	}
}

// --- SetBuildInfo tests ---

func TestSetBuildInfo_ValidTimestamp(t *testing.T) {
	resetState()
	ts := "Mon Jan  2 15:04:05 UTC 2006"
	SetBuildInfo(ts)

	b := Build()
	if b.Timestamp.IsZero() {
		t.Error("Timestamp should not be zero for valid UnixDate format")
	}
	if b.Timestamp.Year() != 2006 {
		t.Errorf("Year = %d, want 2006", b.Timestamp.Year())
	}
}

func TestSetBuildInfo_InvalidTimestamp(t *testing.T) {
	resetState()
	SetBuildInfo("not-a-timestamp")

	b := Build()
	if !b.Timestamp.IsZero() {
		t.Error("Timestamp should be zero for invalid format")
	}
}

func TestSetBuildInfo_Empty(t *testing.T) {
	resetState()
	SetBuildInfo("")

	b := Build()
	if !b.Timestamp.IsZero() {
		t.Error("Timestamp should be zero for empty string")
	}
}

func TestSetBuildInfo_Idempotent(t *testing.T) {
	resetState()
	SetBuildInfo("Mon Jan  2 15:04:05 UTC 2006")
	SetBuildInfo("Tue Jan  3 15:04:05 UTC 2006")

	b := Build()
	if b.Timestamp.Day() != 2 {
		t.Errorf("Day = %d, want 2 (should be idempotent)", b.Timestamp.Day())
	}
}

// --- SetChangelog tests ---

func TestSetChangelog(t *testing.T) {
	resetState()
	SetChangelog("## v1.0.0\n- Initial release")

	a := App()
	if a.Changelog != "## v1.0.0\n- Initial release" {
		t.Errorf("Changelog = %q, want changelog text", a.Changelog)
	}
}

func TestSetChangelog_Idempotent(t *testing.T) {
	resetState()
	SetChangelog("first")
	SetChangelog("second")

	a := App()
	if a.Changelog != "first" {
		t.Errorf("Changelog = %q, want %q (should be idempotent)", a.Changelog, "first")
	}
}

func TestSetChangelogFromFile(t *testing.T) {
	resetState()

	// Create a temp file with changelog content
	dir := t.TempDir()
	path := filepath.Join(dir, "CHANGELOG.md")
	content := "## v1.0.0\n- Feature A\n- Feature B\n"
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	if err := SetChangelogFromFile(path); err != nil {
		t.Fatalf("SetChangelogFromFile() error = %v", err)
	}

	a := App()
	if a.Changelog != content {
		t.Errorf("Changelog = %q, want %q", a.Changelog, content)
	}
}

func TestSetChangelogFromFile_MissingFile(t *testing.T) {
	resetState()

	err := SetChangelogFromFile("/nonexistent/path/CHANGELOG.md")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestSetChangelogFromFile_Idempotent(t *testing.T) {
	resetState()

	dir := t.TempDir()
	path1 := filepath.Join(dir, "CHANGELOG1.md")
	path2 := filepath.Join(dir, "CHANGELOG2.md")
	os.WriteFile(path1, []byte("first"), 0644)
	os.WriteFile(path2, []byte("second"), 0644)

	SetChangelogFromFile(path1)
	SetChangelogFromFile(path2)

	a := App()
	if a.Changelog != "first" {
		t.Errorf("Changelog = %q, want %q (should be idempotent)", a.Changelog, "first")
	}
}

// --- String() method tests ---

func TestVersionString(t *testing.T) {
	tests := []struct {
		name string
		ver  Version
		want string
	}{
		{
			name: "basic version",
			ver:  Version{Major: 1, Minor: 2, Patch: 3},
			want: " 1.2.3",
		},
		{
			name: "with suffix",
			ver:  Version{Major: 1, Minor: 0, Patch: 0, Prefix: "beta"},
			want: "beta 1.0.0",
		},
		{
			name: "zero version",
			ver:  Version{},
			want: " 0.0.0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.ver.String()
			if got != tt.want {
				t.Errorf("String() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestAppInfoString(t *testing.T) {
	a := AppInfo{Name: "myapp", Description: "My application"}
	got := a.String()
	if !strings.Contains(got, "myapp") {
		t.Errorf("String() = %q, should contain app name", got)
	}
	if !strings.Contains(got, "My application") {
		t.Errorf("String() = %q, should contain description", got)
	}
}

func TestBuildInfoString(t *testing.T) {
	ts := time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC)
	b := BuildInfo{
		Timestamp: ts,
		Git:       GitInfo{Commit: "abc", Branch: "main", Repo: "repo"},
	}
	got := b.String()
	if !strings.Contains(got, "Timestamp") {
		t.Errorf("String() = %q, should contain Timestamp label", got)
	}
	if !strings.Contains(got, "Git") {
		t.Errorf("String() = %q, should contain Git label", got)
	}
}

func TestGitInfoString(t *testing.T) {
	g := GitInfo{Commit: "abc123", Branch: "develop", Repo: "github.com/user/repo"}
	got := g.String()
	if !strings.Contains(got, "abc123") {
		t.Errorf("String() = %q, should contain commit", got)
	}
	if !strings.Contains(got, "develop") {
		t.Errorf("String() = %q, should contain branch", got)
	}
	if !strings.Contains(got, "github.com/user/repo") {
		t.Errorf("String() = %q, should contain repo", got)
	}
}

// --- Print test ---

func TestPrint(t *testing.T) {
	resetState()
	SetAppInfo("testapp", "Test application")
	SetVersion("1.2.3")
	SetGitInfo("abc", "main", "repo")

	// Capture stdout
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	Print()

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	if !strings.Contains(output, "Running:") {
		t.Errorf("Print() output should contain 'Running:', got %q", output)
	}
	if !strings.Contains(output, "Version:") {
		t.Errorf("Print() output should contain 'Version:', got %q", output)
	}
	if !strings.Contains(output, "Build:") {
		t.Errorf("Print() output should contain 'Build:', got %q", output)
	}
}

// --- LoadFromFile tests ---

func TestLoadFromFile_AllKeys(t *testing.T) {
	resetState()

	dir := t.TempDir()
	path := filepath.Join(dir, ".version")
	content := `VERSION=2.3.4
GIT_COMMIT=deadbeef
GIT_BRANCH=release
GIT_REPO=github.com/test/repo
BUILD_TIMESTAMP=Mon Jan  2 15:04:05 UTC 2006
`
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	if err := LoadFromFile(path); err != nil {
		t.Fatalf("LoadFromFile() error = %v", err)
	}

	v := Get()
	if v.Raw != "2.3.4" {
		t.Errorf("Raw = %q, want %q", v.Raw, "2.3.4")
	}
	if v.Major != 2 || v.Minor != 3 || v.Patch != 4 {
		t.Errorf("version = %d.%d.%d, want 2.3.4", v.Major, v.Minor, v.Patch)
	}

	g := Git()
	if g.Commit != "deadbeef" {
		t.Errorf("Commit = %q, want %q", g.Commit, "deadbeef")
	}
	if g.Branch != "release" {
		t.Errorf("Branch = %q, want %q", g.Branch, "release")
	}
	if g.Repo != "github.com/test/repo" {
		t.Errorf("Repo = %q, want %q", g.Repo, "github.com/test/repo")
	}

	b := Build()
	if b.Timestamp.IsZero() {
		t.Error("Timestamp should not be zero")
	}
}

func TestLoadFromFile_WithComments(t *testing.T) {
	resetState()

	dir := t.TempDir()
	path := filepath.Join(dir, ".version")
	content := `# Generated by go-version
VERSION=1.0.0

# Git info
GIT_COMMIT=abc123
`
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	if err := LoadFromFile(path); err != nil {
		t.Fatalf("LoadFromFile() error = %v", err)
	}

	v := Get()
	if v.Raw != "1.0.0" {
		t.Errorf("Raw = %q, want %q", v.Raw, "1.0.0")
	}

	g := Git()
	if g.Commit != "abc123" {
		t.Errorf("Commit = %q, want %q", g.Commit, "abc123")
	}
}

func TestLoadFromFile_SpacesAroundEquals(t *testing.T) {
	resetState()

	dir := t.TempDir()
	path := filepath.Join(dir, ".version")
	content := `VERSION = 1.5.0
GIT_COMMIT = abc123
`
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	if err := LoadFromFile(path); err != nil {
		t.Fatalf("LoadFromFile() error = %v", err)
	}

	v := Get()
	if v.Raw != "1.5.0" {
		t.Errorf("Raw = %q, want %q", v.Raw, "1.5.0")
	}
}

func TestLoadFromFile_MissingFile(t *testing.T) {
	resetState()

	err := LoadFromFile("/nonexistent/path/.version")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestLoadFromFile_EmptyFile(t *testing.T) {
	resetState()

	dir := t.TempDir()
	path := filepath.Join(dir, ".version")
	os.WriteFile(path, []byte(""), 0644)

	if err := LoadFromFile(path); err != nil {
		t.Fatalf("LoadFromFile() error = %v", err)
	}

	v := Get()
	if v.Raw != "" {
		t.Errorf("Raw = %q, want empty for empty file", v.Raw)
	}
}

func TestLoadFromFile_IgnoresUnknownKeys(t *testing.T) {
	resetState()

	dir := t.TempDir()
	path := filepath.Join(dir, ".version")
	content := `VERSION=1.0.0
UNKNOWN_KEY=some_value
ANOTHER=thing
`
	os.WriteFile(path, []byte(content), 0644)

	if err := LoadFromFile(path); err != nil {
		t.Fatalf("LoadFromFile() error = %v", err)
	}

	v := Get()
	if v.Raw != "1.0.0" {
		t.Errorf("Raw = %q, want %q", v.Raw, "1.0.0")
	}
}

func TestLoadFromFile_IgnoresLinesWithoutEquals(t *testing.T) {
	resetState()

	dir := t.TempDir()
	path := filepath.Join(dir, ".version")
	content := `VERSION=1.0.0
this line has no equals sign
GIT_COMMIT=abc
`
	os.WriteFile(path, []byte(content), 0644)

	if err := LoadFromFile(path); err != nil {
		t.Fatalf("LoadFromFile() error = %v", err)
	}

	v := Get()
	if v.Raw != "1.0.0" {
		t.Errorf("Raw = %q, want %q", v.Raw, "1.0.0")
	}
	if Git().Commit != "abc" {
		t.Errorf("Commit = %q, want %q", Git().Commit, "abc")
	}
}

func TestLoadFromFile_DoesNotOverwriteExisting(t *testing.T) {
	resetState()

	// Set version first
	SetVersion("1.0.0")
	SetGitInfo("original", "main", "repo1")

	dir := t.TempDir()
	path := filepath.Join(dir, ".version")
	content := `VERSION=2.0.0
GIT_COMMIT=overwritten
GIT_BRANCH=dev
GIT_REPO=repo2
`
	os.WriteFile(path, []byte(content), 0644)

	if err := LoadFromFile(path); err != nil {
		t.Fatalf("LoadFromFile() error = %v", err)
	}

	// Version.Raw is set, so LoadFromFile should skip it
	v := Get()
	if v.Raw != "1.0.0" {
		t.Errorf("Raw = %q, want %q (should not overwrite)", v.Raw, "1.0.0")
	}

	// Git info was already set, should not overwrite
	g := Git()
	if g.Commit != "original" {
		t.Errorf("Commit = %q, want %q (should not overwrite)", g.Commit, "original")
	}
}

// --- LoadFromGit tests ---

func TestLoadFromGit(t *testing.T) {
	resetState()

	// This test runs actual git commands, so it requires being in a git repo.
	// The go-version repo itself is a git repo, so this should work.
	err := LoadFromGit()
	if err != nil {
		t.Skipf("skipping: not in a git repository (%v)", err)
	}

	g := Git()
	if g.Commit == "" {
		t.Error("expected non-empty commit from git")
	}
	if g.Branch == "" {
		t.Error("expected non-empty branch from git")
	}
}

func TestLoadFromGit_DoesNotOverwriteExisting(t *testing.T) {
	resetState()

	SetGitInfo("preset-commit", "preset-branch", "preset-repo")
	SetVersion("9.9.9")

	err := LoadFromGit()
	if err != nil {
		t.Skipf("skipping: not in a git repository (%v)", err)
	}

	g := Git()
	if g.Commit != "preset-commit" {
		t.Errorf("Commit = %q, want %q (should not overwrite)", g.Commit, "preset-commit")
	}
	if g.Branch != "preset-branch" {
		t.Errorf("Branch = %q, want %q", g.Branch, "preset-branch")
	}
}

// --- Getter tests ---

func TestGet_ReturnsVersion(t *testing.T) {
	resetState()
	SetVersion("3.2.1")

	v := Get()
	if v.Major != 3 || v.Minor != 2 || v.Patch != 1 {
		t.Errorf("Get() = %d.%d.%d, want 3.2.1", v.Major, v.Minor, v.Patch)
	}
}

func TestBuild_ReturnsTimestampAndGit(t *testing.T) {
	resetState()
	SetBuildInfo("Mon Jan  2 15:04:05 UTC 2006")
	SetGitInfo("abc", "main", "repo")

	b := Build()
	if b.Timestamp.IsZero() {
		t.Error("Timestamp should not be zero")
	}
	if b.Git.Commit != "abc" {
		t.Errorf("Git.Commit = %q, want %q", b.Git.Commit, "abc")
	}
}

func TestGit_ReturnsBuildGit(t *testing.T) {
	resetState()
	SetGitInfo("abc", "main", "repo")

	g := Git()
	b := Build()
	if g.Commit != b.Git.Commit {
		t.Errorf("Git() and Build().Git should return same data")
	}
}

func TestApp_ReturnsAppInfo(t *testing.T) {
	resetState()
	SetAppInfo("myapp", "desc")
	SetChangelog("log")

	a := App()
	if a.Name != "myapp" || a.Description != "desc" || a.Changelog != "log" {
		t.Errorf("App() = %+v, unexpected values", a)
	}
}

// --- Table-driven SetVersion test ---

func TestSetVersion_Table(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantRaw string
		major   int
		minor   int
		patch   int
		prefix  string
	}{
		{"semver", "1.2.3", "1.2.3", 1, 2, 3, ""},
		{"v-prefix", "v1.2.3", "v1.2.3", 1, 2, 3, ""},
		{"suffix", "1.2.3-alpha", "1.2.3-alpha", 1, 2, 3, "alpha"},
		{"complex suffix", "v0.1.0-dev.42", "v0.1.0-dev.42", 0, 1, 0, "dev.42"},
		{"zeros", "0.0.0", "0.0.0", 0, 0, 0, ""},
		{"large numbers", "100.200.300", "100.200.300", 100, 200, 300, ""},
		{"empty", "", "", 0, 0, 0, ""},
		{"one part", "1", "1", 0, 0, 0, ""},
		{"two parts", "1.2", "1.2", 0, 0, 0, ""},
		{"four parts", "1.2.3.4", "1.2.3.4", 1, 2, 3, ""},
		{"just v", "v", "v", 0, 0, 0, ""},
		{"v only prefix", "v1", "v1", 0, 0, 0, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resetState()
			SetVersion(tt.input)

			v := Get()
			if v.Raw != tt.wantRaw {
				t.Errorf("Raw = %q, want %q", v.Raw, tt.wantRaw)
			}
			if v.Major != tt.major {
				t.Errorf("Major = %d, want %d", v.Major, tt.major)
			}
			if v.Minor != tt.minor {
				t.Errorf("Minor = %d, want %d", v.Minor, tt.minor)
			}
			if v.Patch != tt.patch {
				t.Errorf("Patch = %d, want %d", v.Patch, tt.patch)
			}
			if v.Prefix != tt.prefix {
				t.Errorf("Prefix = %q, want %q", v.Prefix, tt.prefix)
			}
		})
	}
}

// --- Integration test: full workflow ---

func TestFullWorkflow(t *testing.T) {
	resetState()

	// Simulate a typical application startup
	SetAppInfo("myservice", "My microservice")
	SetVersion("v2.1.0-rc.1")
	SetGitInfo("abc123def", "release/2.1", "github.com/org/myservice")
	SetBuildInfo("Mon Jan  2 15:04:05 UTC 2006")
	SetChangelog("## 2.1.0-rc.1\n- Added feature X")

	v := Get()
	if v.Major != 2 || v.Minor != 1 || v.Patch != 0 {
		t.Errorf("version = %d.%d.%d, want 2.1.0", v.Major, v.Minor, v.Patch)
	}
	if v.Prefix != "rc.1" {
		t.Errorf("Prefix = %q, want %q", v.Prefix, "rc.1")
	}

	a := App()
	if a.Name != "myservice" {
		t.Errorf("Name = %q, want %q", a.Name, "myservice")
	}
	if !strings.Contains(a.Changelog, "feature X") {
		t.Errorf("Changelog should contain 'feature X'")
	}

	g := Git()
	if g.Commit != "abc123def" {
		t.Errorf("Commit = %q, want %q", g.Commit, "abc123def")
	}

	b := Build()
	if b.Timestamp.Year() != 2006 {
		t.Errorf("Year = %d, want 2006", b.Timestamp.Year())
	}
}

func TestFullWorkflow_LoadFromFile(t *testing.T) {
	resetState()

	dir := t.TempDir()
	path := filepath.Join(dir, ".version")

	// Simulate a file generated by the go-version CLI
	content := fmt.Sprintf(`# Generated by go-version
VERSION=v3.0.0
GIT_COMMIT=fedcba98
GIT_BRANCH=main
GIT_REPO=github.com/org/app
BUILD_TIMESTAMP=%s
`, time.Now().UTC().Format(time.UnixDate))

	os.WriteFile(path, []byte(content), 0644)

	SetAppInfo("app", "My app")
	if err := LoadFromFile(path); err != nil {
		t.Fatalf("LoadFromFile() error = %v", err)
	}

	v := Get()
	if v.Major != 3 || v.Minor != 0 || v.Patch != 0 {
		t.Errorf("version = %d.%d.%d, want 3.0.0", v.Major, v.Minor, v.Patch)
	}

	g := Git()
	if g.Commit != "fedcba98" {
		t.Errorf("Commit = %q, want %q", g.Commit, "fedcba98")
	}

	b := Build()
	if b.Timestamp.IsZero() {
		t.Error("Timestamp should not be zero")
	}
}
