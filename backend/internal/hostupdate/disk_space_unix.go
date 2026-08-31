//go:build !windows

package hostupdate

import "syscall"

func availableDiskBytes(path string) (int64, error) {
	var disk syscall.Statfs_t
	if err := syscall.Statfs(path, &disk); err != nil {
		return 0, err
	}
	return int64(disk.Bavail) * int64(disk.Bsize), nil
}
