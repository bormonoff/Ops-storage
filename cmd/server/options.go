package main

import (
	"flag"
	"ops-storage/internal/server/logger"
	"os"
	"strconv"
)

type options struct {
	endpoint string

	// recover options
	storeInterval int
	filePath      string
	restore       bool
}

func Parse(opts *options) {
	flag.StringVar(&opts.endpoint, "a", "localhost:8080", "server address. localhost:8080 by default")
	flag.IntVar(&opts.storeInterval, "i", 3, "store interval in seconds")
	flag.StringVar(&opts.filePath, "f", "/tmp/metrics-db.json", "path to a dump file")
	flag.BoolVar(&opts.restore, "r", true, "Should server load data from the previous process")
	flag.Parse()

	if endpoint := os.Getenv("ADDRESS"); endpoint != "" {
		opts.endpoint = endpoint
	}

	if interval := os.Getenv("STORE_INTERVAL"); interval != "" {
		interval, err := strconv.Atoi(interval)
		if err != nil {
			logger.MainLog.Warnf("Can't parse STORE_INTERVAL env. Use default interval. Err: %s", err.Error())
		}
		opts.storeInterval = interval
	}

	if path := os.Getenv("FILE_STORAGE_PATH"); path != "" {
		opts.filePath = path
	}

	if restore := os.Getenv("RESTORE"); restore != "" {
		restore, err := strconv.ParseBool(restore)
		if err != nil {
			logger.MainLog.Warnf("Can't parse RESTORE env. Use default restore flag. Err: %s", err.Error())
		}
		opts.restore = restore
	}

	logger.MainLog.Infof("Options sucessfully parsed")
	logger.MainLog.Infof("Options.endpoint = %s", opts.endpoint)
	logger.MainLog.Infof("Options.storeInterval = %d", opts.storeInterval)
	logger.MainLog.Infof("Options.filePath = %s", opts.filePath)
	logger.MainLog.Infof("Options.restore = %t", opts.restore)
}
