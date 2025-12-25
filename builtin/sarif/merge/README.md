# SARIF Merge

Merge multiple SARIF 2.1.0 reports into a single unified report.

## Usage

```yaml
- name: Merge SARIF reports
  using: builtin/sarif/merge
  with:
    reports:
      - ${{ outputs.trivy_sarif }}
      - ${{ outputs.semgrep_sarif }}
```

## Parameters

### Required
- `reports` - Array of SARIF reports (as JSON strings or objects)

## Description

Combines multiple SARIF reports from different security tools into a single unified report. All runs from each input report are merged into the output report, preserving tool information and results.

## Output

Returns a JSON string containing the merged SARIF 2.1.0 report with:
- Version: `2.1.0`
- Schema: SARIF JSON schema reference
- Runs: Combined array of all runs from input reports

## Use Cases

- Unified security reporting from multiple tools
- Aggregating vulnerability scans, code analysis, and secret detection
- Creating comprehensive security dashboards
- Feeding combined results to issue tracking systems

## Examples

See [examples.yml](examples.yml) for usage examples including:
- Merge two reports
- Merge multiple tools (7+ reports)
- Merge and save to file

## Links

- [SARIF Specification](https://sarifweb.azurewebsites.net/)
- [SARIF 2.1.0 Schema](https://raw.githubusercontent.com/oasis-tcs/sarif-spec/master/Schemata/sarif-schema-2.1.0.json)
