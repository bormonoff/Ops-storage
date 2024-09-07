package storage

import (
	"strconv"
)

type BaseStorage interface {
	Insert(counterType string, name string, val string) error
	GetMetric(counterType string, name string) (string, error)
	GetAllMetrics() *map[string]string
}

func (s *storage) GetAllMetrics() *map[string]string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	res := make(map[string]string)

	for id, val := range s.IncrementalCounters {
		res[id] = strconv.FormatInt(int64(val), 10)
	}
	for id, val := range s.GaugeCounters {
		res[id] = strconv.FormatFloat(float64(val), 'f', -1, 64)
	}
	s.storeToFile()
	return &res
}

func (s *storage) Insert(counterType string, name string, val string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	switch counterType {
	case "gauge":
		newVal, err := strconv.ParseFloat(val, 64)
		if err != nil {
			return ErrIvalidMetric
		}
		s.GaugeCounters[name] = gaugeCounter(newVal)
	case "counter":
		newVal, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			return ErrIvalidMetric
		}

		_, ok := s.IncrementalCounters[name]
		if ok {
			s.IncrementalCounters[name] += incrementalCounter(newVal)
		} else {
			s.IncrementalCounters[name] = incrementalCounter(newVal)
		}
	default:
		return ErrIvalidMetric
	}

	return nil
}

func (s *storage) GetMetric(counterType string, name string) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result string

	switch counterType {
	case "gauge":
		val, ok := s.GaugeCounters[name]
		if !ok {
			return "", ErrNotFound
		}
		result = strconv.FormatFloat(float64(val), 'f', -1, 64)
	case "counter":
		val, ok := s.IncrementalCounters[name]
		if !ok {
			return "", ErrNotFound
		}
		result = strconv.FormatInt(int64(val), 10)
	default:
		return "", ErrNotFound
	}

	return result, nil
}
