package contract

import (
	"encoding/json"
	"sync"
	"time"
)

type Item struct {
	sync.RWMutex
	key       string
	value     interface{}
	ttl       time.Duration
	created   time.Time
	updated   time.Time
	visited   time.Time
	frequency uint
	elements  map[Mode]*Elements
}

func NewItem(key string, value interface{}, ttl time.Duration) *Item {
	var (
		date = time.Now()
	)
	return &Item{
		key:       key,
		value:     value,
		ttl:       ttl,
		created:   date,
		updated:   date,
		visited:   date,
		frequency: 0,
		elements:  make(map[Mode]*Elements, 0),
	}
}

// Key returns the key of this cached item.
func (item *Item) Key() string {
	// immutable
	return item.key
}

// Value returns the value of this cached item.
func (item *Item) Value() interface{} {
	item.RLock()
	defer item.RUnlock()

	item.visited = time.Now()
	item.frequency++
	return item.value
}

// TTL returns this item's expiration duration.
func (item *Item) TTL() time.Duration {
	item.RLock()
	defer item.RUnlock()

	return item.ttl
}

// Created returns when this item was added to the cache.
func (item *Item) Created() time.Time {
	// immutable
	return item.created
}

// Updated returns when this item was last updated.
func (item *Item) Updated() time.Time {
	item.RLock()
	defer item.RUnlock()

	return item.updated
}

// Visited returns when this item was last visted.
func (item *Item) Visited() time.Time {
	item.RLock()
	defer item.RUnlock()

	return item.visited
}

// Frequency returns how often this item has been visted.
func (item *Item) Frequency() uint {
	item.RLock()
	defer item.RUnlock()

	return item.frequency
}

// Elements get item node in list.
func (item *Item) Elements(mode Mode) *Elements {
	item.RLock()
	defer item.RUnlock()

	return item.elements[mode]
}

func (item *Item) SetTTL(ttl time.Duration) *Item {
	item.Lock()
	defer item.Unlock()

	var (
		date = time.Now()
	)

	item.ttl = ttl
	item.updated = date
	item.frequency++
	return item
}

func (item *Item) SetValue(value interface{}) *Item {
	item.Lock()
	defer item.Unlock()

	var (
		date = time.Now()
	)

	item.value = value
	item.updated = date
	// item.frequency++
	return item
}

// SetElement set item node in list.
func (item *Item) SetElements(mode Mode, elements *Elements) *Item {
	item.Lock()
	defer item.Unlock()

	item.elements[mode] = elements
	return item
}

// IsValid returns value is not out of date
func (item *Item) IsValid() bool {
	item.RLock()
	defer item.RUnlock()

	if item.ttl == 0 {
		return true
	}
	return item.visited.Add(item.ttl).After(time.Now())
}

// Expired returns valid expired interval
func (item *Item) Expired() (bool, time.Duration) {
	item.RLock()
	defer item.RUnlock()

	if item.ttl == 0 {
		return true, 0
	}
	return false, item.ttl - time.Since(item.visited)
}

func (item *Item) MarshalJSON() ([]byte, error) {
	item.Lock()
	defer item.Unlock()

	return json.Marshal(
		map[string]interface{}{
			"key":       item.key,
			"value":     item.value,
			"ttl":       item.ttl,
			"created":   item.created.Unix(),
			"updated":   item.updated.Unix(),
			"visited":   item.visited.Unix(),
			"frequency": item.frequency,
		},
	)
}

func (item *Item) UnmarshalJSON(data []byte) error {
	item.Lock()
	defer item.Unlock()

	var (
		temp map[string]interface{}
		err  error
	)

	if err = json.Unmarshal(data, &temp); err != nil {
		return err
	}

	item.key = temp["key"].(string)
	item.value = temp["value"]
	item.ttl = time.Duration(temp["ttl"].(float64))
	item.created = time.Unix(int64(temp["created"].(float64)), 0)
	item.updated = time.Unix(int64(temp["updated"].(float64)), 0)
	item.visited = time.Unix(int64(temp["visited"].(float64)), 0)
	item.frequency = uint(temp["frequency"].(float64))
	item.elements = make(map[Mode]*Elements, 0)
	return nil
}
