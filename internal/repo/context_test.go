package repo

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestGitContextCancelledAndOrdinaryErrors(t *testing.T) {
	t.Parallel()
	plain := t.TempDir()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := GitContext(ctx, plain, "rev-parse", "--git-dir"); err != context.Canceled {
		t.Fatalf("a cancelled context must be returned as it is: %v", err)
	}
	if _, _, err := LocateContext(ctx, plain); err != context.Canceled {
		t.Fatalf("cancellation is not a missing repository: %v", err)
	}
	// Ordinary failures keep their text, with or without a context.
	_, err := Git(plain, "rev-parse", "--git-dir")
	_, ctxErr := GitContext(context.Background(), plain, "rev-parse", "--git-dir")
	if err == nil || !strings.HasPrefix(err.Error(), "git rev-parse: ") || !strings.Contains(err.Error(), "not a git repository") || ctxErr == nil || ctxErr.Error() != err.Error() {
		t.Fatalf("ordinary Git error text changed: %v / %v", err, ctxErr)
	}
	if _, _, err := Locate(plain); err == nil || !strings.HasPrefix(err.Error(), "this command requires a Git repository; coordination state lives in its common directory (git rev-parse: ") || errors.Is(err, context.Canceled) {
		t.Fatalf("ordinary Locate error text changed: %v", err)
	}
}

// The commands that predate cancellation wait for Git and whatever it started.
// A helper that holds Git's output open for longer than the kill delay must
// not turn a successful command into an error with nothing after the colon.
func TestUncancellableGitWaitsForHeldPipes(t *testing.T) {
	real, err := exec.LookPath("git")
	if err != nil {
		t.Skip("needs git")
	}
	shim := t.TempDir()
	waitDelay = 50 * time.Millisecond
	t.Cleanup(func() { waitDelay = 2 * time.Second })
	script := "#!/bin/sh\nsleep 0.3 &\nexec \"" + real + "\" \"$@\"\n"
	if err := os.WriteFile(filepath.Join(shim, "git"), []byte(script), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", shim+string(os.PathListSeparator)+os.Getenv("PATH"))
	if WaitDelay(context.Background()) != 0 || WaitDelay(t.Context()) == 0 {
		t.Fatal("only a cancellable context bounds the wait")
	}
	if out, err := Git(t.TempDir(), "version"); err != nil || !strings.HasPrefix(out, "git version") {
		t.Fatalf("out=%q err=%v", out, err)
	}
	// When the bounded wait does expire, the reason is stated.
	if _, err := GitContext(t.Context(), t.TempDir(), "version"); err == nil || strings.HasSuffix(err.Error(), ": ") {
		t.Fatalf("a cancellable read should report the expired wait: %v", err)
	}
}
