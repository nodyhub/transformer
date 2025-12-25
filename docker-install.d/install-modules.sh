#!/bin/bash
set -euo pipefail

# Script to selectively install security scanning modules
# Usage: install-modules.sh [MODULES]
#   MODULES: comma-separated list of modules (e.g., "trivy,semgrep,nmap")
#            or "all" to install all modules
#            or empty to install all modules

MODULES="${1:-}"

# If MODULES is empty or "all", install all tools
if [ -z "$MODULES" ] || [ "$MODULES" = "all" ]; then
    for script in /tmp/install.d/[0-9][0-9]-*.sh; do
        echo "===== Running $(basename "$script") ====="
        bash "$script" || { echo "FAILED: $(basename "$script")"; exit 1; }
        echo "===== SUCCESS: $(basename "$script") ====="
    done
else
    # Install only specified modules
    IFS="," read -ra MODULE_LIST <<< "$MODULES"
    for module in "${MODULE_LIST[@]}"; do
        module=$(echo "$module" | xargs)
        script=$(ls /tmp/install.d/*-${module}.sh 2>/dev/null | head -1)
        if [ -n "$script" ]; then
            echo "===== Running $(basename "$script") ====="
            bash "$script" || { echo "FAILED: $(basename "$script")"; exit 1; }
            echo "===== SUCCESS: $(basename "$script") ====="
        else
            echo "WARNING: Module '$module' not found, skipping..."
        fi
    done
fi
