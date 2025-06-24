package maps_test

import (
	"fmt"
	"testing"

	"github.com/13770129/containers/maps"
)

// Test OrderedMap with the general test suite to ensure it implements AbstractMap correctly
func TestOrderedMapString(t *testing.T) {
	factory := func() maps.AbstractMap[string, string] {
		return maps.NewOrderedMap[string, string]()
	}

	testData := []TestCase[string, string]{
		{"alpha", "first"},
		{"beta", "second"},
		{"gamma", "third"},
		{"delta", "fourth"},
	}

	testSuite(t, factory, testData)
}

func TestOrderedMapInt(t *testing.T) {
	factory := func() maps.AbstractMap[int, int] {
		return maps.NewOrderedMap[int, int]()
	}

	testData := []TestCase[int, int]{
		{10, 100},
		{20, 200},
		{30, 300},
		{40, 400},
	}

	testSuite(t, factory, testData)
}

func TestOrderedMapMixed(t *testing.T) {
	factory := func() maps.AbstractMap[string, int] {
		return maps.NewOrderedMap[string, int]()
	}

	testData := []TestCase[string, int]{
		{"first", 1},
		{"second", 2},
		{"third", 3},
		{"fourth", 4},
	}

	testSuite(t, factory, testData)
}

// OrderedMap-specific tests that verify insertion order preservation
func TestOrderedMapInsertionOrder(t *testing.T) {
	t.Run("BasicInsertionOrder", func(t *testing.T) {
		om := maps.NewOrderedMap[string, int]()

		// Insert keys in a specific order
		insertOrder := []string{"charlie", "alpha", "delta", "beta"}
		for i, key := range insertOrder {
			om.Store(key, i)
		}

		// Verify iteration follows insertion order exactly
		var iterationOrder []string
		om.Range(func(key string, value int) bool {
			iterationOrder = append(iterationOrder, key)
			return true
		})

		if len(iterationOrder) != len(insertOrder) {
			t.Fatalf("Expected %d keys, got %d", len(insertOrder), len(iterationOrder))
		}

		for i, expectedKey := range insertOrder {
			if iterationOrder[i] != expectedKey {
				t.Errorf("Position %d: expected %s, got %s", i, expectedKey, iterationOrder[i])
			}
		}
	})

	t.Run("UpdatePreservesOrder", func(t *testing.T) {
		om := maps.NewOrderedMap[string, int]()

		// Insert initial values
		keys := []string{"first", "second", "third"}
		for i, key := range keys {
			om.Store(key, i)
		}

		// Update middle value - should preserve position
		om.Store("second", 999)

		// Verify order unchanged, but value updated
		expectedOrder := []string{"first", "second", "third"}
		expectedValues := []int{0, 999, 2}

		i := 0
		om.Range(func(key string, value int) bool {
			if key != expectedOrder[i] {
				t.Errorf("Position %d: expected key %s, got %s", i, expectedOrder[i], key)
			}
			if value != expectedValues[i] {
				t.Errorf("Position %d: expected value %d, got %d", i, expectedValues[i], value)
			}
			i++
			return true
		})
	})

	t.Run("DeletionMaintainsOrder", func(t *testing.T) {
		om := maps.NewOrderedMap[int, string]()

		// Insert sequence: 10, 20, 30, 40, 50
		for i := 1; i <= 5; i++ {
			om.Store(i*10, fmt.Sprintf("value%d", i*10))
		}

		// Delete middle elements (20, 40)
		om.Delete(20)
		om.Delete(40)

		// Verify remaining elements maintain relative order
		expectedKeys := []int{10, 30, 50}
		var actualKeys []int
		om.Range(func(key int, value string) bool {
			actualKeys = append(actualKeys, key)
			return true
		})

		if len(actualKeys) != len(expectedKeys) {
			t.Fatalf("Expected %d keys after deletion, got %d", len(expectedKeys), len(actualKeys))
		}

		for i, expected := range expectedKeys {
			if actualKeys[i] != expected {
				t.Errorf("Position %d: expected key %d, got %d", i, expected, actualKeys[i])
			}
		}
	})

	t.Run("ClearAndReinsertNewOrder", func(t *testing.T) {
		om := maps.NewOrderedMap[string, int]()

		// Initial insertion order
		firstOrder := []string{"zebra", "apple", "mango"}
		for i, key := range firstOrder {
			om.Store(key, i)
		}

		// Clear and insert in different order
		om.Clear()
		secondOrder := []string{"apple", "zebra", "mango"}
		for i, key := range secondOrder {
			om.Store(key, i+10)
		}

		// Verify new insertion order is preserved
		var actualOrder []string
		om.Range(func(key string, value int) bool {
			actualOrder = append(actualOrder, key)
			return true
		})

		for i, expected := range secondOrder {
			if actualOrder[i] != expected {
				t.Errorf("Position %d: expected %s, got %s", i, expected, actualOrder[i])
			}
		}
	})
}

