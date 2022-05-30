package main

import (
	"sort"
	"time"

	"github.com/VividCortex/ewma"
	"github.com/lxn/walk"

	//lint:ignore ST1001 standard behavior lxn/walk
	. "github.com/lxn/walk/declarative"
)

type Data struct {
	Timestamp time.Time
	Value     float64
}

type Dataset struct {
	Current []Data
	Backlog []Data
	Max     float64
	Min     float64
}

type DatabaseTable struct {
	checked       bool
	Index         int
	Count         int
	Countdown     time.Time
	Max           float64
	Min           float64
	MovingAverage ewma.MovingAverage
	Data          []Data
}

type DatabaseTableModel struct {
	walk.TableModelBase
	walk.SorterBase
	sortColumn int
	sortOrder  walk.SortOrder
	items      []*DatabaseTable
}

func NewDatabaseTableModel(dbBacklog *[]Database) *DatabaseTableModel {
	m := new(DatabaseTableModel)
	m.ResetRows(dbBacklog)
	return m
}

// Called by the TableView from SetModel and every time the model publishes a
// RowsReset event.
func (m *DatabaseTableModel) RowCount() int {
	return len(m.items)
}

// Called by the TableView when it needs the text to display for a given cell.
func (m *DatabaseTableModel) Value(row, col int) interface{} {
	item := m.items[row]

	switch col {
	case 0:
		return item.Index
	case 1:
		return item.Count
	case 2:
		return item.Countdown
	case 3:
		return item.Min
	case 4:
		return item.Max
	}

	panic("unexpected col")
}

// Called by the TableView to retrieve if a given row is checked.
func (m *DatabaseTableModel) Checked(row int) bool {
	return m.items[row].checked
}

// Called by the TableView when the user toggled the check box of a given row.
func (m *DatabaseTableModel) SetChecked(row int, checked bool) error {
	m.items[row].checked = checked

	return nil
}

// Called by the TableView to sort the model.
func (m *DatabaseTableModel) Sort(col int, order walk.SortOrder) error {
	m.sortColumn, m.sortOrder = col, order

	sort.SliceStable(m.items, func(i, j int) bool {
		a, b := m.items[i], m.items[j]

		c := func(ls bool) bool {
			if m.sortOrder == walk.SortAscending {
				return ls
			}

			return !ls
		}

		switch m.sortColumn {
		case 0:
			return c(a.Index < b.Index)

		case 1:
			return c(a.Count < b.Count)

		case 2:
			return c(a.Countdown.After(b.Countdown))

		case 3:
			return c(a.Min < b.Min)

		case 4:
			return c(a.Max < b.Max)
		}

		panic("unreachable")
	})

	return m.SorterBase.Sort(col, order)
}

func (m *DatabaseTableModel) ResetRows(dbBacklog *[]Database) {
	dbb := *dbBacklog
	for i := 0; i < len(dbb); i++ {
		temp := DatabaseTable{
			Index:         i,
			Count:         dbb[i].Count,
			Countdown:     dbb[i].Countdown,
			Max:           dbb[i].All.Max,
			Min:           dbb[i].All.Min,
			Data:          dbb[i].All.Backlog,
			MovingAverage: dbb[i].e,
		}
		m.items = append(m.items, &temp)
	}

	// Notify TableView and other interested parties about the reset.
	m.PublishRowsReset()

	m.Sort(m.sortColumn, m.sortOrder)
}

func SaveLogs(owner *GUI) (int, error) {
	var dlg *walk.Dialog
	model := NewDatabaseTableModel(&dbBacklog)
	var tv *walk.TableView

	// find the maximum
	var maxCount int
	for _, v := range model.items {
		if maxCount < v.Count {
			maxCount = v.Count
		}
	}
	// set the checkmark for the maximum
	for _, v := range model.items {
		if maxCount == v.Count {
			v.checked = true
		}
	}

	return Dialog{
		AssignTo: &dlg,
		Title:    "Which logs should be saved",
		Icon:     2,
		MinSize: Size{
			Width:  400,
			Height: 300,
		},
		Layout: VBox{},
		Children: []Widget{
			Label{
				Text: "Which data sets should be saved?",
			},
			TableView{
				AssignTo:         &tv,
				AlternatingRowBG: true,
				CheckBoxes:       true,
				ColumnsOrderable: true,
				Columns: []TableViewColumn{
					{Title: "#", Width: 40},
					{Title: "entries"},
					{Title: "started", Format: "2006-01-02 15:04:05"},
					{Title: "min"},
					{Title: "max"},
				},
				StyleCell: func(style *walk.CellStyle) {
					if model.items[style.Row()].checked {
						if style.Row()%2 == 0 {
							style.BackgroundColor = walk.RGB(159, 215, 255)
						} else {
							style.BackgroundColor = walk.RGB(143, 199, 239)
						}
					}
				},
				Model: model,
			},
			PushButton{
				Text: "Save",
				OnClicked: func() {
					for _, dbt := range model.items {
						if dbt.checked {
							SaveToFile(dbt)
							walk.MsgBox(nil, "", "File saved", walk.MsgBoxIconInformation)
						}
					}
				},
			},
		},
	}.Run(*owner)
}
