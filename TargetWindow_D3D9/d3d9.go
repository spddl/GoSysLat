package TargetWindow_D3D9

import (
	"log"
	"runtime"
	"runtime/debug"
	"syscall"
	"time"

	"github.com/gonutz/d3d9"
	"github.com/gonutz/w32/v2"
)

type Target struct {
	device      *d3d9.Device
	verticesLen uint
	IsActive    bool
	WhiteBox    bool
}

var (
	colorWhite        = d3d9.ColorRGB(255, 255, 255)
	colorBlack        = d3d9.ColorRGB(0, 0, 0)
	previousPlacement w32.WINDOWPLACEMENT
)

func (t *Target) Start(currentFileName string, fullscreen bool, fWidth, fHeight int, pushMode bool) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	classNamePtr, _ := syscall.UTF16PtrFromString("fullscreen_window_class")
	w32.RegisterClassEx(&w32.WNDCLASSEX{
		Cursor: w32.LoadCursor(0, w32.MakeIntResource(w32.IDC_ARROW)),
		WndProc: syscall.NewCallback(func(window w32.HWND, msg uint32, w, l uintptr) uintptr {
			switch msg {
			case w32.WM_KEYDOWN:
				switch w {
				case w32.VK_ESCAPE:
					w32.SendMessage(window, w32.WM_CLOSE, 0, 0)
					w32.PostQuitMessage(0)
				case w32.VK_F11:
					toggleFullscreen(window) // BUG: stays, Composed: Copy with GPU GDI
				}
				return 1
			case w32.WM_DESTROY:
				w32.PostQuitMessage(0)
				return 0
			default:
				return w32.DefWindowProc(window, msg, w, l)
			}
		}),
		ClassName: classNamePtr,
		Icon:      w32.ExtractIcon(currentFileName, 0),
	})

	windowNamePtr, _ := syscall.UTF16PtrFromString("TargetWindow D3D9")
	windowHandle := w32.CreateWindow(
		classNamePtr,
		windowNamePtr,
		w32.WS_OVERLAPPEDWINDOW|w32.WS_VISIBLE,
		w32.CW_USEDEFAULT,
		w32.CW_USEDEFAULT,
		640,
		480,
		0,
		0,
		0,
		nil,
	)

	var err error
	d3d, err := d3d9.Create(d3d9.SDK_VERSION)
	defer d3d.Release()
	check(err)

	if fullscreen { // Hardware: Legacy Flip
		var width, height uint32
		if fWidth == 0 || fHeight == 0 {
			MaxWidth, MaxHeight, _ := checkMaxRes()
			if fWidth == 0 {
				width = MaxWidth
			} else {
				width = uint32(fWidth)
			}
			if fHeight == 0 {
				height = MaxHeight
			} else {
				height = uint32(fHeight)
			}
		} else {
			width = uint32(fWidth)
			height = uint32(fHeight)
		}
		for {
			t.device, _, err = d3d.CreateDevice(
				d3d9.ADAPTER_DEFAULT,
				d3d9.DEVTYPE_HAL,
				d3d9.HWND(windowHandle),
				d3d9.CREATE_SOFTWARE_VERTEXPROCESSING,
				d3d9.PRESENT_PARAMETERS{
					Windowed:         0,
					HDeviceWindow:    d3d9.HWND(windowHandle),
					SwapEffect:       d3d9.SWAPEFFECT_DISCARD,
					BackBufferFormat: d3d9.FMT_X8R8G8B8,
					BackBufferWidth:  width,
					BackBufferHeight: height,
				},
			)
			if err == nil {
				break
			}
			time.Sleep(time.Second)
		}
	} else { // window mode // Composed: Copy with GPU GDI
		t.device, _, err = d3d.CreateDevice(
			d3d9.ADAPTER_DEFAULT,
			d3d9.DEVTYPE_HAL,
			d3d9.HWND(windowHandle),
			d3d9.CREATE_HARDWARE_VERTEXPROCESSING,
			d3d9.PRESENT_PARAMETERS{
				Windowed:      1,
				HDeviceWindow: d3d9.HWND(windowHandle),
				SwapEffect:    d3d9.SWAPEFFECT_DISCARD,
			},
		)
		check(err)
	}
	defer t.device.Release()
	check(t.device.SetRenderState(d3d9.RS_CULLMODE, uint32(d3d9.CULL_NONE)))

	decl, err := t.device.CreateVertexDeclaration([]d3d9.VERTEXELEMENT{
		{
			Stream:     0,
			Offset:     0,
			Type:       d3d9.DECLTYPE_FLOAT2,
			Method:     d3d9.DECLMETHOD_DEFAULT,
			Usage:      d3d9.DECLUSAGE_POSITION,
			UsageIndex: 0,
		},
		d3d9.DeclEnd(),
	})
	check(err)
	defer decl.Release()
	check(t.device.SetVertexDeclaration(decl))

	vertices := []float32{
		-1, -0.8,
		-1, 0.8,
		1, -0.8,

		1, 0.8,
		1, -0.8,
		-1, 0.8,

		-0.8, -1,
		-0.8, 1,
		0.8, 1,

		0.8, 1,
		0.8, -1,
		-0.8, -1,
	}

	t.verticesLen = uint(len(vertices) / 3)
	vb, err := t.device.CreateVertexBuffer(uint(len(vertices)*4), d3d9.USAGE_WRITEONLY, 0, d3d9.POOL_DEFAULT, 0)
	check(err)
	defer vb.Release()
	vbMem, err := vb.Lock(0, 0, d3d9.LOCK_DISCARD)
	check(err)
	vbMem.SetFloat32s(0, vertices)
	check(vb.Unlock())

	check(t.device.SetStreamSource(0, vb, 0, 2*4))

	if !pushMode { // the best I can achieve
		// create a timer that ticks every 10ms and register a callback for it
		// https://docs.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-settimer
		w32.SetTimer(windowHandle, 1, 10, 0) // USER_TIMER_MINIMUM
	}

	var msg w32.MSG
	for w32.GetMessage(&msg, 0, 0, 0) != 0 {
		if msg.Message == 0x0012 { // WM_QUIT
			t.IsActive = false
		} else {
			w32.TranslateMessage(&msg)
			t.IsActive = true
			w32.DispatchMessage(&msg)
		}
	}
}

