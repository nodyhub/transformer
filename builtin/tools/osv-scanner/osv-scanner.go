package osvscanner

import (
	"context"
	"fmt"
	"os/exec"
	"strings"

	"github.com/nodyhub/transformer/registry"
)

func init() {
	registry.Register("builtin/tools/osv-scanner", OSVScanner)
}

// OSVScanner executes OSV-Scanner vulnerability scanner
// Supported inputs:
//   - target: directory, lockfile, or SBOM to scan (required)
//   - format: output format (json, table, markdown, sarif)
//   - output: output file path
//   - call-analysis: enable call analysis for Go
//   - recursive: scan directories recursively
//   - skip-git: skip scanning git repositories
//   - experimental-offline: scan in offline mode
//   - additional-args: array of additional arguments
func OSVScanner(ctx context.Context, with interface{}) (interface{}, error) {
	withMap, ok := with.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("expected 'with' to be a map")
	}

	// Build command arguments
	args := []string{"osv-scanner"}

	// Format
	if format, ok := withMap["format"].(string); ok && format != "" {
		args = append(args, "--format", format)
	} else {
		args = append(args, "--format", "json")
	}

	// Output file
	if output, ok := withMap["output"].(string); ok && output != "" {
		args = append(args, "--output", output)
	}

	// Call analysis
	if callAnalysis, ok := withMap["call-analysis"].(bool); ok && callAnalysis {
		args = append(args, "--call-analysis")
	}

	// Recursive
	if recursive, ok := withMap["recursive"].(bool); ok && recursive {
		args = append(args, "--recursive")
	}

	// Skip git
	if skipGit, ok := withMap["skip-git"].(bool); ok && skipGit {
		args = append(args, "--skip-git")
	}

	// Experimental offline
	if offline, ok := withMap["experimental-offline"].(bool); ok && offline {
		args = append(args, "--experimental-offline")
	}

	// Additional arguments
	if additionalArgs, ok := withMap["additional-args"].([]interface{}); ok {
		for _, arg := range additionalArgs {
			if argStr, ok := arg.(string); ok {
				args = append(args, argStr)
			}
		}
	}

	// Target (required) - can be a directory, lockfile, or SBOM
	if target, ok := withMap["target"].(string); ok && target != "" {
		// Check if it's a lockfile or directory scan
		if strings.Contains(target, "package-lock.json") ||
			strings.Contains(target, "go.mod") ||
			strings.Contains(target, "requirements.txt") ||
			strings.Contains(target, "Cargo.lock") {
			args = append(args, "--lockfile", target)
		} else {
			args = append(args, target)
		}
	} else {
		return nil, fmt.Errorf("'target' field is required")
	}

	// Execute command
	cmd := exec.CommandContext(ctx, args[0], args[1:]...)
	output, err := cmd.CombinedOutput()

	result := map[string]interface{}{
		"command": strings.Join(args, " "),
		"output":  string(output),
	}

	if err != nil {
		result["error"] = err.Error()
		result["exit_code"] = cmd.ProcessState.ExitCode()
	} else {
		result["exit_code"] = 0
	}

	// If output file was specified, return the path
	if outputPath, ok := withMap["output"].(string); ok && outputPath != "" {
		result["output_file"] = outputPath
	}

	return result, nil
}
