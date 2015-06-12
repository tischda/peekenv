// +build windows

package main

import (
	"syscall"
	"unsafe"
)

// Code Copyright 2015 The Go Authors extracted from:
// https://github.com/golang/sys/blob/master/windows/registry/zsyscall_windows.go

var (
	// Advanced Services (advapi32.dll) provide access to the Windows registry
	modadvapi32       = syscall.NewLazyDLL("advapi32.dll")
	procRegEnumValueW = modadvapi32.NewProc("RegEnumValueW")
)

// https://msdn.microsoft.com/en-us/library/windows/desktop/ms724865(v=vs.85).aspx
func regEnumValue(key syscall.Handle, index uint32, name *uint16, nameLen *uint32, reserved *uint32, valtype *uint32, buf *byte, buflen *uint32) (regerrno error) {
	r0, _, _ := syscall.Syscall9(procRegEnumValueW.Addr(), 8, uintptr(key), uintptr(index), uintptr(unsafe.Pointer(name)), uintptr(unsafe.Pointer(nameLen)), uintptr(unsafe.Pointer(reserved)), uintptr(unsafe.Pointer(valtype)), uintptr(unsafe.Pointer(buf)), uintptr(unsafe.Pointer(buflen)), 0)
	if r0 != 0 {
		regerrno = syscall.Errno(r0)
	}
	return
}
