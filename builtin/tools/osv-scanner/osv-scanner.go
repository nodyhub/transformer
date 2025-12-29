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

	args, err := buildOSVScannerArgs(withMap)
	if err != nil {
		return nil, err
	}

	cmd := exec.CommandContext(ctx, args[0], args[1:]...) // #nosec
	var stdout, stderr strings.Builder
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err = cmd.Run()

	result := map[string]interface{}{
		"command": strings.Join(args, " "),
		"output":  stdout.String(),
		"stderr":  stderr.String(),
	}

	if err != nil {
		result["error"] = err.Error()
		if cmd.ProcessState != nil {
			result["exit_code"] = cmd.ProcessState.ExitCode()
		} else {
			result["exit_code"] = -1
		}
	} else {
		result["exit_code"] = 0
	}

	if outputPath, ok := withMap["output"].(string); ok && outputPath != "" {
		result["output_file"] = outputPath
	}

	return result, nil
}

func buildOSVScannerArgs(m map[string]interface{}) ([]string, error) {
	args := []string{"osv-scanner"}

	// Format
	if format, ok := m["format"].(string); ok && format != "" {
		args = append(args, "--format", format)
	} else {
		args = append(args, "--format", "json")
	}

	// Output file
	if output, ok := m["output"].(string); ok && output != "" {
		args = append(args, "--output", output)
	}

	// Call analysis
	if callAnalysis, ok := m["call-analysis"].(bool); ok && callAnalysis {
		args = append(args, "--call-analysis")
	}

	// Recursive
	if recursive, ok := m["recursive"].(bool); ok && recursive {
		args = append(args, "--recursive")
	}

	// Skip git
	if skipGit, ok := m["skip-git"].(bool); ok && skipGit {
		args = append(args, "--skip-git")
	}

	// Experimental offline
	if offline, ok := m["experimental-offline"].(bool); ok && offline {
		args = append(args, "--experimental-offline")
	}

	// Additional arguments
	if additionalArgs, ok := m["additional-args"].([]interface{}); ok {
		for _, arg := range additionalArgs {
			if argStr, ok := arg.(string); ok {
				args = append(args, argStr)
			}
		}
	}

	// Target (required) - can be a directory, lockfile, or SBOM
	if target, ok := m["target"].(string); ok && target != "" {
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

	return args, nil
}
