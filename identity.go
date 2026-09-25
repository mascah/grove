package grove

import (
	"crypto/sha256"
	"fmt"
	"io/fs"
	"maps"
	"runtime/debug"
	"slices"
)

// version and commit are a release build's stamp, set by the linker:
//
//	-ldflags "-X github.com/mascah/grove.version=v0.1.0 -X github.com/mascah/grove.commit=COMMIT"
//
// An unstamped build falls back to Go's build information.
var version, commit string

// Build names one executable and the workflow content it ships. Version and
// Commit are "" when neither a stamp nor Go's build information has them.
type Build struct {
	Version  string // the stamped release, else the main module's version
	Commit   string // the stamped commit, else vcs.revision
	VCS      string // vcs.revision where it differs from the stamped commit
	Modified bool   // vcs.modified: Go's checkout had uncommitted changes
	Guides   string // sha256 over the work and shaping guides and the record model
	Content  string // sha256 over the guides, the reviewer and init's other entrypoint templates
}

// Identity is this executable's Build, the one account of it that version
// and an attempt's launch record both print.
func Identity() Build {
	info, ok := debug.ReadBuildInfo()
	return identify(info, ok, version, commit)
}

func identify(info *debug.BuildInfo, ok bool, stampedVersion, stampedCommit string) Build {
	b := Build{Version: stampedVersion, Commit: stampedCommit}
	if ok {
		if b.Version == "" {
			b.Version = info.Main.Version
		}
		for _, s := range info.Settings {
			switch {
			case s.Key == "vcs.revision" && b.Commit == "":
				b.Commit = s.Value
			case s.Key == "vcs.revision" && s.Value != b.Commit:
				b.VCS = s.Value
			case s.Key == "vcs.modified":
				b.Modified = s.Value == "true"
			}
		}
	}
	guides := sha256.New()
	content := sha256.New()
	shipped := Entrypoints() // the reviewer among them
	for _, name := range []string{"work", "shape", "model"} {
		source, err := fs.ReadFile(Guides, GuideFiles[name])
		if err != nil {
			panic(err) // every file is embedded
		}
		guides.Write(source)
		shipped[GuideFiles[name]] = string(source)
	}
	for _, name := range slices.Sorted(maps.Keys(shipped)) {
		fmt.Fprintf(content, "%s\x00%d\x00%s", name, len(shipped[name]), shipped[name])
	}
	b.Guides = fmt.Sprintf("sha256:%x", guides.Sum(nil)[:6])
	b.Content = fmt.Sprintf("sha256:%x", content.Sum(nil)[:6])
	return b
}

// String is the version line: a missing version or commit is said to be
// unknown, so an unstamped build never reads as a known release.
func (b Build) String() string {
	line := "grove " + b.Version
	if b.Version == "" {
		line = "grove (version unknown)"
	}
	switch {
	case b.Commit == "":
		line += " (commit unknown"
	default:
		line += " (" + b.Commit
	}
	if b.VCS != "" {
		line += ", vcs " + b.VCS
	}
	if b.Modified {
		line += ", modified"
	}
	return line + ") guides " + b.Guides + " content " + b.Content
}
