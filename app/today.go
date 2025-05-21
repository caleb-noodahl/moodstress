package app

import (
	"fmt"
	"log"
	"math"
	"strconv"
	"time"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
	"github.com/noodahl-org/moodstress/app/models"
	"github.com/samber/lo"
)

var blockStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("9")). // red
	Background(lipgloss.Color("9"))  // red

type TodayAggregate struct {
	total    int
	positive int
	negative int
	score    int
}

func (a *App) Today() {
	n := time.Now()
	year, month, day := n.Date()
	key := fmt.Sprintf("intra_%v.%v.%v", month, day, year)

	log.Printf("today: key:%s", key)
	columns := []table.Column{}
	rows := []table.Row{}
	agg := TodayAggregate{}
	metrics, err := models.FetchIntradayMetrics(key, a.db)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("metrics: %v", metrics)
	if len(metrics) > 0 {
		//list
		m := metrics[0]
		columns = append(columns, table.Column{Title: "time", Width: 20})
		columns = append(columns, lo.Map(m.Metrics, func(m models.Metric, _ int) table.Column {
			pos := lo.Ternary(m.Charge, "+", "-")
			return table.Column{Title: fmt.Sprintf("(%s%s)", string(m.Rune), pos), Width: 4}
		})...)
		rows = append(rows, lo.Map(metrics, func(m models.IntradayMetric, _ int) table.Row {
			t := time.Unix(m.Time, 0)
			result := []string{fmt.Sprintf("%v.%v.%v %v:00", t.Month(), t.Day(), t.Year(), t.Hour())}
			result = append(result, lo.Map(m.Metrics, func(m models.Metric, _ int) string {
				return m.Value
			})...)
			return result
		})...)
		//this can be done cleaner
		agg.total = lo.SumBy(metrics, func(m models.IntradayMetric) int {
			return lo.SumBy(m.Metrics, func(metric models.Metric) int {
				val, _ := strconv.Atoi(metric.Value)
				switch metric.Charge {
				case true:
					agg.positive += val
					return val
				case false:
					agg.negative += val
					return val * -1
				}
				return 0
			})
		})
		if agg.positive+agg.negative > 0 {
			percentage := float64(agg.positive) / float64(agg.positive+agg.negative)
			agg.score = int(math.Round(percentage * 100))
		} else {
			agg.score = 50 // neutral when no data
		}

		log.Printf("agg: %v", agg)
		rows = append(rows, []table.Row{
			{},
			{"positive", fmt.Sprintf("%v", agg.positive)},
			{"negative", fmt.Sprintf("%v", agg.negative)},
			{},
			{"score", fmt.Sprintf("%v", agg.score), "/", "100"},
		}...)

		a.showForm = false
		a.showTable = true
		a.table = table.New(
			table.WithFocused(true),
			table.WithColumns(columns),
			table.WithRows(rows),
		)
	}

}
