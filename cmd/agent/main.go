package main

import (
	"ops-storage/internal/agent/app"
	"ops-storage/internal/agent/logger"
)

func main() {
	logger.Initialize()
	logger.Log.Info("The logger has been successfully initialized")

	opts := options{}
	Parse(&opts)

	config := app.Config(app.Config{
		PollInterval:   opts.pollInterval,
		ReportInterval: opts.reportInterval,
		HasCompression: opts.compress,
	})
	config.SetAddr(opts.endpoint)

	app := app.New(config)

	app.Run()
}
