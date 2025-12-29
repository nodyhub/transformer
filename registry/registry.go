// Package registry provides a registry for step implementations.
package registry

import (
	"context"
	"fmt"
	"sync"
)

// StepFunc is the function signature for step implementations
type StepFunc func(ctx context.Context, with interface{}) (interface{}, error)

var (
	stepFunctions = make(map[string]StepFunc)
	mu            sync.RWMutex
)

// Register registers a step function with the given name
func Register(name string, fn StepFunc) {
	mu.Lock()
	defer mu.Unlock()

	if fn == nil {
		panic("registry: Register function is nil")
	}
	if _, exists := stepFunctions[name]; exists {
		panic(fmt.Sprintf("registry: Register called twice for step %q", name))
	}
	stepFunctions[name] = fn
}

// Get retrieves a registered step function by name
func Get(name string) (StepFunc, bool) {
	mu.RLock()
	defer mu.RUnlock()
	fn, ok := stepFunctions[name]
	return fn, ok
}

func Call(ctx context.Context, name string, with interface{}) (interface{}, error) {
	fn, ok := Get(name)
	if !ok {
		return nil, fmt.Errorf("unknown step function: %s", name)
	}
	return fn(ctx, with)
}

// List returns all registered step names
func List() []string {
	mu.RLock()
	defer mu.RUnlock()

	names := make([]string, 0, len(stepFunctions))
	for name := range stepFunctions {
		names = append(names, name)
	}
	return names
}

func Remove(name string) {
	mu.Lock()
	defer mu.Unlock()
	delete(stepFunctions, name)
}

func RemoveAll() {
	mu.Lock()
	defer mu.Unlock()
	stepFunctions = make(map[string]StepFunc)
}
