package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"path/filepath"
	"strconv"

	"github.com/VividCortex/ewma"
	"github.com/eclesh/welford"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/go-echarts/go-echarts/v2/types"
)

func SaveToFile(dbt *DatabaseTable) {
	stats := welford.New()

	// deletes the first high values in descending order
	i := 0
	var lastValue float64
	for ; i < len(dbt.Data); i++ {
		if lastValue == 0.0 || lastValue > dbt.Data[i].Value {
			lastValue = dbt.Data[i].Value
		} else {
			break
		}
	}

	var filename = dbt.Data[0].Timestamp.Format("2006-01-02T15-04-05")
	var projectName string
	if cliMode {
		// filename = dbt.Data[0].Timestamp.Format("2006-01-02T15-04-05")
		if f.Name != "" {
			filename = f.Name + "_" + filename
			projectName = f.Name
		} else {
			projectName = filename
		}
	} else {
		guiFileName := gui.FileName.Text()
		projectName = guiFileName
		// filename = dbt.Data[0].Timestamp.Format("2006-01-02T15-04-05")
		if guiFileName != "" && guiFileName != filename {
			filename = gui.FileName.Text() + "_" + filename
		}
	}
	// CSV
	csvFile, err := os.Create(filepath.Join(f.LogFolder, filename+".csv"))
	if err != nil {
		log.Println(err)
	}
	_, err = csvFile.Write(append([]byte("time,latency"), []byte{13, 10}...))
	if err != nil {
		log.Println(err)
	}

	// https://en.wikipedia.org/wiki/Moving_average#Exponential_moving_average
	SimpleEWMA := ewma.NewMovingAverage()

	var highestLatencyRound5 int
	round5Map := make(map[int]int)
	var XValuesLineChartTime = []string{}
	var YValuesLineChartTime = []opts.LineData{}
	var benchmarkTime float64
	var rawData []float64
	for _, v := range dbt.Data {
		// CSV
		temp := []byte(v.Timestamp.Format("2006-01-02T15:04:05.999999999") + ",")
		temp = append(temp, []byte(float64ToString(v.Value))...) // in ms
		temp = append(temp, []byte{13, 10}...)                   // \r\n
		_, err := csvFile.Write(temp)
		if err != nil {
			log.Println(err)
		}

		// summary
		rawData = append(rawData, v.Value)
		benchmarkTime += v.Value
		SimpleEWMA.Add(v.Value)
		stats.Add(v.Value)

		// chart, bar
		r5 := int(round5(v.Value))
		val, ok := round5Map[r5]
		if ok {
			round5Map[r5] = val + 1
		} else {
			round5Map[r5] = 1
			if highestLatencyRound5 < r5 {
				highestLatencyRound5 = r5
			}
		}

		// chart, line
		XValuesLineChartTime = append(XValuesLineChartTime, v.Timestamp.Sub(dbt.Data[0].Timestamp).String())
		YValuesLineChartTime = append(YValuesLineChartTime, opts.LineData{
			Value: v.Value,
		})
	}
	csvFile.Close()

	var newFile bool
	if _, err := os.Stat(filepath.Join(f.LogFolder, "summary.csv")); errors.Is(err, os.ErrNotExist) {
		newFile = true
	}

	header := "name,start,duration,count,max,min,average,ewma average,stdev\n"
	summary := fmt.Sprintf(`%s,%s,%s,%d,%.3f,%.3f,%.3f,%.3f,%.3f`,
		projectName,
		dbt.Data[0].Timestamp.Format("2006-01-02T15:04:05.999999999"),
		dbt.Data[len(dbt.Data)-1].Timestamp.Sub(dbt.Data[0].Timestamp).String(),
		len(dbt.Data),
		stats.Max(),
		stats.Min(),
		stats.Mean(),
		SimpleEWMA.Value(),
		stats.Stddev(),
	)

	summaryFile, err := os.OpenFile(filepath.Join(f.LogFolder, "summary.csv"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		log.Println(err)
	}

	if newFile {
		if _, err := summaryFile.WriteString(header); err != nil {
			log.Println(err)
		}
	}

	if _, err := summaryFile.WriteString(summary + "\n"); err != nil {
		log.Println(err)
	}
	summaryFile.Close()

	// Chart
	chartXValues := []string{}
	for i := 5; i <= highestLatencyRound5; i += 5 {
		if i == 5 {
			chartXValues = append(chartXValues, "5ms")
		} else {
			chartXValues = append(chartXValues, strconv.Itoa(i))
		}
	}

	BarItems := make([]opts.BarData, highestLatencyRound5/5+1)
	for i := 1; i < len(BarItems); i++ {
		val, ok := round5Map[i*5]
		if !ok {
			BarItems[i-1] = opts.BarData{
				Value: 0,
			}
		} else {
			BarItems[i-1] = opts.BarData{
				Value: val,
			}
		}
	}

	page := components.NewPage().
		SetLayout(components.PageFlexLayout).
		AddCharts(
			createBarChartRound5(
				filename,
				fmt.Sprintf("start: %s, count: %d, max: %.3f, min: %.3f, average: %.3f, ewma average: %.3f", dbt.Data[0].Timestamp.Format("2006-01-02T15:04:05"), stats.Count(), stats.Max(), stats.Min(), stats.Mean(), SimpleEWMA.Value()),
				chartXValues,
				BarItems,
			),
			createLineChartTime(
				"time series",
				fmt.Sprintf("duration: %s, max: %.3f, min: %.3f, average: %.3f, ewma average: %.3f", dbt.Data[len(dbt.Data)-1].Timestamp.Sub(dbt.Data[0].Timestamp).String(), stats.Max(), stats.Min(), stats.Mean(), SimpleEWMA.Value()),
				XValuesLineChartTime,
				YValuesLineChartTime,
			),
		)

	chartFile, err := os.Create(filepath.Join(f.LogFolder, filename+".html"))
	if err != nil {
		log.Println(err)
	}
	page.Render(io.MultiWriter(chartFile))
	chartFile.Close()
}

func createBarChartRound5(title, subtitle string, xValues []string, yValues []opts.BarData) *charts.Bar {
	bar := charts.NewBar()
	bar.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{
			PageTitle: "GoSysLat",
			Theme:     types.ThemeWesteros,
		}),
		charts.WithTitleOpts(opts.Title{
			Title:    title,
			Subtitle: subtitle,
		}),
		charts.WithDataZoomOpts(opts.DataZoom{
			Type:       "inside",
			Start:      0,
			End:        100,
			XAxisIndex: []int{0},
		}),
		charts.WithDataZoomOpts(opts.DataZoom{
			Type:       "slider",
			Start:      0,
			End:        100,
			XAxisIndex: []int{0},
		}),
	)
	bar.SetXAxis(xValues).
		AddSeries("Data", yValues)

	return bar
}

