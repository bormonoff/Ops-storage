package app

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"math/rand"
	"strconv"
	"sync"
	"time"

	"ops-storage/internal/agent/collector"
	"ops-storage/internal/agent/handlers"
	"ops-storage/internal/agent/logger"
)

type Config struct {
	serverAddr     string
	PollInterval   int
	ReportInterval int

	HasCompression bool
}

func (c *Config) SetAddr(url string) {
	c.serverAddr = fmt.Sprintf("http://%v", url)
}

type app struct {
	mu        sync.Mutex
	collector collector.Collection
	config    Config
}

func New(config Config) *app {
	app := app{
		collector: collector.NewCollection(),
		config:    config,
	}
	return &app
}

func (app *app) updateData() {
	for {
		app.mu.Lock()
		app.collector.RefreshStats()
		app.mu.Unlock()

		time.Sleep(time.Duration(app.config.PollInterval) * time.Second)
	}
}

type updateJSONValidator struct {
	MType   string      `json:"type"`
	Name    string      `json:"id"`
	Counter json.Number `json:"delta,omitempty"`
	Gauge   json.Number `json:"value,omitempty"`
}

func (v updateJSONValidator) String() string {
	return fmt.Sprint(v.MType, v.Name, v.Counter, v.Gauge)
}

func updateCounters(c *collector.Collection) []updateJSONValidator {
	// +2 is becouse collector has pollcount and randomval additional fields
	totalUpdateCount := 2 + len(c.RuntimeStats.UintStats) + len(c.RuntimeStats.FloatStats)
	var toUpdate = make([]updateJSONValidator, totalUpdateCount)

	idx := -1
	for id, val := range c.RuntimeStats.UintStats {
		idx++
		toUpdate[idx] = updateJSONValidator{
			MType: "gauge",
			Name:  string(id),
			Gauge: json.Number(strconv.FormatUint(val, 10)),
		}
	}
	for id, val := range c.RuntimeStats.FloatStats {
		idx++
		toUpdate[idx] = updateJSONValidator{
			MType: "gauge",
			Name:  string(id),
			Gauge: json.Number(strconv.FormatFloat(val, 'f', 1, 64)),
		}
	}

	idx++
	toUpdate[idx] = updateJSONValidator{
		MType: "counter",
		Name:  "PollCount",
		Gauge: json.Number(strconv.Itoa(c.PollCount)),
	}
	idx++
	toUpdate[idx] = updateJSONValidator{
		MType: "gauge",
		Name:  "RandomValue",
		Gauge: json.Number(strconv.Itoa(int(rand.Int()))),
	}

	return toUpdate
}

func (app *app) sendData() {
	headers := map[string]string{
		"Content-Type": "application/json",
	}

	if app.config.HasCompression {
		headers["Content-Encoding"] = "gzip"
	}

	url := fmt.Sprintf("%v/update", app.config.serverAddr)
	for {
		app.mu.Lock()

		counters := updateCounters(&app.collector)
		app.mu.Unlock()

		for _, tmp := range counters {
			body, err := json.Marshal(tmp)
			if err != nil {
				logger.Log.Errorf(err.Error())
				continue
			}

			if app.config.HasCompression {
				var buf bytes.Buffer
				encoder, err := gzip.NewWriterLevel(&buf, gzip.BestSpeed)
				if err != nil {
					logger.Log.Errorf("Compressor isn't initialized: %s", err.Error())
					continue
				}

				_, err = encoder.Write(body)
				if err != nil {
					logger.Log.Errorf("Can't wrine body: %s", err.Error())
					continue
				}
				encoder.Close()
				body = buf.Bytes()
			}
			err = handlers.SendPostRequest(url, headers, body)
			if err != nil {
				logger.Log.Errorf("Can't update metric: %s", err.Error())
			}
			logger.Log.Infof("Metric %s has been successfully updated", tmp.Name)
		}
		time.Sleep(time.Duration(app.config.ReportInterval) * time.Second)
	}
}

func (app *app) Run() {
	wg := new(sync.WaitGroup)
	wg.Add(2)

	go app.updateData()
	go app.sendData()

	wg.Wait()
}
