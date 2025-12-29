# Trufflehog

Trufflehog finds and verifies secrets in code, commit history, and filesystems.

## Usage

```yaml
- name: Scan for secrets
  uses: builtin/tools/trufflehog
  with:
    type: git
    target: https://github.com/user/repo.git
    json: true
```

## Parameters

### Required
- `target` - Git repository URL, filesystem path, or docker image

### Optional
- `type` - Scan type: `git` (default), `github`, `gitlab`, `filesystem`, `docker`, `s3`
- `format` - Output format: `json`, `yaml`
- `output` - Output file path
- `since-commit` - Scan commits since this commit
- `branch` - Scan specific branch
- `max-depth` - Maximum commit depth (default: 1000000)
- `include-paths` - File with path patterns to include
- `exclude-paths` - File with path patterns to exclude
- `only-verified` - Boolean to show only verified secrets
- `json` - Boolean for JSON output
- `additional-args` - Array of additional arguments

## Examples

See [examples.yml](examples.yml) for usage examples including:
- Scan Git repository
- Scan recent commits
- Scan specific branch
- Scan filesystem
- Scan Docker image

## Output

Returns a map containing:
- `command` - Full command executed
- `output` - Combined stdout/stderr
- `exit_code` - Exit code (0 = success)
- `error` - Error message if failed
- `output_file` - Path to output file if specified

## Links

- [Trufflehog Documentation](https://github.com/trufflesecurity/trufflehog#readme)
- [Trufflehog GitHub](https://github.com/trufflesecurity/trufflehog)
