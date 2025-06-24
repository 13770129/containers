package maps

type AbstractMap[K, V any] interface {
	Clear()
	CompareAndDelete(key K, old V) (deleted bool)
	CompareAndSwap(key K, old, new V) (swapped bool)
	Delete(key K)
	Len() (len int)
	Load(key K) (value V, ok bool)
	LoadAndDelete(key K) (value V, loaded bool)
	LoadOrStore(key K, value V) (actual V, loaded bool)
	Range(f func(key K, value V) bool)
	Keys(f func(key K) bool)
	Values(f func(value V) bool)
	Store(key K, value V)
	Swap(key K, value V) (previous V, loaded bool)
}

type DefaultAbstractMap[K, V any] struct {
	impl AbstractMap[K, V]
}

func NewDefaultAbstractMap[K, V any](impl AbstractMap[K, V]) *DefaultAbstractMap[K, V] {
	return &DefaultAbstractMap[K, V]{
		impl: impl,
	}
}

func (m *DefaultAbstractMap[K, V]) Clear() {
	var keys []K
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

func (m *DefaultAbstractMap[K, V]) Delete(key K) {
	panic("not implemented")
}

func (m *DefaultAbstractMap[K, V]) Len() int {
	var len int
	for range m.impl.Range {
		len += 1
	}
	return len
}

func (m *DefaultAbstractMap[K, V]) Load(key K) (value V, ok bool) {
	panic("not implemented")
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

func (m *DefaultAbstractMap[K, V]) Range(f func(key K, value V) bool) {
	panic("not implemented")
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
