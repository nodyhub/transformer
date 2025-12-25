# Nmap

Nmap is a network exploration tool and security/port scanner.

## Usage

```yaml
- name: Port scan
  using: builtin/tools/nmap
  with:
    target: 192.168.1.1
    ports: 1-1000
    output: nmap-results
```

## Parameters

### Required
- `target` - Target host(s), IP address(es), or network(s)

### Optional
- `ports` - Port specification (e.g., `22,80,443` or `1-1000`)
- `scan-type` - Scan type: `-sS` (SYN), `-sT` (TCP), `-sU` (UDP), `-sV` (version), `-sC` (scripts), `-A` (aggressive)
- `output` - Output file path (base name, extensions added automatically)
- `format` - Output format: `normal`, `xml`, `grepable`, `all`
- `timing` - Timing template: `0-5`, or `paranoid`, `sneaky`, `polite`, `normal`, `aggressive`, `insane`
- `script` - NSE script(s) to run
- `os-detection` - Boolean to enable OS detection (`-O`)
- `version-detection` - Boolean to enable version detection (`-sV`)
- `aggressive` - Boolean to enable aggressive scan (`-A`)
- `additional-args` - Array of additional arguments

## Examples

See [examples.yml](examples.yml) for usage examples including:
- Basic port scan
- Network range scan
- Service detection
- OS detection
- Aggressive scan
- Vulnerability scanning
- Custom NSE scripts
- Stealth scan

## Timing Templates

- `0` / `paranoid` - Very slow, IDS evasion
- `1` / `sneaky` - Slow, IDS evasion
- `2` / `polite` - Slower, less bandwidth
- `3` / `normal` - Default
- `4` / `aggressive` - Faster, assumes fast network
- `5` / `insane` - Very fast, may miss results

## Output

Returns a map containing:
- `command` - Full command executed
- `output` - Combined stdout/stderr
- `exit_code` - Exit code (0 = success)
- `error` - Error message if failed
- `output_file` - Path to output file if specified

## Links

- [Nmap Documentation](https://nmap.org/docs.html)
- [Nmap NSE Scripts](https://nmap.org/nsedoc/)
- [Nmap Website](https://nmap.org/)
