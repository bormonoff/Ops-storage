package main

import (
	"flag"
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
}
