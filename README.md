# Heap Cache Package

A Go package that provides a simple implementation of a heap cache - Least Recently Used (LRU). The cache is designed to store key-value pairs with a limited capacity, and it automatically evicts the least recently used items when the capacity is exceeded.

## Features

- Efficient LRU caching mechanism.
- Customizable expiration time for cached items.
- Thread-safe operations using a read-write lock.

## Installation

To use this package in your Go project, you can simply import it and install the necessary dependencies:

```bash
go get github.com/your-username/lru-cache-package
