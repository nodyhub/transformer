// Package registry provides a registry for step implementations.
package registry

import (
	"context"
	"fmt"
	"log/slog"
	"slices"
	"strings"
	"sync"
)

// StepFunc is the function signature for step implementations
type StepFunc func(ctx context.Context, with interface{}) (interface{}, error)

// ModuleHook is the function signature for module hooks if the step requires any
// modules to be set up before execution.
type ModuleHook func(ctx context.Context) error

var (
	stepFunctions = make(map[string]StepFunc)
	modulesHooks  = make(map[string][]ModuleHook)
	mu            sync.RWMutex
)

// Register registers a step function with the given name
func Register(name string, fn StepFunc, moduleHooks ...ModuleHook) {
	mu.Lock()
	defer mu.Unlock()

	if fn == nil {
		panic("registry: Register function is nil")
	}
	if _, exists := stepFunctions[name]; exists {
		panic(fmt.Sprintf("registry: Register called twice for step %q", name))
	}
	stepFunctions[name] = fn
	modulesHooks[name] = moduleHooks
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

func RunModuleHooks(ctx context.Context, name string) error {
	mu.RLock()
	defer mu.RUnlock()

	found := false
	for hookPath, hooks := range modulesHooks {
		if strings.HasPrefix(hookPath, name) {
			found = true

			if len(hooks) == 0 {
				continue
			}

			slog.Debug("running module hooks", slog.String("path", hookPath), slog.Int("count", len(hooks)))
			for _, hook := range hooks {
				if err := hook(ctx); err != nil {
					return fmt.Errorf("error running hook for %s: %w", hookPath, err)
				}
			}

		}
	}

	if !found {
		return fmt.Errorf("no hooks found for path: %s", name)
	}

	return nil
}

func ListModules() []string {
	mu.RLock()
	defer mu.RUnlock()

	names := make([]string, 0, len(modulesHooks))
	for name := range modulesHooks {
		names = append(names, name)
	}
	slices.Sort(names)

	return names
}
