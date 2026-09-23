package tui

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
)

// TestTerminal runs testdata/terminal.py against a freshly built grove. Only
// a real pseudo-terminal shows whether raw mode and the alternate screen are
// restored, the result stays off the interface's stream, and a blocked Git
// child dies with the session. Under -race the binary is built with it too.
func TestTerminal(t *testing.T) {
	t.Parallel()
	if testing.Short() {
		t.Skip("builds the binary and drives nine pseudo-terminal sessions")
	}
	python, err := exec.LookPath("python3")
	if err != nil || runtime.GOOS == "windows" {
		t.Skip("needs python3 with the Unix pty and termios modules")
	}
	grove := filepath.Join(t.TempDir(), "grove")
	build := []string{"build", "-o", grove}
	if raceEnabled {
		build = append(build, "-race")
	}
	if out, err := exec.Command("go", append(build, "../../cmd/grove")...).CombinedOutput(); err != nil {
		t.Fatalf("go build: %v\n%s", err, out)
	}
	out, err := exec.Command(python, "testdata/terminal.py", grove).CombinedOutput()
	if t.Log(string(out)); err != nil {
		t.Fatalf("terminal checks failed: %v", err)
	}
}

// A screen that stops accepting output ends the session once, with the cause.
func TestWatchedScreenStopsTheSession(t *testing.T) {
	t.Parallel()
	file, err := os.Create(filepath.Join(t.TempDir(), "screen"))
	if err != nil {
		t.Fatal(err)
	}
	stops := 0
	w := &watched{File: file, stop: func() { stops++ }}
	if _, err := w.Write([]byte("drawn")); err != nil || stops != 0 || w.err != nil {
		t.Fatalf("a good write: %v, %d stops", err, stops)
	}
	file.Close()
	for range 2 {
		if _, err := w.Write([]byte("lost")); err == nil {
			t.Fatal("writing to a closed screen should fail")
		}
	}
	if stops != 1 || !errors.Is(w.err, os.ErrClosed) {
		t.Fatalf("stopped %d times with %v", stops, w.err)
	}
}
