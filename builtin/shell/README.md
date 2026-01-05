# Shell

Execute shell commands in a `sh` shell environment.

## Usage

```yaml
- id: run_cmd
  uses: builtin/shell
  with:
    command: "date"
    stdout: "stdout"   # optional: "stdout", "stderr", "discard", or file path
    stderr: "stderr"   # optional: "stdout", "stderr", "discard", or file path
```

**Parameters:**
- `command` (required): Shell command(s) to execute
- `stdout` (optional): Output destination - "stdout", "discard", or file path (default: captured)
- `stderr` (optional): Error destination - "stderr", "stdout", "discard", or file path (default: captured)

**Output:**
```json
{
  "output": "command output...",
  "exit_code": 0
}
```

**Supported substitutions:**
- `${{ args.0 }}` - CLI argument by index
- `${{ args }}` - All CLI arguments
- `${{ env.PATH }}` - Environment variables
- `${{ outputs.step_id }}` - Previous step outputs