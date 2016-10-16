//+build windows

package main

import (
	"fmt"
	"os"
	"syscall"
	"unsafe"
)

// https://github.com/StalkR/goircbot/blob/ffd0c37e5d201730e4dac1ed1e013b315acab949/lib/disk/space_windows.go#L13
func GetAvailableDiskSpace() (uint64, error) {
	path, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		return 0, err
	}
	kernel32, err := syscall.LoadLibrary("Kernel32.dll")
	if err != nil {
		fmt.Println(err)
		return 0, err
	}
	defer syscall.FreeLibrary(kernel32)
	GetDiskFreeSpaceEx, err := syscall.GetProcAddress(syscall.Handle(kernel32), "GetDiskFreeSpaceExW")
	if err != nil {
		fmt.Println(err)
		return 0, err
	}
	var freeBytesAvailable int64
	var totalNumberOfBytes int64
	var totalNumberOfFreeBytes int64
	r1, _, e1 := syscall.Syscall6(uintptr(GetDiskFreeSpaceEx), 4,
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(path))),
		uintptr(unsafe.Pointer(&freeBytesAvailable)),
		uintptr(unsafe.Pointer(&totalNumberOfBytes)),
		uintptr(unsafe.Pointer(&totalNumberOfFreeBytes)), 0, 0)
	if r1 == 0 {
		if e1 != 0 {
			err = error(e1)
			return 0, err
		} else {
			err = syscall.EINVAL
			return 0, err
		}
	}
	return uint64(freeBytesAvailable), nil
}
