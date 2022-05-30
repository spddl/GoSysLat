package USBController

import (
	"log"
	"strings"
	"syscall"
	"unicode"
)

const (
	EV_RXCHAR = 0x0001
	EV_ERR    = 0x0080
)

var COMDevices []Device

// var HIDDevices []Device

func init() {
	var comHandle DevInfo
	COMDevices, comHandle = FindDevices(GUID_DEVINTERFACE_COMPORT)
	SetupDiDestroyDeviceInfoList(comHandle)
	for i := range COMDevices {
		COMIndex := strings.Index(COMDevices[i].FriendlyName, "COM")
		COMDevices[i].COMid = ReturnNumber(COMDevices[i].FriendlyName[COMIndex+3:])
	}

	// var hidHandle DevInfo
	// HIDDevices, hidHandle = FindDevices(GUID_DEVINTERFACE_MOUSE)
	// SetupDiDestroyDeviceInfoList(hidHandle)
	// for i := range HIDDevices {
	// 	log.Printf("%#v\n", HIDDevices[i])
	// }
}

func ReturnNumber(s string) string {
	r := []rune{}
	for _, c := range s {
		if unicode.IsDigit(c) {
			r = append(r, c)
		} else {
			return string(r)
		}
	}
	return string(r)
}

func OpenComPort(portSpecifier string) syscall.Handle {
	portSpec, _ := syscall.UTF16PtrFromString(portSpecifier)
	hPort, err := syscall.CreateFile(
		portSpec,
		syscall.GENERIC_READ|syscall.GENERIC_WRITE,
		syscall.FILE_SHARE_WRITE|syscall.FILE_SHARE_READ,
		nil,
		syscall.OPEN_EXISTING,
		0, 0)
	if err != nil {
		log.Printf("0x%X, %s", hPort, err)
	}

	var dcb c_DCB
	if err = GetCommState(hPort, &dcb); err != nil {
		log.Println(err)
		syscall.CloseHandle(hPort)
		return hPort // INVALID_HANDLE_VALUE
	}

	dcb.BaudRate = 9600

	dcb.ByteSize = 8
	dcb.Parity = 0   // NOPARITY
	dcb.StopBits = 0 // ONESTOPBIT
	if err = SetCommState(hPort, &dcb); err != nil {
		log.Println(err)
		syscall.CloseHandle(hPort)
		return hPort
	}

	// Read this carefully because timeouts are important
	// https://docs.microsoft.com/en-us/windows/win32/api/winbase/ns-winbase-commtimeouts
	var timeouts c_COMMTIMEOUTS // https://docs.microsoft.com/en-us/windows/win32/api/winbase/ns-winbase-commtimeouts
	err = GetCommTimeouts(hPort, &timeouts)
	if err != nil {
		log.Println(err)
	}

	timeouts.ReadIntervalTimeout = 100
	timeouts.ReadTotalTimeoutConstant = 100
	timeouts.WriteTotalTimeoutConstant = 100
	err = SetCommTimeouts(hPort, &timeouts)
	if err != nil {
		log.Println(err)
	}

	err = SetCommMask(hPort, EV_RXCHAR|EV_ERR) // receive character event
	if err != nil {
		log.Println(err)
		return hPort
	}

	return hPort
}

func CloseComPort(hPort syscall.Handle) {
	if hPort != 0 && hPort != 0xFFFFFFFFFFFFFFFF {
		if err := PurgeComm(hPort, PURGE_RXABORT); err != nil {
			log.Println(err)
			return
		}

		if err := syscall.CloseHandle(hPort); err != nil {
			log.Println(err)
		}
	}
}

func ReadByte(hPort *syscall.Handle) (byte, bool) {
	var nNumberOfBytesToRead uint32 = 0
	buf := make([]byte, 1)
	if err := syscall.ReadFile(*hPort, buf, &nNumberOfBytesToRead, nil); err != nil {
		log.Println(err)
		panic(err)
	}
	return buf[0], nNumberOfBytesToRead != 0
	// return buf[:nNumberOfBytesToRead], nNumberOfBytesToRead != 0

}

func WriteByte(hPort *syscall.Handle, buf []byte) bool {
	// log.Println("WriteByte", buf[0])
	var done uint32
	err := syscall.WriteFile(*hPort, buf, &done, nil)
	if err != nil {
		log.Println(err)
	}
	return err != nil
}
