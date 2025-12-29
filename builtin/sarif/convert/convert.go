package convert

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/nodyhub/transformer/registry"
)

func init() {
	registry.Register("builtin/sarif/convert", Convert)
}

// Convert converts tool output to SARIF 2.1.0 format
func Convert(ctx context.Context, with interface{}) (interface{}, error) {
	withMap, ok := with.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("expected 'with' to be a map")
	}

	tool, err := extractToolName(withMap)
	if err != nil {
		return nil, err
	}

	input, err := extractInput(withMap)
	if err != nil {
		return nil, err
	}

	// If input is a shell wrapper result, extract 'output' and parse as JSON if possible
	if inputMap, ok := input.(map[string]interface{}); ok {
		if out, ok := inputMap["output"].(string); ok && out != "" {
			var parsed interface{}
			if err := json.Unmarshal([]byte(out), &parsed); err == nil {
				input = parsed
			} else {
				input = out
			}
		}
	}

	// Check if input is already SARIF 2.1.0
	if sarifData, isSARIF := checkIfAlreadySARIF(input); isSARIF {
		return sarifData, nil
	}

	return routeToConverter(tool, input)
}

func extractToolName(withMap map[string]interface{}) (string, error) {
	toolRaw, ok := withMap["tool"]
	if !ok {
		return "", fmt.Errorf("'tool' field is required")
	}
	tool, ok := toolRaw.(string)
	if !ok {
		return "", fmt.Errorf("'tool' must be a string")
	}
	return tool, nil
}

func extractInput(withMap map[string]interface{}) (interface{}, error) {
	inputRaw, ok := withMap["input"]
	if !ok {
		return "", fmt.Errorf("'input' field is required")
	}

	// check if inputRaw is a filepath or raw content
	inputStr, isString := inputRaw.(string)
	if !isString {
		return inputRaw, nil
	}

	if _, err := os.Stat(inputStr); err != nil {
		// not a file, treat as raw content
		return inputStr, nil
	}

	// it's a file, read content
	absPath, err := filepath.Abs(inputStr)
	if err != nil {
		return nil, fmt.Errorf("failed to get absolute path of input file: %v", err)
	}
	content, err := os.ReadFile(absPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read input file: %v", err)
	}

	return string(content), nil
}

func checkIfAlreadySARIF(input interface{}) (interface{}, bool) {
	inputStr, ok := input.(string)
	if !ok {
		return nil, false
	}

	var jsonData interface{}
	if err := json.Unmarshal([]byte(inputStr), &jsonData); err != nil {
		return nil, false
	}
	jsonMap, ok := jsonData.(map[string]interface{})
	if !ok {
		return nil, false
	}

	version, ok := jsonMap["version"].(string)
	if !ok || !strings.HasPrefix(version, "2.1.") {
		return nil, false
	}

	if _, ok := jsonMap["runs"]; !ok {
		return nil, false
	}

	return jsonData, true
}

func routeToConverter(tool string, input interface{}) (interface{}, error) {
	switch strings.ToLower(tool) {
	case "trivy":
		return convertTrivy(input)
	case "semgrep":
		return convertSemgrep(input)
	case "nuclei":
		return convertNuclei(input)
	case "codeql":
		return convertCodeQL(input)
	case "trufflehog":
		return convertTrufflehog(input)
	case "nmap":
		return convertNmap(input)
	case "gosec":
		return convertGosec(input)
	case "osv":
		return convertOSV(input)
	default:
		return nil, fmt.Errorf("unsupported tool: %s (supported: trivy, semgrep, nuclei, codeql, trufflehog, nmap, gosec, osv)", tool)
	}
}
