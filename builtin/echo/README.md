# Echo

Echo messages to standard output.

## Parameters

- **message** (required): String or list of strings to print

## Usage

```yml
- name: Print message
  using: builtin/echo
  with:
    message: Hello World
```

## Substitutions

Supports variable substitution:
- `${{ args.0 }}` - CLI argument by index
- `${{ args }}` - All CLI arguments
- `${{ env.HOME }}` - Environment variables
- `${{ outputs.step_id }}` - Previous step outputs

## Return Value

Returns printed content as a string (without trailing newline). Accessible via step `id`:

```yml
- id: greeting
  using: builtin/echo
  with:
    message: Hello

- name: Use output
  using: builtin/echo
  with:
    message: Previous output was ${{ outputs.greeting }}
```

## Examples

See [example.yml](example.yml) for comprehensive usage examples.