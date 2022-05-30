package USBController

import (
	"syscall"
	"unsafe"
)

var (
	// Library
	modkernel32 = syscall.NewLazyDLL("kernel32.dll")

	// Functions
	procPurgeComm       = modkernel32.NewProc("PurgeComm")
	procSetCommMask     = modkernel32.NewProc("SetCommMask")
	procGetCommState    = modkernel32.NewProc("GetCommState")
	procSetCommState    = modkernel32.NewProc("SetCommState")
	procGetCommTimeouts = modkernel32.NewProc("GetCommTimeouts")
	procSetCommTimeouts = modkernel32.NewProc("SetCommTimeouts")
)

type c_COMMTIMEOUTS struct {
	ReadIntervalTimeout         uint32
	ReadTotalTimeoutMultiplier  uint32
	ReadTotalTimeoutConstant    uint32
	WriteTotalTimeoutMultiplier uint32
	WriteTotalTimeoutConstant   uint32
}

type purgeFlag int

const (
	PURGE_TXABORT purgeFlag = 0x01
	PURGE_RXABORT purgeFlag = 0x02
	PURGE_TXCLEAR purgeFlag = 0x04
	PURGE_RXCLEAR purgeFlag = 0x08
)

type c_DCB struct {
	DCBlength  uint32
	BaudRate   uint32
	Pad_cgo_0  [4]byte
	WReserved  uint16
	XonLim     uint16
	XoffLim    uint16
	ByteSize   uint8
	Parity     uint8
	StopBits   uint8
	XonChar    int8
	XoffChar   int8
	ErrorChar  int8
	EofChar    int8
	EvtChar    int8
	WReserved1 uint16
}

func PurgeComm(handle syscall.Handle, purge purgeFlag) error {
	// BOOL PurgeComm( HANDLE hFile, DWORD dwFlags )
	var err error
	r0, _, e1 := syscall.Syscall(procPurgeComm.Addr(), 2, uintptr(handle), uintptr(purge), 0)
	if r0 == 0 {
		if e1 != 0 {
			return error(e1)
		} else {
			return syscall.EINVAL
		}
	}
	return err
}

func GetCommState(handle syscall.Handle, dcb *c_DCB) (err error) {
	r1, _, e1 := syscall.Syscall(procGetCommState.Addr(), 2, uintptr(handle), uintptr(unsafe.Pointer(dcb)), 0)
	if r1 == 0 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func SetCommState(handle syscall.Handle, dcb *c_DCB) (err error) {
	r1, _, e1 := syscall.Syscall(procSetCommState.Addr(), 2, uintptr(handle), uintptr(unsafe.Pointer(dcb)), 0)
	if r1 == 0 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func SetCommMask(handle syscall.Handle, mask uint32) (err error) {
	r0, _, e1 := syscall.Syscall(procSetCommMask.Addr(), 2, uintptr(handle), uintptr(mask), 0)
	if r0 == 0 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = syscall.EINVAL
		}
	}

	return
}

func GetCommTimeouts(handle syscall.Handle, timeouts *c_COMMTIMEOUTS) (err error) {
	r1, _, e1 := syscall.Syscall(procGetCommTimeouts.Addr(), 2, uintptr(handle), uintptr(unsafe.Pointer(timeouts)), 0)
	if r1 == 0 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func SetCommTimeouts(handle syscall.Handle, timeouts *c_COMMTIMEOUTS) (err error) {
	r1, _, e1 := syscall.Syscall(procSetCommTimeouts.Addr(), 2, uintptr(handle), uintptr(unsafe.Pointer(timeouts)), 0)
	if r1 == 0 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}
