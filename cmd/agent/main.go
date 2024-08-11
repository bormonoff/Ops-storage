package main

import (
	"ops-storage/internal/agent/app"
)

func main() {
	opts := options{}
	Parse(&opts)

	config := app.Config(app.Config{
		PollInterval:   opts.pollInterval,
		ReportInterval: opts.reportInterval,
	})
	config.SetAddr(opts.endpoint)

	app := app.NewApp(config)

	app.Run()
}
