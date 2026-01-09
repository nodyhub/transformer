# Contributing to Transformer

Thank you for your interest in contributing to Transformer! We welcome contributions in the form of bug reports, feature requests, and pull requests.

## Implementing Custom Step Functions

One of the best ways to extend Transformer is by implementing custom step functions. Here's how to create and share your own steps.

### Step Function Signature

All step functions must implement this interface:

```go
type StepFunc func(ctx context.Context, with interface{}) (interface{}, error)
```

- **ctx**: Context for cancellation and timeouts
- **with**: Input parameters from your YAML workflow
- **Return**: A map or struct containing output data, or an error

### Example: Custom Step Function

Create a new package in your module:

```go
package greeting

import (
	"context"
	"fmt"
	"github.com/nodyhub/transformer/registry"
)

func init() {
	registry.Register("custom/greet", Greet)
}

// Greet is a custom step function
func Greet(ctx context.Context, with interface{}) (interface{}, error) {
	withMap, ok := with.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("expected 'with' parameter to be a map")
	}
	
	// Extract parameters
	name, ok := withMap["name"].(string)
	if !ok {
		return nil, fmt.Errorf("'name' field is required")
	}
	
	greeting := fmt.Sprintf("Hello, %s!", name)
	
	return map[string]interface{}{
		"message": greeting,
	}, nil
}
```

### Using Your Custom Step in YAML

```yaml
steps:
  - id: greet_step
    uses: custom/greet
    with:
      name: World
  
  - uses: builtin/echo
    with:
      message: ${{ outputs.greet_step.message }}
```

### Optional: Module Hooks for Dependencies

If your step function requires system dependencies, register module hooks:

```go
func init() {
	registry.Register("custom/myfunction", MyFunction, moduleHook)
}

// moduleHook installs required system dependencies
func moduleHook(ctx context.Context) error {
	// Install dependencies via shell commands
	// This is called when: transformer modules install custom/myfunction
	return nil
}
```

Users can then install your module with:

```bash
transformer modules install custom/myfunction
```

### Guidelines for Custom Step Functions

1. **Input Validation**: Always validate and type-check the `with` parameters
2. **Error Handling**: Return meaningful error messages
3. **Output Structure**: Use maps with string keys for consistent output
4. **Context Awareness**: Respect the context for cancellation
5. **Documentation**: Include clear examples in your module's README
6. **Testing**: Write unit tests following the pattern in [builtin/](./builtin/)

### Testing Your Step Function

Create tests following this pattern:

```go
package greeting

import (
	"context"
	"testing"
)

func TestGreet(t *testing.T) {
	result, err := Greet(context.Background(), map[string]interface{}{
		"name": "Alice",
	})
	
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	
	output := result.(map[string]interface{})
	if output["message"] != "Hello, Alice!" {
		t.Errorf("expected 'Hello, Alice!', got %v", output["message"])
	}
}
```

## Built-in Examples

For more complex examples, explore the built-in step functions:

- [builtin/echo](./builtin/echo/) - Simple output function
- [builtin/shell](./builtin/shell/) - Execute shell commands with I/O redirection
- [builtin/file/read](./builtin/file/read/) - Read file contents
- [builtin/file/write](./builtin/file/write/) - Write content to files

## Submitting Your Contribution

1. **Fork** the repository
2. **Create a branch** for your feature (`git checkout -b feature/my-step`)
3. **Implement** your step function with tests
4. **Document** with examples and a README
5. **Test** your implementation thoroughly
6. **Submit** a pull request with a clear description

## Questions?

Feel free to open an issue to discuss your ideas before implementing!
