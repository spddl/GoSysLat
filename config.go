package main

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"syscall"
)

type Config struct {
	ComPort     int  `json:"ComPort"`
	OSDX        int  `json:"OSD_x"`
	OSDY        int  `json:"OSD_y"`
	OSDHeight   int  `json:"OSD_height"`
	OSDWidth    int  `json:"OSD_width"`
	PushMode    bool `json:"Push_Mode"`
	Display     bool `json:"Display"`
	DetectLight bool `json:"Detect_Light"`
}

var (
	executable string
	c          Config
)

func init() {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	executable = filepath.Base(ex)
	c.Load()
}

func (c *Config) read() []byte {
	stream, _ := syscall.UTF16PtrFromString(executable + ":Stream")
	hStream, err := syscall.CreateFile(stream, syscall.GENERIC_READ, syscall.FILE_SHARE_READ, nil, syscall.OPEN_EXISTING, 0, 0)
	if err != nil { // config not found
		return []byte{}
	}
	defer syscall.CloseHandle(hStream)

	var done uint32 = 0
	buf := make([]byte, 128)
	err = syscall.ReadFile(hStream, buf, &done, nil)
	if err != nil {
		log.Println(err)
	}
	return buf[:done]
}

func (c *Config) write(data []byte) {
	stream, _ := syscall.UTF16PtrFromString(executable + ":Stream")
	syscall.DeleteFile(stream)
	hStream, err := syscall.CreateFile(stream, syscall.GENERIC_WRITE, syscall.FILE_SHARE_WRITE, nil, syscall.OPEN_ALWAYS, 0, 0)
	if err != nil {
		log.Println(err)
	}
	defer syscall.CloseHandle(hStream)

	var done uint32
	err = syscall.WriteFile(hStream, data, &done, nil)
	if err != nil {
		log.Println(err)
	}
}

func (c *Config) Load() {
	// more < GoSysLat.exe:Stream:$DATA
	data := c.read()
	if len(data) != 0 {
		err := json.Unmarshal(data, c)
		if err != nil {
			log.Println(err)
		}
	} else {
		c.Reset()
	}

	if cliMode {
		log.Printf("f %#v\n", f)
		if f.RTSSOSDX != 0 {
			c.OSDX = f.RTSSOSDX
		}
		if f.RTSSOSDY != 0 {
			c.OSDY = f.RTSSOSDY
		}
		if f.RTSSOSDHeight != 0 {
			c.OSDHeight = f.RTSSOSDHeight
		}
		if f.RTSSOSDWidth != 0 {
			c.OSDWidth = f.RTSSOSDWidth
		}
		log.Printf("c %#v\n", c)
	}
}

func (c *Config) Save() {
	data, err := json.Marshal(c)
	if err != nil {
		log.Println(err)
	}
	c.write(data)
}

func (c *Config) Reset() {
	*c = Config{
		// default settings
		ComPort:     1,
		Display:     true,
		DetectLight: true,
		OSDHeight:   10,
		OSDWidth:    10,
	}
}

func (c *Config) RefreshOSDStrings() {
	RTSSOSDBlack = "<P=" + strconv.Itoa(c.OSDX) + "," + strconv.Itoa(c.OSDY) + "><C=000000><B=" + strconv.Itoa(c.OSDWidth) + "," + strconv.Itoa(c.OSDHeight) + "><C><P=0,0>"
	RTSSOSDWhite = "<P=" + strconv.Itoa(c.OSDX) + "," + strconv.Itoa(c.OSDY) + "><C=FFFFFF><B=" + strconv.Itoa(c.OSDWidth) + "," + strconv.Itoa(c.OSDHeight) + "><C><P=0,0>"
}
