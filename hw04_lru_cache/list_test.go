package hw04lrucache

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestList(t *testing.T) {
	t.Run("empty list", func(t *testing.T) {
		l := NewList()

		require.Equal(t, l.Len(), 0)
		require.Nil(t, l.Front())
		require.Nil(t, l.Back())
	})

	t.Run("complex", func(t *testing.T) {
		l := NewList()

		l.PushFront(10) // [10]
		l.PushBack(20)  // [10, 20]
		l.PushBack(30)  // [10, 20, 30]
		require.Equal(t, 3, l.Len())

		middle := l.Front().Next // 20
		l.Remove(middle)         // [10, 30]
		require.Equal(t, 2, l.Len())

		for i, v := range [...]int{40, 50, 60, 70, 80} {
			if i%2 == 0 {
				l.PushFront(v)
			} else {
				l.PushBack(v)
			}
		} // [80, 60, 40, 10, 30, 50, 70]

		require.Equal(t, 7, l.Len())
		require.Equal(t, 80, l.Front().Value)
		require.Equal(t, 70, l.Back().Value)

		l.MoveToFront(l.Front()) // [80, 60, 40, 10, 30, 50, 70]
		l.MoveToFront(l.Back())  // [70, 80, 60, 40, 10, 30, 50]

		elems := make([]int, 0, l.Len())
		for i := l.Front(); i != nil; i = i.Next {
			elems = append(elems, i.Value.(int))
		}
		require.Equal(t, []int{70, 80, 60, 40, 10, 30, 50}, elems)
	})
}

func TestPushFront(t *testing.T) {
	t.Run("empty list", func(t *testing.T) {
		l := NewList()
		val := 10
		front := l.PushFront(val) // [10]
		require.Equal(t, 1, l.Len())
		require.Equal(t, val, l.Front().Value)
		require.Equal(t, val, l.Back().Value)
		require.Equal(t, front, l.Front())
		require.Nil(t, front.Next)
		require.Nil(t, front.Prev)
	})

	t.Run("list with single element", func(t *testing.T) {
		l := NewList()
		val1 := 10
		val2 := 20
		l.PushFront(val1)          // [10]
		front := l.PushFront(val2) // [20, 10]
		require.Equal(t, 2, l.Len())
		require.Equal(t, val2, l.Front().Value)
		require.Equal(t, val1, l.Back().Value)
		require.Equal(t, val1, l.Front().Next.Value)
		require.Equal(t, val2, l.Back().Prev.Value)
		require.Equal(t, front, l.Front())
	})

	t.Run("list with multiple elements", func(t *testing.T) {
		l := NewList()
		val1 := 10
		val2 := 20
		val3 := 30
		l.PushFront(val1)          // [10]
		l.PushFront(val2)          // [20, 10]
		front := l.PushFront(val3) // [30, 20, 10]
		require.Equal(t, 3, l.Len())
		require.Equal(t, val3, l.Front().Value)
		require.Equal(t, val2, l.Front().Next.Value)
		require.Equal(t, val1, l.Front().Next.Next.Value)
		require.Equal(t, val1, l.Back().Value)
		require.Equal(t, val2, l.Back().Prev.Value)
		require.Equal(t, val3, l.Back().Prev.Prev.Value)
		require.Equal(t, front, l.Front())
	})
}

func TestPushBack(t *testing.T) {
	t.Run("empty list", func(t *testing.T) {
		l := NewList()
		val := 10
		back := l.PushBack(val) // [10]
		require.Equal(t, 1, l.Len())
		require.Equal(t, val, l.Front().Value)
		require.Equal(t, val, l.Back().Value)
		require.Equal(t, back, l.Back())
		require.Nil(t, l.Back().Next)
		require.Nil(t, l.Back().Prev)
	})

	t.Run("list with single element", func(t *testing.T) {
		l := NewList()
		val1 := 10
		val2 := 20
		l.PushFront(val1)        // [10]
		back := l.PushBack(val2) // [10, 20]
		require.Equal(t, 2, l.Len())
		require.Equal(t, back, l.Back())
		require.Equal(t, val2, l.Back().Value)
		require.Equal(t, val1, l.Front().Value)
		require.Equal(t, val1, l.Back().Prev.Value)
		require.Equal(t, val2, l.Front().Next.Value)
		require.Nil(t, l.Back().Next)
		require.Nil(t, l.Front().Prev)
	})

	t.Run("list with multiple elements", func(t *testing.T) {
		l := NewList()
		val1 := 10
		val2 := 20
		val3 := 30
		l.PushFront(val1)        // [10]
		l.PushFront(val2)        // [20, 10]
		back := l.PushBack(val3) // [20, 10, 30]
		require.Equal(t, 3, l.Len())
		require.Equal(t, back, l.Back())
		require.Equal(t, val3, l.Back().Value)
		require.Equal(t, val1, l.Back().Prev.Value)
		require.Equal(t, val2, l.Front().Value)
		require.Equal(t, val1, l.Front().Next.Value)
		require.Nil(t, l.Back().Next)
		require.Nil(t, l.Front().Prev)
	})
}

