package error

import (
	"errors"
)

var (
	ErrIvalidMetric = errors.New("invalid metric type or value")
	ErrNotFound     = errors.New("metric isn't found")
	ErrInternal     = errors.New("internal error")
)
