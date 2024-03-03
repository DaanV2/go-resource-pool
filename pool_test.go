package pools_test

import (
	"sync"
	"testing"
	"time"

	pools "github.com/DaanV2/go-resource-pool"
	"github.com/stretchr/testify/require"
)

type TestItem struct {
	ID int
	Update bool
}

func Test_ResourcePool(t *testing.T) {
	t.Run("Should be able to get a resource from the pool and update it", func(t *testing.T) {
		original := &TestItem{
			ID: 1,
			Update: false,
		}

		pool := pools.NewResourcePool([]*TestItem{original}, 10)

		item, returnFn := pool.Loan(0)
		item.Update = true
		returnFn(item)

		require.Equal(t, original.Update, true)
	})

	t.Run("Should be able to get a resource from the pool using a string key", func(t *testing.T) {
		original := &TestItem{
			ID: 1,
			Update: false,
		}

		pool := pools.NewResourcePool([]*TestItem{original}, 10)

		item, returnFn := pool.LoanForString("dir/file.txt")
		item.Update = true
		returnFn(item)

		require.Equal(t, original.Update, true)
	})

	t.Run("Should get an item, update it and return it to the pool, while another has to wait", func(t *testing.T) {
		original := &TestItem{
			ID: 1,
			Update: false,
		}

		pool := pools.NewResourcePool([]*TestItem{original}, 10)

		wait := sync.WaitGroup{}
		wait.Add(2)

		// Wait 1 second, then update
		go func() {
			defer wait.Done()

			item, returnFn := pool.Loan(0)
			item.Update = true
			item.ID = 2

			time.Sleep(5 * time.Second)
			returnFn(item)
			
		}()

		go func() {
			defer wait.Done()
			time.Sleep(2 * time.Second)

			item, returnFn := pool.Loan(0)
			require.True(t, item.Update)
			require.Equal(t, item.ID, 2)

			item.ID = 3

			returnFn(item)
		}()

		wait.Wait()
		require.Equal(t, original.ID, 3)
	})

	t.Run("Larger pool size test", func(t *testing.T) {
		items := make([]*TestItem, 0, 100)

		for i := 0; i < 100; i++ {
			items = append(items, &TestItem{
				ID: i,
				Update: false,
			})
		}

		pool := pools.NewResourcePool(items, 100)

		wait := sync.WaitGroup{}
		len := pool.Len()
		wait.Add(len)

		for i := 0; i < len; i++ {
			go func(index int) {
				defer wait.Done()

				item, returnFn := pool.Loan(uint64(index))
				item.Update = true
				returnFn(item)
			}(i)
		}

		wait.Wait()

		for i := 0; i < len; i++ {
			require.True(t, items[i].Update)
		}
	})
}

func Fuzz_ResourcePool_KeyInt(f *testing.F) {
	f.Fuzz(func(t *testing.T, key uint64) {
		original := &TestItem{
			ID: 1,
			Update: false,
		}

		pool := pools.NewResourcePool([]*TestItem{original}, 10)

		item, returnFn := pool.Loan(key)
		item.Update = true
		returnFn(item)

		require.Equal(t, original.Update, true)
	})
}

func Fuzz_ResourcePool_KeyString(f *testing.F) {
	f.Fuzz(func(t *testing.T, key uint64) {
		original := &TestItem{
			ID: 1,
			Update: false,
		}

		pool := pools.NewResourcePool([]*TestItem{original}, 10)

		item, returnFn := pool.Loan(key)
		item.Update = true
		returnFn(item)

		require.Equal(t, original.Update, true)
	})
}

func Fuzz_ResourcePool_KeyBytes(f *testing.F) {
	f.Fuzz(func(t *testing.T, key []byte) {
		original := &TestItem{
			ID: 1,
			Update: false,
		}

		pool := pools.NewResourcePool([]*TestItem{original}, 10)

		item, returnFn := pool.LoanForBytes(key)
		item.Update = true
		returnFn(item)

		require.Equal(t, original.Update, true)
	})
}

func Fuzz_ResourcePool_KeyUint64(f *testing.F) {
	f.Fuzz(func(t *testing.T, key uint64) {
		original := &TestItem{
			ID: 1,
			Update: false,
		}

		pool := pools.NewResourcePool([]*TestItem{original}, 10)

		item, returnFn := pool.LoanForUint64(key)
		item.Update = true
		returnFn(item)

		require.Equal(t, original.Update, true)
	})
}