package contract

import (
	"sort"
	"sync"
)

type Items struct {
	sync.Mutex
	mode  Mode
	items []*Item
}

func NewItems() *Items {
	return &Items{}
}

func (is *Items) Initialize(items map[string]*Item) *Items {
	is.Lock()
	defer is.Unlock()

	is.items = make([]*Item, 0)
	for _, item := range items {
		is.items = append(is.items, item)
	}
	return is
}

func (is *Items) List(mode Mode) []*Item {
	is.Lock()
	defer is.Unlock()

	is.mode = mode
	sort.Sort(is)

	return is.items
}

func (is *Items) Mode() Mode {
	return is.mode
}

func (is *Items) lfu(i, j *Item) bool {
	return i.frequency < j.frequency
}

func (is *Items) lru(i, j *Item) bool {
	return i.visited.Before(j.visited)
}

func (is *Items) Len() int {
	return len(is.items)
}

func (is *Items) Less(i, j int) bool {
	switch is.mode {
	case LFU:
		return is.lfu(is.items[i], is.items[j])
	case LRU:
		return is.lru(is.items[i], is.items[j])
	}
	return false
}

func (is *Items) Swap(i, j int) {
	is.items[i], is.items[j] = is.items[j], is.items[i]
}
