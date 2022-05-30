package USBController

import (
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	// Library
	modSetupapi = windows.NewLazyDLL("setupapi.dll")

	// Functions
	procSetupDiGetClassDevsW = modSetupapi.NewProc("SetupDiGetClassDevsW")
)

type Device struct {
	Idata        DevInfoData
	FriendlyName string
	COMid        string
}

var GUID_DEVINTERFACE_COMPORT, _ = windows.GUIDFromString("{86E0D1E0-8089-11D0-9CE4-08003E301F73}")
var GUID_DEVINTERFACE_MOUSE, _ = windows.GUIDFromString("{378DE44C-56EF-11D1-BC8C-00A0C91405DD}")

func FindDevices(classGuid windows.GUID) ([]Device, DevInfo) {
	var allDevices []Device
	handle, err := SetupDiGetClassDevs(&classGuid, nil, 0, uint32(DIGCF_PRESENT|DIGCF_DEVICEINTERFACE))

	if err != nil {
		panic(err)
	}

	var index = 0
	for {
		idata, err := SetupDiEnumDeviceInfo(handle, index)
		if err != nil { // ERROR_NO_MORE_ITEMS
			break
		}
		index++

		dev := Device{
			Idata: *idata,
		}

		val, err := SetupDiGetDeviceRegistryProperty(handle, idata, SPDRP_FRIENDLYNAME)
		if err == nil {
			dev.FriendlyName = val.(string)
		}

		allDevices = append(allDevices, dev)
	}
	return allDevices, handle
}

func SetupDiGetClassDevs(classGuid *windows.GUID, enumerator *uint16, hwndParent uintptr, flags uint32) (handle DevInfo, err error) {
	r0, _, e1 := syscall.Syscall6(procSetupDiGetClassDevsW.Addr(), 4, uintptr(unsafe.Pointer(classGuid)), uintptr(unsafe.Pointer(enumerator)), uintptr(hwndParent), uintptr(flags), 0, 0)
	handle = DevInfo(r0)
	if handle == DevInfo(windows.InvalidHandle) {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func GetDeviceProperty(dis DevInfo, devInfoData *DevInfoData, devPropKey DEVPROPKEY) ([]byte, error) {
	var propt, size uint32
	buf := make([]byte, 16)
	run := true
	for run {
		err := SetupDiGetDeviceProperty(dis, devInfoData, &devPropKey, &propt, &buf[0], uint32(len(buf)), &size, 0)
		switch {
		case size > uint32(len(buf)):
			buf = make([]byte, size+16)
		case err != nil:
			return buf, err
		default:
			run = false
		}
	}

	return buf, nil
}
