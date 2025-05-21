package app

import (
	"encoding/json"
	"log"
	"runtime"
	"sort"
	"strconv"
	"time"

	"github.com/NimbleMarkets/ntcharts/linechart/timeserieslinechart"
	"github.com/charmbracelet/lipgloss"
	"github.com/noodahl-org/moodstress/app/models"
	"github.com/samber/lo"
)

var green = lipgloss.NewStyle().Foreground(lipgloss.Color("10"))
var red = lipgloss.NewStyle().Foreground(lipgloss.Color("9"))

func (a *App) Graph() {
	iter, err := a.db.NewIter(nil)
	if err != nil {
		_, _, line, _ := runtime.Caller(1)
		log.Fatalf("%s %v", err, line)
	}
	defer iter.Close()

	var dataPoints []timeserieslinechart.TimePoint

	for iter.First(); iter.Valid(); iter.Next() {
		key := iter.Key()
		log.Printf("key: %s", key)

		results := []models.IntradayMetric{}
		if err := json.Unmarshal(iter.Value(), &results); err != nil {
			_, _, line, _ := runtime.Caller(1)
			log.Printf("%s %v", err, line)
			continue
		}

		for _, result := range results {
			t := time.Unix(result.Time, 0)
			score := lo.SumBy(result.Metrics, func(metric models.Metric) int {
				val, _ := strconv.Atoi(metric.Value)
				if metric.Charge {
					return val
				} else {
					return val * -1
				}
			})

			value := float64(score)
			dataPoints = append(dataPoints, timeserieslinechart.TimePoint{
				Time:  t,
				Value: value,
			})
		}
	}

	if err := iter.Error(); err != nil {
		_, _, line, _ := runtime.Caller(1)
		log.Fatalf("%s %v", err, line)
	}

	sort.Slice(dataPoints, func(i, j int) bool {
		return dataPoints[i].Time.Before(dataPoints[j].Time)
	})

	a.timeseries = timeserieslinechart.New(50, 8,
		timeserieslinechart.WithStyle(red),
		timeserieslinechart.WithDataSetTimeSeries("health score", dataPoints),
	)

	// Draw the chart
	a.timeseries.DrawBrailleAll()
	a.timeseries.DrawXYAxisAndLabel()

	a.showBarChart = false
	a.showForm = false
	a.showTable = false
	a.showTimeseriesChart = true
}
