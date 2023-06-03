package contract

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSort(t *testing.T) {
	var (
		list   = make(map[string]*Item, 0)
		item_1 = NewItem("test_1", struct{}{}, 0)
		item_2 = NewItem("test_2", struct{}{}, 0)
	)

	list["test_1"] = item_1
	list["test_2"] = item_2

	var (
		items = NewItems().Initialize(list)
	)

	item_1.Value()
	assert.Equal(t, items.List(LFU), []*Item{item_2, item_1})
	item_1.Value()
	time.Sleep(time.Second)
	item_2.Value()
	assert.Equal(t, items.List(LRU), []*Item{item_1, item_2})
	assert.Equal(t, items.List(LFU), []*Item{item_2, item_1})
}

func TestSerialize(t *testing.T) {
	var (
		item_1 = NewItem("test", struct{}{}, time.Second)
		item_2 = NewItem("", struct{}{}, 0)
	)
	item_1.Value()

	var (
		val_1, val_2 []byte
		err          error
	)

	val_1, err = json.Marshal(item_1)
	assert.Nil(t, err, nil)
	err = json.Unmarshal(val_1, &item_2)
	assert.Nil(t, err, nil)

	val_2, _ = json.Marshal(item_2)
	assert.Equal(t, val_1, val_2)
}
