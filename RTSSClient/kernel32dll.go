package RTSSClient

import (
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	// Library
	kernel32 = syscall.NewLazyDLL("kernel32.dll")

	// Functions
	openFileMapping = kernel32.NewProc("OpenFileMappingW")
)

func OpenFileMapping(dwDesiredAccess uint32, bInheritHandle bool, lpName string) (syscall.Handle, error) {
	namep, _ := windows.UTF16PtrFromString(lpName)
	var inheritHandle uint32
	if bInheritHandle {
		inheritHandle = 1
	}

	ret, _, err := openFileMapping.Call(uintptr(dwDesiredAccess), uintptr(inheritHandle), uintptr(unsafe.Pointer(namep)))

	return syscall.Handle(ret), err
}
