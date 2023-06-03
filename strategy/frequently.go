package strategy

import (
	"container/list"

	"github.com/kisekivul/cache/contract"
)

var (
	_ contract.Strategy = &Frequently{}
)

type Frequently struct {
	list *list.List
}

func (f *Frequently) Initialize(items []*contract.Item) contract.Strategy {
	f.list = list.New()

	for _, item := range items {
		f.Append(item)
	}
	return f
}

func (f *Frequently) Mode() contract.Mode {
	return contract.LFU
}

func (f *Frequently) Append(item *contract.Item) *contract.Item {
	var (
		elements  *contract.Elements
		element   *list.Element
		frequency = item.Frequency()
	)

	if element = f.list.Front(); element != nil {
		elements = element.Value.(*contract.Elements)
	}
LOOP:
	for {
		switch {
		case element == nil:
			elements = contract.NewElements(
				map[string]interface{}{
					"frequency": item.Frequency(),
					"list":      list.New(),
				},
			)
			element = f.list.PushFront(elements)
			break LOOP
		case frequency < elements.Get("frequency").(uint):
			elements = contract.NewElements(
				map[string]interface{}{
					"frequency": item.Frequency(),
					"list":      list.New(),
				},
			)
			element = f.list.InsertBefore(elements, element)
			break LOOP
		case frequency == elements.Get("frequency").(uint):
			break LOOP
		case frequency > elements.Get("frequency").(uint):
			if element = element.Next(); element != nil {
				elements = element.Value.(*contract.Elements)
				continue
			}

			elements = contract.NewElements(
				map[string]interface{}{
					"frequency": item.Frequency(),
					"list":      list.New(),
				},
			)
			element = f.list.PushBack(elements)
			break LOOP
		}
	}

	item.SetElements(
		f.Mode(),
		contract.NewElements(
			map[string]interface{}{
				"elements": element,
				"element":  elements.Get("list").(*list.List).PushBack(item),
			},
		),
	)
	return item
}

func (f *Frequently) Update(item *contract.Item) *contract.Item {
	f.Remove(item)
	return f.Append(item)
}

func (f *Frequently) Remove(item *contract.Item) *contract.Item {
	var (
		elements *contract.Elements
		element  *list.Element
		temp     *list.List
	)

	if f.list.Len() == 0 {
		return item
	}

	if elements = item.Elements(contract.LFU); elements != nil {
		element = elements.Get("elements").(*list.Element)
		temp = element.Value.(*contract.Elements).Get("list").(*list.List)
		item = temp.Remove(elements.Get("element").(*list.Element)).(*contract.Item)
		if temp.Len() == 0 {
			f.list.Remove(element)
		}
	}
	return item
}

func (f *Frequently) Execute() *contract.Item {
	var (
		element *list.Element
		items   *list.List
	)

	if f.list.Len() == 0 {
		return nil
	}

	items = f.list.Front().Value.(*contract.Elements).Get("list").(*list.List)
	element = items.Front()
	return items.Remove(element).(*contract.Item)
}
