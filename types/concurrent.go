package types

import "sync"

type ConcurrentCollection[T any] struct {
	mutex      sync.Mutex
	collection []T
}

func (c *ConcurrentCollection[T]) Add(value T) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.collection = append(c.collection, value)
}

func (c *ConcurrentCollection[T]) Get() []T {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	return c.collection
}

func (c *ConcurrentCollection[T]) Append(value T) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.collection = append(c.collection, value)
}
