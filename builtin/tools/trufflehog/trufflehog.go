package trufflehog

import (
	"context"
	"fmt"
	"os/exec"
	"strings"

	"github.com/nodyhub/transformer/registry"
)

func init() {
	registry.Register("builtin/tools/trufflehog", Trufflehog)
}

// Trufflehog executes Trufflehog secret scanner
// Supported inputs:
//   - target: git repository URL, filesystem path, or docker image (required)
//   - type: scan type (git, github, gitlab, filesystem, docker, s3)
//   - format: output format (json, yaml)
//   - output: output file path
//   - since-commit: scan commits since this commit
//   - branch: scan specific branch
//   - max-depth: maximum commit depth (default: 1000000)
//   - include-paths: file with path patterns to include
//   - exclude-paths: file with path patterns to exclude
//   - only-verified: only show verified secrets
//   - json: output as JSON
//   - additional-args: array of additional arguments
func Trufflehog(ctx context.Context, with interface{}) (interface{}, error) {
	withMap, ok := with.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("expected 'with' to be a map")
	}

	args, err := buildTrufflehogArgs(withMap)
	if err != nil {
		return nil, err
	}

	cmd := exec.CommandContext(ctx, args[0], args[1:]...) // #nosec
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

func buildTrufflehogArgs(m map[string]interface{}) ([]string, error) {
	target, ok := m["target"].(string)
	if !ok || target == "" {
		return nil, fmt.Errorf("'target' field is required")
	}

	args := []string{"trufflehog"}

	// Scan type (default: git)
	scanType := "git"
	if t, ok := m["type"].(string); ok && t != "" {
		scanType = t
	}
	args = append(args, scanType)

	// JSON output
	if jsonOutput, ok := m["json"].(bool); ok && jsonOutput {
		args = append(args, "--json")
	}

	// Only verified
	if onlyVerified, ok := m["only-verified"].(bool); ok && onlyVerified {
		args = append(args, "--only-verified")
	}

	// Since commit
	if sinceCommit, ok := m["since-commit"].(string); ok && sinceCommit != "" {
		args = append(args, "--since-commit", sinceCommit)
	}

	// Branch
	if branch, ok := m["branch"].(string); ok && branch != "" {
		args = append(args, "--branch", branch)
	}

	// Max depth
	if maxDepth, ok := m["max-depth"].(int); ok {
		args = append(args, "--max-depth", fmt.Sprintf("%d", maxDepth))
	}

	// Include paths
	if includePaths, ok := m["include-paths"].(string); ok && includePaths != "" {
		args = append(args, "--include-paths", includePaths)
	}

	// Exclude paths
	if excludePaths, ok := m["exclude-paths"].(string); ok && excludePaths != "" {
		args = append(args, "--exclude-paths", excludePaths)
	}

	// Additional arguments
	if additionalArgs, ok := m["additional-args"].([]interface{}); ok {
		for _, arg := range additionalArgs {
			if argStr, ok := arg.(string); ok {
				args = append(args, argStr)
			}
		}
	}

	// Add target
	args = append(args, target)
	return args, nil
}
