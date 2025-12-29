package nuclei

import (
	"context"
	"fmt"
	"os/exec"
	"strings"

	"github.com/nodyhub/transformer/registry"
)

func init() {
	registry.Register("builtin/tools/nuclei", Nuclei)
}

// Nuclei executes Nuclei vulnerability scanner
// Supported inputs:
//   - target: target URL(s) or file containing URLs (required)
//   - templates: template or template directory to run
//   - workflows: workflow or workflow directory to run
//   - severity: filter by severity (info, low, medium, high, critical)
//   - tags: filter by tags (comma-separated)
//   - exclude-tags: exclude tags (comma-separated)
//   - output: output file path
//   - format: output format (json, sarif, markdown)
//   - silent: display results only
//   - stats: display scan statistics
//   - update-templates: update templates before scanning
//   - additional-args: array of additional arguments
func Nuclei(ctx context.Context, with interface{}) (interface{}, error) {
	withMap, ok := with.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("expected 'with' to be a map")
	}

	args, err := buildNucleiArgs(withMap)
	if err != nil {
		return nil, err
	}

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

	if outputPath, ok := withMap["output"].(string); ok && outputPath != "" {
		result["output_file"] = outputPath
	}

	return result, nil
}

func buildNucleiArgs(m map[string]interface{}) ([]string, error) {
	args := []string{"nuclei"}

	// Target (required)
	if target, ok := m["target"].(string); ok && target != "" {
		if strings.HasPrefix(target, "http://") || strings.HasPrefix(target, "https://") {
			args = append(args, "-u", target)
		} else {
			args = append(args, "-l", target)
		}
	} else if targets, ok := m["target"].([]interface{}); ok {
		for _, t := range targets {
			if targetStr, ok := t.(string); ok {
				args = append(args, "-u", targetStr)
			}
		}
	} else {
		return nil, fmt.Errorf("'target' field is required")
	}

	// Templates
	if templates, ok := m["templates"].(string); ok && templates != "" {
		args = append(args, "-t", templates)
	}

	// Workflows
	if workflows, ok := m["workflows"].(string); ok && workflows != "" {
		args = append(args, "-w", workflows)
	}

	// Severity
	if severity, ok := m["severity"].(string); ok && severity != "" {
		args = append(args, "-severity", severity)
	}

	// Tags
	if tags, ok := m["tags"].(string); ok && tags != "" {
		args = append(args, "-tags", tags)
	}

	// Exclude tags
	if excludeTags, ok := m["exclude-tags"].(string); ok && excludeTags != "" {
		args = append(args, "-exclude-tags", excludeTags)
	}

	// Output file
	if output, ok := m["output"].(string); ok && output != "" {
		args = append(args, "-o", output)
	}

	// Format
	if format, ok := m["format"].(string); ok && format != "" {
		args = append(args, "-"+format)
	}

	// Silent mode
	if silent, ok := m["silent"].(bool); ok && silent {
		args = append(args, "-silent")
	}

	// Stats
	if stats, ok := m["stats"].(bool); ok && stats {
		args = append(args, "-stats")
	}

	// Update templates
	if updateTemplates, ok := m["update-templates"].(bool); ok && updateTemplates {
		args = append(args, "-update-templates")
	}

	// Additional arguments
	if additionalArgs, ok := m["additional-args"].([]interface{}); ok {
		for _, arg := range additionalArgs {
			if argStr, ok := arg.(string); ok {
				args = append(args, argStr)
			}
		}
	}

	return args, nil
}
