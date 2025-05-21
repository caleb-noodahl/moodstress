package models

import (
	"encoding/json"
	"fmt"
	"log"
	"runtime"
	"time"

	"github.com/cockroachdb/pebble"
)

type Metric struct {
	Name    string `json:"name"`
	Rune    rune   `json:"rune"`
	Value   string `json:"value"`
	Charge  bool   `json:"charge"`
	Anchors string `json:"anchors"`
}

type IntradayMetric struct {
	Time    int64    `json:"time"`
	Metrics []Metric `json:"metrics"`
	Note    string   `json:"notes"`
}

type BinaryMetric struct {
	Key     string
	Metrics []Metric `json:"metrics"`
}

func FetchIntradayMetrics(key string, db *pebble.DB) ([]IntradayMetric, error) {
	log.Printf("fetch_intraday_metrics: %s", key)
	results := []IntradayMetric{}
	record, closer, err := db.Get([]byte(key))
	if err != nil {
		log.Printf("err:%s", err)
		if err == pebble.ErrNotFound {
			record = []byte("[]")
		} else {
			_, _, line, _ := runtime.Caller(1)
			return nil, fmt.Errorf("%s %v", err, line)
		}
	} else {
		defer closer.Close()
	}
	if err := json.Unmarshal(record, &results); err != nil {
		return nil, err
	}
	return results, nil

}

func NewSCHEMAP() *IntradayMetric {
	return &IntradayMetric{
		Time: time.Now().UTC().Unix(),
		Metrics: []Metric{
			{
				Name:    "stress",
				Rune:    rune('S'),
				Value:   "",
				Charge:  false,
				Anchors: "1=calm 2=mild 3=moderate 4=high 5=severe",
			},
			{
				Name:    "clarity",
				Rune:    rune('C'),
				Value:   "",
				Charge:  true,
				Anchors: "1=foggy 2=unclear 3=neutral 4=clear 5=sharp",
			},
			{
				Name:    "hormonal",
				Rune:    rune('H'),
				Value:   "",
				Charge:  false,
				Anchors: "1=calm 2=mild 3=moderate 4=high 5=severe",
			},
			{
				Name:    "energy",
				Rune:    rune('E'),
				Value:   "",
				Charge:  true,
				Anchors: "1=exhausted 2=tired 3=moderate 4=good 5=high",
			},
			{
				Name:    "mood",
				Rune:    rune('M'),
				Value:   "",
				Charge:  true,
				Anchors: "1=depressed 2=sad 3=flat 4=positive 5=excellent",
			},
			{
				Name:    "attention",
				Rune:    rune('A'),
				Value:   "",
				Charge:  true,
				Anchors: "1=none 2=distracted 3=moderate 4=good 5=laser",
			},
			{
				Name:    "pain",
				Rune:    rune('P'),
				Value:   "",
				Charge:  false,
				Anchors: "1=none 2=mild 3=moderate 4=high 5=severe",
			},
		},
	}
}
