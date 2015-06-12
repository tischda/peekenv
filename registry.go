// +build windows

package main

import (
	"log"
	"syscall"
	"unsafe"
)

type realRegistry struct{}

// do not reorder
var hKeyTable = []syscall.Handle{
	syscall.HKEY_CLASSES_ROOT,
	syscall.HKEY_CURRENT_USER,
	syscall.HKEY_LOCAL_MACHINE,
	syscall.HKEY_USERS,
	syscall.HKEY_PERFORMANCE_DATA,
	syscall.HKEY_CURRENT_CONFIG,
	syscall.HKEY_DYN_DATA,
}

// Read string from Windows registry (no expansion).
// Thanks to http://npf.io/2012/11/go-win-stuff/
func (realRegistry) GetString(path regPath, valueName string) (value string, err error) {
	handle := openKey(path, syscall.KEY_QUERY_VALUE)
	defer syscall.RegCloseKey(handle)

	var typ uint32
	var bufSize uint32

	// https://msdn.microsoft.com/en-us/library/windows/desktop/ms724911(v=vs.85).aspx
	err = syscall.RegQueryValueEx(
		handle,
		syscall.StringToUTF16Ptr(valueName),
		nil,
		&typ,
		nil,
		&bufSize)

	if err != nil {
		return
	}

	data := make([]uint16, bufSize/2+1)

	err = syscall.RegQueryValueEx(
		handle,
		syscall.StringToUTF16Ptr(valueName),
		nil,
		&typ,
		(*byte)(unsafe.Pointer(&data[0])),
		&bufSize)

	if err != nil {
		return
	}
	return syscall.UTF16ToString(data), nil
}

// Enumerates the values for the specified registry key index. The function
// returns an array of valueNames.
func (realRegistry) EnumValues(path regPath) []string {
	var values []string
	name, err := getNextEnumValue(path, uint32(0))
	for i := 1; err == nil; i++ {
		values = append(values, name)
		name, err = getNextEnumValue(path, uint32(i))
	}
	return values
}

// Enumerates the values for the specified registry key. The function
// returns one indexed value name for the key each time it is called.
func getNextEnumValue(path regPath, index uint32) (string, error) {
	handle := openKey(path, syscall.KEY_QUERY_VALUE)
	defer syscall.RegCloseKey(handle)

	var nameLen uint32 = 16383
	name := make([]uint16, nameLen)

	// https://msdn.microsoft.com/en-us/library/windows/desktop/ms724872(v=vs.85).aspx
	err := regEnumValue(
		handle,
		index,
		&name[0],
		&nameLen,
		nil,
		nil,
		nil,
		nil)

	return syscall.UTF16ToString(name), err
}

// Opens a Windows registry key and returns a handle. You must close
// the handle with `defer syscall.RegCloseKey(handle)` in the calling code.
func openKey(path regPath, desiredAccess uint32) syscall.Handle {
	var handle syscall.Handle

	// https://msdn.microsoft.com/en-us/library/windows/desktop/ms724897(v=vs.85).aspx
	err := syscall.RegOpenKeyEx(
		hKeyTable[path.hKeyIdx],
		syscall.StringToUTF16Ptr(path.lpSubKey),
		0,
		desiredAccess,
		&handle)

	if err != nil {
		log.Fatalln("Cannot open registry path:", path)
	}
	return handle
}
