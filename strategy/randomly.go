package strategy

import (
	"container/list"
	"math/rand"
	"time"

	"github.com/kisekivul/cache/contract"
)

var (
	_ contract.Strategy = &Randomly{}
)

type Randomly struct {
	list *list.List
}

func (r *Randomly) Initialize(items []*contract.Item) contract.Strategy {
	r.list = list.New()
	for _, item := range items {
		r.Append(item)
	}
	return r
}

func (r *Randomly) Mode() contract.Mode {
	return contract.DEF
}

func (r *Randomly) Append(item *contract.Item) *contract.Item {
	return item.SetElements(
		r.Mode(),
		contract.NewElements(
			map[string]interface{}{
				"element": r.list.PushBack(item),
			},
		),
	)
}

func (r *Randomly) Update(item *contract.Item) *contract.Item {
	return item
}

func (r *Randomly) Remove(item *contract.Item) *contract.Item {
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

func (r *Randomly) Execute() *contract.Item {
	var (
		element          *list.Element
		item             *contract.Item
		length, position int
	)

	if length = r.list.Len(); length == 0 {
		return nil
	}

	rand.Seed(time.Now().UnixNano())
	if position = rand.Intn(r.list.Len()); position < length/2 {
		element = r.list.Front()
		for i := 1; i <= position; i++ {
			element = element.Next()
		}
	} else {
		element = r.list.Back()
		for i := length - 1; i > position; i-- {
			element = element.Prev()
		}
	}

	if element != nil {
		item = r.list.Remove(element).(*contract.Item)
	}
	return item
}
