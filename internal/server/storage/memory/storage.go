package memory

import (
	"sync"
)

type (
	incrementalCounter int64
	gaugeCounter       float64
)

type memStorage struct {
	IncrementalCounters map[string]incrementalCounter `json:"incrementalCounters,omitempty"`
	GaugeCounters       map[string]gaugeCounter       `json:"gaugeCounters,omitempty"`

	recFilePath string

	mu sync.RWMutex
}

var (
	store = New()
)

func New() memStorage {
	return memStorage{
		IncrementalCounters: make(map[string]incrementalCounter),
		GaugeCounters:       make(map[string]gaugeCounter),
	}
}

func Instance() *memStorage {
	return &store
}
