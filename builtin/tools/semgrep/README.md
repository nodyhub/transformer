# Semgrep

Semgrep is a fast, open-source static analysis tool for finding bugs and enforcing code standards.

## Usage

```yaml
- name: Security scan with Semgrep
  uses: builtin/tools/semgrep
  with:
    config: p/security-audit
    format: json
    output: semgrep-results.json
```

## Parameters

### Optional
- `target` - Path to scan (default: `.`)
- `config` - Rules config: `auto`, `p/ci`, `p/security-audit`, or custom path
- `format` - Output format: `json` (default), `sarif`, `text`, `gitlab-sast`, `junit-xml`
- `output` - Output file path
- `severity` - Filter by severity: `INFO`, `WARNING`, `ERROR`
- `exclude` - Exclude patterns (string or array)
- `max-memory` - Maximum memory in MB (default: 5000)
- `metrics` - Send anonymous metrics: `on`, `off`
- `additional-args` - Array of additional arguments

## Examples

See [examples.yml](examples.yml) for usage examples including:
- Security audit
- CI/CD scan
- Custom rules with exclusions
- Multiple rulesets

## Output

Returns a map containing:
- `command` - Full command executed
- `output` - Combined stdout/stderr
- `exit_code` - Exit code (0 = success)
- `error` - Error message if failed
- `output_file` - Path to output file if specified

## Links

- [Semgrep Documentation](https://semgrep.dev/docs/)
- [Semgrep Registry](https://semgrep.dev/explore)
- [Semgrep GitHub](https://github.com/semgrep/semgrep)
