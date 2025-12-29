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

	action := getString(withMap, "action", "database-analyze")
	database := getString(withMap, "database", "")
	if database == "" {
		return nil, fmt.Errorf("'database' field is required")
	}

	var args []string
	args = append(args, "codeql", action)

	switch action {
	case "database-create":
		a, err := buildCreateArgs(withMap, database)
		if err != nil {
			return nil, err
		}
		args = append(args, a...)
	case "database-analyze", "analyze":
		args = append(args, buildAnalyzeArgs(withMap, database)...)
	}

	args = append(args, buildCommonArgs(withMap)...)

	// Execute command
	cmd := exec.CommandContext(ctx, args[0], args[1:]...)
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

	if outputPath := getString(withMap, "output", ""); outputPath != "" {
		result["output_file"] = outputPath
	}

	return result, nil
}

func getString(m map[string]interface{}, key, def string) string {
	if v, ok := m[key]; ok {
		if s, ok := v.(string); ok && s != "" {
			return s
		}
	}
	return def
}

func getInt(m map[string]interface{}, key string) (int, bool) {
	if v, ok := m[key]; ok {
		if i, ok := v.(int); ok {
			return i, true
		}
	}
	return 0, false
}

func buildCreateArgs(m map[string]interface{}, database string) ([]string, error) {
	args := []string{database}
	language := getString(m, "language", "")
	if language == "" {
		return nil, fmt.Errorf("'language' field is required for database-create")
	}
	args = append(args, "--language="+language)
	if sourceRoot := getString(m, "source-root", ""); sourceRoot != "" {
		args = append(args, "--source-root="+sourceRoot)
	}
	return args, nil
}

func buildAnalyzeArgs(m map[string]interface{}, database string) []string {
	args := []string{database}
	if query := getString(m, "query", ""); query != "" {
		args = append(args, query)
	}
	if format := getString(m, "format", ""); format != "" {
		args = append(args, "--format="+format)
	} else {
		args = append(args, "--format=sarif-latest")
	}
	if output := getString(m, "output", ""); output != "" {
		args = append(args, "--output="+output)
	}
	return args
}

func buildCommonArgs(m map[string]interface{}) []string {
	args := []string{}
	if threads, ok := getInt(m, "threads"); ok {
		args = append(args, fmt.Sprintf("--threads=%d", threads))
	}
	if ram, ok := getInt(m, "ram"); ok {
		args = append(args, fmt.Sprintf("--ram=%d", ram))
	}
	if additionalArgs, ok := m["additional-args"].([]interface{}); ok {
		for _, arg := range additionalArgs {
			if argStr, ok := arg.(string); ok {
				args = append(args, argStr)
			}
		}
	}
	return args
}
