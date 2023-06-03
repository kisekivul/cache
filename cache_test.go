package cache

import (
	"os"
	"testing"
	"time"

	"github.com/kisekivul/cache/contract"
	"github.com/kisekivul/cache/storage"
	"github.com/stretchr/testify/assert"
)

func TestCache(t *testing.T) {
	var (
		cache contract.Cache
		fs    contract.Storage
		err   error
	)

	fs, err = (&storage.File{}).Initialize(time.Second, "./test.txt")
	assert.Nil(t, err, nil)

	defer func() {
		fs.Exit()
		os.Remove("./test.txt")
	}()

	cache = NewCache("test")
	cache.Initialize(contract.LRU, 2, fs)

	cache.Add("test_1", "test_1", time.Minute)
	cache.Add("test_2", "test_2", time.Minute)
	cache.Get("test_2")
	cache.Add("test_3", "test_3", time.Minute)

	assert.Equal(t, 2, len(cache.List()))
}
