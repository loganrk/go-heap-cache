package cache

import (
	"container/list"
	"errors"
	"fmt"
	"sync"
	"time"
)

// LruCache represents a least-recently-used (LRU) cache.
type LruCache struct {
	queue    *list.List               // Doubly linked list to maintain item order.
	items    map[string]*list.Element // Map to store cached items by key.
	expire   int64                    // Expiration time in seconds.
	capacity int                      // Maximum number of items the cache can hold.
	mu       sync.RWMutex             // Read-write lock for concurrent access.
}

// Retrieves a cached item by key and updates its position in the cache.
func (c *LruCache) Get(key string) (any, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if item, ok := c.items[key]; ok {
		itemNode := item.Value.(Node)

		if itemNode.expireAt == NO_EXPIRE || itemNode.expireAt > time.Now().Unix() {
			c.queue.MoveToFront(item)

			return itemNode.data, nil
		} else {

			return nil, errors.New("heapCache :: Cache expired")
		}
	}

	return nil, errors.New("heapCache :: Cache not found")
}

// Adds or updates a cached item with the given key and value.
func (c *LruCache) Set(key string, value any) {
	var expireAt int64

	if c.expire != NO_EXPIRE {
		expireAt = int64(time.Now().Unix()) + c.expire
	} else {
		expireAt = NO_EXPIRE
	}
	c.set(key, value, expireAt)

}

// Adds or updates a cached item with the given key, value and custom expiration time specified by the caller.
func (c *LruCache) SetWithExpire(key string, value any, expiry int64) {
	var expireAt int64

	if expiry != NO_EXPIRE {
		expireAt = int64(time.Now().Unix()) + expiry
	} else {
		expireAt = NO_EXPIRE
	}

	c.set(key, value, expireAt)

}

// Returns a map of all non-expired items in the cache.
func (c *LruCache) GetAll() map[string]any {
	c.mu.RLock()
	defer c.mu.RUnlock()

	items := make(map[string]any, len(c.items))

	for key, item := range c.items {
		itemNode := item.Value.(Node)

		if itemNode.expireAt == NO_EXPIRE || itemNode.expireAt > time.Now().Unix() {
			items[key] = item
		}
	}

	return items
}

// Returns the current count of items in the cache.
func (c *LruCache) Count() int {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return len(c.items)
}

// Removes expired items from the cache and returns a status message.
func (c *LruCache) Delete(key string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if item, ok := c.items[key]; ok {
		c.queue.Remove(item)
		delete(c.items, key)

		return nil
	}

	return errors.New("heapCache :: Cache not found")
}

// Removes a cached item with the specified key from the cache.
func (c *LruCache) DeleteExpired() (string, error) {
	var deletedItemsCount int

	c.mu.Lock()
	defer c.mu.Unlock()

	if len(c.items) == 0 {

		return "", errors.New("heapCache :: No items available")
	}

	for key, item := range c.items {
		itemNode := item.Value.(Node)

		if itemNode.expireAt != NO_EXPIRE && itemNode.expireAt < time.Now().Unix() {
			c.queue.Remove(item)
			delete(c.items, key)
			deletedItemsCount++
		}
	}

	if deletedItemsCount > 0 {

		return fmt.Sprintf("heapCache :: %d item(s) are deleted", deletedItemsCount), nil
	} else {

		return "", errors.New("heapCache :: No expired items are found")
	}
}

// set is a private helper function for setting a cached item
func (c *LruCache) set(key string, value any, expireAt int64) {

	c.mu.Lock()
	defer c.mu.Unlock()

	if item, ok := c.items[key]; !ok {

		if c.capacity == len(c.items) {
			tailElement := c.queue.Back()
			c.queue.Remove(tailElement)
			if selectedNode, ok := tailElement.Value.(Node); ok {
				delete(c.items, selectedNode.key) // Retrieve the key for removal from the map.
			}
		}

		c.items[key] = c.queue.PushFront(Node{
			data:     value,
			expireAt: expireAt,
			key:      key,
		})

	} else {
		item.Value = Node{
			data:     value,
			expireAt: expireAt,
			key:      key,
		}

		c.items[key] = item
		c.queue.MoveToFront(item)
	}

}

// returns the least recently used (LRU) node in the cache, representing the oldest item.
func (c *LruCache) headNode() (Node, error) {
	headElement := c.queue.Front()

	if headElement == nil {

		return Node{}, errors.New("heapCache :: head node is not available")
	}

	if tailNode, ok := headElement.Value.(Node); ok {
		return tailNode, nil
	}

	return Node{}, errors.New("heapCache :: head node is not able get")
}

// returns the most recently used (MRU) node in the cache, representing the newest item.
func (c *LruCache) tailNode() (Node, error) {
	tailElement := c.queue.Back()

	if tailElement == nil {

		return Node{}, errors.New("heapCache :: tail node is not available")
	}

	if tailNode, ok := tailElement.Value.(Node); ok {
		return tailNode, nil
	}

	return Node{}, errors.New("heapCache :: tail node is not able get")

}
