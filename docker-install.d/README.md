# Docker Installation Scripts

This directory contains installation scripts for Docker image setup.

## Files

- `00-system-utilities.sh` - Installs essential system libraries and utilities (ca-certificates, wget, curl, jq, unzip, etc.)
- `install-modules.sh` - Selectively installs security scanning modules

## Usage

The installation scripts are integrated into the Docker build process via the Dockerfile. The `00-system-utilities.sh` script runs first to install base system dependencies.

## Selective Module Installation

You can selectively install modules using the `install-modules.sh` script:

```bash
# Install all modules
./install-modules.sh all

# Install specific modules (comma-separated)
./install-modules.sh "trivy,semgrep"

# Install single module
./install-modules.sh "nmap"
```

## Adding New Modules

To add new installation modules:

1. Create a script named `XX-modulename.sh` in this directory
2. Follow the naming convention: two-digit prefix (01-99) followed by module name
3. Include the bash shebang and strict error handling:
   ```bash
   #!/bin/bash
   set -euo pipefail
   ```
4. Install dependencies and tools as needed
5. The script will be automatically discovered by `install-modules.sh`

## Adding New Modules

To add new installation modules:

1. Create a script named `XX-modulename.sh` in this directory
2. Follow the naming convention: two-digit prefix (01-99) followed by module name
3. Include the bash shebang and strict error handling:
   ```bash
   #!/bin/bash
   set -euo pipefail
   ```
4. Install dependencies and tools as needed
5. The script will be automatically discovered by `install-modules.sh`
