package cache_test

import (
	"testing"
	"time"

	"github.com/jcnnll/cache"
)

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
