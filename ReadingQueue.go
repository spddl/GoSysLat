package main

import (
	"context"
	"fmt"
	"os"
	"sync"
	"syscall"
	"time"

	"github.com/spddl/RTSSClient"
	"github.com/spddl/USBController"
)

type Queue struct {
	ctx    context.Context
	cancel context.CancelFunc
	wg     *sync.WaitGroup
	ticker *time.Ticker
	sendch chan int32
}

func (q *Queue) StartQueue() {
	q.ctx, q.cancel = context.WithCancel(context.Background())
	q.wg = &sync.WaitGroup{}
	q.wg.Add(2) // Reading & (Ticker || TickerCli)
	q.ticker = time.NewTicker(time.Second)
	q.sendch = make(chan int32)
}

func (q *Queue) Reading(hPort *syscall.Handle) {
	defer q.wg.Done()

	sysLatResultsBytes := []byte{}

	// timerA := time.Now()
	// timerB := time.Now()
	// timerC := time.Now()

	for {
		select {
		case <-q.ctx.Done():
			return
		default:
		}

		serialReadData, ok := USBController.ReadByte(hPort) // the function runs every 100ms on timeout
		if !ok {                                            // and starts from the beginning if no data was received
			continue
		}

		// log.Printf("serialReadData - %d\n", serialReadData)

		switch serialReadData {
		case 1: // A
			// fmt.Printf("- From C to A => %s\n", time.Since(timerC))
			// timerA = time.Now()
			q.sendch <- -1
		case 2: // B
			// fmt.Printf("- From A to B => \t\t%s\n", time.Since(timerA))
			// timerB = time.Now()
			q.sendch <- -2
		case 3: // C
			// log.Println("C")
			// fmt.Printf("- From B to C => \t\t\t\t%s\n", time.Since(timerB))
			// timerC = time.Now()

			// log.Println("len(sysLatResultsBytes)", len(sysLatResultsBytes), sysLatResultsBytes, string(sysLatResultsBytes))
			// sysLatResultsBytes = []byte{}

			if len(sysLatResultsBytes) != 0 {
				data := make([]byte, 8)
				copy(data, sysLatResultsBytes)
				intVar := ByteArrayToInt(data)
				q.sendch <- int32(intVar)
				sysLatResultsBytes = []byte{}
			} else {
				if db.Count != 0 {
					dbBacklog = append(dbBacklog, db)
					gui.ResetData()
				}
				db.Countdown = time.Time{}
			}
		default:
			sysLatResultsBytes = append(sysLatResultsBytes, serialReadData)
		}
	}
}

func (q *Queue) Ticker() {
	defer q.wg.Done()
	for {
		select {
		case <-q.ticker.C:
			countCurrent := len(db.All.Current)
			if countCurrent != 0 {
				// log.Println(countCurrent, "records in one sec")

				// calculates the average of all data accumulated within one second
				var sum float64 = 0
				for i := 0; i < countCurrent; i++ {
					sum += db.All.Current[i].Value
				}
				SecValue := sum / float64(countCurrent)

				// and adds it to the "Second.Current"
				db.Second.Backlog = append(db.Second.Backlog, Data{
					Timestamp: time.Now(),
					Value:     SecValue,
				})
				gui.SetValue(gui.SecValue, float64ToString(Round(SecValue)))
				lenSecondBacklog := len(db.Second.Backlog)
				if lenSecondBacklog > 1 {
					gui.SetDeltaValue(gui.SecValueDelta, db.Second.Backlog[lenSecondBacklog-2].Value, SecValue)
				}

				// All "Current" data is pushed into the "Backlog" for later use
				db.All.Backlog = append(db.All.Backlog, db.All.Current...)
				db.All.Current = []Data{}

				countSecond := len(db.Second.Backlog)
				if countSecond >= 60 { // one minute has passed and we can start calculating the average in the last minute
					last60Values := db.Second.Backlog[countSecond-60:]
					var sum float64 = 0
					for i := 0; i < 60; i++ {
						sum += last60Values[i].Value
					}
					MinuteValue := sum / float64(len(last60Values))

					// and adds it to the "Minute.Current"
					db.Minute.Current = append(db.Minute.Current, Data{
						Timestamp: time.Now(),
						Value:     MinuteValue,
					})
					gui.SetValue(gui.MinuteValue, float64ToString(Round(MinuteValue)))

					lenMinuteCurrent := len(db.Minute.Current)
					if lenMinuteCurrent > 1 {
						gui.SetDeltaValue(gui.MinuteValueDelta, db.Minute.Current[lenMinuteCurrent-2].Value, MinuteValue)
					}

					db.Second.Backlog = append(db.Second.Backlog, db.Second.Current...)
					db.Second.Current = []Data{}
				}
			}
		case <-q.ctx.Done():
			q.ticker.Stop()
			return
		}
	}
}

