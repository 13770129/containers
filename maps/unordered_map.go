package maps

type UnorderedMap[K comparable, V any] struct {
	*DefaultAbstractMap[K, V]
	m map[K]V
}

func NewUnorderedMap[K comparable, V any]() *UnorderedMap[K, V] {
	m := &UnorderedMap[K, V]{
		m: map[K]V{},
	}
	m.DefaultAbstractMap = NewDefaultAbstractMap(m)
	return m
}

func (m *UnorderedMap[K, V]) Delete(key K) {
	delete(m.m, key)
}

func (m *UnorderedMap[K, V]) Len() int {
	return len(m.m)
}

func (m *UnorderedMap[K, V]) Load(key K) (value V, ok bool) {
	value, ok = m.m[key]
	return value, ok
}

func (m *UnorderedMap[K, V]) Range(f func(key K, value V) bool) {
	for k, v := range m.m {
		if !f(k, v) {
			break
		}
	}
}

func (m *UnorderedMap[K, V]) Store(key K, value V) {
	m.m[key] = value
}
