package maps_test

import (
	"fmt"
	"slices"
	"testing"

	"github.com/13770129/containers/maps"
)

// MapFactory creates new map instances for testing different implementations.
// This abstraction enables testing multiple map implementations with identical test logic.
type MapFactory[K, V any] func() maps.AbstractMap[K, V]

// testSuite executes comprehensive tests against any AbstractMap implementation.
// The factory parameter enables different implementations to be tested with the same suite.
func testSuite[K comparable, V comparable](t *testing.T, factory MapFactory[K, V], testData []TestCase[K, V]) {
	t.Helper()

	t.Run("BasicOperations", func(t *testing.T) {
		testBasicOperations(t, factory, testData)
	})

	t.Run("AdvancedOperations", func(t *testing.T) {
		testAdvancedOperations(t, factory, testData)
	})

	t.Run("IterationOperations", func(t *testing.T) {
		testIterationOperations(t, factory, testData)
	})

	t.Run("EdgeCases", func(t *testing.T) {
		testEdgeCases(t, factory, testData)
	})
}

// TestCase represents a key-value pair for testing.
// Generic design allows the same test structure to work with different data types.
type TestCase[K, V any] struct {
	Key   K
	Value V
}

// testBasicOperations verifies fundamental map operations: Store, Load, Delete, Len.
func testBasicOperations[K comparable, V comparable](t *testing.T, factory MapFactory[K, V], testData []TestCase[K, V]) {
	t.Run("StoreAndLoad", func(t *testing.T) {
		m := factory()

		for _, tc := range testData {
			m.Store(tc.Key, tc.Value)

			if value, ok := m.Load(tc.Key); !ok {
				t.Errorf("Expected to find key %v after storing", tc.Key)
			} else if value != tc.Value {
				t.Errorf("Expected value %v for key %v, got %v", tc.Value, tc.Key, value)
			}
		}

		expectedLen := len(testData)
		if actualLen := m.Len(); actualLen != expectedLen {
			t.Errorf("Expected length %d, got %d", expectedLen, actualLen)
		}
	})

	t.Run("LoadNonExistentKey", func(t *testing.T) {
		m := factory()

		if len(testData) > 0 {
			value, ok := m.Load(testData[0].Key)
			if ok {
				t.Errorf("Expected ok=false for non-existent key, got ok=true with value %v", value)
			}
		}
	})

	t.Run("Delete", func(t *testing.T) {
		m := factory()

		for _, tc := range testData {
			m.Store(tc.Key, tc.Value)
		}

		for _, tc := range testData {
			m.Delete(tc.Key)

			if value, ok := m.Load(tc.Key); ok {
				t.Errorf("Expected key %v to be deleted, but found value %v", tc.Key, value)
			}
		}

		if length := m.Len(); length != 0 {
			t.Errorf("Expected empty map after deleting all keys, got length %d", length)
		}
	})

	t.Run("Clear", func(t *testing.T) {
		m := factory()

		for _, tc := range testData {
			m.Store(tc.Key, tc.Value)
		}

		m.Clear()

		if length := m.Len(); length != 0 {
			t.Errorf("Expected length 0 after Clear(), got %d", length)
		}

		for _, tc := range testData {
			if value, ok := m.Load(tc.Key); ok {
				t.Errorf("Expected key %v to be cleared, but found value %v", tc.Key, value)
			}
		}
	})
}

