package trivy

import (
	"context"
	"fmt"
	"os/exec"
	"strings"

	"github.com/nodyhub/transformer/registry"
)

func init() {
	registry.Register("builtin/tools/trivy", Trivy)
}

// Trivy executes Trivy vulnerability scanner
// Supported inputs:
//   - target: image name, filesystem path, or repository (required)
//   - type: scan type (image, fs, repo, config, sbom) default: image
//   - format: output format (json, sarif, table, cyclonedx, spdx) default: json
//   - severity: comma-separated severities (UNKNOWN,LOW,MEDIUM,HIGH,CRITICAL)
//   - output: output file path
//   - scanners: comma-separated scanners (vuln, misconfig, secret, license)
//   - exit-code: exit code when vulnerabilities are found (0-255)
//   - ignore-unfixed: ignore unfixed vulnerabilities
//   - additional-args: array of additional arguments
func Trivy(ctx context.Context, with interface{}) (interface{}, error) {
	withMap, ok := with.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("expected 'with' to be a map")
	}

	// Required: target
	target, ok := withMap["target"].(string)
	if !ok || target == "" {
		return nil, fmt.Errorf("'target' field is required")
	}

	// Build command arguments
	args := []string{"trivy"}

	// Scan type (default: image)
	scanType := "image"
	if t, ok := withMap["type"].(string); ok && t != "" {
		scanType = t
	}
	args = append(args, scanType)

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

	// Severity
	if severity, ok := withMap["severity"].(string); ok && severity != "" {
		args = append(args, "--severity", severity)
	}

	// Scanners
	if scanners, ok := withMap["scanners"].(string); ok && scanners != "" {
		args = append(args, "--scanners", scanners)
	}

	// Exit code
	if exitCode, ok := withMap["exit-code"].(int); ok {
		args = append(args, "--exit-code", fmt.Sprintf("%d", exitCode))
	}

	// Ignore unfixed
	if ignoreUnfixed, ok := withMap["ignore-unfixed"].(bool); ok && ignoreUnfixed {
		args = append(args, "--ignore-unfixed")
	}

	// Additional arguments
	if additionalArgs, ok := withMap["additional-args"].([]interface{}); ok {
		for _, arg := range additionalArgs {
			if argStr, ok := arg.(string); ok {
				args = append(args, argStr)
			}
		}
	}

	// Add target
	args = append(args, target)

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
