package RTSSClient

const (
	MAX_PATH = 260
)

type (
	DWORD         uint32
	LONG          int32
	LARGE_INTEGER int64
	LPBYTE        *byte
)

type RTSS_SHARED_MEMORY struct {
	DwSignature DWORD
	// signature allows applications to verify status of shared memory

	// The signature can be set to:
	// 'RTSS'	- statistics server's memory is initialized and contains
	//			valid data
	// 0xDEAD	- statistics server's memory is marked for deallocation and
	//			no longer contain valid data
	// otherwise	the memory is not initialized
	DwVersion DWORD
	// structure version ((major<<16) + minor)
	// must be set to 0x0002xxxx for v2.x structure

	DwAppEntrySize DWORD
	// size of RTSS_SHARED_MEMORY_OSD_ENTRY for compatibility with future versions
	DwAppArrOffset DWORD
	// offset of arrOSD array for compatibility with future versions
	DwAppArrSize DWORD
	// size of arrOSD array for compatibility with future versions

	DwOSDEntrySize DWORD
	// size of RTSS_SHARED_MEMORY_APP_ENTRY for compatibility with future versions
	DwOSDArrOffset DWORD
	// offset of arrApp array for compatibility with future versions
	DwOSDArrSize DWORD
	// size of arrOSD array for compatibility with future versions

	DwOSDFrame DWORD
	// Global OSD frame ID. Increment it to force the server to update OSD for all currently active 3D
	// applications.

	// next fields are valid for v2.14 and newer shared memory format only

	DwBusy LONG // long? int32 or int64
	// set bit 0 when you're writing to shared memory and reset it when done

	// WARNING: do not forget to reset it, otherwise you'll completely lock OSD updates for all clients

	// next fields are valid for v2.15 and newer shared memory format only

	DwDesktopVideoCaptureFlags DWORD
	dwDesktopVideoCaptureStat  DWORD // DWORD dwDesktopVideoCaptureStat[5];
	// shared copy of desktop video capture flags and performance stats for 64-bit applications

	// next fields are valid for v2.16 and newer shared memory format only

	DwLastForegroundApp DWORD
	// last foreground application entry index
	DwLastForegroundAppProcessID DWORD
	// last foreground application process ID

	// OSD slot descriptor structure

	// WARNING: next fields should never (!!!) be accessed directly, use the offsets to access them in order to provide
	// compatibility with future versions

	arrOSD [8]RTSS_SHARED_MEMORY_OSD_ENTRY
	// array of OSD slots
	arrApp [256]RTSS_SHARED_MEMORY_APP_ENTRY
	// array of application descriptors

	// next fields are valid for v2.9 and newer shared memory format only

	// WARNING: due to design flaw there is no offset available for this field, so it must be calculated manually as
	// dwAppArrOffset + dwAppArrSize * dwAppEntrySize

	// VIDEO_CAPTURE_PARAM autoVideoCaptureParam;
}

type RTSS_SHARED_MEMORY_OSD_ENTRY struct {
	SzOSD [256]byte // 	char	szOSD[256];
	// OSD slot text
	SzOSDOwner [256]byte // char	szOSDOwner[256];
	// OSD slot owner ID

	// next fields are valid for v2.7 and newer shared memory format only

	SzOSDEx [4096]byte // char	szOSDEx[4096];
	// extended OSD slot text

	// next fields are valid for v2.12 and newer shared memory format only

	Buffer [262144]byte // BYTE	buffer[262144];
	// OSD slot data buffer
}

