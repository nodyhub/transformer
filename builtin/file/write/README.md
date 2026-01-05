# File Write

Write files to disk.

## Usage

```yaml
- id: write_file
  uses: builtin/file/write
  with:
    path: "file.txt"
    content: "content to write"
    mode: "write"      # optional: "write" (default) or "append"
    perm: 0644         # optional: file permissions (default: 0644)
```

**Parameters:**
- `path` (required): File path to write to
- `content` (required): Content to write (string, bytes, or any value)
- `mode` (optional): Write mode - "write" (default, truncate) or "append"
- `perm` (optional): File permissions as integer (default: 0644)

**Output:**
```json
{
  "path": "file.txt",
  "bytes": 1234
}
```