func TestRemove(t *testing.T) {
	t.Run("last element", func(t *testing.T) {
		l := NewList()
		val := 10
		toRemove := l.PushFront(val) // [10]
		l.Remove(toRemove)           // []
		require.Equal(t, l.Len(), 0)
		require.Nil(t, l.Front())
		require.Nil(t, l.Back())
	})

	t.Run("first of two", func(t *testing.T) {
		l := NewList()
		val1 := 10
		val2 := 20
		l.PushFront(val1)             // [10]
		toRemove := l.PushFront(val2) // [20, 10]
		l.Remove(toRemove)            // [10]
		require.Equal(t, l.Len(), 1)
		require.Equal(t, l.Front().Value, val1)
		require.Equal(t, l.Back().Value, val1)
		require.Nil(t, l.Front().Next)
		require.Nil(t, l.Front().Prev)
		require.Nil(t, l.Back().Next)
		require.Nil(t, l.Back().Prev)
	})

	t.Run("last of two", func(t *testing.T) {
		l := NewList()
		val1 := 10
		val2 := 20
		toRemove := l.PushFront(val1) // [10]
		l.PushFront(val2)             // [20, 10]
		l.Remove(toRemove)            // [20]
		require.Equal(t, l.Len(), 1)
		require.Equal(t, l.Front().Value, val2)
		require.Equal(t, l.Back().Value, val2)
		require.Nil(t, l.Front().Next)
		require.Nil(t, l.Front().Prev)
		require.Nil(t, l.Back().Next)
		require.Nil(t, l.Back().Prev)
	})

	t.Run("in the middle", func(t *testing.T) {
		l := NewList()
		val1 := 10
		val2 := 20
		val3 := 30
		l.PushFront(val1)             // [10]
		toRemove := l.PushFront(val2) // [20, 10]
		l.PushFront(val3)             // [30, 20, 10]
		l.Remove(toRemove)            // [30, 10]
		require.Equal(t, l.Len(), 2)
		require.Equal(t, l.Front().Value, val3)
		require.Equal(t, l.Back().Value, val1)
		require.Equal(t, l.Back(), l.Front().Next)
		require.Nil(t, l.Front().Prev)
		require.Equal(t, l.Front(), l.Back().Prev)
		require.Nil(t, l.Back().Next)
	})

	t.Run("first of many", func(t *testing.T) {
		l := NewList()
		val1 := 10
		val2 := 20
		val3 := 30
		l.PushFront(val1)             // [10]
		l.PushFront(val2)             // [20, 10]
		toRemove := l.PushFront(val3) // [30, 20, 10]
		l.Remove(toRemove)            // [20, 10]
		require.Equal(t, l.Len(), 2)
		require.Equal(t, l.Front().Value, val2)
		require.Equal(t, l.Back().Value, val1)
		require.Equal(t, l.Back(), l.Front().Next)
		require.Nil(t, l.Front().Prev)
		require.Equal(t, l.Front(), l.Back().Prev)
		require.Nil(t, l.Back().Next)
	})

	t.Run("last of many", func(t *testing.T) {
		l := NewList()
		val1 := 10
		val2 := 20
		val3 := 30
		toRemove := l.PushFront(val1) // [10]
		l.PushFront(val2)             // [20, 10]
		l.PushFront(val3)             // [30, 20, 10]
		l.Remove(toRemove)            // [30, 20]
		require.Equal(t, l.Len(), 2)
		require.Equal(t, l.Front().Value, val3)
		require.Equal(t, l.Back().Value, val2)
		require.Equal(t, l.Back(), l.Front().Next)
		require.Nil(t, l.Front().Prev)
		require.Equal(t, l.Front(), l.Back().Prev)
		require.Nil(t, l.Back().Next)
	})
}

