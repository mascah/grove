package create

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"
)

// The test re-executes itself as a helper that takes the lock and waits to be
// killed; the parent must then acquire the lock without any cleanup step.
func TestLockReleasedWhenHolderIsKilled(t *testing.T) {
	if path := os.Getenv("GROVE_LOCK_HELPER"); path != "" {
		if _, err := lock(path); err != nil {
			os.Exit(2)
		}
		fmt.Println("locked")
		time.Sleep(time.Minute)
		return
	}
	path := filepath.Join(t.TempDir(), "lock")
	helper := exec.Command(os.Args[0], "-test.run=^TestLockReleasedWhenHolderIsKilled$")
	helper.Env = append(os.Environ(), "GROVE_LOCK_HELPER="+path)
	stdout, err := helper.StdoutPipe()
	if err != nil {
		t.Fatal(err)
	}
	if err := helper.Start(); err != nil {
		t.Fatal(err)
	}
	if _, err := io.ReadFull(stdout, make([]byte, len("locked\n"))); err != nil {
		t.Fatalf("helper never reported the lock: %v", err)
	}
	done := make(chan error, 1)
	go func() {
		unlock, err := lock(path)
		if err == nil {
			unlock()
		}
		done <- err
	}()
	select {
	case <-done:
		t.Fatal("lock acquired while the helper still held it")
	case <-time.After(200 * time.Millisecond):
	}
	if err := helper.Process.Kill(); err != nil {
		t.Fatal(err)
	}
	helper.Wait()
	select {
	case err := <-done:
		if err != nil {
			t.Fatal(err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("lock still held after the holder was killed")
	}
}
