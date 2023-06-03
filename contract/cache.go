package contract

import (
	"time"
)

type Cache interface {
	Initialize(mode Mode, limit uint, storage Storage) Cache
	Get(key string) (interface{}, error)
	Exist(key string) bool
	Add(key string, val interface{}, ttl time.Duration) error
	TTL(key string, ttl time.Duration) error
	Update(key string, val interface{}) error
	Delete(key string) error
	Keys() []string
	List() map[interface{}]interface{}
	Refresh()
	Flush()
}
