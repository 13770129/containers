package maps_test

import (
	"fmt"
	"testing"

	"github.com/13770129/containers/maps"
)

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