func TestMoveToFront(t *testing.T) {
	t.Run("single element", func(t *testing.T) {
		l := NewList()
		val := 10
		toMove := l.PushFront(val) // [10]
		l.MoveToFront(toMove)      // [10]
		require.Equal(t, l.Len(), 1)
		require.Equal(t, l.Front(), l.Back())
		require.Nil(t, l.Front().Next)
		require.Nil(t, l.Front().Prev)
		require.Nil(t, l.Back().Next)
		require.Nil(t, l.Back().Prev)
	})

	t.Run("first of two", func(t *testing.T) {
		l := NewList()
		val1 := 10
		val2 := 20
		l.PushFront(val1)           // [10]
		toMove := l.PushFront(val2) // [20, 10]
		l.MoveToFront(toMove)       // [20, 10]
		require.Equal(t, l.Len(), 2)
		require.Equal(t, l.Front().Value, val2)
		require.Equal(t, l.Back().Value, val1)
		require.Equal(t, l.Back(), l.Front().Next)
		require.Nil(t, l.Front().Prev)
		require.Equal(t, l.Front(), l.Back().Prev)
		require.Nil(t, l.Back().Next)
	})

	t.Run("last of two", func(t *testing.T) {
		l := NewList()
		val1 := 10
		val2 := 20
		toMove := l.PushFront(val1) // [10]
		l.PushFront(val2)           // [20, 10]
		l.MoveToFront(toMove)       // [10, 20]
		require.Equal(t, l.Len(), 2)
		require.Equal(t, l.Front().Value, val1)
		require.Equal(t, l.Back().Value, val2)
		require.Equal(t, l.Back(), l.Front().Next)
		require.Nil(t, l.Front().Prev)
		require.Equal(t, l.Front(), l.Back().Prev)
		require.Nil(t, l.Back().Next)
	})

	t.Run("from the middle", func(t *testing.T) {
		l := NewList()
		val1 := 10
		val2 := 20
		val3 := 30
		l.PushFront(val1)           // [10]
		toMove := l.PushFront(val2) // [20, 10]
		l.PushFront(val3)           // [30, 20, 10]
		l.MoveToFront(toMove)       // [20, 30, 10]
		require.Equal(t, l.Len(), 3)
		require.Equal(t, l.Front().Value, val2)
		require.Equal(t, l.Back().Value, val1)
		require.Equal(t, l.Front().Next.Value, val3)
		require.Equal(t, l.Back().Prev.Value, val3)
	})

	t.Run("first of many", func(t *testing.T) {
		l := NewList()
		val1 := 10
		val2 := 20
		val3 := 30
		l.PushFront(val1)           // [10]
		l.PushFront(val2)           // [20, 10]
		toMove := l.PushFront(val3) // [30, 20, 10]
		l.MoveToFront(toMove)       // [30, 20, 10]
		require.Equal(t, l.Len(), 3)
		require.Equal(t, l.Front().Value, val3)
		require.Equal(t, l.Back().Value, val1)
		require.Equal(t, l.Front().Next.Value, val2)
		require.Equal(t, l.Back().Prev.Value, val2)
	})

	t.Run("last of many", func(t *testing.T) {
		l := NewList()
		val1 := 10
		val2 := 20
		val3 := 30
		toMove := l.PushFront(val1) // [10]
		l.PushFront(val2)           // [20, 10]
		l.PushFront(val3)           // [30, 20, 10]
		l.MoveToFront(toMove)       // [10, 30, 20]
		require.Equal(t, l.Len(), 3)
		require.Equal(t, l.Front().Value, val1)
		require.Equal(t, l.Back().Value, val2)
		require.Equal(t, l.Front().Next.Value, val3)
		require.Equal(t, l.Back().Prev.Value, val3)
	})
}
