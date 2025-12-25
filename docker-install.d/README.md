# Security Tool Installation Scripts

This directory contains modular installation scripts for security scanning tools used by the transformer.

## Naming Convention

Scripts are prefixed with two digits (01-99) to control execution order:

- `01-trivy.sh` - Trivy vulnerability scanner
- `02-semgrep.sh` - Semgrep static analysis
- `03-nuclei.sh` - Nuclei web scanner
- `04-codeql.sh` - CodeQL code analysis (x86_64 only)
- `05-trufflehog.sh` - Trufflehog secret scanner
- `06-nmap.sh` - Nmap network scanner
- `07-gosec.sh` - Gosec Go security scanner
- `08-osv-scanner.sh` - OSV vulnerability scanner

## Script Requirements

Each script must:

1. Start with the bash shebang: `#!/bin/bash`
2. Use strict error handling: `set -euo pipefail`
3. Install its own dependencies (e.g., `ca-certificates`, `wget`, `tar`, etc.)
4. Echo what it's installing for build visibility
5. Verify the installation with a `--version` check
6. Be idempotent (safe to run multiple times)
7. Clean up temporary files and APT cache

## Script Template

```bash
#!/bin/bash
set -euo pipefail

echo "Installing ToolName dependencies..."
apt-get update
apt-get install -y --no-install-recommends ca-certificates wget tar
rm -rf /var/lib/apt/lists/*

echo "Installing ToolName..."
TOOL_VERSION="1.0.0"
wget -qO tool.tar.gz "https://example.com/tool_${TOOL_VERSION}_linux_amd64.tar.gz"
tar -xzf tool.tar.gz -C /usr/local/bin
rm tool.tar.gz
chmod +x /usr/local/bin/tool
tool --version
```

## Adding a New Tool

1. Create a new script with the next available number:
   ```bash
   touch install.d/09-newtool.sh
   chmod +x install.d/09-newtool.sh
   ```

2. Follow the template above

3. Update the tool list in this README

4. Rebuild the Docker image:
   ```bash
   docker build -t transformer:latest .
   ```

## Updating Tool Versions

Edit the `*_VERSION` variable at the top of each script:

```bash
# Before
TRIVY_VERSION="0.48.0"

# After
TRIVY_VERSION="0.49.0"
```

Then rebuild the image.

## Selective Module Installation

By default, all security tools are installed. You can build a smaller image by installing only specific modules:

```bash
# Install only Trivy and Semgrep
docker build --build-arg MODULES="trivy,semgrep" -t transformer:lite .

# Install only Nmap
docker build --build-arg MODULES="nmap" -t transformer:nmap .

# Install all tools (default)
docker build -t transformer:full .
# or
docker build --build-arg MODULES="all" -t transformer:full .
```

Available modules:
- `trivy` - Vulnerability scanner
- `semgrep` - Static analysis
- `nuclei` - Web scanner
- `codeql` - Code analysis
- `truffleEach module installs its own required dependencies
- `nmap` - Network scanner
- `gosec` - Go security scanner
- `osv-scanner` - Vulnerability scanner

**Note:** The base system dependencies (00-base.sh) are always installed.

## Disabling a Tool

To temporarily disable a tool installation without deleting the script:

1. Rename it to remove the `.sh` extension:
   ```bash
   mv install.d/07-gosec.sh install.d/07-gosec.sh.disabled
   ```

2. Or move it to a backup directory:
   ```bash
   mkdir -p install.d/.disabled
   mv install.d/07-gosec.sh install.d/.disabled/
   ```

3. Or use the MODULES build argument (recommended):
   ```bash
   docker build --build-arg MODULES="trivy,semgrep" -t transformer:custom .
   ```

## Execution Order

Scripts are executed in alphabetical/numerical order. When using the MODULES argument, only the base (00-base.sh) and specified modules are executed:

```dockerfile
RUN chmod +x /tmp/install.d/*.sh && \
    for script in /tmp/install.d/*.sh; do \
        echo "Running $(basename "$script")..."; \
        bash "$script"; \
    done
```

The `00-base.sh` script must run first to install system dependencies required by other tools.

## Tool Versions

Current versions (update this when changing version variables):

| Tool | Version | Release Notes | Status |
|------|---------|---------------|--------|
| Trivy | Latest (apt) | https://github.com/aquasecurity/trivy/releases | ✅ Active |
| Semgrep | Latest (pip) | https://github.com/semgrep/semgrep/releases | ✅ Active |
| Nuclei | 3.3.6 | https://github.com/projectdiscovery/nuclei/releases | ✅ Active |
| CodeQL | 2.15.5 | https://github.com/github/codeql-cli-binaries/releases | ⚠️ x86_64 only |
| Trufflehog | 3.82.13 | https://github.com/trufflesecurity/trufflehog/releases | ✅ Active |
| Nmap | Latest (apt) | https://nmap.org/download.html | ✅ Active |
| Gosec | 2.18.2 | https://github.com/securego/gosec/releases | ✅ Active |
| OSV-Scanner | 1.9.1 | https://github.com/google/osv-scanner/releases | ✅ Active |

## Architecture-Specific Tools

### CodeQL (x86_64 only)

CodeQL installation is automatically skipped on ARM64 systems because GitHub doesn't provide ARM64 Linux binaries. The script detects the architecture at runtime:

- **x86_64/amd64**: CodeQL is installed normally
- **ARM64/aarch64**: Installation is skipped with an informational message

**Running on ARM64:**
- The Docker build will succeed but CodeQL won't be available
- Use x86_64 hardware/VMs for CodeQL support
- Consider multi-platform builds: `docker buildx build --platform linux/amd64,linux/arm64`

**Check for ARM64 support:**
```bash
curl -sL https://api.github.com/repos/github/codeql-cli-binaries/releases/latest | jq -r '.assets[].name' | grep arm64
```

## Troubleshooting

### Script fails during build

Check the Docker build output to see which script failed:
```bash
docker build -t transformer:latest . 2>&1 | grep -A 10 "Running.*\.sh"
```

### Tool not found after installation

Verify the tool was installed correctly:
```bash
docker run --rm transformer which <toolname>
docker run --rm transformer <toolname> --version
```

### Permission errors

Ensure all scripts are executable:
```bash
chmod +x install.d/*.sh
```

## Best Practices

- **Pin versions**: Use specific versions instead of `latest` for reproducibility
- **Verify checksums**: Add SHA256 verification for downloaded binaries when available
- **Minimize layers**: Keep related commands in the same script
- **Clean up**: Remove downloaded archives and temp files
- **Test locally**: Build and test the Docker image before committing changes

## Example: Adding Bandit (Python Security)

```bash
cat > install.d/09-bandit.sh << 'EOF'
#!/bin/bash
set -euo pipefail

echo "Installing Bandit..."
pip3 install --no-cache-dir bandit --break-system-packages
bandit --version
EOF

chmod +x install.d/09-bandit.sh
docker build -t transformer:latest .
```
