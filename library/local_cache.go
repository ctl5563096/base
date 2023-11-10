package library

import (
	"sync"
	"time"
)

type LruNode struct {
	Prev    *LruNode
	Next    *LruNode
	Timeout int64
	Value   interface{}
	Key     string
}

// LocalCache  is an LRU cache. It is not safe for concurrent access.
type LocalCache struct {
	maxSize int                      // Maximum number of items in cache
	size    int                      // Number of items in cache
	head    *LruNode                 // Head of linked list
	tail    *LruNode                 // Tail of linked list
	cache   map[interface{}]*LruNode // Keyed list of items
	locker  *sync.RWMutex            // Mutex for cache
}

// NewLocalCache creates a new LocalCache of the given size.
// 初始化实例
// cache := NewLocalCache(1024)
// 向缓存存放数据
// cache.Put("key", "value", time.Second * 5)
// 获取缓存数据
// val,exist = cache.Get("key")
// val 获取到的缓存值
// exist 是否存在
func NewLocalCache(maxSize int) *LocalCache {
	return &LocalCache{
		maxSize: maxSize,
		size:    0,
		head:    nil,
		tail:    nil,
		cache:   make(map[interface{}]*LruNode),
		locker:  new(sync.RWMutex),
	}
}

// Put adds a value to the cache.expire 有效期,单位秒
func (c *LocalCache) Put(key string, val interface{}, expire time.Duration) error {

	// 确认容量是否超出,如果超出了则进行清理

	locker := c.locker

	locker.Lock()
	defer locker.Unlock()
	c.ifFullRemoveLast()

	ts := time.Now().Add(expire).Unix()

	node := &LruNode{
		Prev:    nil,
		Next:    nil,
		Timeout: ts,
		Value:   val,
		Key:     key,
	}
	c.addToHead(node)
	return nil
}

// Get looks up a key's value from the cache.
func (c *LocalCache) Get(key string) (interface{}, bool) {
	locker := c.locker

	locker.RLock()
	node, ok := c.cache[key]
	if !ok {
		locker.RUnlock()
		return nil, false
	}
	locker.RUnlock()

	locker.Lock()
	defer locker.Unlock()
	if node.Timeout < time.Now().Unix() {
		c.delete(node)
		return nil, false
	}
	c.moveToHead(node)
	return node.Value, true
}

func (c *LocalCache) delete(node *LruNode) {
	if node == nil {
		return
	}

	key := node.Key

	if c.head == node {
		c.head = node.Next
	}
	if c.tail == node {
		c.tail = node.Prev
	}

	if node.Prev != nil {
		node.Prev.Next = node.Next
	}
	if node.Next != nil {
		node.Next.Prev = node.Prev
	}

	node.Next, node.Prev = nil, nil

	delete(c.cache, key)
	c.size -= 1
}

func (c *LocalCache) moveToHead(node *LruNode) {
	if node == nil {
		return
	}
	// 已经在队列头了
	if c.head == node {
		return
	}
	// 在队列尾了
	if c.tail == node {
		c.tail = node.Prev
	}

	//
	prev := node.Prev
	if prev != nil {
		prev.Next = node.Next
	}

	head := c.head
	node.Prev, node.Next = nil, head
	head.Prev = node
	c.head = node
}

func (c *LocalCache) addToHead(node *LruNode) {
	if node == nil {
		return
	}

	head := c.head

	if head == nil {
		c.head, c.tail = node, node
	} else {

		node.Next = head
		node.Prev = nil

		head.Prev = node

		c.head = node
	}
	c.cache[node.Key] = node
	c.size++
}

func (c *LocalCache) ifFullRemoveLast() {
	if c.size+1 > c.maxSize {
		c.delete(c.tail)
	}
}
