package contract

type Mode string

const (
	DEF Mode = "def"
	LRU Mode = "lru" // least recently used  最近最少使用
	LFU Mode = "lfu" // least frequently used 最不经常使用
)

type Strategy interface {
	Initialize([]*Item) Strategy
	Mode() Mode
	Append(*Item) *Item
	Update(*Item) *Item
	Remove(*Item) *Item
	Execute() *Item
}
