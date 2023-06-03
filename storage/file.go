package storage

import (
	"encoding/json"
	"os"
	"sync"
	"time"

	"github.com/kisekivul/cache/contract"
)

type File struct {
	sync.RWMutex
	path     string
	interval time.Duration
	latest   time.Time
	file     *os.File
}

func (f *File) Initialize(interval time.Duration, params interface{}) (contract.Storage, error) {
	var (
		err error
	)

	if f.interval = interval; f.interval == 0 {
		f.interval = time.Second
	}

	f.path = params.(string)
	if f.file, err = os.OpenFile(f.path, os.O_RDWR|os.O_CREATE, 0766); err != nil {
		return f, err
	}
	return f, nil
}

func (f *File) Load() (map[string]*contract.Item, error) {
	f.Lock()
	defer f.Unlock()

	var (
		val   []byte
		err   error
		items map[string]*contract.Item
	)

	if val, err = os.ReadFile(f.path); err != nil {
		return nil, err
	}

	if val == nil {
		return nil, nil
	}

	if err = json.Unmarshal(val, &items); err != nil {
		return nil, err
	}
	return items, nil
}

func (f *File) Save(items map[string]*contract.Item) error {
	if f.interval > time.Since(f.latest) {
		return nil
	}

	f.Lock()
	defer f.Unlock()

	if info, _ := f.file.Stat(); info == nil {
		return nil
	}

	var (
		val []byte
		err error
	)

	if err = f.file.Truncate(0); err != nil {
		return err
	}
	f.file.Seek(0, 0)

	if val, err = json.Marshal(items); err != nil {
		return err
	}

	if _, err = f.file.Write(val); err != nil {
		return err
	}
	f.latest = time.Now()
	return nil
}

func (f *File) Exit() {
	f.Lock()
	defer f.Unlock()

	f.file.Close()
}
