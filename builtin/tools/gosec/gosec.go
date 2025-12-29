package gosec

import (
	"context"
	"fmt"
	"strings"

	"github.com/nodyhub/transformer/registry"
)

func init() {
	registry.Register("builtin/tools/gosec", Gosec)
}

// Gosec executes Gosec Go security scanner
// Supported inputs:
//   - target: path to scan (default: "./...")
//   - format: output format (json, yaml, csv, junit-xml, html, sonarqube, golint, sarif, text)
//   - output: output file path
//   - severity: filter by severity (low, medium, high)
//   - confidence: filter by confidence (low, medium, high)
//   - exclude: exclude directories (comma-separated)
//   - exclude-generated: exclude generated files
//   - nosec: ignore #nosec comments
//   - tests: scan test files
//   - additional-args: array of additional arguments
func Gosec(ctx context.Context, with interface{}) (interface{}, error) {
	withMap, ok := with.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("expected 'with' to be a map")
	}

	args := buildGosecArgs(withMap)

	// Build the command string
	cmdStr := strings.Join(args, " ")

	shellInput := map[string]interface{}{
		"command": cmdStr,
	}
	// Optionally pass through stdout/stderr/output if present in withMap
	if v, ok := withMap["stdout"]; ok {
		shellInput["stdout"] = v
	}
	if v, ok := withMap["stderr"]; ok {
		shellInput["stderr"] = v
	}

	// Call builtin/shell
	return registry.Call(ctx, "builtin/shell", shellInput)
}

func buildGosecArgs(m map[string]interface{}) []string {
	args := []string{"gosec"}

	// Format
	if format, ok := m["format"].(string); ok && format != "" {
		args = append(args, "-fmt", format)
	} else {
		args = append(args, "-fmt", "json")
	}

	// Output file
	if output, ok := m["output"].(string); ok && output != "" {
		args = append(args, "-out", output)
	}

	// Severity
	if severity, ok := m["severity"].(string); ok && severity != "" {
		args = append(args, "-severity", severity)
	}

	// Confidence
	if confidence, ok := m["confidence"].(string); ok && confidence != "" {
		args = append(args, "-confidence", confidence)
	}

	// Exclude directories
	if exclude, ok := m["exclude"].(string); ok && exclude != "" {
		args = append(args, "-exclude", exclude)
	}

	// Exclude generated files
	if excludeGenerated, ok := m["exclude-generated"].(bool); ok && excludeGenerated {
		args = append(args, "-exclude-generated")
	}

	// Ignore #nosec comments
	if nosec, ok := m["nosec"].(bool); ok && nosec {
		args = append(args, "-nosec")
	}

	// Include tests
	if tests, ok := m["tests"].(bool); ok && tests {
		args = append(args, "-tests")
	}

	// no-fail: default true, can be overridden
	noFail := true
	if nf, ok := m["no-fail"].(bool); ok {
		noFail = nf
	}
	if noFail {
		args = append(args, "-no-fail")
	}

	// Additional arguments
	if additionalArgs, ok := m["additional-args"].([]interface{}); ok {
		for _, arg := range additionalArgs {
			if argStr, ok := arg.(string); ok {
				args = append(args, argStr)
			}
		}
	}

	// Target (default: ./...)
	target := "./..."
	if t, ok := m["target"].(string); ok && t != "" {
		target = t
	}
	args = append(args, target)
	return args
}
