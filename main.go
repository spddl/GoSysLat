package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/MakeNowJust/hotkey"
	"github.com/VividCortex/ewma"
	"github.com/lxn/walk"
	"golang.org/x/sys/windows"

	"github.com/spddl/RTSSClient"
	"github.com/spddl/TargetWindow_D3D9"
	"github.com/spddl/TargetWindow_OpenGL"

	"github.com/spddl/USBController"
)

type Database struct {
	Count     int
	Countdown time.Time
	All       Dataset
	Second    Dataset
	Minute    Dataset
	e         ewma.MovingAverage
}

var (
	db              Database
	dbBacklog       []Database
	hPort           syscall.Handle
	gui             GUI
	RTSSOSDBlack    string
	RTSSOSDWhite    string
	sigs            chan os.Signal
	queue           Queue
	currentFileName string
	toggle          chan bool
	trigger         bool
)

func init() {
	log.SetFlags(log.Lmicroseconds | log.Lshortfile) // https://ispycode.com/GO/Logging/Setting-output-flags

	// https://docs.microsoft.com/en-us/windows/win32/api/processthreadsapi/nf-processthreadsapi-getpriorityclass
	windows.SetPriorityClass(windows.CurrentProcess(), 0x00000100) // REALTIME_PRIORITY_CLASS, if it is started as admin otherwise it is only high

	// https://en.wikipedia.org/wiki/Moving_average#Exponential_moving_average
	db.e = ewma.NewMovingAverage(60)

	e, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	currentFileName = filepath.Base(e)
}

func cleanUp() {
	RTSSClient.UpdateOSD("") // needs testing, maybe you can hide the rectangle like this
	if targetOGL.IsActive {
		targetOGL.Close()
	}
	if targetD3D9.IsActive {
		targetD3D9.Close()
	}
	if trigger {
		USBController.WriteByte(&hPort, []byte{5})
	}
	queue.cancel()
	log.Println("queue.wg.Wait() // block")
	queue.wg.Wait() // block

	USBController.CloseComPort(hPort) // close handle to ComPort
}

func prettyPrint(i interface{}) string {
	s, _ := json.MarshalIndent(i, "", "\t")
	return string(s)
}

func main() {
	sigs = make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		cleanUp()
		os.Exit(0)
	}()

	c.RefreshOSDStrings()
	if !cliMode {
		gui.Start()
		gui.Show()

		hkey := hotkey.New()
		hkey.Register(hotkey.Ctrl, hotkey.DELETE, func() {
			gui.ResetData()
			gui.ResetGUI()
		})
		hkey.Register(hotkey.Ctrl, hotkey.RETURN, func() {
			var dbTemp = []Database{db}
			gui.ResetData()
			gui.ResetGUI()
			dbt := NewDatabaseTableModel(&dbTemp)
			for _, values := range dbt.items {
				if len(values.Data) != 0 {
					SaveToFile(values)
				}
			}
		})

		/// the syslat is not designed for this.
		// var quit chan struct{}
		// hkey.Register(hotkey.Ctrl, 0x52, func() {
		// 	USBController.WriteByte(&hPort, []byte{5})
		// 	trigger = !trigger
		// 	if trigger {
		// 		mouseHook = SetWindowsHookEx(WH_MOUSE_LL, (HOOKPROC)(func(nCode int, wparam WPARAM, lparam LPARAM) LRESULT {
		// 			if nCode == 0 && wparam == WM_LBUTTONDOWN {
		// 				log.Println("WM_LBUTTONDOWN", trigger)
		// 				if trigger {
		// 					go USBController.WriteByte(&hPort, []byte{0})
		// 				}
		// 			}
		// 			return CallNextHookEx(mouseHook, nCode, wparam, lparam)
		// 		}), 0, 0)
		// 		quit = make(chan struct{})
		// 		go func() {
		// 			var msg MSG
		// 			for {
		// 				select {
		// 				case <-quit: // BUG: the loop will not be killed
		// 					return
		// 				default:
		// 					// if GetMessage(&msg, 0, 0, 0) != 0 {
		// 					// 	TranslateMessage(&msg)
		// 					// 	DispatchMessage(&msg)
		// 					// }
		// 					for PeekMessage(&msg, 0, 0, 0, PM_NOREMOVE|PM_QS_INPUT) {

		// 					}
		// 				}
		// 			}
		// 		}()
		// 	} else {
		// 		close(quit)
		// 		UnhookWindowsHookEx(mouseHook)
		// 		mouseHook = 0
		// 	}
		// })
	}

	if cliMode {
		c.ComPort = f.Port
	}

	// the first valid port as default port
	if c.ComPort == 0 || !TestValidComPort(IntToString(c.ComPort)) {
		c.Reset()
		c.ComPort = StringToInt(USBController.COMDevices[0].COMid)
		c.Save()
	}

	hPort = USBController.OpenComPort(fmt.Sprintf("COM%d", c.ComPort))
	if hPort == 0 || hPort == 0xFFFFFFFFFFFFFFFF {
		walk.MsgBox(nil, "Error", fmt.Sprintf("Can't Open COM%d port", c.ComPort), walk.MsgBoxIconError)
		c.Reset()
		// the first valid port as default port
		c.ComPort = StringToInt(USBController.COMDevices[0].COMid)
		c.Save()
	}
	if c.Display {
		USBController.WriteByte(&hPort, []byte{1})
	} else {
		USBController.WriteByte(&hPort, []byte{2})
	}
	if c.DetectLight {
		USBController.WriteByte(&hPort, []byte{3})
	} else {
		USBController.WriteByte(&hPort, []byte{4})
	}

	queue.StartQueue()
	go queue.Reading(&hPort)
	if cliMode {
		go queue.TickerCli()
	} else {
		go queue.Ticker()
	}
	go queue.DataProcessing()

	if cliMode {
		if f.D3D9 {
			targetD3D9 = TargetWindow_D3D9.Target{}
			go func() {
				targetD3D9.Start(currentFileName, f.Fullscreen, f.Width, f.Height, c.PushMode)
				sigs <- syscall.SIGINT
			}()
		} else if f.OGL {
			targetOGL = TargetWindow_OpenGL.Target{}
			go func() {
				targetOGL.Start(f.Fullscreen, f.Width, f.Height, c.PushMode)
				sigs <- syscall.SIGINT
			}()
		}
	} else {
		gui.Run()
		sigs <- syscall.SIGINT
	}

	<-sigs
	cleanUp()
}

func TestValidComPort(port string) bool {
	for _, com := range USBController.COMDevices {
		if com.COMid == port {
			return true
		}
	}
	return false
}
