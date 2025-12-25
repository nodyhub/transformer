# SARIF Suppress

Filter out findings from a SARIF report based on suppression rules.

## Usage

```yaml
- name: Suppress false positives
  using: builtin/sarif/suppress
  with:
    input: ${{ outputs.merged_sarif }}
    path_patterns:
      - "vendor/**"
      - "**/node_modules/**"
    severities:
      - note
```

## Parameters

### Required
- `input` - SARIF report (as JSON string or object)

### Optional
- `rule_ids` - Array of exact rule IDs to suppress (e.g., `CVE-2024-1234`)
- `rule_patterns` - Array of regex patterns for rule IDs (e.g., `CVE-2024-.*`)
- `paths` - Array of exact file paths to suppress
- `path_patterns` - Array of glob patterns for file paths (supports `*`, `?`, `**/`)
- `severities` - Array of severity levels to suppress: `error`, `warning`, `note`
- `tools` - Array of tool names to suppress findings from

## Glob Pattern Syntax

- `*` - Match any characters except `/`
- `?` - Match any single character except `/`
- `**/` - Match zero or more directories
- Examples:
  - `vendor/**` - All files under vendor directory
  - `**/*.test.*` - All test files
  - `src/*/generated/*.go` - Generated Go files in src subdirectories

## Use Cases

- Suppress known false positives
- Exclude vendor/dependency code from scans
- Filter out low-severity findings
- Suppress specific CVEs or vulnerability IDs
- Exclude test files from production scans
- Filter findings by tool name

## Output

Returns a JSON string containing the filtered SARIF 2.1.0 report with suppressed findings removed.

## Examples

See [examples.yml](examples.yml) for usage examples including:
- Suppress by path patterns (vendor, tests)
- Suppress by severity level
- Suppress by rule ID
- Suppress by CVE patterns
- Combine multiple suppression criteria

## Links

- [SARIF Specification](https://sarifweb.azurewebsites.net/)
- [SARIF 2.1.0 Schema](https://raw.githubusercontent.com/oasis-tcs/sarif-spec/master/Schemata/sarif-schema-2.1.0.json)
