package memory

import (
	"strconv"

	serror "ops-storage/internal/server/storage/error"
)

func (s *memStorage) GetAll() (*map[string]string, error) {
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
	return &res, nil
}

func (s *memStorage) Insert(counterType string, name string, val string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	switch counterType {
	case "gauge":
		newVal, err := strconv.ParseFloat(val, 64)
		if err != nil {
			return serror.ErrIvalidMetric
		}
		s.GaugeCounters[name] = gaugeCounter(newVal)
	case "counter":
		newVal, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			return serror.ErrIvalidMetric
		}

		_, ok := s.IncrementalCounters[name]
		if ok {
			s.IncrementalCounters[name] += incrementalCounter(newVal)
		} else {
			s.IncrementalCounters[name] = incrementalCounter(newVal)
		}
	default:
		return serror.ErrIvalidMetric
	}

	return nil
}

func (s *memStorage) Get(counterType string, name string) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result string

	switch counterType {
	case "gauge":
		val, ok := s.GaugeCounters[name]
		if !ok {
			return "", serror.ErrNotFound
		}
		result = strconv.FormatFloat(float64(val), 'f', -1, 64)
	case "counter":
		val, ok := s.IncrementalCounters[name]
		if !ok {
			return "", serror.ErrNotFound
		}
		result = strconv.FormatInt(int64(val), 10)
	default:
		return "", serror.ErrNotFound
	}

	return result, nil
}

func (s *memStorage) IsStorageAlive() bool {
	if s.GaugeCounters != nil && s.IncrementalCounters != nil {
		return true
	}
	return false
}