func createLineChartTime(title, subtitle string, xValues []string, yValues []opts.LineData) *charts.Line {
	line := charts.NewLine()
	line.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{
			Theme: types.ThemeWesteros,
		}),
		charts.WithTitleOpts(opts.Title{
			Title:    title,
			Subtitle: subtitle,
		}),
		charts.WithDataZoomOpts(opts.DataZoom{
			Type:       "inside",
			Start:      0,
			End:        100,
			XAxisIndex: []int{0},
		}),
		charts.WithDataZoomOpts(opts.DataZoom{
			Type:       "slider",
			Start:      0,
			End:        100,
			XAxisIndex: []int{0},
		}),
	)

	line.SetXAxis(xValues).
		AddSeries("Data", yValues).
		SetSeriesOptions(
			charts.WithLabelOpts(
				opts.Label{
					Show: true,
				},
			),
			charts.WithAreaStyleOpts(
				opts.AreaStyle{
					Opacity: 0.2,
				},
			),
			charts.WithLineChartOpts(
				opts.LineChart{
					Smooth: true,
				},
			),
		)

	return line
}

func round5(x float64) float64 {
	// x = 1 => 5
	// x = 6 => 10
	return math.Ceil(x/5) * 5
}

func round(val float64) float64 {
	return math.Round(val*100) / 100
}

func roundInt(val float64) int {
	return int(math.Round(val))
}
