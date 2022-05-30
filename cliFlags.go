package main

import (
	"flag"
	"fmt"
	"log"
	"time"
)

var cliMode bool

type Flags struct {
	Name          string
	LogFolder     string
	Port          int
	Delay         time.Duration
	Count         int
	RTSS          bool
	RTSSOSDX      int
	RTSSOSDY      int
	RTSSOSDWidth  int
	RTSSOSDHeight int
	D3D9          bool
	OGL           bool
	Fullscreen    bool
	Width         int
	Height        int
	Print         bool
}

var f Flags

// go build; .\GoSysLat.exe -d3d9 -port 4 -count 10000
// go build; .\GoSysLat.exe -d3d9 -fullscreen -port 4 -time 30s -name Test
func init() {
	flag.StringVar(&f.Name, "name", "", "will be used later for the name")
	flag.StringVar(&f.LogFolder, "logs", "", "Log folder")
	flag.IntVar(&f.Port, "port", 0, "")
	flag.IntVar(&f.Count, "count", -1, "after which amount should be canceled")
	timePtr := flag.String("time", "", "time span: 5h30m40s")
	flag.BoolVar(&f.RTSS, "rtss", false, "RTSS enable support")
	flag.IntVar(&f.RTSSOSDX, "rtssosdx", 0, "RTSSOSDX")
	flag.IntVar(&f.RTSSOSDY, "rtssosdy", 0, "RTSSOSDY")
	flag.IntVar(&f.RTSSOSDWidth, "rtssosdwidth", 0, "RTSSOSDWidth")
	flag.IntVar(&f.RTSSOSDHeight, "rtssosdheight", 0, "RTSSOSDHeight")
	flag.BoolVar(&f.D3D9, "d3d9", false, "D3D9 enable support")
	flag.BoolVar(&f.OGL, "ogl", false, "OpenGL enable support")
	flag.BoolVar(&f.Fullscreen, "fullscreen", false, "activate fullscreen")
	flag.IntVar(&f.Width, "width", 0, "the width of the D3D9 application")
	flag.IntVar(&f.Height, "height", 0, "the height of the D3D9 application")
	flag.BoolVar(&f.Print, "print", false, "show the values in the console")
	flag.Parse()

	if f.Name != "" ||
		f.Port != 0 ||
		f.Count != -1 ||
		*timePtr != "" ||
		f.RTSS || f.D3D9 || f.OGL {
		// f.D3D9Fullscreen, f.LogFolder, f.Width & f.Height is optional

		if !f.RTSS && !f.D3D9 && !f.OGL {
			panic("use -rtss, -d3d9 or -ogl")
		}
		if f.Count == -1 && *timePtr == "" {
			fmt.Println("with -count or -time you can define an end so that a log is created")
		}

		var err error
		f.Delay, err = time.ParseDuration(*timePtr) // Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h".
		if err != nil {
			log.Println(err)
		}

		// log.Println(prettyPrint(f))

		cliMode = true
	}
}
