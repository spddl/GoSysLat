package main

import (
	"log"
	"time"

	"github.com/lxn/walk"
	"github.com/spddl/TargetWindow_D3D9"
	"github.com/spddl/TargetWindow_OpenGL"

	//lint:ignore ST1001 standard behavior lxn/walk
	. "github.com/lxn/walk/declarative"
)

type GUI struct {
	*walk.MainWindow
	Data             *Config
	FileName         *walk.LineEdit
	Count            *walk.Label
	Countdown        *walk.Label
	ValueMin         *walk.Label
	ValueMax         *walk.Label
	Value            *walk.Label
	ValueDelta       *walk.Label
	SecValue         *walk.Label
	SecValueDelta    *walk.Label
	MinuteValue      *walk.Label
	MinuteValueDelta *walk.Label
	EwmaValue        *walk.Label
	EwmaValueDelta   *walk.Label
}

var (
	greenColor      = walk.RGB(0, 0xff, 0)
	redColor        = walk.RGB(0xff, 0, 0)
	backgroundColor = walk.RGB(0xf0, 0xf0, 0xf0)
)

var (
	targetD3D9 TargetWindow_D3D9.Target
	targetOGL  TargetWindow_OpenGL.Target
)

func (g *GUI) Start() {
	g.Data = &c

	var openSettings *walk.Action

	if err := (MainWindow{
		AssignTo: &g.MainWindow,
		Title:    "GoSysLat",
		Icon:     2,
		Size: Size{
			Width:  250,
			Height: 1,
		},
		DoubleBuffering: true,
		MenuItems: []MenuItem{
			Action{
				AssignTo: &openSettings,
				Text:     "&Settings",
				OnTriggered: func() {
					RunSettingsDialog(g)
				},
			},
			Menu{
				Text: "&Target Window",
				Items: []MenuItem{
					Action{
						AssignTo: &openSettings,
						Text:     "Open in Window (D3D9)",
						OnTriggered: func() {
							targetD3D9 = TargetWindow_D3D9.Target{}
							go func() {
								targetD3D9.Start(currentFileName, false, f.Width, f.Height, c.PushMode)
								targetD3D9.Close()
							}()
						},
					},
					Action{
						AssignTo: &openSettings,
						Text:     "Open in Fullscreen (D3D9)",
						OnTriggered: func() {
							targetD3D9 = TargetWindow_D3D9.Target{}
							go func() {
								targetD3D9.Start(currentFileName, true, f.Width, f.Height, c.PushMode)
								targetD3D9.Close()
							}()
						},
					},
					Separator{},
					Action{
						AssignTo: &openSettings,
						Text:     "Open in Window (OGL)",
						OnTriggered: func() {
							targetOGL = TargetWindow_OpenGL.Target{}
							go targetOGL.Start(false, f.Width, f.Height, c.PushMode)
						},
					},
					Action{
						AssignTo: &openSettings,
						Text:     "Open in Fullscreen (OGL)",
						OnTriggered: func() {
							targetOGL = TargetWindow_OpenGL.Target{}
							go targetOGL.Start(true, f.Width, f.Height, c.PushMode)
						},
					},
				},
			},
			Menu{
				Text: "&Help",
				Items: []MenuItem{
					Action{
						Text:        "How to start?",
						OnTriggered: g.howToStart_Triggered,
					},
					Action{
						Text:        "About",
						OnTriggered: g.aboutAction_Triggered,
					},
				},
			},
		},
		Layout: VBox{
			Alignment: AlignHNearVNear,
		},
		Children: []Widget{
			Composite{
				Layout: HBox{},
				Children: []Widget{
					Label{
						Text: "Project name",
					},
					LineEdit{
						AssignTo: &g.FileName,
						Text:     "",
					},
				},
			},
			Composite{
				Layout: Grid{Columns: 3},
				Children: []Widget{

					Label{
						Text: "Count:",
					},
					Label{
						ColumnSpan: 2,
						Text:       "-",
						AssignTo:   &g.Count,
						Font:       Font{Family: "Segoe UI", PointSize: 20},
					},

					Label{
						Text: "Current Value (in ms):",
					},
					Label{
						Text:     "-",
						AssignTo: &g.Value,
						Font:     Font{Family: "Segoe UI", PointSize: 20},
						MinSize: Size{
							Width: 150,
						},
					},
					Label{
						Text:       "",
						AssignTo:   &g.ValueDelta,
						TextColor:  backgroundColor,
						Background: SolidColorBrush{Color: backgroundColor},
					},

					Label{
						Text: "Average seconds:",
					},
					Label{
						Text:     "-",
						AssignTo: &g.SecValue,
						Font:     Font{PointSize: 20},
					},
					Label{
						Text:       "",
						AssignTo:   &g.SecValueDelta,
						TextColor:  backgroundColor,
						Background: SolidColorBrush{Color: backgroundColor},
					},

					Label{
						Text: "Average minute:",
					},
					Label{
						Text:     "-",
						AssignTo: &g.MinuteValue,
						Font:     Font{PointSize: 20},
					},
					Label{
						Text:       "",
						AssignTo:   &g.MinuteValueDelta,
						TextColor:  backgroundColor,
						Background: SolidColorBrush{Color: backgroundColor},
					},

					Label{
						Text:        "EWMA Value:",
						ToolTipText: "SimpleEWMA",
					},
					Label{
						Text:     "-",
						AssignTo: &g.EwmaValue,
						Font:     Font{PointSize: 20},
					},
					Label{
						Text:       "",
						AssignTo:   &g.EwmaValueDelta,
						TextColor:  backgroundColor,
						Background: SolidColorBrush{Color: backgroundColor},
					},
				},
			},

			Composite{
				Layout: Grid{Columns: 2},
				Children: []Widget{
					Label{
						Text: "Countdown:",
					},
					Label{
						Text:     "",
						AssignTo: &g.Countdown,
					},

					Label{
						Text: "Min:",
					},
					Label{
						Text:     "-",
						AssignTo: &g.ValueMin,
					},

					Label{
						Text: "Max:",
					},
					Label{
						Text:     "-",
						AssignTo: &g.ValueMax,
					},
				},
			},
			VSpacer{},
			PushButton{
				Text: "Reset (Ctrl + Del)",
				OnClicked: func() {
					g.ResetData()
					g.ResetGUI()
				},
			},
			PushButton{
				Text: "Save Logs (Ctrl + Return)",
				OnClicked: func() {
					SaveLogs(&gui)
				},
			},
		},
	}).Create(); err != nil {
		log.Fatal(err)
	}
}

