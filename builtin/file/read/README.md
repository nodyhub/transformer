# File Read

Read files from disk.

## Usage

```yaml
- id: read_file
  uses: builtin/file/read
  with:
    path: "file.txt"
```

**Parameters:**
- `path` (required): File path to read from

**Output:**
```json
{
  "path": "file.txt",
  "content": "file content...",
  "bytes": 1234
}
```