func (t *Target) Close() {
	t.IsActive = false
}

func (t *Target) SetWhite() {
	t.WhiteBox = true
	t.draw()
}
func (t *Target) SetBlack() {
	t.WhiteBox = false
	t.draw()
}

func (t *Target) draw() {
	if !t.IsActive {
		return
	}

	if t.WhiteBox {
		t.device.Clear(nil, d3d9.CLEAR_TARGET, colorWhite, 0, 0)
	} else {
		t.device.Clear(nil, d3d9.CLEAR_TARGET, colorBlack, 0, 0)
	}
	check(t.device.BeginScene()) // invalid call
	check(t.device.DrawPrimitive(d3d9.PT_TRIANGLELIST, 0, t.verticesLen))
	check(t.device.EndScene())
	t.device.Present(nil, nil, 0, nil)
}

func checkMaxRes() (uint32, uint32, uint32) {
	var maxScreenW uint32
	var maxScreenH uint32
	var maxRefreshRate uint32
	d3d, err := d3d9.Create(d3d9.SDK_VERSION)
	check(err)
	defer d3d.Release()

	adapterCount := d3d.GetAdapterCount()
	for adapter := uint(0); adapter < adapterCount; adapter++ {
		displayMode, err := d3d.GetAdapterDisplayMode(adapter)
		check(err)

		if maxRefreshRate < displayMode.RefreshRate {
			maxRefreshRate = displayMode.RefreshRate
			if displayMode.Width > maxScreenW {
				maxScreenW = displayMode.Width
			}
			if displayMode.Height > maxScreenH {
				maxScreenH = displayMode.Height
			}
		}
	}
	return maxScreenW, maxScreenH, maxRefreshRate
}

func toggleFullscreen(window w32.HWND) {
	style := w32.GetWindowLong(window, w32.GWL_STYLE)
	if style&w32.WS_OVERLAPPEDWINDOW != 0 {
		// go into full-screen
		var monitorInfo w32.MONITORINFO
		monitor := w32.MonitorFromWindow(window, w32.MONITOR_DEFAULTTOPRIMARY)
		if w32.GetWindowPlacement(window, &previousPlacement) &&
			w32.GetMonitorInfo(monitor, &monitorInfo) {
			w32.SetWindowLong(
				window,
				w32.GWL_STYLE,
				style & ^w32.WS_OVERLAPPEDWINDOW,
			)
			w32.SetWindowPos(
				window,
				0,
				int(monitorInfo.RcMonitor.Left),
				int(monitorInfo.RcMonitor.Top),
				int(monitorInfo.RcMonitor.Right-monitorInfo.RcMonitor.Left),
				int(monitorInfo.RcMonitor.Bottom-monitorInfo.RcMonitor.Top),
				w32.SWP_NOOWNERZORDER|w32.SWP_FRAMECHANGED,
			)
		}
		w32.ShowCursor(false)
	} else {
		// go into windowed mode
		w32.SetWindowLong(
			window,
			w32.GWL_STYLE,
			style|w32.WS_OVERLAPPEDWINDOW,
		)
		w32.SetWindowPlacement(window, &previousPlacement)
		w32.SetWindowPos(window, 0, 0, 0, 0, 0,
			w32.SWP_NOMOVE|w32.SWP_NOSIZE|w32.SWP_NOZORDER|
				w32.SWP_NOOWNERZORDER|w32.SWP_FRAMECHANGED,
		)
		w32.ShowCursor(true)
	}
}

func check(err error) {
	if err != nil {
		debug.PrintStack()
		log.Fatalln(err)
		panic(err)
	}
}
