# File Builtin

Write files to disk.

## Usage

### Write

Write or append content to a file.

```yaml
- id: write_report
  uses: builtin/file/write
  with:
    path: "report.sarif"
    content: ${{ outputs.merged_sarif }}
    mode: "write"  # or "append"
    perm: 0644     # optional, default 0644
```

**Parameters:**
- `path` (required): File path to write to
- `content` (required): Content to write (string, bytes, or any value that can be converted to string)
- `mode` (optional): Write mode - "write" (default, truncate) or "append"
- `perm` (optional): File permissions as integer (default: 0644)

**Output:**
```json
{
  "path": "report.sarif",
  "bytes": 1234
}
```

## Examples

### Write SARIF report
```yaml
- id: save_report
  uses: builtin/file/write
  with:
    path: "security-report.sarif"
    content: ${{ outputs.merged_sarif }}
```

### Append to log file
```yaml
- id: append_log
  uses: builtin/file/write
  with:
    path: "logs/output.log"
    content: "New log entry\n"
    mode: "append"
```

### Write to nested directory
```yaml
- id: save_config
  uses: builtin/file/write
  with:
    path: "output/reports/config.json"
    content: ${{ outputs.config }}
```
