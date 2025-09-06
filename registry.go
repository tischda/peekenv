//go:build windows
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
func (realRegistry) GetString(path regKey, valueName string) (value string, err error) {
	handle := openKey(path, syscall.KEY_QUERY_VALUE)
	defer syscall.RegCloseKey(handle)

	var typ uint32
	var bufSize uint32

	name, err := syscall.UTF16PtrFromString(valueName)
	if err != nil {
		return "", err
	}

	// First call: Get the required buffer size
	// Pass nil for data buffer to get size in bufSize
	err = syscall.RegQueryValueEx(
		handle,
		name,
		nil,
		&typ,
		nil,      // nil data buffer
		&bufSize) // receives required size

	if err != nil {
		return "", err
	}

	// Allocate buffer with the exact size needed
	// Add 1 to handle potential rounding for UTF16
	data := make([]uint16, bufSize/2+1)

	// Second call: Actually get the data with properly sized buffer
	err = syscall.RegQueryValueEx(
		handle,
		name,
		nil,
		&typ,
		(*byte)(unsafe.Pointer(&data[0])), // properly sized buffer
		&bufSize)

	if err != nil {
		return "", err
	}
	return syscall.UTF16ToString(data), nil
}

// Enumerates the values for the specified registry key index. The function
// returns an array of valueNames.
func (realRegistry) EnumValues(path regKey) []string {
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
func getNextEnumValue(path regKey, index uint32) (string, error) {
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

// Opens a Windows registry key and returns a handle. You must close the
// handle with `defer syscall.RegCloseKey(handle)` in the calling code.
func openKey(path regKey, desiredAccess uint32) syscall.Handle {
	var handle syscall.Handle

	subkey, err := syscall.UTF16PtrFromString(path.lpSubKey)
	if err != nil {
		log.Fatalln("Error on registry path.subKey:", path.lpSubKey, err)
	}

	err = syscall.RegOpenKeyEx(
		hKeyTable[path.hKeyIdx],
		subkey,
		0,
		desiredAccess,
		&handle)

	if err != nil {
		log.Fatalln("Cannot open registry path:", path, err)
	}
	return handle
}
