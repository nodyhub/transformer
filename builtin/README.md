# Built-in Functions

Transformer includes several built-in functions for common operations. Each function can be used as a step in your workflow by referencing it with the `uses` field.

## Echo

[`builtin/echo`](./echo/) - Print messages to standard output.

Print strings, lists, or interpolated messages. Useful for logging, displaying results, and debugging.

```yaml
- uses: builtin/echo
  with:
    message: "Hello, World!"
```

**Output:** `{ "output": "..." }`

## Shell

[`builtin/shell`](./shell/) - Execute shell commands.

Run arbitrary shell commands with support for output/error redirection. Execute scripts, run tools, or perform system operations.

```yaml
- uses: builtin/shell
  with:
    command: "date"
```

**Output:** `{ "output": "...", "exit_code": 0 }`

## File Operations

### Read

[`builtin/file/read`](./file/read/) - Read files from disk.

Load file contents for processing. Returns the file path, content, and byte count.

```yaml
- uses: builtin/file/read
  with:
    path: "data.txt"
```

**Output:** `{ "path": "...", "content": "...", "bytes": 1234 }`

### Write

[`builtin/file/write`](./file/write/) - Write files to disk.

Save content to files with optional append mode and permission control. Returns the file path and bytes written.

```yaml
- uses: builtin/file/write
  with:
    path: "output.txt"
    content: "data"
```

**Output:** `{ "path": "...", "bytes": 1234 }`

## Variable Substitution

All built-in functions support variable substitution in their parameters:

- **Previous outputs:** `${{ outputs.step_id.property }}`
- **CLI arguments:** `${{ args.0 }}` or `${{ args }}`
- **Environment variables:** `${{ env.HOME }}`

## Adding Custom Functions

To add custom functions, implement the function signature:

```go
func MyFunction(ctx context.Context, with interface{}) (interface{}, error) {
  // Implementation
}
```

Register it in your `init()`:

```go
func init() {
  registry.Register("custom/myfunction", MyFunction)
}
```
