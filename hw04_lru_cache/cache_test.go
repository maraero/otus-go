package hw04lrucache

import (
	"math/rand"
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCache(t *testing.T) {
	t.Run("empty cache", func(t *testing.T) {
		c := NewCache(10)
		_, ok := c.Get("aaa")
		require.False(t, ok)
		_, ok = c.Get("bbb")
		require.False(t, ok)
	})

	t.Run("add single value", func(t *testing.T) {
		c := NewCache(10)
		success := c.Set("aaa", 10) // [{aaa, 10}] 
		require.False(t, success)
		val, ok := c.Get("aaa") // [{aaa, 10}]
		require.True(t, ok)
		require.Equal(t, val, 10)
	})

	t.Run("update single value", func(t *testing.T) {
		c := NewCache(10)
		c.Set("aaa", 10) // [{aaa, 10}]
		c.Set("aaa", 20) // [{aaa, 20}]
		val, ok := c.Get("aaa") // [{aaa, 20}]
		require.True(t, ok)
		require.Equal(t, val, 20)
	})

	t.Run("update value", func(t *testing.T) {
		c := NewCache(10)
		c.Set("aaa", 10) // [{aaa, 10}]
		c.Set("bbb", 20) // [{bbb, 20}, {aaa: 10}]
		c.Set("aaa", 30) // [{aaa, 30}, {bbb: 20}]
		val, ok := c.Get("aaa") // [{aaa, 30}, {bbb: 20}]
		require.True(t, ok)
		require.Equal(t, val, 30)
	})

	t.Run("exceed capacity", func(t *testing.T) {
		c := NewCache(1)
		c.Set("aaa", 10) // [{aaa, 10}]
		c.Set("bbb", 20) // [{bbb, 20}]
		_, ok := c.Get("aaa") // [{bbb, 20}] 
		require.False(t, ok)
		val, ok := c.Get("bbb") // [{bbb, 20}]
		require.True(t, ok)
		require.Equal(t, val, 20)
	})

	t.Run("exceed capacity in the cache with values update", func(t *testing.T) {
		t.Skip()
		c := NewCache(2)
		c.Set("a", 10) // [{a: 10}]
		c.Set("b", 20) // [{b: 20}, {a: 10}]
		c.Set("b", 30) // [{b: 30}, {a: 10}]
		c.Set("a", 40) // [{a: 40}, {b: 30}]
		c.Set("b", 50) // [{b: 50}, {a: 40}]
		c.Set("c", 60) // [{c: 60}, {b: 50}]
		_, ok := c.Get("a")
		require.False(t, ok)
		val, ok := c.Get("b")
		require.True(t, ok)
		require.Equal(t, val, 50)
		val, ok = c.Get("c")
		require.True(t, ok)
		require.Equal(t, val, 60)
	})

	t.Run("simple", func(t *testing.T) {
		t.Skip()
		c := NewCache(5)

		wasInCache := c.Set("aaa", 100) // [{aaa, 100}]
		require.False(t, wasInCache)

		wasInCache = c.Set("bbb", 200) // [{bbb, 200}, {aaa, 100}]
		require.False(t, wasInCache)

		val, ok := c.Get("aaa") // [{aaa: 100}, {bbb:200}]
		require.True(t, ok)
		require.Equal(t, 100, val)

		val, ok = c.Get("bbb") // [{bbb, 200}, {aaa, 100}]
		require.True(t, ok)
		require.Equal(t, 200, val)

		wasInCache = c.Set("aaa", 300) // [{aaa: 300}, {bbb:200}]
		require.True(t, wasInCache)

		val, ok = c.Get("aaa") // [{aaa: 300}, {bbb:200}]
		require.True(t, ok)
		require.Equal(t, 300, val)

		val, ok = c.Get("ccc") // [{aaa: 300}, {bbb:200}]
		require.False(t, ok)
		require.Nil(t, val)
	})

	t.Run("clear cache", func(t *testing.T) {
		c := NewCache(2)
		c.Set("a", 10)
		c.Set("b", 20)
		c.Clear()
		_, ok := c.Get("a")
		require.False(t, ok)
		_, ok = c.Get("b")
		require.False(t, ok)
	})
}

func TestCacheMultithreading(t *testing.T) {
	t.Skip() // Remove me if task with asterisk completed.

	c := NewCache(10)
	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Set(Key(strconv.Itoa(i)), i)
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Get(Key(strconv.Itoa(rand.Intn(1_000_000))))
		}
	}()

	wg.Wait()
}
