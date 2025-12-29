# CodeQL

CodeQL is GitHub's semantic code analysis engine for discovering vulnerabilities across a codebase.

## Usage

```yaml
- name: Analyze code with CodeQL
  uses: builtin/tools/codeql
  with:
    action: database-analyze
    database: /path/to/codeql-db
    query: security-extended
    format: sarif-latest
    output: codeql-results.sarif
```

## Parameters

### Required
- `database` - Path to CodeQL database or source directory

### Optional
- `action` - Action to perform: `database-create`, `database-analyze` (default), `analyze`
- `language` - Language for database creation: `go`, `javascript`, `python`, `java`, `cpp`, `csharp`, `ruby` (required for database-create)
- `source-root` - Source root directory (for database-create)
- `query` - Query pack or QL file to run (e.g., `security-extended`, `security-and-quality`)
- `format` - Output format: `sarif-latest` (default), `csv`, `json`
- `output` - Output file path
- `threads` - Number of threads (default: 0 = auto)
- `ram` - Amount of RAM in MB
- `additional-args` - Array of additional arguments

## Examples

See [examples.yml](examples.yml) for usage examples including:
- Create database
- Analyze database
- Complete workflow (Go)
- Python analysis
- Java analysis with build command
- Custom query

## Supported Languages

- **Go** (`go`)
- **JavaScript/TypeScript** (`javascript`)
- **Python** (`python`)
- **Java/Kotlin** (`java`)
- **C/C++** (`cpp`)
- **C#** (`csharp`)
- **Ruby** (`ruby`)

## Query Packs

Common query packs available:
- `security-extended` - Security-focused queries
- `security-and-quality` - Security + quality queries
- Language-specific: `go-security-extended.qls`, `python-security-and-quality.qls`, etc.

## Output Formats

- `sarif-latest` - SARIF format (recommended for integration)
- `csv` - Comma-separated values
- `json` - JSON format

## Output

Returns a map containing:
- `command` - Full command executed
- `output` - Combined stdout/stderr
- `exit_code` - Exit code (0 = success)
- `error` - Error message if failed
- `output_file` - Path to output file if specified

## Notes

- **x86_64 only**: CodeQL currently only supports x86_64 architecture (no ARM64 Linux binaries)
- Database creation can be time-consuming for large codebases
- Ensure sufficient RAM for analysis (recommended: 4GB+)
- For compiled languages (Java, C++), you may need to specify build commands

## Links

- [CodeQL Documentation](https://codeql.github.com/docs/)
- [CodeQL Query Reference](https://codeql.github.com/codeql-query-help/)
- [CodeQL CLI](https://codeql.github.com/docs/codeql-cli/)
- [CodeQL GitHub](https://github.com/github/codeql)