type RTSS_SHARED_MEMORY_APP_ENTRY struct {
	// application identification related fields

	DwProcessID DWORD
	// process ID
	SzName [MAX_PATH]byte
	// process executable name
	DwFlags DWORD
	// application specific flags

	// instantaneous framerate related fields

	DwTime0 DWORD
	// start time of framerate measurement period (in milliseconds)

	// Take a note that this field must contain non-zero value to calculate
	// framerate properly!
	DwTime1 DWORD
	// end time of framerate measurement period (in milliseconds)
	DwFrames DWORD
	// amount of frames rendered during (dwTime1 - dwTime0) period
	DwFrameTime DWORD
	// frame time (in microseconds)

	// to calculate framerate use the following formulas:

	// 1000.0f * dwFrames / (dwTime1 - dwTime0) for framerate calculated once per second
	// or
	// 1000000.0f / dwFrameTime for framerate calculated once per frame

	// framerate statistics related fields

	DwStatFlags DWORD
	// bitmask containing combination of STATFLAG_... flags
	DwStatTime0 DWORD
	// statistics record period start time
	DwStatTime1 DWORD
	// statistics record period end time
	DwStatFrames DWORD
	// total amount of frames rendered during statistics record period
	DwStatCount DWORD
	// amount of min/avg/max measurements during statistics record period
	DwStatFramerateMin DWORD
	// minimum instantaneous framerate measured during statistics record period
	DwStatFramerateAvg DWORD
	// average instantaneous framerate measured during statistics record period
	DwStatFramerateMax DWORD
	// maximum instantaneous framerate measured during statistics record period

	// OSD related fields

	DwOSDX DWORD
	// OSD X-coordinate (coordinate wrapping is allowed, i.e. -5 defines 5
	// pixel offset from the right side of the screen)
	DwOSDY DWORD
	// OSD Y-coordinate (coordinate wrapping is allowed, i.e. -5 defines 5
	// pixel offset from the bottom side of the screen)
	DwOSDPixel DWORD
	// OSD pixel zooming ratio
	DwOSDColor DWORD
	// OSD color in RGB format
	DwOSDFrame DWORD
	// application specific OSD frame ID. Don't change it directly!

	DwScreenCaptureFlags DWORD
	SzScreenCapturePath  [MAX_PATH]byte // szScreenCapturePath[MAX_PATH] string

	// next fields are valid for v2.1 and newer shared memory format only

	DwOSDBgndColor DWORD
	// OSD background color in RGB format

	// next fields are valid for v2.2 and newer shared memory format only

	DwVideoCaptureFlags   DWORD
	SzVideoCapturePath    [MAX_PATH]byte // szVideoCapturePath[MAX_PATH] string
	DwVideoFramerate      DWORD
	DwVideoFramesize      DWORD
	DwVideoFormat         DWORD
	DwVideoQuality        DWORD
	DwVideoCaptureThreads DWORD

	DwScreenCaptureQuality DWORD
	DwScreenCaptureThreads DWORD

	// next fields are valid for v2.3 and newer shared memory format only

	DwAudioCaptureFlags DWORD

	// next fields are valid for v2.4 and newer shared memory format only

	DwVideoCaptureFlagsEx DWORD

	// next fields are valid for v2.5 and newer shared memory format only

	DwAudioCaptureFlags2 DWORD

	DwStatFrameTimeMin   DWORD
	DwStatFrameTimeAvg   DWORD
	DwStatFrameTimeMax   DWORD
	DwStatFrameTimeCount DWORD

	DwStatFrameTimeBuf          DWORD // dwStatFrameTimeBuf[1024] DWORD
	DwStatFrameTimeBufPos       DWORD
	DwStatFrameTimeBufFramerate DWORD

	// next fields are valid for v2.6 and newer shared memory format only

	QwAudioCapturePTTEventPush    LARGE_INTEGER
	QwAudioCapturePTTEventRelease LARGE_INTEGER

	QwAudioCapturePTTEventPush2    LARGE_INTEGER
	QwAudioCapturePTTEventRelease2 LARGE_INTEGER

	// next fields are valid for v2.8 and newer shared memory format only

	DwPrerecordSizeLimit DWORD
	DwPrerecordTimeLimit DWORD

	// next fields are valid for v2.13 and newer shared memory format only

	QwStatTotalTime                LARGE_INTEGER
	DwStatFrameTimeLowBuf          [1024]DWORD // dwStatFrameTimeLowBuf[1024] DWORD
	DwStatFramerate1Dot0PercentLow DWORD
	DwStatFramerate0Dot1PercentLow DWORD
}
