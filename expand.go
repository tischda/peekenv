package main

import (
	"log"
	"syscall"
	"unsafe"
)

var modkernel32 = syscall.NewLazyDLL("kernel32.dll")
var procExpandEnvironmentStringsW = modkernel32.NewProc("ExpandEnvironmentStringsW")

// expandVariable returns the resolved value of environment variables containing
// other variables such as %APPDATA% or %USERPROFILE%.
//
// Parameters:
//   - v: The environment variable to expand.
func expandVariable(v string) string {
	src, err := syscall.UTF16PtrFromString(v)
	if err != nil {
		log.Fatalln("String with NULL passed to StringToUTF16Ptr")
	}
	buf := make([]uint16, 32767) // Maximum environment variable size on Windows
	dst := &buf[0]
	size := uintptr(len(buf))

	n, _, _ := procExpandEnvironmentStringsW.Call(
		uintptr(unsafe.Pointer(src)),
		uintptr(unsafe.Pointer(dst)),
		size,
	)
	if n != 0 && n <= size {
		v = syscall.UTF16ToString(buf[:n-1])
	}
	return v
}
