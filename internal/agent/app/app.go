package app

import (
	"fmt"
	"math/rand"
	"strconv"
	"sync"
	"time"

	"ops-storage/internal/agent/collector"
	"ops-storage/internal/agent/handlers"
)

type Config struct {
	serverAddr     string
	PollInterval   int
	ReportInterval int
}

func (c *Config) SetAddr(url string) {
	c.serverAddr = fmt.Sprintf("http://%v", url)
}

type app struct {
	mu        sync.Mutex
	collector collector.Collection
	config    Config
}

func NewApp(config Config) *app {
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

func (app *app) SendData() {
	headers := map[string]string{
		"Content-Type": "text/plain",
	}
	for {
		app.mu.Lock()

		for id, val := range app.collector.RuntimeStats.UintStats {
			url := fmt.Sprintf("%v/update/gauge/%v/%v",
				app.config.serverAddr, string(id), strconv.FormatUint(val, 10))
			handlers.SendPostRequest(url, headers)
		}

		for id, val := range app.collector.RuntimeStats.FloatStats {
			url := fmt.Sprintf("%v/update/gauge/%v/%v",
				app.config.serverAddr, string(id), strconv.FormatFloat(val, 'f', 1, 64))
			handlers.SendPostRequest(url, headers)
		}

		pollUrl := fmt.Sprintf("%v/update/counter/PollCount/%v",
			app.config.serverAddr, strconv.Itoa(app.collector.PollCount))
		handlers.SendPostRequest(pollUrl, headers)

		randValue := fmt.Sprintf("%v/update/gauge/RandomValue/%v",
			app.config.serverAddr, strconv.Itoa(int(rand.Int())))
		handlers.SendPostRequest(randValue, headers)

		app.mu.Unlock()

		time.Sleep(time.Duration(app.config.ReportInterval) * time.Second)
	}
}

func (app *app) Run() {
	wg := new(sync.WaitGroup)
	wg.Add(2)

	go app.updateData()
	go app.SendData()

	wg.Wait()
}
