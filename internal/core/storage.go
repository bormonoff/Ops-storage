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

var (
	storage = Storage{
		incrementalCounters: make(map[string]incrementalCounter),
		gaugeCounters:       make(map[string]gaugeCounter),
	}

	ErrIvalidMetric = errors.New("invalid metric type or value")
)

func Update(counterType string, name string, val string) error {
	switch counterType {
	case "gauge":
		newVal, err := strconv.ParseFloat(val, 64)
		if err != nil {
			return ErrIvalidMetric
		}
		storage.gaugeCounters[name] = gaugeCounter(newVal)
	case "counter":
		newVal, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			return ErrIvalidMetric
		}

		_, ok := storage.incrementalCounters[name]
		if ok {
			storage.incrementalCounters[name] += incrementalCounter(newVal)
		} else {
			storage.incrementalCounters[name] = incrementalCounter(newVal)
		}
	default:
		return ErrIvalidMetric
	}


	return nil
}
