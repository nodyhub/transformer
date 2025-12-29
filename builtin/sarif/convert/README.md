# SARIF Convert

Convert security tool output to SARIF 2.1.0 format.

## Usage

```yaml
- name: Convert Trivy to SARIF
  uses: builtin/sarif/convert
  with:
    tool: trivy
    input: ${{ outputs.trivy_scan }}
```

## Parameters

### Required
- `tool` - The security tool name: `trivy`, `semgrep`, `nuclei`, `codeql`, `trufflehog`, `nmap`, `gosec`, `osv`
- `input` - The tool's output (string, JSON, or raw output)

## Supported Tools

- **Trivy**: Converts JSON vulnerability output
- **Semgrep**: Converts JSON results
- **Nuclei**: Converts JSONL output
- **CodeQL**: Converts SARIF output (pass-through) or JSON format
- **Trufflehog**: Converts JSONL secret detection output
- **Nmap**: Converts JSON network scan output
- **Gosec**: Converts JSON security scan output for Go code
- **OSV-Scanner**: Converts JSON vulnerability output from OSV database

## Severity Mapping

Tool severities are automatically mapped to SARIF levels:
- `CRITICAL`, `HIGH` → `error`
- `MEDIUM` → `warning`
- `LOW`, `INFO` → `note`

## Output

Returns a SARIF 2.1.0 report structure (map) containing:
- Version: `2.1.0`
- Schema: SARIF JSON schema reference
- Runs: Array of tool runs with results, locations, and metadata

## Examples

See [examples.yml](examples.yml) for usage examples including:
- Convert Trivy vulnerability scan
- Convert Semgrep code analysis
- Convert Nuclei security scan
- Convert Trufflehog secret detection
- Convert Nmap network scan
- Convert Gosec Go security scan
- Convert OSV-Scanner vulnerability scan
- Convert CodeQL analysis

## Links

- [SARIF Specification](https://sarifweb.azurewebsites.net/)
- [SARIF 2.1.0 Schema](https://raw.githubusercontent.com/oasis-tcs/sarif-spec/master/Schemata/sarif-schema-2.1.0.json)
