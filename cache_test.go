package cache_test

import (
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/jcnnll/cache"
)

func TestDelete(t *testing.T) {
	c := cache.New[string](10)
	c.Set("key", "value", 10*time.Second)

	val, ok := c.Get("key")
	if !ok || val != "value" {
		t.Fatal("expected key to exists before deletion")
	}

	c.Delete("key")

	_, ok = c.Get("key")
	if ok {
		t.Fatal("expected key to be deleted")
	}
}

func TestTTLExtension(t *testing.T) {
	c := cache.New[string](10)

	// Set with short TTL
	c.Set("key", "value", 1*time.Second)
	time.Sleep(500 * time.Millisecond)
	c.Set("key", "value", 10*time.Second)

	// Wait for originall TTL to pass
	time.Sleep(1 * time.Second)

	// Key should still be present
	_, ok := c.Get("key")
	if !ok {
		t.Fatal("expected key to still exist after TTL extension")
	}
}

func TestConcurrency(t *testing.T) {
	const goroutines = 100
	const operations = 1000

	c := cache.New[int](1000)
	var wg sync.WaitGroup
	wg.Add(goroutines)

	for g := range goroutines {
		go func(id int) {
			defer wg.Done()

			// Each goroutine does a mix of operations
			for i := range operations {
				key := strconv.Itoa(id*operations + i)

				switch i % 3 {
				case 0:
					c.Set(key, id, 5*time.Second)
				case 1:
					c.Get(key)
				case 2:
					c.Delete(key)
				}
			}
		}(g)
	}

	wg.Wait()
	// if we get here without any deadlocks or panics, we pass
}

func TestCApacityBoundary(t *testing.T) {
	c := cache.New[string](1)

	c.Set("first", "value1", 10*time.Second)
	c.Set("second", "value2", 10*time.Second)

	// Fiest key should be evicted
	if _, ok := c.Get("first"); ok {
		t.Fatal("expected 'first' to be evicted")
	}

	if val, ok := c.Get("second"); !ok || val != "value2" {
		t.Fatal("expected 'second' to exist")
	}
}

func TestGenericTypes(t *testing.T) {
	// Test with integers
	intCache := cache.New[int](10)
	intCache.Set("answer", 42, 5*time.Second)
	if val, ok := intCache.Get("answer"); !ok || val != 42 {
		t.Fatalf("expected integer value 42, got %v", val)
	}

	// Test with booleans
	boolCache := cache.New[bool](10)
	boolCache.Set("flag", true, 5*time.Second)
	if val, ok := boolCache.Get("flag"); !ok || !val {
		t.Fatal("expected boolean value true")
	}

	// Test with structs
	type Person struct {
		Name string
		Age  int
	}
	structCache := cache.New[Person](10)
	alice := Person{Name: "Alice", Age: 30}
	structCache.Set("alice", alice, 5*time.Second)
	if val, ok := structCache.Get("alice"); !ok || val.Name != "Alice" || val.Age != 30 {
		t.Fatal("expected struct to match original")
	}
}

func TestZeroTTL(t *testing.T) {
	c := cache.New[string](10)
	c.Set("forever", "value", 0) // Zero duration

	time.Sleep(100 * time.Millisecond) // Brief wait

	if val, ok := c.Get("forever"); !ok || val != "value" {
		t.Fatal("expected item with zero TTL to remain cached")
	}
}

func TestLRUOrdering(t *testing.T) {
	c := cache.New[string](3)

	c.Set("a", "A", 10*time.Second)
	c.Set("b", "B", 10*time.Second)
	c.Set("c", "C", 10*time.Second)

	// Access in reverse to change LRU status
	c.Get("c")
	c.Get("b")
	c.Get("a")

	// Add another item to evict c
	c.Set("d", "D", 10*time.Second)

	if _, ok := c.Get("c"); ok {
		t.Fatal("expected 'c' to be evicted")
	}

	// a,b,d sould still exist
	if _, ok := c.Get("a"); !ok {
		t.Fatal("expected 'a' to exist")
	}
	if _, ok := c.Get("b"); !ok {
		t.Fatal("expected 'b' to exist")
	}
	if _, ok := c.Get("d"); !ok {
		t.Fatal("expected 'd' to exist")
	}
}

func TestUpdate(t *testing.T) {
	c := cache.New[string](10)
	c.Set("key", "original", 10*time.Second)
	c.Set("key", "updated", 10*time.Second)

	val, ok := c.Get("key")
	if !ok || val != "updated" {
		t.Fatalf("expected 'udated', got '%v'", val)
	}
}

func TestSetGet(t *testing.T) {
	c := cache.New[string](1)

	c.Set("a", "A", 5*time.Second)
	val, ok := c.Get("a")
	if !ok || val != "A" {
		t.Fatal("expected to retrieve 'A'")
	}
}

func TestTTLExpiry(t *testing.T) {
	c := cache.New[string](1)
	c.Set("expire", "soon", 1*time.Second)

	time.Sleep(2 * time.Second)
	_, ok := c.Get("expire")
	if ok {
		t.Fatal("expected key to expire")
	}
}

func TestLRUEviction(t *testing.T) {
	c := cache.New[string](2)

	c.Set("a", "A", 10*time.Second)
	c.Set("b", "B", 10*time.Second)
	c.Get("a")                      // Access 'a' to mark as recent
	c.Set("c", "C", 10*time.Second) // Evict 'b'

	if _, ok := c.Get("b"); ok {
		t.Fatal("expected 'b' to be evicted")
	}
	if _, ok := c.Get("a"); !ok {
		t.Fatal("expected 'a' to remain")
	}
	if _, ok := c.Get("c"); !ok {
		t.Fatal("expected 'c' to remain")
	}
}
