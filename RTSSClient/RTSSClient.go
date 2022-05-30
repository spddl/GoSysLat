package RTSSClient

import (
	"log"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

const (
	slotOwnerOSD   = "GoSyslat______________"
	clientPriority = 1

	Black = "<P=10,10><C=000000><B=10,10><C><P=0,0>"
	White = "<P=10,10><C=FFFFFF><B=10,10><C><P=0,0>"
)

type unsafeSlice struct {
	Data unsafe.Pointer
	Len  int
	Cap  int
}

const (
	STANDARD_RIGHTS_REQUIRED = 0x000F0000
	SECTION_QUERY            = 0x0001
	SECTION_MAP_WRITE        = 0x0002
	SECTION_MAP_READ         = 0x0004
	SECTION_MAP_EXECUTE      = 0x0008
	SECTION_EXTEND_SIZE      = 0x0010

	SECTION_ALL_ACCESS  = STANDARD_RIGHTS_REQUIRED | SECTION_QUERY | SECTION_MAP_WRITE | SECTION_MAP_READ | SECTION_MAP_EXECUTE | SECTION_EXTEND_SIZE
	FILE_MAP_ALL_ACCESS = SECTION_ALL_ACCESS
)

func UpdateOSD(lpText string) {
	var bResult bool
	hMapFile, _ := OpenFileMapping(FILE_MAP_ALL_ACCESS, false, "RTSSSharedMemoryV2")
	if hMapFile == 0 {
		// when RTSS has not been started
		// log.Println("could not open RTSSSharedMemoryV2")
		// log.Println(err)
		return
	}
	defer syscall.CloseHandle(hMapFile)

	if hMapFile != 0 {
		pMapAddr, err := syscall.MapViewOfFile(hMapFile, FILE_MAP_ALL_ACCESS, 0, 0, 0)
		if pMapAddr == 0 || err != nil {
			log.Println(err)
			return
		}
		defer syscall.UnmapViewOfFile(pMapAddr)

		pMem := (*RTSS_SHARED_MEMORY)(unsafe.Pointer(pMapAddr))
		// log.Println("pMem", prettyPrint(pMem))
		if (pMem.DwSignature == 1381258067 /*'RTSS'*/) && (pMem.DwVersion >= 0x00020000) {
			for dwPass := 0; dwPass < 2; dwPass++ {
				// 1st pass : find previously captured OSD slot
				// 2nd pass : otherwise find the first unused OSD slot and capture it

				// If the caller is "SysLat" allow it to take over the first OSD slot
				var dwEntry DWORD
				if clientPriority == 0 {
					dwEntry = 0
				} else {
					dwEntry = 1
				}

				for ; dwEntry < pMem.DwOSDArrSize; dwEntry++ {
					// allow primary OSD clients (e.g. EVGA Precision / MSI Afterburner) to use the first slot exclusively, so third party
					// applications start scanning the slots from the second one - CHANGED THIS TO 0 SO I CAN BE PRIMARY BECAUSE I NEED THE CORNERS

					pEntry := (*RTSS_SHARED_MEMORY_OSD_ENTRY)(unsafe.Pointer(pMapAddr + uintptr(pMem.DwOSDArrOffset) + uintptr(dwEntry)*uintptr(pMem.DwOSDEntrySize)))
					if dwPass != 0 {
						if pEntry.SzOSDOwner != [256]byte{} {
							strcpy(&pEntry.SzOSDOwner, slotOwnerOSD, len(slotOwnerOSD))
						}
					}

					// remember that strcmp returns 0 if the strings match... so the following if statement basically says if the strings match
					pEntry_szOSDOwner := windows.ByteSliceToString(pEntry.SzOSDOwner[:])
					if pEntry_szOSDOwner == slotOwnerOSD {
						lpTextPtr, _ := windows.BytePtrFromString(lpText)
						if pMem.DwVersion >= 0x00020007 {
							// use extended text slot for v2.7 and higher shared memory, it allows displaying 4096 symbols
							// instead of 256 for regular text slot
							if pMem.DwVersion >= 0x0002000e {
								// OSD locking is supported on v2.14 and higher shared memory

								DwBusy := (*LONG)(unsafe.Pointer(&pMem.DwBusy))
								if *DwBusy != 0 {
									// bit 0 of this variable will be set if OSD is locked by renderer and cannot be refreshed
									// at the moment
									*DwBusy = 0
								} else {
									// maybe strncpy_s is better, but it does not seem to be necessary until now
									strncpy(&pEntry.SzOSDEx[0], lpTextPtr, int32(unsafe.Sizeof(pEntry.SzOSDEx)))
								}

								// DWORD dwBusy = _interlockedbittestandset(&pMem.dwBusy, 0); // also not necessary
								// https://cpp.hotexamples.com/de/examples/-/-/_interlockedbittestandset/cpp-_interlockedbittestandset-function-examples.html
								// https://docs.microsoft.com/en-us/windows/win32/api/winnt/nf-winnt-_interlockedbittestandset
							} else {
								strncpy(&pEntry.SzOSDEx[0], lpTextPtr, int32(unsafe.Sizeof(pEntry.SzOSDEx)))
							}

						} else {
							strncpy(&pEntry.SzOSDEx[0], lpTextPtr, int32(unsafe.Sizeof(pEntry.SzOSDEx)))
						}

						pMem.DwOSDFrame++
						bResult = true
						break
					}
				}
				if bResult {
					break
				}
			}
		}
	}
}

func BytePtrToString(p *byte) string {
	if p == nil {
		log.Println("p == nil")
		return ""
	}
	if *p == 0 {
		log.Println("*p == 0")
		return ""
	}

	// Find NUL terminator.
	n := 0
	for ptr := unsafe.Pointer(p); *(*byte)(ptr) != 0; n++ {
		ptr = unsafe.Pointer(uintptr(ptr) + 1)
	}

	n += 100
	log.Println("n", n)

	var s []byte
	h := (*unsafeSlice)(unsafe.Pointer(&s))
	h.Data = unsafe.Pointer(p)
	h.Len = n
	h.Cap = n

	return string(s)
}

// BytePtrToString converts byte pointer to a Go string.
func BytePtrToStringg(p *byte) string {
	a := (*[10000]uint8)(unsafe.Pointer(p))
	i := 0
	for a[i] != 0 {
		i++
	}
	return string(a[:i])
}

func strcpy(dest *[256]byte, src string, n int) {
	maxn := len(dest) - 1
	if n > maxn {
		n = maxn
	}

	for i := 0; i < n; i++ {
		dest[i] = src[i]
	}
}

func strncpy(dest, src *byte, len int32) *byte {
	// Copy up to the len or first NULL bytes - whichever comes first.
	var (
		pSrc  = src
		pDest = dest
		i     int32
	)
	for i < len && *pSrc != 0 {
		*pDest = *pSrc
		i++
		pSrc = (*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(src)) + uintptr(i)))
		pDest = (*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(dest)) + uintptr(i)))
	}

	// The rest of the dest will be padded with zeros to the len.
	for i < len {
		*pDest = 0
		i++
		pDest = (*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(dest)) + uintptr(i)))
	}

	return dest
}

func GetSharedMemoryVersion() DWORD {
	hMapFile, err := OpenFileMapping(FILE_MAP_ALL_ACCESS, false, "RTSSSharedMemoryV2")
	if hMapFile == 0 { // when RTSS has not been started
		log.Println("could not open RTSSSharedMemoryV2")
		log.Println(err)
		return DWORD(0)
	}
	defer syscall.CloseHandle(hMapFile)

	if hMapFile != 0 {
		pMapAddr, err := syscall.MapViewOfFile(hMapFile, FILE_MAP_ALL_ACCESS, 0, 0, 0)
		if pMapAddr == 0 || err != nil {
			log.Println(err)
			return DWORD(0)
		}
		defer syscall.UnmapViewOfFile(pMapAddr)

		pMem := (*RTSS_SHARED_MEMORY)(unsafe.Pointer(pMapAddr))
		if (pMem.DwSignature == 1381258067 /*'RTSS'*/) && (pMem.DwVersion >= 0x00020000) {
			return pMem.DwVersion
		} else {
			log.Printf("RTSS DwSignature %#v, DwVersion %#v\n", pMem.DwSignature, pMem.DwVersion)
		}
	}
	return DWORD(0)
}
