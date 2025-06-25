package maps

type UnorderedMap[Key comparable, Value any] struct {
	*DefaultAbstractMap[Key, Value]
	m map[Key]Value
}

func NewUnorderedMap[Key comparable, Value any]() *UnorderedMap[Key, Value] {
	um := &UnorderedMap[Key, Value]{
		m: map[Key]Value{},
	}
	um.DefaultAbstractMap = NewDefaultAbstractMap(um)
	return um
}

func (um *UnorderedMap[Key, Value]) Delete(key Key) {
	delete(um.m, key)
}

func (um *UnorderedMap[Key, Value]) Len() int {
	return len(um.m)
}

func (um *UnorderedMap[Key, Value]) Load(key Key) (value Value, ok bool) {
	value, ok = um.m[key]
	return value, ok
}

func (um *UnorderedMap[Key, Value]) Range(f func(key Key, value Value) bool) {
	for k, v := range um.m {
		if !f(k, v) {
			break
		}
	}
}

func (um *UnorderedMap[Key, Value]) Store(key Key, value Value) {
	um.m[key] = value
}
