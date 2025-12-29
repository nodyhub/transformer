# Security Tool Wrappers

This package provides builtin wrappers for popular security scanning tools, making them easy to use in transformer workflows.

## Available Tools

- **builtin/tools/trivy** - Trivy vulnerability scanner
- **builtin/tools/semgrep** - Semgrep static analysis
- **builtin/tools/nuclei** - Nuclei vulnerability scanner
- **builtin/tools/codeql** - CodeQL semantic code analysis
- **builtin/tools/gosec** - Gosec Go security scanner
- **builtin/tools/trufflehog** - Trufflehog secret scanner
- **builtin/tools/osv-scanner** - OSV vulnerability scanner
- **builtin/tools/nmap** - Nmap network scanner

## Usage Examples

### Trivy

Scan a container image for vulnerabilities:

```yaml
- name: Scan Docker image
  uses: builtin/tools/trivy
  with:
    type: image
    target: myapp:latest
    format: json
    output: trivy-report.json
    severity: HIGH,CRITICAL
```

Scan filesystem:

```yaml
- name: Scan filesystem
  uses: builtin/tools/trivy
  with:
    type: fs
    target: /path/to/code
    format: sarif
    output: trivy.sarif
    scanners: vuln,secret,misconfig
```

### Semgrep

Run security audit:

```yaml
- name: Semgrep security scan
  uses: builtin/tools/semgrep
  with:
    target: .
    config: p/security-audit
    format: json
    output: semgrep-results.json
    severity: ERROR
```

Custom rules:

```yaml
- name: Semgrep custom rules
  uses: builtin/tools/semgrep
  with:
    target: src/
    config: /path/to/rules.yml
    format: sarif
    output: semgrep.sarif
    exclude:
      - "**/test/**"
      - "**/vendor/**"
```

### Nuclei

Scan web application:

```yaml
- name: Nuclei web scan
  uses: builtin/tools/nuclei
  with:
    target: https://example.com
    templates: /root/nuclei-templates/
    severity: critical,high
    output: nuclei-results.json
    format: json
```

Multiple targets:

```yaml
- name: Scan multiple hosts
  uses: builtin/tools/nuclei
  with:
    target: targets.txt
    tags: cve,exposure
    silent: true
    output: findings.json
```

### CodeQL

Create database and analyze:

```yaml
- name: Create CodeQL database
  uses: builtin/tools/codeql
  with:
    action: database-create
    database: /tmp/codeql-db
    language: go
    source-root: .

- name: Analyze with CodeQL
  uses: builtin/tools/codeql
  with:
    action: database-analyze
    database: /tmp/codeql-db
    query: security-extended
    format: sarif-latest
    output: codeql.sarif
```

Analyze with custom query:

```yaml
- name: Run custom CodeQL query
  uses: builtin/tools/codeql
  with:
    action: database-analyze
    database: /tmp/codeql-db
    query: /path/to/custom.ql
    format: sarif-latest
    output: results.sarif
    threads: 4
    ram: 4096
```

### Gosec

Scan Go code:

```yaml
- name: Gosec security scan
  uses: builtin/tools/gosec
  with:
    target: ./...
    format: sarif
    output: gosec.sarif
    severity: medium
    tests: true
```

### Trufflehog

Scan Git repository:

```yaml
- name: Scan for secrets
  uses: builtin/tools/trufflehog
  with:
    type: git
    target: https://github.com/user/repo.git
    json: true
    only-verified: true
    since-commit: HEAD~50
```

Scan filesystem:

```yaml
- name: Scan filesystem for secrets
  uses: builtin/tools/trufflehog
  with:
    type: filesystem
    target: /path/to/code
    json: true
```

### OSV-Scanner

Scan dependencies:

```yaml
- name: OSV vulnerability scan
  uses: builtin/tools/osv-scanner
  with:
    target: .
    format: sarif
    output: osv.sarif
    recursive: true
    call-analysis: true
```

Scan specific lockfile:

```yaml
- name: Scan package-lock.json
  uses: builtin/tools/osv-scanner
  with:
    target: package-lock.json
    format: json
    output: osv-results.json
```

### Nmap

Port scan:

```yaml
- name: Nmap port scan
  uses: builtin/tools/nmap
  with:
    target: 192.168.1.0/24
    ports: 1-1000
    timing: 4
    output: nmap-results
    format: xml
```

Service detection:

```yaml
- name: Service detection scan
  uses: builtin/tools/nmap
  with:
    target: example.com
    version-detection: true
    script: default,vuln
    output: nmap-scan
    format: all
```

## Complete Security Pipeline Example

```yaml
- name: Trivy container scan
  id: trivy
  uses: builtin/tools/trivy
  with:
    type: image
    target: myapp:latest
    format: json

- name: Convert Trivy to SARIF
  id: trivy-sarif
  uses: builtin/sarif/convert
  with:
    tool: trivy
    input: ${{ outputs.trivy }}

- name: Semgrep code scan
  id: semgrep
  uses: builtin/tools/semgrep
  with:
    target: .
    config: p/security-audit
    format: sarif
    output: semgrep.sarif

- name: Create CodeQL database
  id: codeql-db
  uses: builtin/tools/codeql
  with:
    action: database-create
    database: /tmp/codeql-db
    language: go
    source-root: .

- name: CodeQL analysis
  id: codeql
  uses: builtin/tools/codeql
  with:
    action: database-analyze
    database: /tmp/codeql-db
    query: security-extended
    format: sarif-latest
    output: codeql.sarif

- name: Secret scanning
  id: secrets
  uses: builtin/tools/trufflehog
  with:
    type: filesystem
    target: .
    json: true

- name: Merge SARIF reports
  uses: builtin/sarif/merge
  with:
    inputs:
      - trivy.sarif
      - semgrep.sarif
      - codeql.sarif
    output: merged-report.sarif
```

## Return Value

All tool wrappers return a map with:

- `command`: The full command that was executed
- `output`: The command output (stdout + stderr)
- `exit_code`: The exit code (0 for success)
- `error`: Error message (if any)
- `output_file`: Path to output file (if specified in `with.output`)

## Common Parameters

Most tools support these common parameters:

- `target`: What to scan (required for most tools)
- `format`: Output format (json, sarif, etc.)
- `output`: Output file path
- `additional-args`: Array of additional CLI arguments

## Notes

- Tools must be installed in the Docker container (see install.d/ scripts)
- Use tool-specific formats for best results
- Convert tool outputs to SARIF for unified reporting
- Check exit codes for scan failures
- Some tools may require additional configuration files
