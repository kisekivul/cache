package contract

import (
	"errors"
)

var (
	ErrKeyNotFound     = errors.New("key not found")
	ErrStrategyNotSet  = errors.New("strategy not set")
	ErrEvictItemFailed = errors.New("evict item failed")
)
