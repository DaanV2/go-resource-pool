package pools_test

import (
	"math/rand"
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

		item, returnFn := pool.LoanByString("dir/file.txt")
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

func Test_ResourcePool_Deadlocks(t *testing.T) {
	t.Run("Larger pool size test nr2.", func(t *testing.T) {
		amount := 10
		items := make([]*TestItem, 0, amount)

		for i := 0; i < amount; i++ {
			items = append(items, &TestItem{
				ID: i,
				Update: false,
			})
		}

		pool := pools.NewResourcePool(items, amount)

		wait := sync.WaitGroup{}
		wait.Add(100)

		seed := 123456789

		// 1000 go routines, access a random item from the pool, sleep and update it
		// This should not cause any deadlocks
		for i := 0; i < 100; i++ {
			go func(index int) {
				defer wait.Done()
				cs := seed * index
				cs = cs ^ seed
				rnd := rand.New(rand.NewSource(int64(cs)))

				for j := 0; j < amount; j++ {
					index := rnd.Uint64() % uint64(amount)
					item, returnFn := pool.Loan(index)
					item.Update = true
					time.Sleep(100 * time.Millisecond)
					returnFn(item)
				}
			}(i)
		}

		wait.Wait()
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

		item, returnFn := pool.LoanByBytes(key)
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

		item, returnFn := pool.LoanByUint64(key)
		item.Update = true
		returnFn(item)

		require.Equal(t, original.Update, true)
	})
}