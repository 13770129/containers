package unordered_map

import "github.com/13770129/containers/abstract_map"

type UnorderedMap[K comparable, V any] struct {
	abstract_map.DefaultAbstractMap[K, V]
	m map[K]V
}
