package strategy

import (
	"container/list"

	"github.com/kisekivul/cache/contract"
)

var (
	_ contract.Strategy = &Recently{}
)

type Recently struct {
	list *list.List
}

func (r *Recently) Initialize(items []*contract.Item) contract.Strategy {
	r.list = list.New()
	for _, item := range items {
		r.Append(item)
	}
	return r
}

func (r *Recently) Mode() contract.Mode {
	return contract.LRU
}

func (r *Recently) Append(item *contract.Item) *contract.Item {
	return item.SetElements(
		r.Mode(),
		contract.NewElements(
			map[string]interface{}{
				"element": r.list.PushBack(item),
			},
		),
	)
}

func (r *Recently) Update(item *contract.Item) *contract.Item {
	r.list.MoveToBack(item.Elements(r.Mode()).Get("element").(*list.Element))
	return item
}

func (r *Recently) Remove(item *contract.Item) *contract.Item {
	var (
		elements *contract.Elements
	)

	if r.list.Len() == 0 {
		return item
	}

	if elements = item.Elements(contract.LRU); elements != nil {
		item = r.list.Remove(elements.Get("element").(*list.Element)).(*contract.Item)
	}
	return item
}

func (r *Recently) Execute() *contract.Item {
	var (
		element *list.Element
		item    *contract.Item
	)

	if r.list.Len() == 0 {
		return nil
	}

	if element = r.list.Front(); element != nil {
		item = r.list.Remove(element).(*contract.Item)
	}
	return item
}
