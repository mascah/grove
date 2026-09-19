//go:build windows

package create

import "errors"

// D-003 defers Windows locking (LockFileEx) until it is required.
func lock(string) (func(), error) {
	return nil, errors.New("record creation is not supported on Windows yet")
}
