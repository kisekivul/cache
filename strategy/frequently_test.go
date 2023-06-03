package strategy

import (
	"testing"

	"github.com/kisekivul/cache/contract"
	"github.com/stretchr/testify/assert"
)

func TestFrequentlyy(t *testing.T) {
	var (
		strategy contract.Strategy
		item     *contract.Item
		items    = []*contract.Item{
			contract.NewItem("test_1", struct{}{}, 0),
			contract.NewItem("test_2", struct{}{}, 0),
			contract.NewItem("test_3", struct{}{}, 0),
			contract.NewItem("test_4", struct{}{}, 0),
			contract.NewItem("test_5", struct{}{}, 0),
		}
	)

	items[0].Value()
	items[0].Value()
	items[2].Value()
	items[3].Value()
	items[4].Value()
	items[4].Value()
	items[4].Value()

	strategy = (&Frequently{}).Initialize(items)

	assert.Equal(t, items[1], strategy.Execute())

	item = contract.NewItem("test_6", struct{}{}, 0)
	items = append(items, item)
	strategy.Append(item)
	assert.Equal(t, items[5], strategy.Execute())

	items = append(items[:1], items[2:len(items)-1]...)
	assert.Equal(t, 4, len(items))

	strategy.Remove(items[len(items)-1])
}
