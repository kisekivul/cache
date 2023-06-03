package strategy

import (
	"container/list"
	"testing"

	"github.com/kisekivul/cache/contract"
	"github.com/stretchr/testify/assert"
)

func TestRandomly(t *testing.T) {
	var (
		strategy contract.Strategy
		item     *contract.Item
		items    = []*contract.Item{
			contract.NewItem("test_1", struct{}{}, 0),
		}
	)
	strategy = (&Randomly{}).Initialize(items)

	item = strategy.Execute()
	assert.Equal(t, true, item != nil)
	item = strategy.Execute()
	assert.Equal(t, true, item == nil)

	strategy.Append(contract.NewItem("test_2", struct{}{}, 0))
	item = strategy.Execute()
	assert.Equal(t, true, item != nil)

	item = strategy.Append(contract.NewItem("test_3", struct{}{}, 0))
	assert.Equal(t, item, strategy.Execute())
	assert.Equal(t, item, item.Elements(strategy.Mode()).Get("element").(*list.Element).Value)
}
