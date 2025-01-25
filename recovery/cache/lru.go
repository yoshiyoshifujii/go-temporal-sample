package cache

import (
	"container/list"
	"errors"
	"sync"
	"time"
)

var (
	ErrCacheFull = errors.New("cache capacity is fully occupied with pinned elements")
)

type lru struct {
	mut      sync.Mutex
	byAccess *list.List
	byKey    map[string]*list.Element
	maxSize  int
	ttl      time.Duration
	pin      bool
	rmFunc   RemovedFunc
}

func New(maxSize int, opts *Options) Cache {
	if opts == nil {
		opts = &Options{}
	}
	return &lru{
		byAccess: list.New(),
		byKey:    make(map[string]*list.Element, opts.InitialCapacity),
		maxSize:  maxSize,
		ttl:      opts.TTL,
		pin:      opts.Pin,
		rmFunc:   opts.RemovedFunc,
	}
}

func NewLRU(maxSize int) Cache {
	return New(maxSize, nil)
}

func NewLRUWithInitialCapacity(initialCapacity, maxSize int) Cache {
	return New(maxSize, &Options{
		InitialCapacity: initialCapacity,
	})
}

func (c *lru) Get(key string) interface{} {
	c.mut.Lock()
	defer c.mut.Unlock()

	elt := c.byKey[key]
	if elt == nil {
		return nil
	}

	entry := elt.Value.(*cacheEntry)

	if c.pin {
		entry.refCount++
	}

	if entry.refCount == 0 && !entry.expiration.IsZero() && time.Now().After(entry.expiration) {
		if c.rmFunc != nil {
			go c.rmFunc(entry.value)
		}
		c.byAccess.Remove(elt)
		delete(c.byKey, entry.key)
		return nil
	}

	c.byAccess.MoveToFront(elt)
	return entry.value
}

func (c *lru) Put(key string, value interface{}) interface{} {
	if c.pin {
		panic("Cannot use Put API in Pin mode. Use Delete and PutIfNotExist if necessary.")
	}
	val, _ := c.putInternal(key, value, true)
	return val
}

func (c *lru) PutIfNotExist(key string, value interface{}) (interface{}, error) {
	existing, err := c.putInternal(key, value, false)
	if err != nil {
		return nil, err
	}

	if existing == nil {
		return value, nil
	}

	return existing, nil
}

func (c *lru) Delete(key string) {
	c.mut.Lock()
	defer c.mut.Unlock()

	elt := c.byKey[key]
	if elt != nil {
		entry := c.byAccess.Remove(elt).(*cacheEntry)
		if c.rmFunc != nil {
			go c.rmFunc(entry.value)
		}
		delete(c.byKey, key)
	}
}

func (c *lru) Release(key string) {
	c.mut.Lock()
	defer c.mut.Unlock()

	elt := c.byKey[key]
	entry := elt.Value.(*cacheEntry)
	entry.refCount--
}

func (c *lru) Size() int {
	c.mut.Lock()
	defer c.mut.Unlock()

	return len(c.byKey)
}

func (c *lru) putInternal(key string, value interface{}, allowUpdate bool) (interface{}, error) {
	c.mut.Lock()
	defer c.mut.Unlock()

	elt := c.byKey[key]
	if elt != nil {
		entry := elt.Value.(*cacheEntry)
		existing := entry.value
		if allowUpdate {
			entry.value = value
		}
		if c.ttl != 0 {
			entry.expiration = time.Now().Add(c.ttl)
		}
		c.byAccess.MoveToFront(elt)
		if c.pin {
			entry.refCount++
		}
		return existing, nil
	}

	entry := &cacheEntry{
		key:   key,
		value: value,
	}

	if c.pin {
		entry.refCount++
	}

	if c.ttl != 0 {
		entry.expiration = time.Now().Add(c.ttl)
	}

	c.byKey[key] = c.byAccess.PushFront(entry)
	if len(c.byKey) > c.maxSize {
		oldest := c.byAccess.Back().Value.(*cacheEntry)

		if oldest.refCount > 0 {
			c.byAccess.Remove(c.byAccess.Front())
			delete(c.byKey, key)
			return nil, ErrCacheFull
		}

		c.byAccess.Remove(c.byAccess.Back())
		if c.rmFunc != nil {
			go c.rmFunc(oldest.value)
		}
		delete(c.byKey, oldest.key)
	}

	return nil, nil
}

type cacheEntry struct {
	key        string
	expiration time.Time
	value      interface{}
	refCount   int
}
