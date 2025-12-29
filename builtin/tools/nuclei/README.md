# Nuclei

Nuclei is a fast vulnerability scanner that uses templates to find security issues.

## Usage

```yaml
- name: Nuclei vulnerability scan
  uses: builtin/tools/nuclei
  with:
    target: https://example.com
    severity: high,critical
    output: nuclei-results.json
```

## Parameters

### Required
- `target` - Target URL(s), file containing URLs, or array of URLs

### Optional
- `templates` - Template or template directory to run
- `workflows` - Workflow or workflow directory to run
- `severity` - Filter by severity: `info`, `low`, `medium`, `high`, `critical`
- `tags` - Filter by tags (comma-separated)
- `exclude-tags` - Exclude tags (comma-separated)
- `output` - Output file path
- `format` - Output format: `json`, `sarif`, `markdown`
- `silent` - Boolean to display results only
- `stats` - Boolean to display scan statistics
- `update-templates` - Boolean to update templates before scanning
- `additional-args` - Array of additional arguments

## Examples

See [examples.yml](examples.yml) for usage examples including:
- Scan web application
- Scan multiple targets
- CVE detection
- Custom workflows

## Output

Returns a map containing:
- `command` - Full command executed
- `output` - Combined stdout/stderr
- `exit_code` - Exit code (0 = success)
- `error` - Error message if failed
- `output_file` - Path to output file if specified

## Links

- [Nuclei Documentation](https://docs.projectdiscovery.io/nuclei/)
- [Nuclei Templates](https://github.com/projectdiscovery/nuclei-templates)
- [Nuclei GitHub](https://github.com/projectdiscovery/nuclei)
