package storage

import (
	"os"
	"testing"
	"time"

	"github.com/kisekivul/cache/contract"
	"github.com/stretchr/testify/assert"
)

func TestFile(t *testing.T) {
	var (
		f   contract.Storage
		err error
	)

	f, err = (&File{}).Initialize(time.Second, "./test.txt")
	assert.Nil(t, err)

	var (
		item_1 = contract.NewItem("test_1", struct{}{}, 0)
		item_2 = contract.NewItem("test_2", struct{}{}, 0)
	)

	err = f.Save(
		map[string]*contract.Item{
			"test_1": item_1,
			"test_2": item_2,
		},
	)
	assert.Nil(t, err)

	var (
		list map[string]*contract.Item
	)
	list, err = f.Load()
	assert.Nil(t, err)
	assert.Equal(t, 2, len(list))

	err = f.Save(
		map[string]*contract.Item{
			"test_1": item_1,
		},
	)
	assert.Nil(t, err)

	list, err = f.Load()
	assert.Nil(t, err)
	assert.Equal(t, 1, len(list))

	f.Exit()
	os.Remove("./test.txt")
}
