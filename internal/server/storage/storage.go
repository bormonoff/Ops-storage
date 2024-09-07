package storage

import (
	"errors"
	"sync"
)

type (
	incrementalCounter int64
	gaugeCounter       float64
)

type storage struct {
	IncrementalCounters map[string]incrementalCounter `json:"incrementalCounters,omitempty"`
	GaugeCounters       map[string]gaugeCounter       `json:"gaugeCounters,omitempty"`

	recFilePath string

	mu sync.RWMutex
}

var (
	store = New()

	ErrIvalidMetric = errors.New("invalid metric type or value")
	ErrNotFound     = errors.New("metric isn't found")
)

func New() storage {
	return storage{
		IncrementalCounters: make(map[string]incrementalCounter),
		GaugeCounters:       make(map[string]gaugeCounter),
	}
}

func StorageInstace() BaseStorage {
	return &store
}
