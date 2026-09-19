package create

import (
	"os"
	"syscall"
)

// lock takes an exclusive advisory lock that the kernel releases when the
// holding process exits, so a crash never leaves a stale lock behind.
func lock(path string) (func(), error) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0o644)
	if err != nil {
		return nil, err
	}
	if err := syscall.Flock(int(f.Fd()), syscall.LOCK_EX); err != nil {
		f.Close()
		return nil, err
	}
	return func() { f.Close() }, nil
}