// testAdvancedOperations verifies atomic-style operations that require careful implementation.
func testAdvancedOperations[K comparable, V comparable](t *testing.T, factory MapFactory[K, V], testData []TestCase[K, V]) {
	if len(testData) == 0 {
		t.Skip("No test data provided for advanced operations")
		return
	}

	firstCase := testData[0]

	t.Run("LoadOrStore", func(t *testing.T) {
		m := factory()

		// First call should store the value since key doesn't exist
		actual, loaded := m.LoadOrStore(firstCase.Key, firstCase.Value)
		if loaded {
			t.Error("Expected loaded=false for first LoadOrStore call")
		}
		if actual != firstCase.Value {
			t.Errorf("Expected actual value %v, got %v", firstCase.Value, actual)
		}

		// Second call should load the existing value without storing
		differentValue := firstCase.Value
		actual, loaded = m.LoadOrStore(firstCase.Key, differentValue)
		if !loaded {
			t.Error("Expected loaded=true for second LoadOrStore call")
		}
		if actual != firstCase.Value {
			t.Errorf("Expected actual value %v (original), got %v", firstCase.Value, actual)
		}
	})

	t.Run("LoadAndDelete", func(t *testing.T) {
		m := factory()
		m.Store(firstCase.Key, firstCase.Value)

		// Should return existing value and remove it atomically
		value, loaded := m.LoadAndDelete(firstCase.Key)
		if !loaded {
			t.Error("Expected loaded=true for existing key")
		}
		if value != firstCase.Value {
			t.Errorf("Expected value %v, got %v", firstCase.Value, value)
		}

		// Key should no longer exist
		if _, ok := m.Load(firstCase.Key); ok {
			t.Error("Expected key to be deleted after LoadAndDelete")
		}

		// Second call on non-existent key should return not loaded
		_, loaded = m.LoadAndDelete(firstCase.Key)
		if loaded {
			t.Error("Expected loaded=false for non-existent key")
		}
	})

	t.Run("Swap", func(t *testing.T) {
		m := factory()

		// Swap on non-existent key should store new value
		newValue := firstCase.Value
		previous, loaded := m.Swap(firstCase.Key, newValue)
		if loaded {
			t.Error("Expected loaded=false for swap on empty map")
		}
		// Since this is a non-existent key, previous should be the zero value
		var zero V
		if previous != zero {
			t.Errorf("Expected zero value for previous on non-existent key, got %v", previous)
		}

		if value, ok := m.Load(firstCase.Key); !ok || value != newValue {
			t.Errorf("Expected new value %v to be stored", newValue)
		}

		// Swap on existing key should return previous value
		if len(testData) > 1 {
			secondValue := testData[1].Value
			previous, loaded = m.Swap(firstCase.Key, secondValue)
			if !loaded {
				t.Error("Expected loaded=true for swap on existing key")
			}
			if previous != newValue {
				t.Errorf("Expected previous value %v, got %v", newValue, previous)
			}
		}
	})

	t.Run("CompareAndSwap", func(t *testing.T) {
		m := factory()
		m.Store(firstCase.Key, firstCase.Value)

		// Successful compare and swap with correct old value
		if len(testData) > 1 {
			newValue := testData[1].Value
			swapped := m.CompareAndSwap(firstCase.Key, firstCase.Value, newValue)
			if !swapped {
				t.Error("Expected successful compare and swap")
			}

			if value, ok := m.Load(firstCase.Key); !ok || value != newValue {
				t.Errorf("Expected value %v after swap, got %v", newValue, value)
			}

			// Failed compare and swap with incorrect old value
			swapped = m.CompareAndSwap(firstCase.Key, firstCase.Value, firstCase.Value)
			if swapped {
				t.Error("Expected failed compare and swap with wrong old value")
			}
		}
	})

	t.Run("CompareAndDelete", func(t *testing.T) {
		m := factory()
		m.Store(firstCase.Key, firstCase.Value)

		// Failed compare and delete with incorrect value
		if len(testData) > 1 {
			wrongValue := testData[1].Value
			deleted := m.CompareAndDelete(firstCase.Key, wrongValue)
			if deleted {
				t.Error("Expected failed compare and delete with wrong value")
			}

			if _, ok := m.Load(firstCase.Key); !ok {
				t.Error("Expected key to still exist after failed compare and delete")
			}
		}

		// Successful compare and delete with correct value
		deleted := m.CompareAndDelete(firstCase.Key, firstCase.Value)
		if !deleted {
			t.Error("Expected successful compare and delete")
		}

		if _, ok := m.Load(firstCase.Key); ok {
			t.Error("Expected key to be deleted after successful compare and delete")
		}
	})
}

