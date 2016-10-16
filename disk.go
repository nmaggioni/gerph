//+build !windows

package main

import (
	"os"
	"syscall"
)

func GetAvailableDiskSpace() (uint64, error) {
	var stat syscall.Statfs_t
	cwd, err := os.Getwd()
	if err != nil {
		return 0, err
	}
	syscall.Statfs(cwd, &stat)
	return stat.Bavail * uint64(stat.Bsize), nil
}
