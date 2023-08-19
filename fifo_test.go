package cache

import (
	"fmt"
	"testing"
	"time"
)

func Test_FifoCacheSetAndGet(t *testing.T) {
	cache := New(&Config{
		Capacity:       DEFAULT_CAPACITY,
		Expire:         DEFAULT_EXPIRE,
		EvictionPolicy: EVICTION_POLICY_FIFO,
	})
	cache.Set("key1", "value1")

	value, err := cache.Get("key1")
	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}
	if value != "value1" {
		t.Errorf("Expected value 'value1', but got %v", value)
	}
	_, err = cache.Get("non_existing_key")
	if err == nil || err.Error() != "heapCache :: Cache not found" {
		t.Errorf("Expected 'Cache not found' error, but got %v", err)
	}
}

func Test_FifoCacheSetWithExpire(t *testing.T) {
	cache := New(&Config{
		Capacity:       DEFAULT_CAPACITY,
		Expire:         NO_EXPIRE,
		EvictionPolicy: EVICTION_POLICY_FIFO,
	})

	cache.SetWithExpire("key2", "value2", 2)

	value, err := cache.Get("key2")
	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}
	if value != "value2" {
		t.Errorf("Expected value 'value2', but got %v", value)
	}

	cache.SetWithExpire("key3", "value3", 2)
	time.Sleep(3 * time.Second)
	_, err = cache.Get("key3")
	if err == nil || err.Error() != "heapCache :: Cache expired" {
		t.Errorf("Expected 'Cache expired' error, but got %v", err)
	}
}

func Test_FifoCacheGetAll(t *testing.T) {
	cache := New(&Config{
		Capacity:       5,
		Expire:         NO_EXPIRE,
		EvictionPolicy: EVICTION_POLICY_FIFO,
	})
	cache.Set("key1", "value1")
	cache.Set("key2", "value2")
	cache.Set("key3", "value4")
	cache.Set("key4", "value4")

	items := cache.GetAll()
	if len(items) != 4 {
		t.Errorf("Expected 4 item in the cache, but got %d", len(items))
	}

	cache.Set("key5", "value5")
	cache.Set("key6", "value6")
	cache.Set("key7", "value7")

	items = cache.GetAll()
	if len(items) != 5 {
		t.Errorf("Expected 5 item in the cache, but got %d", len(items))
	}
}

func Test_FifoCacheDelete(t *testing.T) {
	cache := New(&Config{
		Capacity:       DEFAULT_CAPACITY,
		Expire:         NO_EXPIRE,
		EvictionPolicy: EVICTION_POLICY_FIFO,
	})

	cache.Set("key5", "value5")

	err := cache.Delete("key5")
	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}

	err = cache.Delete("key5")
	if err == nil || err.Error() != "heapCache :: Cache not found" {
		t.Errorf("Expected 'heapCache :: Cache not found', but got %v", err)
	}

	_, err = cache.Get("key5")
	if err == nil || err.Error() != "heapCache :: Cache not found" {
		t.Errorf("Expected 'Cache not found' error after Delete, but got %v", err)
	}
}

func Test_FifoCacheDeleteExpired(t *testing.T) {
	cache := New(&Config{
		Capacity:       DEFAULT_CAPACITY,
		Expire:         NO_EXPIRE,
		EvictionPolicy: EVICTION_POLICY_FIFO,
	})

	_, err := cache.DeleteExpired()
	if err != nil && err.Error() != "heapCache :: No items available" {
		t.Errorf("Expected 'heapCache :: No items available', but got %v", err)
	}

	cache.SetWithExpire("key5", "value5", NO_EXPIRE)
	cache.SetWithExpire("key6", "value6", 2)

	time.Sleep(3 * time.Second)

	result, _ := cache.DeleteExpired()
	if result != "heapCache :: 1 item(s) are deleted" {
		t.Errorf("Expected '1 item(s) are deleted', but got %s", result)
	}

	cache.SetWithExpire("key7", "value7", NO_EXPIRE)
	cache.Set("key8", "value8")

	_, err = cache.DeleteExpired()
	if err != nil && err.Error() != "heapCache :: No expired items are found" {
		t.Errorf("Expected 'No expired items are found', but got %v", err)
	}

	count := cache.Count()
	if count != 3 {
		t.Errorf("Expected '2', but got %d", count)
	}

}

func Test_FifoCacheSHeadNode(t *testing.T) {
	cache := New(&Config{
		Capacity:       5,
		Expire:         NO_EXPIRE,
		EvictionPolicy: EVICTION_POLICY_FIFO,
	})

	_, err := cache.headNode()
	if err == nil || err.Error() != "heapCache :: head node is not available" {
		t.Errorf("Expected 'head node is not available' error, but got %v", err)
	}

	cache.Set("key1", "value1")
	cache.Set("key2", "value2")

	node, err := cache.headNode()
	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}
	if node.key != "key2" {
		t.Errorf("Expected head node's key to be 'key2', but got %s", node.key)
	}

	cache.Set("key3", "value3")
	cache.Set("key1", "value1")

	node, err = cache.headNode()
	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}

	if node.key != "key1" {
		t.Errorf("Expected head node's key to be 'key1', but got %s", node.key)
	}
}

func Test_FifoCacheSTailNode(t *testing.T) {
	cache := New(&Config{
		Capacity:       5,
		Expire:         NO_EXPIRE,
		EvictionPolicy: EVICTION_POLICY_FIFO,
	})

	_, err := cache.tailNode()
	if err == nil || err.Error() != "heapCache :: tail node is not available" {
		t.Errorf("Expected 'tail node is not available' error, but got %v", err)
	}

	cache.Set("key1", "value1")

	node, err := cache.tailNode()
	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}
	if node.key != "key1" {
		t.Errorf("Expected tail node's key to be 'key1', but got %s", node.key)
	}

	cache.Set("key2", "value2")
	cache.Set("key3", "value3")
	cache.Set("key1", "value1")

	node, err = cache.tailNode()
	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}
	if node.key != "key2" {
		t.Errorf("Expected tail node's key to be 'key2', but got %s", node.key)
	}

	cache.Set("key4", "value4")
	cache.Set("key5", "value5")
	cache.Set("key6", "value6")

	node, err = cache.tailNode()
	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}
	if node.key != "key3" {
		t.Errorf("Expected tail node's key to be 'key3', but got %s", node.key)
	}

}

func Benchmark_FifoCacheSet(b *testing.B) {
	cache := New(&Config{
		Capacity:       DEFAULT_CAPACITY,
		Expire:         DEFAULT_EXPIRE,
		EvictionPolicy: EVICTION_POLICY_FIFO,
	})

	for i := 0; i < b.N; i++ {
		cache.Set("key"+fmt.Sprintf("%d", i), "value"+fmt.Sprintf("%d", i))
	}
}

func Benchmark_FifoCacheGet(b *testing.B) {
	cache := New(&Config{
		Capacity:       DEFAULT_CAPACITY,
		Expire:         DEFAULT_EXPIRE,
		EvictionPolicy: EVICTION_POLICY_FIFO,
	})

	for i := 0; i < b.N; i++ {
		cache.Set("key"+fmt.Sprintf("%d", i), "value"+fmt.Sprintf("%d", i))
	}
	b.ResetTimer() // Reset the timer to exclude setup time
	for i := 0; i < b.N; i++ {
		cache.Get("key" + fmt.Sprintf("%v", i))
	}
}