// testIterationOperations verifies Range, Keys, and Values methods work correctly.
func testIterationOperations[K comparable, V comparable](t *testing.T, factory MapFactory[K, V], testData []TestCase[K, V]) {
	t.Run("Range", func(t *testing.T) {
		m := factory()

		for _, tc := range testData {
			m.Store(tc.Key, tc.Value)
		}

		// Collect all key-value pairs through iteration
		var collected []TestCase[K, V]
		m.Range(func(key K, value V) bool {
			collected = append(collected, TestCase[K, V]{Key: key, Value: value})
			return true
		})

		if len(collected) != len(testData) {
			t.Errorf("Expected %d items from Range, got %d", len(testData), len(collected))
		}

		// Verify each expected pair exists in results (order may vary)
		for _, expected := range testData {
			found := false
			for _, actual := range collected {
				if actual.Key == expected.Key && actual.Value == expected.Value {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("Expected to find key-value pair {%v: %v} in Range results", expected.Key, expected.Value)
			}
		}
	})

	t.Run("RangeEarlyTermination", func(t *testing.T) {
		if len(testData) <= 1 {
			t.Skip("Need at least 2 test cases for early termination test")
			return
		}

		m := factory()
		for _, tc := range testData {
			m.Store(tc.Key, tc.Value)
		}

		// Terminate iteration after first item
		count := 0
		m.Range(func(key K, value V) bool {
			count++
			return false
		})

		if count != 1 {
			t.Errorf("Expected Range to stop after 1 iteration, but got %d iterations", count)
		}
	})

	t.Run("Keys", func(t *testing.T) {
		m := factory()

		for _, tc := range testData {
			m.Store(tc.Key, tc.Value)
		}

		var collectedKeys []K
		m.Keys(func(key K) bool {
			collectedKeys = append(collectedKeys, key)
			return true
		})

		if len(collectedKeys) != len(testData) {
			t.Errorf("Expected %d keys, got %d", len(testData), len(collectedKeys))
		}

		for _, expected := range testData {
			if !slices.Contains(collectedKeys, expected.Key) {
				t.Errorf("Expected to find key %v in Keys results", expected.Key)
			}
		}
	})

	t.Run("Values", func(t *testing.T) {
		m := factory()

		for _, tc := range testData {
			m.Store(tc.Key, tc.Value)
		}

		var collectedValues []V
		m.Values(func(value V) bool {
			collectedValues = append(collectedValues, value)
			return true
		})

		if len(collectedValues) != len(testData) {
			t.Errorf("Expected %d values, got %d", len(testData), len(collectedValues))
		}

		for _, expected := range testData {
			if !slices.Contains(collectedValues, expected.Value) {
				t.Errorf("Expected to find value %v in Values results", expected.Value)
			}
		}
	})
}

// testEdgeCases covers boundary conditions and error scenarios.
func testEdgeCases[K comparable, V comparable](t *testing.T, factory MapFactory[K, V], testData []TestCase[K, V]) {
	t.Run("EmptyMap", func(t *testing.T) {
		m := factory()

		if length := m.Len(); length != 0 {
			t.Errorf("Expected length 0 for empty map, got %d", length)
		}

		// Range over empty map should not invoke callback
		called := false
		m.Range(func(key K, value V) bool {
			called = true
			return true
		})
		if called {
			t.Error("Range function should not be called on empty map")
		}

		// Keys and Values should also not invoke callbacks on empty map
		m.Keys(func(key K) bool {
			called = true
			return true
		})
		if called {
			t.Error("Keys function should not be called on empty map")
		}

		m.Values(func(value V) bool {
			called = true
			return true
		})
		if called {
			t.Error("Values function should not be called on empty map")
		}
	})

	t.Run("OperationsAfterClear", func(t *testing.T) {
		if len(testData) == 0 {
			t.Skip("No test data for operations after clear")
			return
		}

		m := factory()

		// Store data, clear, then verify normal operations work
		for _, tc := range testData {
			m.Store(tc.Key, tc.Value)
		}
		m.Clear()

		firstCase := testData[0]

		m.Store(firstCase.Key, firstCase.Value)
		if value, ok := m.Load(firstCase.Key); !ok || value != firstCase.Value {
			t.Error("Store/Load should work after Clear")
		}

		if length := m.Len(); length != 1 {
			t.Errorf("Expected length 1 after storing one item post-clear, got %d", length)
		}
	})

	t.Run("OverwriteExistingKey", func(t *testing.T) {
		if len(testData) < 2 {
			t.Skip("Need at least 2 test cases for overwrite test")
			return
		}

		m := factory()

		key := testData[0].Key
		firstValue := testData[0].Value
		secondValue := testData[1].Value

		// Store initial value
		m.Store(key, firstValue)
		if length := m.Len(); length != 1 {
			t.Errorf("Expected length 1 after first store, got %d", length)
		}

		// Overwrite with new value
		m.Store(key, secondValue)
		if length := m.Len(); length != 1 {
			t.Errorf("Expected length still 1 after overwrite, got %d", length)
		}

		// Should retrieve the new value
		if value, ok := m.Load(key); !ok || value != secondValue {
			t.Errorf("Expected overwritten value %v, got %v", secondValue, value)
		}
	})
}

// Test functions demonstrating usage with specific implementations.
func TestUnorderedMapString(t *testing.T) {
	factory := func() maps.AbstractMap[string, string] {
		return maps.NewUnorderedMap[string, string]()
	}

	testData := []TestCase[string, string]{
		{"key1", "value1"},
		{"key2", "value2"},
		{"key3", "value3"},
	}

	testSuite(t, factory, testData)
}

func TestUnorderedMapInt(t *testing.T) {
	factory := func() maps.AbstractMap[int, int] {
		return maps.NewUnorderedMap[int, int]()
	}

	testData := []TestCase[int, int]{
		{1, 10},
		{2, 20},
		{3, 30},
	}

	testSuite(t, factory, testData)
}

func TestUnorderedMapMixed(t *testing.T) {
	factory := func() maps.AbstractMap[string, int] {
		return maps.NewUnorderedMap[string, int]()
	}

	testData := []TestCase[string, int]{
		{"one", 1},
		{"two", 2},
		{"three", 3},
	}

	testSuite(t, factory, testData)
}

// Performance benchmarks using the same factory pattern for consistency.
func BenchmarkUnorderedMapOperations(b *testing.B) {
	m := maps.NewUnorderedMap[string, string]()

	b.Run("Store", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			key := fmt.Sprintf("key%d", i)
			m.Store(key, "value")
		}
	})

	// Pre-populate for Load benchmark
	for i := 0; i < 1000; i++ {
		key := fmt.Sprintf("key%d", i)
		m.Store(key, "value")
	}

	b.Run("Load", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			key := fmt.Sprintf("key%d", i%1000)
			m.Load(key)
		}
	})
}
