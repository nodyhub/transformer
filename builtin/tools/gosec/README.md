# Gosec

Gosec inspects Go source code for security problems by scanning the Go AST.

## Usage

```yaml
- name: Gosec security scan
  uses: builtin/tools/gosec
  with:
    target: ./...
    format: sarif
    output: gosec.sarif
```

## Parameters

### Optional
- `target` - Path to scan (default: `./...`)
- `format` - Output format: `json` (default), `yaml`, `csv`, `junit-xml`, `html`, `sonarqube`, `golint`, `sarif`, `text`
- `output` - Output file path
- `severity` - Filter by severity: `low`, `medium`, `high`
- `confidence` - Filter by confidence: `low`, `medium`, `high`
- `exclude` - Exclude directories (comma-separated)
- `exclude-generated` - Boolean to exclude generated files
- `nosec` - Boolean to ignore #nosec comments
- `tests` - Boolean to scan test files
- `additional-args` - Array of additional arguments

## Examples

See [examples.yml](examples.yml) for usage examples including:
- Basic Go scan
- Scan including test files
- Exclude directories
- CI/CD integration

## Output

Returns a map containing:
- `command` - Full command executed
- `output` - Combined stdout/stderr
- `exit_code` - Exit code (0 = success)
- `error` - Error message if failed
- `output_file` - Path to output file if specified

## Links

- [Gosec Documentation](https://github.com/securego/gosec#readme)
- [Gosec Rules](https://github.com/securego/gosec#available-rules)
- [Gosec GitHub](https://github.com/securego/gosec)
