package core

import (
	"errors"
	"strconv"
)

type (
	incrementalCounter int64
	gaugeCounter       float64
)

type Storage struct {
	incrementalCounters map[string]incrementalCounter
	gaugeCounters       map[string]gaugeCounter
}

func createNewStorage() Storage {
	return Storage{
		incrementalCounters: make(map[string]incrementalCounter),
		gaugeCounters:       make(map[string]gaugeCounter),
	}
}

var (
	storage = createNewStorage()

	ErrIvalidMetric = errors.New("invalid metric type or value")
	ErrNotFound     = errors.New("metric isn't found")
)

type BaseStorage interface {
	Insert(counterType string, name string, val string) error
	GetMetric(counterType string, name string) (string, error)
	GetActualMetrics() *map[string]string
}

func GetStorageInstace() BaseStorage {
	return &storage
}

func (s *Storage) GetActualMetrics() *map[string]string {
	res := make(map[string]string)

	for id, val := range s.incrementalCounters {
		res[id] = strconv.FormatInt(int64(val), 10)
	}
	for id, val := range s.gaugeCounters {
		res[id] = strconv.FormatFloat(float64(val), 'f', -1, 64)
	}
	return &res
}

func (s *Storage) Insert(counterType string, name string, val string) error {
	switch counterType {
	case "gauge":
		newVal, err := strconv.ParseFloat(val, 64)
		if err != nil {
			return ErrIvalidMetric
		}
		s.gaugeCounters[name] = gaugeCounter(newVal)
	case "counter":
		newVal, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			return ErrIvalidMetric
		}

		_, ok := s.incrementalCounters[name]
		if ok {
			s.incrementalCounters[name] += incrementalCounter(newVal)
		} else {
			s.incrementalCounters[name] = incrementalCounter(newVal)
		}
	default:
		return ErrIvalidMetric
	}

	return nil
}

func (s *Storage) GetMetric(counterType string, name string) (string, error) {
	var result string

	switch counterType {
	case "gauge":
		val, ok := s.gaugeCounters[name]
		if !ok {
			return "", ErrNotFound
		}
		result = strconv.FormatFloat(float64(val), 'f', -1, 64)
	case "counter":
		val, ok := s.incrementalCounters[name]
		if !ok {
			return "", ErrNotFound
		}
		result = strconv.FormatInt(int64(val), 10)
	default:
		return "", ErrNotFound
	}

	return result, nil
}
