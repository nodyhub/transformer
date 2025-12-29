package semgrep

import (
	"context"
	"fmt"
	"os/exec"
	"strings"

	"github.com/nodyhub/transformer/registry"
)

func init() {
	registry.Register("builtin/tools/semgrep", Semgrep)
}

// Semgrep executes Semgrep static analysis
// Supported inputs:
//   - target: path to scan (default: ".")
//   - config: rules config (auto, p/ci, p/security-audit, or path)
//   - format: output format (json, sarif, text, gitlab-sast, junit-xml)
//   - output: output file path
//   - severity: filter by severity (INFO, WARNING, ERROR)
//   - exclude: exclude patterns (can be array)
//   - max-memory: maximum memory in MB (default: 5000)
//   - metrics: send anonymous metrics (on/off)
//   - additional-args: array of additional arguments
func Semgrep(ctx context.Context, with interface{}) (interface{}, error) {
	withMap, ok := with.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("expected 'with' to be a map")
	}

	args := buildSemgrepArgs(withMap)

	cmd := exec.CommandContext(ctx, args[0], args[1:]...) // #nosec
	var stdout, stderr strings.Builder
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()

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

func buildSemgrepArgs(m map[string]interface{}) []string {
	args := []string{"semgrep", "scan"}

	// Config (default: auto)
	config := "auto"
	if c, ok := m["config"].(string); ok && c != "" {
		config = c
	}
	args = append(args, "--config", config)

	// Format
	if format, ok := m["format"].(string); ok && format != "" {
		args = append(args, "--"+format)
	} else {
		args = append(args, "--json")
	}

	// Output file
	if output, ok := m["output"].(string); ok && output != "" {
		args = append(args, "--output", output)
	}

	// Severity
	if severity, ok := m["severity"].(string); ok && severity != "" {
		args = append(args, "--severity", severity)
	}

	// Exclude patterns
	if exclude, ok := m["exclude"].([]interface{}); ok {
		for _, pattern := range exclude {
			if patternStr, ok := pattern.(string); ok {
				args = append(args, "--exclude", patternStr)
			}
		}
	} else if exclude, ok := m["exclude"].(string); ok && exclude != "" {
		args = append(args, "--exclude", exclude)
	}

	// Max memory
	if maxMemory, ok := m["max-memory"].(int); ok {
		args = append(args, "--max-memory", fmt.Sprintf("%d", maxMemory))
	}

	// Metrics
	if metrics, ok := m["metrics"].(string); ok && metrics != "" {
		args = append(args, "--metrics", metrics)
	}

	// Additional arguments
	if additionalArgs, ok := m["additional-args"].([]interface{}); ok {
		for _, arg := range additionalArgs {
			if argStr, ok := arg.(string); ok {
				args = append(args, argStr)
			}
		}
	}

	// Target (default: current directory)
	target := "."
	if t, ok := m["target"].(string); ok && t != "" {
		target = t
	}
	args = append(args, target)
	return args
}
