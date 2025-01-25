package cache

import "time"

type Cache interface {
	Get(key string) interface{}
	Put(key string, value interface{}) interface{}
	PutIfNotExist(key string, value interface{}) (interface{}, error)
	Delete(key string)
	Release(key string)
	Size() int
}

type Options struct {
	TTL             time.Duration
	InitialCapacity int
	Pin             bool
	RemovedFunc     RemovedFunc
}

type RemovedFunc func(interface{})
