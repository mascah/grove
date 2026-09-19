package repo

import (
	"context"
	"errors"
	"strings"
	"testing"
)

func TestGitContextCancelledAndOrdinaryErrors(t *testing.T) {
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
