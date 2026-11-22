package cache

import (
	"container/list"
	"sync"
	"sync/atomic"
)

type lruEntry struct {
	key  string
	data []byte
	size int
}

// L1MemoryCache is a thread-safe in-memory LRU cache with byte-size accounting.
type L1MemoryCache struct {
	mu         sync.RWMutex
	maxBytes   int64
	currBytes  int64
	items      map[string]*list.Element
	evictList  *list.List
	hits       atomic.Uint64
	misses     atomic.Uint64
	evictions  atomic.Uint64
}

// NewL1MemoryCache creates an L1 in-memory LRU cache with a maximum memory limit in bytes.
func NewL1MemoryCache(maxBytes int64) *L1MemoryCache {
	if maxBytes <= 0 {
		maxBytes = 24 * 1024 * 1024 // 24MB default
	}
	return &L1MemoryCache{
		maxBytes:  maxBytes,
		items:     make(map[string]*list.Element),
		evictList: list.New(),
	}
}

// Get retrieves data from the LRU cache.
func (c *L1MemoryCache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	elem, exists := c.items[key]
	if !exists {
		c.misses.Add(1)
		return nil, false
	}

	c.evictList.MoveToFront(elem)
	c.hits.Add(1)
	entry := elem.Value.(*lruEntry)

	// Return copy to prevent caller mutation
	res := make([]byte, len(entry.data))
	copy(res, entry.data)
	return res, true
}

// Set stores data in the LRU cache, evicting oldest entries if over maxBytes limit.
func (c *L1MemoryCache) Set(key string, data []byte) {
	dataSize := int64(len(data))
	if dataSize > c.maxBytes {
		return
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	if elem, exists := c.items[key]; exists {
		c.evictList.MoveToFront(elem)
		entry := elem.Value.(*lruEntry)
		c.currBytes += dataSize - int64(entry.size)
		entry.data = data
		entry.size = int(dataSize)
	} else {
		entry := &lruEntry{
			key:  key,
			data: data,
			size: int(dataSize),
		}
		elem := c.evictList.PushFront(entry)
		c.items[key] = elem
		c.currBytes += dataSize
	}

	for c.currBytes > c.maxBytes && c.evictList.Len() > 0 {
		c.evictOldest()
	}
}

func (c *L1MemoryCache) evictOldest() {
	elem := c.evictList.Back()
	if elem == nil {
		return
	}
	c.evictList.Remove(elem)
	entry := elem.Value.(*lruEntry)
	delete(c.items, entry.key)
	c.currBytes -= int64(entry.size)
	c.evictions.Add(1)
}

// Stats returns operational metrics of the L1 cache.
func (c *L1MemoryCache) Stats() (items int, bytes int64, hits uint64, misses uint64, evictions uint64) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.items), c.currBytes, c.hits.Load(), c.misses.Load(), c.evictions.Load()
}

// Clear flushes all entries from memory.
func (c *L1MemoryCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items = make(map[string]*list.Element)
	c.evictList.Init()
	c.currBytes = 0
}
