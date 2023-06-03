package contract

import (
	"time"
)

type Storage interface {
	Initialize(interval time.Duration, params interface{}) (Storage, error)
	Load() (map[string]*Item, error)
	Save(items map[string]*Item) error
	Exit()
}