// Test Range iteration behavior and early termination
func TestOrderedMapRangeIteration(t *testing.T) {
	t.Run("EarlyTermination", func(t *testing.T) {
		om := maps.NewOrderedMap[int, string]()

		// Insert multiple items
		for i := 1; i <= 10; i++ {
			om.Store(i, fmt.Sprintf("item%d", i))
		}

		// Range but stop after 3 items
		var visitedKeys []int
		om.Range(func(key int, value string) bool {
			visitedKeys = append(visitedKeys, key)
			return len(visitedKeys) < 3 // Stop after 3 items
		})

		if len(visitedKeys) != 3 {
			t.Errorf("Expected exactly 3 visited keys, got %d", len(visitedKeys))
		}

		// Verify we got the first 3 keys in insertion order
		for i, key := range visitedKeys {
			expectedKey := i + 1
			if key != expectedKey {
				t.Errorf("Position %d: expected key %d, got %d", i, expectedKey, key)
			}
		}
	})

	t.Run("EmptyMapRange", func(t *testing.T) {
		om := maps.NewOrderedMap[string, int]()

		rangeCallCount := 0
		om.Range(func(key string, value int) bool {
			rangeCallCount++
			return true
		})

		if rangeCallCount != 0 {
			t.Errorf("Range should not call function on empty map, but was called %d times", rangeCallCount)
		}
	})
}

// Test Keys and Values iteration methods
func TestOrderedMapKeysAndValues(t *testing.T) {
	t.Run("KeysIterationOrder", func(t *testing.T) {
		om := maps.NewOrderedMap[string, int]()

		insertOrder := []string{"delta", "alpha", "charlie", "beta"}
		for i, key := range insertOrder {
			om.Store(key, i*10)
		}

		var keysOrder []string
		om.Keys(func(key string) bool {
			keysOrder = append(keysOrder, key)
			return true
		})

		// Verify Keys() returns keys in insertion order
		for i, expected := range insertOrder {
			if keysOrder[i] != expected {
				t.Errorf("Keys position %d: expected %s, got %s", i, expected, keysOrder[i])
			}
		}
	})

	t.Run("ValuesIterationOrder", func(t *testing.T) {
		om := maps.NewOrderedMap[int, string]()

		testPairs := []struct {
			key   int
			value string
		}{
			{30, "thirty"},
			{10, "ten"},
			{20, "twenty"},
		}

		for _, pair := range testPairs {
			om.Store(pair.key, pair.value)
		}

		var valuesOrder []string
		om.Values(func(value string) bool {
			valuesOrder = append(valuesOrder, value)
			return true
		})

		// Verify Values() returns values in insertion order
		expectedValues := []string{"thirty", "ten", "twenty"}
		for i, expected := range expectedValues {
			if valuesOrder[i] != expected {
				t.Errorf("Values position %d: expected %s, got %s", i, expected, valuesOrder[i])
			}
		}
	})
}

