package lru

import (
	"container/list"
	"errors"
)

type (
	LRU struct {
		size     int
		list     *list.List
		elements map[interface{}]*list.Element
	}
	item struct {
		key   interface{}
		value interface{}
	}
)

func New(size int) (*LRU, error) {
	if size <= 0 {
		return nil, errors.New("Размер кэша должен быть больше 0")
	}
	l := &LRU{
		size:     size,
		list:     list.New(),
		elements: make(map[interface{}]*list.Element),
	}
	return l, nil
}

// Add добавляет значение value с ключом key в кеш. Если при добавлении
// вытесняется элемент из кеша, то возвращает true
func (l *LRU) Add(key, value interface{}) bool {
	if el, ok := l.elements[key]; ok {
		// Элемент уже в кэше. Обновляем его
		l.list.MoveToFront(el)
		el.Value.(*item).value = value
		return false
	}
	it := &item{
		key: key,
		value: value,
	}
	el := l.list.PushFront(it)
	l.elements[key] = el
	if len(l.elements) > l.size {
		// Вытесняем старый элемент из кэша
		l.RemoveOldest()
		return true
	}
	return false
}

func (l *LRU) Len() int {
	return l.list.Len()
}

func (l *LRU) Keys() []interface{} {
	keys := make([]interface{}, len(l.elements))
	i := 0
	for e := l.list.Front(); e != nil; e = e.Next() {
		keys[i] = e.Value.(*item).key
		i++
	}
	return keys
}

// Purge полностью очищает кэш
func (l *LRU) Purge() {
	l.list.Init()
	l.elements = make(map[interface{}]*list.Element)
}

// Get возращает значение value по ключу key. Если ключ найден,
// освежает его в кэше
func (l *LRU) Get(key interface{}) (value interface{}, ok bool) {
	if el, ok := l.elements[key]; ok {
		l.list.MoveToFront(el)
		return el.Value.(*item).value, true
	}
	return
}

func (l *LRU) Contains(key interface{}) bool {
	_, ok := l.elements[key]
	return ok
}

func (l * LRU) Remove(key interface{}) bool {
	if el, ok := l.elements[key]; ok {
		l.removeElement(el)
		return true
	}
	return false
}

func (l *LRU) removeElement(el *list.Element) {
	l.list.Remove(el)
	delete(l.elements, el.Value.(*item).key)
}

func (l *LRU) GetOldest() (key, value interface{}, ok bool) {
	if old := l.list.Back(); old != nil {
		item := old.Value.(*item)
		return item.key, item.value, true
	}
	return
}

func (l *LRU) RemoveOldest() (key, value interface{}, ok bool) {
	k, v, ok := l.GetOldest()
	if ok {
		ok = l.Remove(k)
		return k, v, ok
	}
	return
}