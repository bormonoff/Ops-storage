package main

import (
	"flag"
	"ops-storage/internal/server/logger"
	"os"
)

type options struct {
	endpoint string
}

func Parse(opts *options) {
	flag.StringVar(&opts.endpoint, "a", "localhost:8080", "server address. localhost:8080 by default")
	flag.Parse()

	if endpoint := os.Getenv("ADDRESS"); endpoint != "" {
		opts.endpoint = endpoint
	}

	logger.MainLog.Infof("Options sucessfully parsed")
	logger.MainLog.Infof("Options.endpoint = %s", opts.endpoint)
}
