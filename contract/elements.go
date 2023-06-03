package contract

import (
	"sync"
)

type Elements struct {
	sync.RWMutex
	dict map[string]interface{}
}

func NewElements(data map[string]interface{}) *Elements {
	var (
		elements = &Elements{
			dict: make(map[string]interface{}, 0),
		}
	)

	for k, v := range data {
		elements.dict[k] = v
	}
	return elements
}

func (es *Elements) Get(key string) interface{} {
	es.RLock()
	defer es.RUnlock()

	return es.dict[key]
}

func (es *Elements) Set(key string, value interface{}) *Elements {
	es.Lock()
	defer es.Unlock()

	es.dict[key] = value
	return es
}
