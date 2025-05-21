package app

import (
	"encoding/json"
	"fmt"
	"log"
	"runtime"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/cockroachdb/pebble"
	"github.com/noodahl-org/moodstress/app/models"
	"github.com/samber/lo"
)

type IntraLogView struct {
	intra   *models.IntradayMetric
	metrics []models.IntradayMetric
}

func NewIntraLogView() *IntraLogView {

	return &IntraLogView{
		intra:   models.NewSCHEMAP(),
		metrics: []models.IntradayMetric{},
	}
}

func (a *App) IntraLog() {
	year, month, day := time.Now().Date()
	key := fmt.Sprintf("intra_%v.%v.%v", month, day, year)

	inputs := lo.Map(a.intralogview.intra.Metrics, func(m models.Metric, i int) huh.Field {
		return huh.NewInput().Title(m.Name + " ").Description(m.Anchors).
			Value(&a.intralogview.intra.Metrics[i].Value)
	})

	var err error
	a.intralogview.metrics, err = models.FetchIntradayMetrics(key, a.db)
	if err != nil {
		log.Fatal(err)
	}

	a.form = huh.NewForm(
		huh.NewGroup(
			inputs...,
		).Title("new intraday log item"),
	)

	a.form.SubmitCmd = func() tea.Msg {
		year, month, day := time.Now().Date()
		key := fmt.Sprintf("intra_%v.%v.%v", month, day, year)

		a.intralogview.metrics = append(a.intralogview.metrics, *a.intralogview.intra)
		data, err := json.Marshal(a.intralogview.metrics)
		if err != nil {
			_, _, line, _ := runtime.Caller(1)
			log.Fatalf("%s %v", err, line)
		}
		a.intralogview.intra.Time = time.Now().Unix()
		if err := a.db.Set([]byte(key), data, &pebble.WriteOptions{Sync: true}); err != nil {
			_, _, line, _ := runtime.Caller(1)
			log.Fatalf("%s %v", err, line)
		}

		a.Init()
		return nil
	}
}
