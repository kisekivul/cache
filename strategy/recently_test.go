package strategy

import (
	"testing"

	"github.com/kisekivul/cache/contract"
	"github.com/stretchr/testify/assert"
)

func TestRecently(t *testing.T) {
	var (
		strategy contract.Strategy
		item     *contract.Item
		items    = []*contract.Item{
			contract.NewItem("test_1", struct{}{}, 0),
			contract.NewItem("test_2", struct{}{}, 0),
			contract.NewItem("test_3", struct{}{}, 0),
			contract.NewItem("test_4", struct{}{}, 0),
		}
	)
	strategy = (&Recently{}).Initialize(items)

	for i := 0; i < len(items); i++ {
		strategy.Update(items[i])
	}
	item = strategy.Execute()
	assert.Equal(t, items[0], item)

	items = items[1:]
	for i := len(items) - 1; i >= 0; i-- {
		strategy.Update(items[i])
	}
	item = strategy.Execute()
	assert.Equal(t, items[len(items)-1], item)
	items = items[:len(items)-1]

	item = contract.NewItem("test_5", struct{}{}, 0)
	strategy.Append(item)
	for i := 0; i < len(items); i++ {
		strategy.Update(items[i])
	}
	assert.Equal(t, item, strategy.Execute())

	assert.Equal(t, items[0], strategy.Remove(items[0]))
	items = items[1:]

	for i := 0; i < len(items); i++ {
		assert.Equal(t, items[i], strategy.Execute())
	}
}
