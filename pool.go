package pools

import (
	"github.com/DaanV2/go-locks"
)

type (
	// ResourcePool is a pool of resources that can be locked and unlocked. It is useful for managing resources that can be used by multiple goroutines.
	ResourcePool[T any] struct {
		items    []*T
		lockPool *locks.Pool
	}

	// ReturnLoan is a function that should be called to return a resource to the pool
	ReturnLoan[T any] func(toReturn *T)
)

// NewResourcePoolWith creates a new ResourcePool with the given items and locks amount. The pool should have ownership of the items.
func NewResourcePool[T any](items []*T, locksAmount int) *ResourcePool[T] {
	return &ResourcePool[T]{
		items:    items,
		lockPool: locks.NewPool(locksAmount),
	}
}

// Loan returns a resource from the pool and a function that should be called to return the resource to the pool.
// This allows the resource to be set to nil, or be easily returned or replaced. At the point of returning the resource,
// the lock is released. And the pool should gain ownership of the resource.
// Example:
//   item, returnFn := pool.Loan(index)
//   defer returnFn(item)
//   // Do something with item
func (pool *ResourcePool[T]) Loan(index uint64) (*T, ReturnLoan[T]) {
	lock := pool.lockPool.GetLock(index)
	
	lock.Lock()
	returnFn := func(item *T) {
		defer lock.Unlock()

		pool.items[index] = item
	}
	
	item := pool.items[index]
	
	return item, returnFn
}

func (pool *ResourcePool[T]) Len() int {
	return len(pool.items)
}

// LoanByUint64 returns a resource from the pool and a function that should be called to return the resource to the pool.
// It uses a uint64 key to get the resource from the pool. This is useful for resources that are identified by a number id.
func (pool *ResourcePool[T]) LoanByUint64(key uint64) (*T, ReturnLoan[T]) {
	index := pool.getIndex(key)

	return pool.Loan(index)
}

// LoanByString returns a resource from the pool and a function that should be called to return the resource to the pool.
// It uses a string key to get the resource from the pool. This is useful for files, or other resources that are identified by a string.
func (pool *ResourcePool[T]) LoanByString(key string) (*T, ReturnLoan[T]) {
	v := locks.KeyForString(key)
	index := pool.getIndex(v)

	return pool.Loan(index)
}

// LoanByString returns a resource from the pool and a function that should be called to return the resource to the pool.
func (pool *ResourcePool[T]) LoanByBytes(key []byte) (*T, ReturnLoan[T]) {
	v := locks.KeyForBytes(key)
	index := pool.getIndex(v)

	return pool.Loan(index)
}

// getIndex returns the index of the item in the pool for the given key
func (pool *ResourcePool[T]) getIndex(key uint64) uint64 {
	v := key % uint64(len(pool.items))
	
	return v
}