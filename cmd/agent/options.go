package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
)

type options struct {
	endpoint       string
	reportInterval int
	pollInterval   int
	compress       bool
}

func Parse(opts *options) {
	flag.StringVar(&opts.endpoint, "a", "localhost:8080", "server address. localhost:8080 by default")
	flag.IntVar(&opts.reportInterval, "r", 10, "report interval, s. 10 seconds by default")
	flag.IntVar(&opts.pollInterval, "p", 2, "poll interval, s. 2 seconds by default")
	flag.BoolVar(&opts.compress, "c", true, "turn on/off the compression")
	flag.Parse()

	parseEnv(opts)
}

func parseEnv(opts *options) {
	if endpoint := os.Getenv("ADDRESS"); endpoint != "" {
		opts.endpoint = endpoint
	}
	if interval := os.Getenv("REPORT_INTERVAL"); interval != "" {
		interval, err := strconv.Atoi(interval)
		if err != nil {
			fmt.Println("REPORT_INTERVAL environment var is not int")
		} else {
			opts.reportInterval = interval
		}
	}
	if interval := os.Getenv("POLL_INTERVAL"); interval != "" {
		interval, err := strconv.Atoi(interval)
		if err != nil {
			fmt.Println("POLL_INTERVAL environment var is not int")
		} else {
			opts.pollInterval = interval
		}
	}
}
