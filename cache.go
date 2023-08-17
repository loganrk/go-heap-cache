package cache

import "container/list"

const (
	DEFAULT_CAPACITY    int   = 100
	DEFAULT_EXPIRE      int64 = 1200 // 20 Minutes
	NO_EXPIRE           int64 = -1
	EVICTION_POLICY_LRU int   = 1
)

type Cache interface {
	Get(key string) (any, error)
	Set(key string, value any)
	SetWithExpire(key string, value any, expiry int64)
	GetAll() map[string]any
	Count() int
	Delete(key string) error
	DeleteExpired() (string, error)
	set(key string, value any, expireAt int64)
	headNode() (Node, error)
	tailNode() (Node, error)
}

type Config struct {
	Capacity       int
	Expire         int64
	EvictionPolicy int
}

// Node represents an item in the cache.
type Node struct {
	data     any   // The cached data.
	expireAt int64 // Expiration time for the item.
	key      string
}

// Creates a new cache instance based on the provided configuration.
// It selects the eviction policy based on the configuration and returns an
// instance of the chosen cache type initialized with the specified capacity and expiration time.
func New(conf *Config) Cache {

	switch conf.EvictionPolicy {

	case EVICTION_POLICY_LRU:

		return &LruCache{
			queue:    list.New(),
			items:    make(map[string]*list.Element, conf.Capacity),
			capacity: conf.Capacity,
			expire:   conf.Expire,
		}

	default:

		return &LruCache{
			queue:    list.New(),
			items:    make(map[string]*list.Element, conf.Capacity),
			capacity: conf.Capacity,
			expire:   conf.Expire,
		}
	}

}
