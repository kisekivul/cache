package strategy

import (
	"github.com/kisekivul/cache/contract"
)

func NewStrategy(mode contract.Mode, data map[string]*contract.Item) contract.Strategy {
	var (
		strategy contract.Strategy
	)

	switch mode {
	case contract.LFU:
		strategy = &Frequently{}
	case contract.LRU:
		strategy = &Recently{}
	default:
		strategy = &Randomly{}
	}
	return strategy.Initialize(contract.NewItems().Initialize(data).List(strategy.Mode()))
}
