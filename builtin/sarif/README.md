# SARIF Builtin

Work with SARIF (Static Analysis Results Interchange Format) reports.

## Sub-packages

### [convert/](convert/) - builtin/sarif

Convert security tool output to SARIF 2.1.0 format.

Supports: Trivy, Semgrep, Nuclei, CodeQL, Trufflehog, Nmap, Gosec, OSV-Scanner

See [convert/README.md](convert/README.md) for details and examples.

### [merge/](merge/) - builtin/sarif/merge

Merge multiple SARIF 2.1.0 reports into a single unified report.

Combines runs from different security tools into one comprehensive report.

See [merge/README.md](merge/README.md) for details and examples.

### [suppress/](suppress/) - builtin/sarif/suppress

Filter out findings from a SARIF report based on suppression rules.

Supports filtering by rule IDs, paths, severities, and tool names with glob patterns.

See [suppress/README.md](suppress/README.md) for details and examples.

## Quick Example

```yaml
# Convert tool outputs to SARIF
- name: Convert Trivy to SARIF
  id: trivy_sarif
  using: builtin/sarif
  with:
    tool: trivy
    input: ${{ outputs.trivy_scan }}

- name: Convert Semgrep to SARIF
  id: semgrep_sarif
  using: builtin/sarif
  with:
    tool: semgrep
    input: ${{ outputs.semgrep_scan }}

# Merge reports
- name: Merge reports
  id: merged
  using: builtin/sarif/merge
  with:
    reports:
      - ${{ outputs.trivy_sarif }}
      - ${{ outputs.semgrep_sarif }}

# Suppress false positives
- name: Suppress vendor files
  using: builtin/sarif/suppress
  with:
    input: ${{ outputs.merged }}
    path_patterns:
      - "vendor/**"
      - "**/node_modules/**"
    severities:
      - note
```

## Links

- [SARIF Specification](https://sarifweb.azurewebsites.net/)
- [SARIF 2.1.0 Schema](https://raw.githubusercontent.com/oasis-tcs/sarif-spec/master/Schemata/sarif-schema-2.1.0.json)
- `paths` (optional): Array of exact file paths to suppress
- `path_patterns` (optional): Array of glob patterns for file paths
  - `*` matches any characters except `/`
  - `**` matches any characters including `/`
  - `?` matches a single character
- `severities` (optional): Array of severity levels to suppress (`error`, `warning`, `note`)
- `tools` (optional): Array of tool names to suppress

**Output:**
Returns a JSON string containing the filtered SARIF report with suppressed findings removed.

**Examples:**

Suppress specific CVEs:
```yaml
- uses: builtin/sarif/suppress
  with:
    input: ${{ outputs.report }}
    rule_ids:
      - CVE-2024-1234
      - GHSA-xxxx-yyyy-zzzz
```

Suppress vendor and test files:
```yaml
- uses: builtin/sarif/suppress
  with:
    input: ${{ outputs.report }}
    path_patterns:
      - "vendor/**"
      - "**/node_modules/**"
      - "test/**"
```

Suppress low severity findings:
```yaml
- uses: builtin/sarif/suppress
  with:
    input: ${{ outputs.report }}
    severities:
      - note
      - warning
```

Combine multiple criteria:
```yaml
- uses: builtin/sarif/suppress
  with:
    input: ${{ outputs.report }}
    rule_patterns:
      - "CVE-2023-.*"
    path_patterns:
      - "test/**"
    severities:
      - note
```

## Usage

## Usage

### Convert tool output to SARIF

```yaml
- id: trivy_scan
  uses: builtin/tools/trivy
  with:
    type: image
    target: nginx:latest
    format: json

- id: trivy_sarif
  uses: builtin/sarif/convert
  with:
    tool: "trivy"
    input: ${{ outputs.trivy_scan }}
```

### Merge SARIF reports

```yaml
- id: semgrep_scan
  uses: builtin/tools/semgrep
  with:
    target: .
    config: auto
    format: json

- id: semgrep_sarif
  uses: builtin/sarif/convert
  with:
    tool: "semgrep"
    input: ${{ outputs.semgrep_scan }}

- id: merge_reports
  uses: builtin/sarif/merge
  with:
    reports:
      - ${{ outputs.trivy_sarif }}
      - ${{ outputs.semgrep_sarif }}

- id: save_report
  uses: builtin/file/write
  with:
    path: "security-report.sarif"
    content: ${{ outputs.merge_reports }}
```

## Examples

### Complete security scanning pipeline

```yaml
steps:
  # Run Trivy container scan
  - id: trivy_scan
    using: builtin/tools/trivy
    with:
      type: image
      target: nginx:latest
      format: json
  
  # Run Semgrep code scan
  - id: semgrep_scan
    using: builtin/tools/semgrep
    with:
      target: .
      config: auto
      format: json
  
  # Run Nuclei web scan
  - id: nuclei_scan
    using: builtin/tools/nuclei
    with:
      target: https://example.com
      format: json
  
  # Convert to SARIF
  - id: trivy_sarif
    using: builtin/sarif/convert
    with:
      tool: trivy
      input: ${{ outputs.trivy_scan }}
  
  - id: semgrep_sarif
    using: builtin/sarif/convert
    with:
      tool: semgrep
      input: ${{ outputs.semgrep_scan }}
  
  - id: nuclei_sarif
    using: builtin/sarif/convert
    with:
      tool: nuclei
      input: ${{ outputs.nuclei_scan }}
  
  # Merge all reports
  - id: unified_report
    using: builtin/sarif/merge
    with:
      reports:
        - ${{ outputs.trivy_sarif }}
        - ${{ outputs.semgrep_sarif }}
        - ${{ outputs.nuclei_sarif }}
  
  # Save to file
  - id: save_report
    using: builtin/file/write
    with:
      path: "security-scan-results.sarif"
      content: ${{ outputs.unified_report }}
  
  # Print confirmation
  - using: builtin/echo
    with:
      message: "Unified SARIF report saved to security-scan-results.sarif"
```

## SARIF Format

The merged report follows the SARIF 2.1.0 specification. Each tool's results are preserved in separate runs within the merged report.

Learn more about SARIF: https://sarifweb.azurewebsites.net/