// Performance benchmarks comparing OrderedMap to UnorderedMap
func BenchmarkOrderedMapVsUnordered(b *testing.B) {
	b.Run("OrderedMapStore", func(b *testing.B) {
		om := maps.NewOrderedMap[string, int]()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			key := fmt.Sprintf("key%d", i)
			om.Store(key, i)
		}
	})

	b.Run("UnorderedMapStore", func(b *testing.B) {
		um := maps.NewUnorderedMap[string, int]()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			key := fmt.Sprintf("key%d", i)
			um.Store(key, i)
		}
	})

	// Benchmark Load operations on pre-populated maps
	b.Run("OrderedMapLoad", func(b *testing.B) {
		om := maps.NewOrderedMap[string, int]()
		// Pre-populate with test data
		for i := 0; i < 10000; i++ {
			key := fmt.Sprintf("key%d", i)
			om.Store(key, i)
		}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			key := fmt.Sprintf("key%d", i%10000)
			om.Load(key)
		}
	})

	b.Run("UnorderedMapLoad", func(b *testing.B) {
		um := maps.NewUnorderedMap[string, int]()
		// Pre-populate with test data
		for i := 0; i < 10000; i++ {
			key := fmt.Sprintf("key%d", i)
			um.Store(key, i)
		}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			key := fmt.Sprintf("key%d", i%10000)
			um.Load(key)
		}
	})

	// Benchmark Range iteration performance
	b.Run("OrderedMapRange", func(b *testing.B) {
		om := maps.NewOrderedMap[int, string]()
		for i := 0; i < 1000; i++ {
			om.Store(i, fmt.Sprintf("value%d", i))
		}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			om.Range(func(key int, value string) bool {
				return true // Process all elements
			})
		}
	})

	b.Run("UnorderedMapRange", func(b *testing.B) {
		um := maps.NewUnorderedMap[int, string]()
		for i := 0; i < 1000; i++ {
			um.Store(i, fmt.Sprintf("value%d", i))
		}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			um.Range(func(key int, value string) bool {
				return true // Process all elements
			})
		}
	})
}

// Test memory efficiency and large dataset handling
func TestOrderedMapLargeDataset(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping large dataset test in short mode")
	}

	t.Run("LargeInsertionAndIteration", func(t *testing.T) {
		om := maps.NewOrderedMap[int, string]()

		// Insert large number of items
		itemCount := 100000
		for i := 0; i < itemCount; i++ {
			om.Store(i, fmt.Sprintf("value_%d", i))
		}

		if om.Len() != itemCount {
			t.Fatalf("Expected length %d, got %d", itemCount, om.Len())
		}

		// Verify iteration order for sample of items
		sampleIndices := []int{0, 1000, 50000, 99999}
		iterationIndex := 0
		nextSampleIndex := 0

		om.Range(func(key int, value string) bool {
			// Check if this is one of our sample indices
			if nextSampleIndex < len(sampleIndices) && iterationIndex == sampleIndices[nextSampleIndex] {
				expectedKey := sampleIndices[nextSampleIndex]
				if key != expectedKey {
					t.Errorf("Sample position %d: expected key %d, got %d", nextSampleIndex, expectedKey, key)
				}
				nextSampleIndex++
			}
			iterationIndex++
			return true
		})

		if iterationIndex != itemCount {
			t.Errorf("Expected to iterate over %d items, but iterated over %d", itemCount, iterationIndex)
		}
	})
}

// Edge case testing for OrderedMap specific scenarios
func TestOrderedMapEdgeCases(t *testing.T) {
	t.Run("StoreNilValue", func(t *testing.T) {
		om := maps.NewOrderedMap[string, *string]()

		om.Store("key", nil)

		value, ok := om.Load("key")
		if !ok {
			t.Error("Expected to find key with nil value")
		}
		if value != nil {
			t.Errorf("Expected nil value, got %v", value)
		}
	})

	t.Run("ZeroValueTypes", func(t *testing.T) {
		om := maps.NewOrderedMap[int, int]()

		// Store zero values
		om.Store(0, 0)

		value, ok := om.Load(0)
		if !ok {
			t.Error("Expected to find zero key")
		}
		if value != 0 {
			t.Errorf("Expected zero value, got %d", value)
		}
	})

	t.Run("StringKeys", func(t *testing.T) {
		om := maps.NewOrderedMap[string, string]()

		testCases := []struct {
			key   string
			value string
		}{
			{"", "empty_key"},          // Empty string key
			{" ", "space_key"},         // Space key
			{"key with spaces", "val"}, // Key with spaces
			{"emoji_key_🔑", "emoji"},   // Unicode key
		}

		for _, tc := range testCases {
			om.Store(tc.key, tc.value)
		}

		// Verify all keys can be retrieved and order is preserved
		i := 0
		om.Range(func(key string, value string) bool {
			expected := testCases[i]
			if key != expected.key || value != expected.value {
				t.Errorf("Position %d: expected {%q: %q}, got {%q: %q}",
					i, expected.key, expected.value, key, value)
			}
			i++
			return true
		})
	})
}
