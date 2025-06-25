package maps

type MapOps[Key, Value any] interface {
	Delete(key Key)
	Load(key Key) (value Value, ok bool)
	Range(f func(key Key, value Value) bool)
	Store(key Key, value Value)
}

type AbstractMap[Key, Value any] interface {
	MapOps[Key, Value]
	Clear()
	CompareAndDelete(key Key, old Value) (deleted bool)
	CompareAndSwap(key Key, old, new Value) (swapped bool)
	Len() int
	LoadAndDelete(key Key) (value Value, loaded bool)
	LoadOrStore(key Key, value Value) (actual Value, loaded bool)
	Keys(f func(key Key) bool)
	Values(f func(value Value) bool)
	Swap(key Key, value Value) (previous Value, loaded bool)
}

func FromGoMaps[Key comparable, Value any](m AbstractMap[Key, Value], gms ...map[Key]Value) {
	for _, gm := range gms {
		for k, v := range gm {
			m.Store(k, v)
		}
	}
}

func FromAbstractMap[Key, Value any](m AbstractMap[Key, Value], ams ...AbstractMap[Key, Value]) {
	for _, am := range ams {
		for k, v := range am.Range {
			m.Store(k, v)
		}
	}
}

type DefaultAbstractMap[Key, Value any] struct {
	impl AbstractMap[Key, Value]
}

func NewDefaultAbstractMap[Key, Value any](impl AbstractMap[Key, Value]) *DefaultAbstractMap[Key, Value] {
	return &DefaultAbstractMap[Key, Value]{
		impl: impl,
	}
}

func (m *DefaultAbstractMap[Key, Value]) Clear() {
	var keys []Key
	for key := range m.impl.Range {
		keys = append(keys, key)
	}
	for _, key := range keys {
		m.impl.Delete(key)
	}
}

func (m *DefaultAbstractMap[K, V]) CompareAndDelete(key K, old V) (deleted bool) {
	value, ok := m.impl.Load(key)
	if !ok {
		return false
	}
	// Compare using interface{} since we can't assume comparable types
	if any(value) == any(old) {
		m.impl.Delete(key)
		return true
	}
	return false
}

func (m *DefaultAbstractMap[K, V]) CompareAndSwap(key K, old, new V) (swapped bool) {
	value, ok := m.impl.Load(key)
	if !ok {
		return false
	}
	// Compare using interface{} since we can't assume comparable types
	if any(value) == any(old) {
		m.impl.Store(key, new)
		return true
	}
	return false
}

func (m *DefaultAbstractMap[K, V]) Len() int {
	var len int
	for range m.impl.Range {
		len += 1
	}
	return len
}

func (m *DefaultAbstractMap[K, V]) LoadAndDelete(key K) (value V, loaded bool) {
	value, loaded = m.impl.Load(key)
	if loaded {
		m.impl.Delete(key)
	}
	return value, loaded
}

func (m *DefaultAbstractMap[K, V]) LoadOrStore(key K, value V) (actual V, loaded bool) {
	actual, loaded = m.impl.Load(key)
	if !loaded {
		m.impl.Store(key, value)
		actual = value
	}
	return actual, loaded
}

func (m *DefaultAbstractMap[K, V]) Keys(f func(key K) bool) {
	for key := range m.impl.Range {
		if !f(key) {
			break
		}
	}
}

func (m *DefaultAbstractMap[K, V]) Values(f func(value V) bool) {
	for _, value := range m.impl.Range {
		if !f(value) {
			break
		}
	}
}

func (m *DefaultAbstractMap[K, V]) Swap(key K, value V) (previous V, loaded bool) {
	previous, loaded = m.impl.Load(key)
	m.impl.Store(key, value)
	return previous, loaded
}