func (g *GUI) aboutAction_Triggered() {
	walk.MsgBox(g, "About", "I started with GoSysLat to better understand how the SysLat works. When I had rebuilt the basic framework I tried to improve the system further, faster with alternative methods to RTSS, e.g. OpenGL and DirectX9. The goal was to get more accurate and better visualized results.", walk.MsgBoxIconInformation)
}

func (g *GUI) howToStart_Triggered() {
	walk.MsgBox(g, "How to start", `First you need a different firmware than the factory one. "github.com/spddl/SysLat_Firmware"

After that you have to select the correct COM port in the settings.

Now you can create a "Target Window" or open it in a game with RTSS. If an area now changes between black and white, the communication with your SysLat works.
Now you can hold the sensor against this area on the monitor and the software will do the rest.`, walk.MsgBoxIconInformation)
}

func (g *GUI) SetValue(value *walk.Label, text string) {
	if g == nil || value == nil {
		return
	}
	g.Synchronize(func() {
		if err := value.SetText(text); err != nil {
			log.Println(err)
		}
	})
}

func (g *GUI) SetDeltaValue(value *walk.Label, oldValue, newValue float64) {
	if g == nil || value == nil {
		return
	}
	deltafloat64 := Round(newValue - oldValue)
	g.Synchronize(func() {
		switch {
		case deltafloat64 == 0:
			value.SetTextColor(backgroundColor)
		case newValue > oldValue:
			value.SetText("( +" + float64ToString(deltafloat64) + " )")
			value.SetTextColor(redColor)
		case newValue < oldValue:
			value.SetText("( " + float64ToString(deltafloat64) + " )")
			value.SetTextColor(greenColor)
		default:
			// will never happen but makes the compiler happy
			value.SetTextColor(backgroundColor)
		}
	})
}

func (g *GUI) ResetData() {
	db.Count = 0
	db.Countdown = time.Time{}
	db.All = Dataset{}
	db.Second = Dataset{}
	db.Minute = Dataset{}
}

func (g *GUI) ResetGUI() {
	g.Synchronize(func() {
		g.Count.SetText("-")
		g.Value.SetText("-")
		g.ValueDelta.SetText("")
		g.SecValue.SetText("-")
		g.SecValueDelta.SetText("")
		g.MinuteValue.SetText("-")
		g.MinuteValueDelta.SetText("")
		g.EwmaValue.SetText("-")
		g.EwmaValueDelta.SetText("")
		g.Countdown.SetText("")
		g.ValueMin.SetText("-")
		g.ValueMax.SetText("-")
	})
}
