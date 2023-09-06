# Heap Cache Package

A Go package that provides a simple implementation of a heap cache - Least Recently Used (LRU). The cache is designed to store key-value pairs with a limited capacity, and it automatically evicts the least recently used items when the capacity is exceeded.

## Features

- Efficient LRU caching mechanism.
- Customizable expiration time for cached items.
- Thread-safe operations using a read-write lock.

## Installation

To use this package in your Go project, you can simply import it and install the necessary dependencies:

```bash
go get github.com/loganrk/go-heap-cache
```

## Quickstart

```go
package main

import (
	"fmt"

	cache "github.com/loganrk/go-heap-cache"
)

func main() {
	// Create a new cache configuration
	conf := &cache.Config{
		Capacity:       cache.DEFAULT_CAPACITY,
		Expire:         cache.DEFAULT_EXPIRE,
		EvictionPolicy: cache.EVICTION_POLICY_LRU,
	}

	// Create a new cache instance based on the configuration
	c := cache.New(conf)

	// Add items to the cache
	c.Set("key1", "value1")
	c.Set("key2", "value2")
	c.SetWithExpire("key3", "value3", 60)

	// Retrieve items from the cache
	item1, err1 := c.Get("key1")
	item2, err2 := c.Get("key2")
	item3, err3 := c.Get("key3")

	if err1 == nil {
		fmt.Println("Item 1:", item1)
	}
	if err2 == nil {
		fmt.Println("Item 2:", item2)
	}
	if err3 == nil {
		fmt.Println("Item 3:", item3)
	}

	// Delete an item from the cache
	err := c.Delete("key2")
	if err == nil {
		fmt.Println("Item 'key2' deleted")
	}

	// Delete expired items from the cache
	status, err := c.DeleteExpired()
	if err == nil {
		fmt.Println(status)
	}

	// Get the count of items in the cache
	count := c.Count()
	fmt.Println("Number of items in cache:", count)
}

```

