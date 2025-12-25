# Shell

Execute shell commands in a `sh` shell environment.

## Parameters

- **command** (required): Shell command(s) to execute
- **stdout** (optional): `"stdout"`, `"discard"`, or file path
- **stderr** (optional): `"stderr"`, `"stdout"`, `"discard"`, or file path

## Usage

```yml
- name: Run command
  using: builtin/shell
  with:
    command: date
```

## Output Redirection

- Not specified: Captured for return value
- `"stdout"` / `"stderr"`: Write to main process
- `"discard"`: Suppress output
- File path: Append to file

```yml
- name: Show in terminal
  using: builtin/shell
  with:
    command: echo "Hello"
    stdout: stdout
    stderr: stderr
```

## Substitutions

Supports variable substitution:
- `${{ args.0 }}` - CLI argument by index
- `${{ args }}` - All CLI arguments
- `${{ env.PATH }}` - Environment variables
- `${{ outputs.step_id }}` - Previous step outputs

## Return Value

Returns captured output (stdout and stderr combined) as a string. Only output not redirected elsewhere is captured. Accessible via step `id`:

```yml
- id: date
  using: builtin/shell
  with:
    command: date

- name: Use output
  using: builtin/echo
  with:
    message: Date is ${{ outputs.date }}
```

## Notes

- Executed in `sh` shell
- Working directory is current directory
- Inherits parent process environment
- Non-zero exit codes fail the step

## Examples

See [example.yml](example.yml) for comprehensive usage examples.