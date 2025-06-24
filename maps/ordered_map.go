package maps

import "container/list"

// entry represents a key-value pair stored in the linked list.
// This structure allows us to store both key and value together,
// enabling efficient iteration while maintaining map semantics.
type entry[K, V any] struct {
	key   K
	value V
}

// OrderedMap implements AbstractMap with insertion order preservation.
// It uses a hybrid approach: a standard Go map for O(1) key lookups
// combined with a doubly-linked list to maintain insertion order.
// This provides O(1) performance for Store, Load, and Delete operations
// while ensuring Range operations iterate in insertion order.
type OrderedMap[K comparable, V any] struct {
	*DefaultAbstractMap[K, V]
	m map[K]*list.Element // Maps keys to their corresponding list elements
	l *list.List          // Doubly-linked list maintaining insertion order
}

// NewOrderedMap creates a new OrderedMap instance.
// The map is initialized empty with no memory pre-allocation,
// allowing it to grow dynamically as items are added.
func NewOrderedMap[K comparable, V any]() *OrderedMap[K, V] {
	om := &OrderedMap[K, V]{
		m: make(map[K]*list.Element),
		l: list.New(),
	}
	// Embed DefaultAbstractMap to inherit common functionality
	// like CompareAndSwap, LoadOrStore, etc.
	om.DefaultAbstractMap = NewDefaultAbstractMap(om)
	return om
}

// Store adds or updates a key-value pair in the map.
// If the key already exists, its value is updated in-place
// without changing its position in the iteration order.
// If the key is new, it's appended to the end of the order.
// Time complexity: O(1)
func (om *OrderedMap[K, V]) Store(key K, value V) {
	if element, exists := om.m[key]; exists {
		// Key exists: update value in-place, preserving order position
		element.Value.(*entry[K, V]).value = value
	} else {
		// New key: create entry and append to end of list
		newEntry := &entry[K, V]{key: key, value: value}
		element := om.l.PushBack(newEntry)
		om.m[key] = element
	}
}

// Load retrieves the value associated with a key.
// Returns the value and true if the key exists,
// or the zero value and false if the key doesn't exist.
// Time complexity: O(1)
func (om *OrderedMap[K, V]) Load(key K) (value V, ok bool) {
	if element, exists := om.m[key]; exists {
		return element.Value.(*entry[K, V]).value, true
	}
	// Return zero value for type V when key not found
	var zero V
	return zero, false
}

// Delete removes a key-value pair from the map.
// If the key exists, it's removed from both the map and the list.
// If the key doesn't exist, this operation is a no-op.
// Time complexity: O(1)
func (om *OrderedMap[K, V]) Delete(key K) {
	if element, exists := om.m[key]; exists {
		// Remove from both data structures atomically
		delete(om.m, key)
		om.l.Remove(element)
	}
}

// Len returns the number of key-value pairs in the map.
// This leverages the built-in map's length for O(1) performance
// rather than counting list elements.
// Time complexity: O(1)
func (om *OrderedMap[K, V]) Len() int {
	return len(om.m)
}

// Range calls the provided function for each key-value pair in insertion order.
// The iteration stops early if the function returns false.
// This method traverses the linked list, ensuring consistent insertion order
// regardless of Go's map iteration randomization.
// Time complexity: O(n) where n is the number of elements
func (om *OrderedMap[K, V]) Range(f func(key K, value V) bool) {
	// Iterate through linked list to maintain insertion order
	for element := om.l.Front(); element != nil; element = element.Next() {
		entry := element.Value.(*entry[K, V])
		if !f(entry.key, entry.value) {
			// Function returned false, stop iteration early
			break
		}
	}
}