func (q *Queue) TickerCli() {
	defer q.wg.Done()
	for {
		select {
		case <-q.ticker.C:
			countCurrent := len(db.All.Current)
			if countCurrent != 0 {
				// calculates the average of all data accumulated within one second
				var sum float64 = 0
				for i := 0; i < countCurrent; i++ {
					sum += db.All.Current[i].Value
				}
				SecValue := sum / float64(countCurrent)

				// and adds it to the "Second.Current"
				db.Second.Backlog = append(db.Second.Backlog, Data{
					Timestamp: time.Now(),
					Value:     SecValue,
				})

				lenSecondBacklog := len(db.Second.Backlog)
				if lenSecondBacklog > 1 {
					gui.SetDeltaValue(gui.SecValueDelta, db.Second.Backlog[lenSecondBacklog-2].Value, SecValue)
				}

				// All "Current" data is pushed into the "Backlog" for later use
				db.All.Backlog = append(db.All.Backlog, db.All.Current...)
				db.All.Current = []Data{}

				countSecond := len(db.Second.Backlog)
				if countSecond >= 60 { // one minute has passed and we can start calculating the average in the last minute
					last60Values := db.Second.Backlog[countSecond-60:]
					var sum float64 = 0
					for i := 0; i < 60; i++ {
						sum += last60Values[i].Value
					}
					MinuteValue := sum / float64(len(last60Values))

					// and adds it to the "Minute.Current"
					db.Minute.Current = append(db.Minute.Current, Data{
						Timestamp: time.Now(),
						Value:     MinuteValue,
					})

					lenMinuteCurrent := len(db.Minute.Current)
					if lenMinuteCurrent > 1 {
						gui.SetDeltaValue(gui.MinuteValueDelta, db.Minute.Current[lenMinuteCurrent-2].Value, MinuteValue)
					}

					db.Second.Backlog = append(db.Second.Backlog, db.Second.Current...)
					db.Second.Current = []Data{}
				}
			}
		case <-q.ctx.Done():
			q.ticker.Stop()
			return
		}
	}
}

func (q *Queue) DataProcessing() {
	for sysLatResult := range q.sendch {
		// log.Println("sysLatResult", sysLatResult, c.DetectLight)

		if !c.DetectLight {
			switch sysLatResult {
			case -1:
				sysLatResult = -2
			case -2:
				sysLatResult = -1
			}
		}

		switch sysLatResult {
		case -1: // White
			if targetOGL.IsActive {
				targetOGL.SetWhite()
			}
			if targetD3D9.IsActive {
				targetD3D9.SetWhite()
			}
			if !cliMode {
				RTSSClient.UpdateOSD(RTSSOSDWhite)
			} else if f.RTSS {
				RTSSClient.UpdateOSD(RTSSOSDWhite)
			}
		case -2: // Black
			if targetOGL.IsActive {
				targetOGL.SetBlack()
			}
			if targetD3D9.IsActive {
				targetD3D9.SetBlack()
			}
			if !cliMode {
				RTSSClient.UpdateOSD(RTSSOSDBlack)
			} else if f.RTSS {
				RTSSClient.UpdateOSD(RTSSOSDBlack)
			}

		default: // Data in µs
			// fmt.Println(sysLatResult)

			timestamp := time.Now()
			valueFloat := float64(sysLatResult) / 1000

			if cliMode {
				db.Count++
				if db.Countdown.IsZero() {
					db.Countdown = timestamp
					queue.ticker.Reset(time.Second)
				}

				db.All.Current = append(db.All.Current, Data{
					Timestamp: timestamp,
					Value:     valueFloat,
				})

				if f.Print {
					db.e.Add(valueFloat)
					// https://docs.microsoft.com/en-us/windows/console/console-virtual-terminal-sequences
					// fmt.Printf("\033c%f", valueFloat)
					fmt.Printf("\033c%f", db.e.Value())
				}

				SinceCountdown := time.Since(db.Countdown)
				if f.Count != -1 && f.Count < db.Count || time.Duration(0) != f.Delay && SinceCountdown > f.Delay {
					dbBacklog = append(dbBacklog, db)
					var dbTemp = []Database{db}
					dbt := NewDatabaseTableModel(&dbTemp)
					for _, values := range dbt.items {
						if len(values.Data) != 0 {
							SaveToFile(values)
							os.Exit(0)
						}
					}
				}

			} else {
				oldEwmaValue := db.e.Value()
				db.e.Add(valueFloat)
				newEwmaValue := db.e.Value()
				gui.SetValue(gui.EwmaValue, float64ToString(Round(newEwmaValue)))
				gui.SetDeltaValue(gui.EwmaValueDelta, oldEwmaValue, newEwmaValue)

				db.Count++
				if db.Countdown.IsZero() {
					gui.ResetGUI()
					db.Countdown = timestamp
					queue.ticker.Reset(time.Second)
					if gui.FileName != nil && gui.FileName.Text() == "" {
						gui.Synchronize(func() {
							gui.FileName.SetText(timestamp.Format("2006-01-02T15-04-05"))
						})
					}
				}

				gui.SetValue(gui.Count, IntToString(db.Count))

				db.All.Current = append(db.All.Current, Data{
					Timestamp: timestamp,
					Value:     valueFloat,
				})

				if valueFloat > db.All.Max {
					db.All.Max = valueFloat
					gui.SetValue(gui.ValueMax, float64ToString(valueFloat))
				}

				if valueFloat < db.All.Min || db.All.Min == 0 {
					db.All.Min = valueFloat
					gui.SetValue(gui.ValueMin, float64ToString(valueFloat))
				}

				gui.SetValue(gui.Value, float64ToString(valueFloat))
				SinceCountdown := time.Since(db.Countdown)
				gui.SetValue(gui.Countdown, SinceCountdown.String())
				if cliMode {
					if f.Count >= db.Count || SinceCountdown >= f.Delay {
						dbBacklog = append(dbBacklog, db)
						var dbTemp = []Database{db}
						dbt := NewDatabaseTableModel(&dbTemp)
						for _, values := range dbt.items {
							if len(values.Data) != 0 {
								SaveToFile(values)
							}
						}
					}
				}
			}
		}
	}
}
