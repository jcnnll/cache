# cache
A generic, thread-safe in-memory LRU (Least Recently Used) cache with per-key TTL (Time-To-Live) support, implemented in Go.

This package allows you to:
- Store and retrieve typed values (`Cache[T]`) safely using Go generics
- Automatically evict items when a custom TTL expires
- Limit the number of entries using an LRU eviction strategy
- Handle concurrent access with internal locking

## Features

- ✅ Go 1.18+ generics
- ✅ Per-key expiration (TTL)
- ✅ LRU eviction policy (evicts least recently accessed items first)
- ✅ Safe for concurrent use
- ✅ Clean, dependency-free implementation

## Installation

```bash
go get github.com/jcnnll/cache
```

## Usage Example

```go
package main

import (
	"fmt"
	"time"

	"github.com/jcnnll/cache"
)

func main() {
	c := cache.New  // max 2 items

	c.Set("a", "Alpha", 5*time.Second)
	c.Set("b", "Beta", 5*time.Second)

	val, ok := c.Get("a")
	if ok {
		fmt.Println("a:", val)
	}

	c.Set("c", "Gamma", 5*time.Second) // Triggers LRU eviction

	if _, ok := c.Get("b"); !ok {
		fmt.Println("b was evicted (LRU)")
	}
}

```

## License

This project is licensed under [The Unlicense](https://unlicense.org/), a public domain dedication.

You are free to use, modify, and distribute this software without restriction.
