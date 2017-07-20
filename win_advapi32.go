// +build windows

package main

import (
	"syscall"
	"unsafe"
)

// Code inspired from:
// https://github.com/golang/sys/blob/master/windows/registry/zsyscall_windows.go

var (
	// Advanced Services (advapi32.dll) provide access to the Windows registry
	modadvapi32       = syscall.NewLazyDLL("advapi32.dll")
	procRegEnumValueW = modadvapi32.NewProc("RegEnumValueW")
)

// https://msdn.microsoft.com/en-us/library/windows/desktop/ms724865(v=vs.85).aspx
func regEnumValue(hKey syscall.Handle, dwIndex uint32, lpValueName *uint16, lpcchValueName *uint32, lpReserved *uint32, lpType *uint32, lpData *byte, lpcbData *uint32) (regerrno error) {
	ret, _, _ := procRegEnumValueW.Call(
		uintptr(hKey),
		uintptr(dwIndex),
		uintptr(unsafe.Pointer(lpValueName)),
		uintptr(unsafe.Pointer(lpcchValueName)),
		uintptr(unsafe.Pointer(lpReserved)),
		uintptr(unsafe.Pointer(lpType)),
		uintptr(unsafe.Pointer(lpData)),
		uintptr(unsafe.Pointer(lpcbData)))

	// If the function fails, the return value is a system error code
	if ret != 0 {
		regerrno = syscall.Errno(ret)
	}
	return
}
