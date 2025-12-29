# OSV-Scanner

OSV-Scanner finds existing vulnerabilities affecting your project's dependencies using the OSV database.

## Usage

```yaml
- name: Scan dependencies
  uses: builtin/tools/osv-scanner
  with:
    target: .
    format: sarif
    output: osv.sarif
```

## Parameters

### Required
- `target` - Directory, lockfile, or SBOM to scan

### Optional
- `format` - Output format: `json` (default), `table`, `markdown`, `sarif`
- `output` - Output file path
- `call-analysis` - Boolean to enable call analysis for Go
- `recursive` - Boolean to scan directories recursively
- `skip-git` - Boolean to skip scanning git repositories
- `experimental-offline` - Boolean to scan in offline mode
- `additional-args` - Array of additional arguments

## Examples

See [examples.yml](examples.yml) for usage examples including:
- Scan all dependencies
- Scan specific lockfile
- Go with call analysis
- Python requirements
- Rust Cargo.lock
- Offline scan

## Supported Lockfiles

- Node.js: `package-lock.json`, `yarn.lock`, `pnpm-lock.yaml`
- Go: `go.mod`, `go.sum`
- Python: `requirements.txt`, `Pipfile.lock`, `poetry.lock`
- Rust: `Cargo.lock`
- Ruby: `Gemfile.lock`
- Java: `pom.xml`, `gradle.lockfile`
- .NET: `packages.lock.json`

## Output

Returns a map containing:
- `command` - Full command executed
- `output` - Combined stdout/stderr
- `exit_code` - Exit code (0 = success)
- `error` - Error message if failed
- `output_file` - Path to output file if specified

## Links

- [OSV-Scanner Documentation](https://google.github.io/osv-scanner/)
- [OSV Database](https://osv.dev/)
- [OSV-Scanner GitHub](https://github.com/google/osv-scanner)
