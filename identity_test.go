package grove

import (
	"bytes"
	"context"
	"errors"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"runtime/debug"
	"slices"
	"strings"
	"testing"

	"github.com/mascah/grove/internal/repo"
)

func TestIdentify(t *testing.T) {
	vcs := func(version string, settings ...string) *debug.BuildInfo {
		info := &debug.BuildInfo{Main: debug.Module{Version: version}}
		for i := 0; i < len(settings); i += 2 {
			info.Settings = append(info.Settings, debug.BuildSetting{Key: settings[i], Value: settings[i+1]})
		}
		return info
	}
	const rev = "0123456789abcdef0123456789abcdef01234567"
	for _, c := range []struct {
		name             string
		info             *debug.BuildInfo
		ok               bool
		version, stamped string
		want             string
	}{
		{"release from an archive", vcs("(devel)"), true, "v0.1.0", rev, "grove v0.1.0 (" + rev + ")"},
		{"release from a clean clone", vcs("(devel)", "vcs.revision", rev, "vcs.modified", "false"), true, "v0.1.0", rev, "grove v0.1.0 (" + rev + ")"},
		{"release whose stamp is not Go's checkout", vcs("(devel)", "vcs.revision", "fedcba", "vcs.modified", "true"), true, "v0.1.0", rev, "grove v0.1.0 (" + rev + ", vcs fedcba, modified)"},
		{"release from a dirty tree", vcs("(devel)", "vcs.revision", rev, "vcs.modified", "true"), true, "v0.1.0", rev, "grove v0.1.0 (" + rev + ", modified)"},
		{"release without its revision", vcs("(devel)"), true, "v0.1.0", "", "grove v0.1.0 (revision unknown)"},
		{"development build", vcs("(devel)", "vcs.revision", rev, "vcs.modified", "false"), true, "", "", "grove (devel) (" + rev + ")"},
		{"modified development build", vcs("(devel)", "vcs.revision", rev, "vcs.modified", "true"), true, "", "", "grove (devel) (" + rev + ", modified)"},
		{"archive without a stamp", vcs("(devel)"), true, "", "", "grove (devel) (revision unknown)"},
		{"go install of a pseudo-version", vcs("v0.0.0-20260925210455-47852e36169d"), true, "", "", "grove v0.0.0-20260925210455-47852e36169d (revision unknown)"},
		{"no build information", nil, false, "", "", "grove (version unknown) (revision unknown)"},
	} {
		got := identify(c.info, c.ok, c.version, c.stamped).String()
		if !strings.HasPrefix(got, c.want+" guides sha256:") {
			t.Errorf("%s: %q, want the prefix %q", c.name, got, c.want)
		}
	}
}

// A shipped template moves the content digest and leaves the guides digest,
// which covers only the guides and the record model. Neither depends on the
// build. The test changes package variables, so it runs alone.
func TestIdentityDigestsCoverTheirContent(t *testing.T) {
	base := identify(nil, false, "", "")
	if now := Identity(); base.Guides != now.Guides || base.Content != now.Content {
		t.Fatalf("the digests depend on the build: %+v and %+v", base, now)
	}
	for name, template := range map[string]*string{"reviewer": &Reviewer, "grove-work skill": &workLoad} {
		saved := *template
		*template += "\n"
		b := identify(nil, false, "", "")
		*template = saved
		if b.Content == base.Content || b.Guides != base.Guides {
			t.Errorf("a %s change: %+v from %+v", name, b, base)
		}
	}
}

// TestStampedArchiveBuild builds a copy of the source without .git, as a
// source archive holds it, once with the release stamp and once without, and
// runs each outside any checkout. Skipped under -short: it builds twice.
func TestStampedArchiveBuild(t *testing.T) {
	t.Parallel()
	if testing.Short() {
		t.Skip("builds the binary twice from a copy of the source")
	}
	listed, err := repo.Git(".", "ls-files", "-z", "--cached", "--others", "--exclude-standard")
	if err != nil {
		t.Fatal(err)
	}
	src := t.TempDir()
	for _, name := range strings.Split(strings.TrimSuffix(listed, "\x00"), "\x00") {
		data, err := os.ReadFile(name)
		if errors.Is(err, fs.ErrNotExist) {
			continue // deleted in the working tree
		} else if err != nil {
			t.Fatal(err)
		}
		if err := os.MkdirAll(filepath.Join(src, filepath.Dir(name)), 0o755); err == nil {
			err = os.WriteFile(filepath.Join(src, name), data, 0o644)
		}
		if err != nil {
			t.Fatal(err)
		}
	}
	env := slices.DeleteFunc(os.Environ(), func(kv string) bool {
		return slices.ContainsFunc(repo.GitLocation, func(v string) bool { return strings.HasPrefix(kv, v+"=") })
	})
	run := func(dir, name string, args ...string) string {
		t.Helper()
		cmd := exec.Command(name, args...)
		cmd.Dir, cmd.Env = dir, env
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("%s %v: %v\n%s", name, args, err, out)
		}
		return string(out)
	}
	const rev = "0123456789abcdef0123456789abcdef01234567"
	bin := t.TempDir()
	stamped, unstamped := filepath.Join(bin, "grove"), filepath.Join(bin, "grove-devel")
	run(src, "go", "build", "-trimpath", "-ldflags", "-X github.com/mascah/grove.version=v9.9.9 -X github.com/mascah/grove.revision="+rev, "-o", stamped, "./cmd/grove")
	run(src, "go", "build", "-o", unstamped, "./cmd/grove")
	here := Identity()
	digests := " guides " + here.Guides + " content " + here.Content + "\n"
	elsewhere := t.TempDir()
	if got := run(elsewhere, stamped, "version"); got != "grove v9.9.9 ("+rev+")"+digests {
		t.Fatalf("stamped: %q", got)
	}
	if got := run(elsewhere, unstamped, "version"); got != "grove (devel) (revision unknown)"+digests {
		t.Fatalf("unstamped: %q", got)
	}
	if binary, err := os.ReadFile(stamped); err != nil || bytes.Contains(binary, []byte(src)) {
		t.Fatalf("the trimmed build names the path it was built in (%v)", err)
	}
	// The stamped build serves a project with no Grove source in reach.
	project := filepath.Join(elsewhere, "project")
	if out, err := repo.Command(context.Background(), elsewhere, "init", "-q", project).CombinedOutput(); err != nil {
		t.Fatalf("git init: %v\n%s", err, out)
	}
	run(project, stamped, "init")
	run(project, stamped, "check")
	if got := run(elsewhere, stamped, "guide", "work"); !strings.HasPrefix(got, "# Executing assigned Grove work\n") {
		t.Fatalf("guide work: %.80q", got)
	}
}
