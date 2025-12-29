# Trivy

Trivy is a comprehensive and versatile vulnerability scanner for containers and other artifacts.

## Usage

```yaml
- name: Scan Docker image
  uses: builtin/tools/trivy
  with:
    type: image
    target: myapp:latest
    format: json
    output: trivy-report.json
```

## Parameters

### Required
- `target` - Image name, filesystem path, or repository to scan

### Optional
- `type` - Scan type: `image` (default), `fs`, `repo`, `config`, `sbom`
- `format` - Output format: `json` (default), `sarif`, `table`, `cyclonedx`, `spdx`
- `output` - Output file path
- `severity` - Comma-separated severities: `UNKNOWN`, `LOW`, `MEDIUM`, `HIGH`, `CRITICAL`
- `scanners` - Comma-separated scanners: `vuln`, `misconfig`, `secret`, `license`
- `exit-code` - Exit code when vulnerabilities are found (0-255)
- `ignore-unfixed` - Boolean to ignore unfixed vulnerabilities
- `additional-args` - Array of additional arguments

## Examples

See [examples.yml](examples.yml) for usage examples including:
- Scan container image
- Scan filesystem
- Scan Git repository
- Advanced scan with custom config

## Output

Returns a map containing:
- `command` - Full command executed
- `output` - Combined stdout/stderr
- `exit_code` - Exit code (0 = success)
- `error` - Error message if failed
- `output_file` - Path to output file if specified

## Links

- [Trivy Documentation](https://aquasecurity.github.io/trivy/)
- [Trivy GitHub](https://github.com/aquasecurity/trivy)
