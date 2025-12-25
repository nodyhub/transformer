package codeql

import (
	"context"
	"fmt"
	"os/exec"
	"strings"

	"github.com/nodyhub/transformer/registry"
)

func init() {
	registry.Register("builtin/tools/codeql", CodeQL)
}

// CodeQL executes GitHub CodeQL code analysis
// Supported inputs:
//   - database: path to CodeQL database or source to create database from (required)
//   - action: action to perform (database-create, database-analyze, analyze)
//   - language: language for database creation (go, javascript, python, java, cpp, csharp, ruby)
//   - source-root: source root directory for database creation
//   - query: query pack or QL file to run
//   - format: output format (sarif-latest, csv, json)
//   - output: output file path
//   - threads: number of threads (default: 0 = auto)
//   - ram: amount of RAM in MB
//   - additional-args: array of additional arguments
func CodeQL(ctx context.Context, with interface{}) (interface{}, error) {
	withMap, ok := with.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("expected 'with' to be a map")
	}

	// Build command arguments
	args := []string{"codeql"}

	// Action (default: database-analyze if database exists)
	action := "database-analyze"
	if a, ok := withMap["action"].(string); ok && a != "" {
		action = a
	}
	args = append(args, action)

	// Database path (required)
	database, ok := withMap["database"].(string)
	if !ok || database == "" {
		return nil, fmt.Errorf("'database' field is required")
	}

	// Handle different actions
	switch action {
	case "database-create":
		args = append(args, database)

		// Language (required for database-create)
		if language, ok := withMap["language"].(string); ok && language != "" {
			args = append(args, "--language="+language)
		} else {
			return nil, fmt.Errorf("'language' field is required for database-create")
		}

		// Source root
		if sourceRoot, ok := withMap["source-root"].(string); ok && sourceRoot != "" {
			args = append(args, "--source-root="+sourceRoot)
		}

	case "database-analyze", "analyze":
		args = append(args, database)

		// Query pack or file
		if query, ok := withMap["query"].(string); ok && query != "" {
			args = append(args, query)
		}

		// Format
		if format, ok := withMap["format"].(string); ok && format != "" {
			args = append(args, "--format="+format)
		} else {
			args = append(args, "--format=sarif-latest")
		}

		// Output file
		if output, ok := withMap["output"].(string); ok && output != "" {
			args = append(args, "--output="+output)
		}
	}

	// Threads
	if threads, ok := withMap["threads"].(int); ok {
		args = append(args, fmt.Sprintf("--threads=%d", threads))
	}

	// RAM
	if ram, ok := withMap["ram"].(int); ok {
		args = append(args, fmt.Sprintf("--ram=%d", ram))
	}

	// Additional arguments
	if additionalArgs, ok := withMap["additional-args"].([]interface{}); ok {
		for _, arg := range additionalArgs {
			if argStr, ok := arg.(string); ok {
				args = append(args, argStr)
			}
		}
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
