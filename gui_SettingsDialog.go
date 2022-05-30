package main

import (
	"fmt"
	"log"

	"github.com/lxn/walk"
	"github.com/spddl/USBController"

	//lint:ignore ST1001 standard behavior lxn/walk
	. "github.com/lxn/walk/declarative"
)

type Port struct {
	Id      int
	Name    string
	COMName string
}

func KnownPorts() []*Port {
	var result []*Port
	for _, v := range USBController.COMDevices {
		p := new(Port)
		p.Name = v.FriendlyName
		p.Id = StringToInt(v.COMid)
		p.COMName = "COM" + v.COMid
		result = append(result, p)
	}
	return result
}

func RunSettingsDialog(owner *GUI) (int, error) {
	var dlg *walk.Dialog
	var db *walk.DataBinder
	var ComPortCB *walk.ComboBox
	var acceptPB, cancelPB *walk.PushButton
	var OSDXS, OSDYS, OSDWidthS, OSDHeightS *walk.Slider
	var OSDXNE, OSDYNE, OSDWidthNE, OSDHeightNE *walk.NumberEdit
	var DetectLightL *walk.Label
	var DisplayCB, DetectLightCB *walk.CheckBox

	return Dialog{
		AssignTo:      &dlg,
		Title:         "Settings",
		Icon:          2,
		DefaultButton: &acceptPB,
		CancelButton:  &cancelPB,
		DataBinder: DataBinder{
			AssignTo:       &db,
			Name:           "data",
			DataSource:     owner.Data,
			ErrorPresenter: ToolTipErrorPresenter{},
		},
		MinSize: Size{
			Width:  400,
			Height: 300,
		},
		Layout: VBox{},
		Children: []Widget{
			Composite{
				Layout: Grid{Columns: 2},
				Children: []Widget{
					Label{
						Text: "COM Port:",
					},
					ComboBox{
						AssignTo:      &ComPortCB,
						Value:         c.ComPort,
						BindingMember: "Id",
						DisplayMember: "Name",
						Model:         KnownPorts(),
						OnCurrentIndexChanged: func() { // BUG: called 3 times
							newCOM := KnownPorts()[ComPortCB.CurrentIndex()].Id
							if c.ComPort != newCOM {
								queue.cancel()
								queue.wg.Wait()
								USBController.CloseComPort(hPort) // close handle to ComPort

								log.Println("OpenComPort", KnownPorts()[ComPortCB.CurrentIndex()].COMName)
								hPort = USBController.OpenComPort(KnownPorts()[ComPortCB.CurrentIndex()].COMName)
								if hPort == 0 || hPort == 0xFFFFFFFFFFFFFFFF {
									log.Printf("Can't Open COM%d Port. %x", c.ComPort, hPort)
									walk.MsgBox(nil, "Error", fmt.Sprintf("Can't Open COM%d port\nchange the COM port under settings and restart the tool", c.ComPort), walk.MsgBoxIconError)
									return
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
								go queue.Ticker()

								c.ComPort = KnownPorts()[ComPortCB.CurrentIndex()].Id
								c.Save()
								go queue.DataProcessing()
							}
						},
					},

					GroupBox{
						Title:      "RTSS OSD",
						ColumnSpan: 2,
						Layout:     Grid{Columns: 3},
						Children: []Widget{
							Label{
								Text:        "X:",
								ToolTipText: "X Coordinate",
							},
							Slider{
								AssignTo: &OSDXS,
								MinValue: 0,
								MaxValue: 100,
								Value:    Bind("OSDX"),
								OnValueChanged: func() {
									value := OSDXS.Value()
									OSDXNE.SetValue(float64(value))
									c.OSDX = value
									c.RefreshOSDStrings()
									c.Save()
								},
							},
							NumberEdit{
								AssignTo: &OSDXNE,
								MaxSize: Size{
									Width: 60,
								},
								Value: Bind("OSDX", Range{
									Min: 0.0,
									Max: 100.0,
								}),
								SpinButtonsVisible: true,
								Decimals:           0,
								OnValueChanged: func() {
									value := int(OSDXNE.Value())
									OSDXS.SetValue(value)
									c.OSDX = value
									c.RefreshOSDStrings()
									c.Save()
								},
							},

							Label{
								Text:        "Y:",
								ToolTipText: "Y Coordinate",
							},
							Slider{
								AssignTo: &OSDYS,
								MinValue: 0,
								MaxValue: 100,
								Value:    Bind("OSDY"),
								OnValueChanged: func() {
									value := OSDYS.Value()
									OSDYNE.SetValue(float64(value))
									c.OSDY = value
									c.RefreshOSDStrings()
									c.Save()
								},
							},
							NumberEdit{
								AssignTo: &OSDYNE,
								MaxSize: Size{
									Width: 60,
								},
								Value: Bind("OSDY", Range{
									Min: 0.0,
									Max: 100.0,
								}),
								SpinButtonsVisible: true,
								Decimals:           0,
								OnValueChanged: func() {
									value := int(OSDYNE.Value())
									OSDYS.SetValue(value)
									c.OSDY = value
									c.RefreshOSDStrings()
									c.Save()
								},
							},

							Label{
								Text:        "Width:",
								ToolTipText: "Width of the OSD's square",
							},
							Slider{
								AssignTo: &OSDWidthS,
								MinValue: 0,
								MaxValue: 100,
								Value:    Bind("OSDWidth"),
								OnValueChanged: func() {
									value := OSDWidthS.Value()
									OSDWidthNE.SetValue(float64(value))
									c.OSDWidth = value
									c.RefreshOSDStrings()
									c.Save()
								},
							},
							NumberEdit{
								AssignTo: &OSDWidthNE,
								MaxSize: Size{
									Width: 60,
								},
								Value: Bind("OSDWidth", Range{
									Min: 0.0,
									Max: 100.0,
								}),
								SpinButtonsVisible: true,
								Decimals:           0,
								OnValueChanged: func() {
									value := int(OSDWidthNE.Value())
									OSDWidthS.SetValue(value)
									c.OSDWidth = value
									c.RefreshOSDStrings()
									c.Save()
								},
							},

							Label{
								Text:        "Height:",
								ToolTipText: "Height of the OSD's square",
							},
							Slider{
								AssignTo: &OSDHeightS,
								MinValue: 0,
								MaxValue: 100,
								Value:    Bind("OSDWidth"),
								OnValueChanged: func() {
									value := OSDHeightS.Value()
									OSDHeightNE.SetValue(float64(value))
									c.OSDHeight = value
									c.RefreshOSDStrings()
									c.Save()
								},
							},
							NumberEdit{
								AssignTo: &OSDHeightNE,
								MaxSize: Size{
									Width: 60,
								},
								Value: Bind("OSDHeight", Range{
									Min: 0.0,
									Max: 100.0,
								}),
								SpinButtonsVisible: true,
								Decimals:           0,
								OnValueChanged: func() {
									value := int(OSDHeightNE.Value())
									OSDHeightS.SetValue(value)
									c.OSDHeight = value
									c.RefreshOSDStrings()
									c.Save()
								},
							},
						},
					},

					GroupBox{
						Title: "D3D9 / OpenGL",
						Layout: VBox{
							Alignment: AlignHNearVNear,
						},
						Children: []Widget{
							RadioButtonGroup{
								DataMember: "PushMode",
								Buttons: []RadioButton{
									{
										ToolTipText: "Here the images are updated as repetitively as possible",
										Text:        "Benchmark Mode",
										Value:       false,
										OnClicked: func() {
											c.PushMode = false
											c.Save()
										},
									},
									{
										ToolTipText: "Here the images are only updated when the syslat sends a signal",
										Text:        "Push Mode",
										Value:       true,
										OnClicked: func() {
											c.PushMode = true
											c.Save()
										},
									},
								},
							},
						},
					},
					GroupBox{
						Title: "Device Settings",
						// Layout: VBox{
						// 	Alignment: AlignHNearVNear,
						// },
						Layout: Grid{
							Alignment: AlignHNearVNear,
							Columns:   2,
						},

						Children: []Widget{

							Label{
								Text: "Display:",
							},
							CheckBox{
								AssignTo: &DisplayCB,
								Checked:  Bind("Display"),
								OnCheckedChanged: func() {
									c.Display = DisplayCB.Checked()
									c.Save()

									if DisplayCB.Checked() {
										USBController.WriteByte(&hPort, []byte{1})
									} else {
										USBController.WriteByte(&hPort, []byte{2})
									}
								},
							},

							Label{
								AssignTo: &DetectLightL,
								Text:     "Detect Light:",
							},
							CheckBox{
								AssignTo: &DetectLightCB,
								Checked:  Bind("DetectLight"),
								OnCheckedChanged: func() {
									c.DetectLight = DetectLightCB.Checked()
									c.Save()

									var result bool
									if DetectLightCB.Checked() {
										DetectLightL.SetText("Detect Light: (Black > White)")
										result = USBController.WriteByte(&hPort, []byte{3})
									} else {
										DetectLightL.SetText("Detect no Light: (White > Black)")
										result = USBController.WriteByte(&hPort, []byte{4})
									}
									log.Println(result)
								},
							},
						},
					},
					VSpacer{},
				},
			},
			Composite{
				Layout: HBox{},
				Children: []Widget{
					HSpacer{},
					PushButton{
						AssignTo: &acceptPB,
						Text:     "OK",
						OnClicked: func() {
							dlg.Accept()
						},
					},
				},
			},
		},
	}.Run(*owner)
}
