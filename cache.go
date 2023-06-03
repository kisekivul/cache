package cache

import (
	"sync"
	"time"

	"github.com/kisekivul/cache/contract"
	"github.com/kisekivul/cache/strategy"
)

type Cache struct {
	sync.RWMutex
	// data
	name     string
	limit    uint
	items    map[string]*contract.Item
	strategy contract.Strategy
	storage  contract.Storage
	// keepalive
	trigger *time.Timer
	expired time.Duration
}

func NewCache(name string) contract.Cache {
	return &Cache{
		name:     name,
		items:    make(map[string]*contract.Item),
		strategy: nil,
		storage:  nil,
		trigger:  nil,
		expired:  0,
	}
}

func (c *Cache) Initialize(mode contract.Mode, limit uint, storage contract.Storage) contract.Cache {
	if storage != nil {
		c.storage = storage
		if items, _ := c.storage.Load(); items != nil {
			c.items = items
		}
	}

	if c.limit = limit; c.limit > 0 {
		c.strategy = strategy.NewStrategy(mode, c.items)
	}

	go func() {
		for {
			c.RLock()
			c.save()
			c.RUnlock()

			time.Sleep(10 * time.Second)
		}
	}()
	return c
}

func (c *Cache) count() uint64 {
	return uint64(len(c.items))
}

func (c *Cache) Count() uint64 {
	c.RLock()
	defer c.RUnlock()

	return c.count()
}

func (c *Cache) get(key string) *contract.Item {
	var (
		item *contract.Item
		ex   bool
	)

	if item, ex = c.items[key]; ex {
		return item
	}
	return nil
}

func (c *Cache) Get(key string) (interface{}, error) {
	c.RLock()
	defer c.RUnlock()

	var (
		item  *contract.Item
		value interface{}
	)

	if item = c.get(key); item != nil {
		value = item.Value()
		// trigger strategy
		if c.strategy != nil {
			c.strategy.Update(item)
		}
		return value, nil
	}
	return nil, contract.ErrKeyNotFound
}

func (c *Cache) add(item *contract.Item) {
	c.items[item.Key()] = item
	if c.strategy != nil {
		c.strategy.Append(item)
	}

	if item.TTL() != 0 && (item.TTL() < c.expired || c.expired == 0) {
		c.refresh(item.TTL() - time.Since(item.Visited()))
	}
}

func (c *Cache) Add(key string, val interface{}, ttl time.Duration) error {
	c.Lock()
	defer c.Unlock()

	if c.limit > 0 && uint(len(c.items)) >= c.limit {
		// refresh first
		c.refresh(0)
		// double check
		if c.limit > 0 && uint(len(c.items)) >= c.limit {
			if c.strategy == nil {
				return contract.ErrStrategyNotSet
			}

			// double check
			for uint(len(c.items)) >= c.limit {
				if item := c.strategy.Execute(); item != nil {
					c.delete(item.Key())
				} else {
					break
				}
			}
		}
	}
	// last check
	if c.limit > 0 && uint(len(c.items)) >= c.limit {
		return contract.ErrEvictItemFailed
	}

	c.add(contract.NewItem(key, val, ttl))
	return nil
}

func (c *Cache) TTL(key string, ttl time.Duration) error {
	c.Lock()
	defer c.Unlock()

	var (
		item *contract.Item
		ex   bool
	)

	if item, ex = c.items[key]; ex {
		item.SetTTL(ttl)
	}
	return nil
}

func (c *Cache) Update(key string, val interface{}) error {
	var (
		item *contract.Item
		ex   bool
	)

	if item, ex = c.items[key]; ex {
		item.SetValue(val)
	}
	return nil
}

func (c *Cache) Exist(key string) bool {
	c.RLock()
	defer c.RUnlock()

	var (
		exist bool
	)

	_, exist = c.items[key]
	return exist
}

func (c *Cache) delete(key string) *contract.Item {
	var (
		item *contract.Item
	)

	if item = c.get(key); item != nil {
		delete(c.items, key)
		// trigger strategy

		if c.strategy != nil {
			c.strategy.Remove(item)
		}
	}
	return item
}

func (c *Cache) Delete(key string) error {
	c.Lock()
	defer c.Unlock()

	c.delete(key)
	return nil
}

func (c *Cache) Keys() []string {
	c.RLock()
	defer c.RUnlock()

	var (
		list = make([]string, 0)
	)

	for key := range c.items {
		list = append(list, key)
	}
	return list
}

func (c *Cache) List() map[interface{}]interface{} {
	c.RLock()
	defer c.RUnlock()

	var (
		list = make(map[interface{}]interface{})
	)

	for _, item := range c.items {
		if item.IsValid() {
			list[item.Key()] = item.Value()
		}
	}
	return list
}

func (c *Cache) Flush() {
	c.Lock()
	defer c.Unlock()

	if c.trigger != nil {
		c.trigger.Stop()
	}
	c.items = make(map[string]*contract.Item)
	c.expired = 0

	c.save()
}

func (c *Cache) Refresh() {
	c.Lock()
	defer c.Unlock()

	c.refresh(0)
}

func (c *Cache) refresh(next time.Duration) {
	if next <= 0 {
		if c.trigger != nil {
			c.trigger.Stop()
		}

		for key, item := range c.items {
			var (
				valid, expired = item.Expired()
			)

			switch {
			case valid:
				continue
			case expired <= 0:
				c.delete(key)
			default:
				if expired < next || next == 0 {
					next = expired
				}
			}
		}
	}
	c.expired = next

	c.save()

	if next > 0 {
		c.trigger = time.AfterFunc(
			next,
			func() {
				go c.Refresh()
			},
		)
	}
}

func (c *Cache) save() {
	if c.storage != nil {
		go func() {
			var (
				i   int
				err error
			)

			for i < 3 {
				if err = c.storage.Save(c.items); err == nil {
					break
				}
				i++
				time.Sleep(time.Second)
			}
		}()
	}
}
