package lru

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	if _, err := New(1); err != nil {
		t.Error("Ошибки в создании кэша быть не должно")
	}
	if _, err := New(0); err == nil {
		t.Error("Должна возвращаться ошибка создания кэша")
	}
}

func TestLRU_Add(t *testing.T) {
	l, _ := New(1)
	assert.False(t, l.Add(1, 1), "Вытеснение быть не должно")
	assert.False(t, l.Add(1, 2), "Вытеснение быть не должно")
	assert.True(t, l.Add(2, 1), "Должно произойти вытеснение")
}

func TestLRU_Get(t *testing.T) {
	l, _ := New(1)
	l.Add(1, 1)
	v, ok := l.Get(1)
	assert.True(t, ok)
	assert.Equal(t, v, 1)
}

func TestLRU_Add_Get(t *testing.T) {
	l, _ := New(1)
	l.Add(1, 1)
	l.Add(1, 2)
	v, _ := l.Get(1)
	assert.Equal(t, v, 2)
}

func TestLRU_Len(t *testing.T) {
	l, _ := New(2)
	l.Add(1, 1)
	assert.Equal(t, l.Len(), 1, "Не правильный расчет размера")
}

func TestLRU_Keys(t *testing.T) {
	t.SkipNow()
	l, _ := New(2)
	l.Add(1, 1)
	l.Add(2, 2)
	assert.Equal(t, len(l.Keys()), 2)
}

func TestLRU_Purge(t *testing.T) {
	l, _ := New(2)
	l.Add(1, 1)
	l.Add(2, 2)
	l.Purge()
	assert.Equal(t, l.Len(), 0)
	assert.Equal(t, len(l.elements), 0)
}

func TestLRU_Contains(t *testing.T) {
	l, _ := New(2)
	l.Add(1, 1)
	l.Add(2, 2)
	assert.True(t, l.Contains(1))
	assert.True(t, l.Contains(2))
	assert.False(t, l.Contains(3))
}

func TestLRU_Remove(t *testing.T) {
	l,_ := New(2)
	l.Add(1, 1)
	l.Add(2, 2)
	assert.True(t, l.Contains(2))
	l.Remove(2)
	assert.False(t, l.Contains(2))
}

func TestLRU_GetOldest(t *testing.T) {
	l, _ := New(3)
	k, _, _ := l.GetOldest()
	assert.Nil(t, k)
	l.Add(1, 1)
	l.Add(2, 2)
	l.Add(3, 3)
	k, _, _ = l.GetOldest()
	assert.Equal(t, k, 1)
}

func TestLRU_RemoveOldest(t *testing.T) {
	l, _ := New(3)
	l.Add(1, 1)
	l.Add(2, 2)
	l.Add(3, 3)
	k, v, ok := l.RemoveOldest()
	assert.Equal(t, k, 1)
	assert.Equal(t, v, 1)
	assert.True(t, ok)
}

// Комплексная проверка работы кеша LRU
func TestLRU (t *testing.T) {
	l, err := New(100)
	assert.NoError(t, err)

	for i := 0; i < 200; i++ {
		l.Add(i, i)
	}
	assert.Equal(t, l.Len(), 100)

	assert.Equal(t, len(l.Keys()), 100)
	for i, k := range l.Keys() {
		v, ok := l.Get(k)
		assert.True(t, ok)
		assert.NotNil(t, v)
		assert.Equal(t, k, v)
		// Проверяем, что новые ключи, вытеснили старые
		assert.Equal(t, k, 200-1-i)
	}

	for i := 0; i < 100; i++ {
		_, ok := l.Get(i)
		assert.False(t, ok)
	}
	for i := 100; i < 200; i++ {
		_, ok := l.Get(i)
		assert.True(t, ok)
	}

	l.Purge()
	assert.Equal(t, l.Len(), 0)
	v, ok := l.Get(150)
	assert.False(t, ok)
	assert.Nil(t, v)
}