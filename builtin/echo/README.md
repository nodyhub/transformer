# Echo

Echo messages to standard output.

## Usage

```yaml
- id: print_message
  uses: builtin/echo
  with:
    message: "Hello World"
```

**Parameters:**
- `message` (required): String or list of strings to print

**Output:**
```json
{
  "output": "Hello World"
}
```

**Supported substitutions:**
- `${{ args.0 }}` - CLI argument by index
- `${{ args }}` - All CLI arguments
- `${{ env.HOME }}` - Environment variables
- `${{ outputs.step_id.output }}` - Previous step output
