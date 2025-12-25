package convert

import (
	"context"
	"encoding/json"
	"fmt"
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

func extractInput(withMap map[string]interface{}) (string, error) {
	inputRaw, ok := withMap["input"]
	if !ok {
		return "", fmt.Errorf("'input' field is required")
	}

	switch v := inputRaw.(type) {
	case string:
		return v, nil
	case []byte:
		return string(v), nil
	default:
		return fmt.Sprintf("%v", v), nil
	}
}

func checkIfAlreadySARIF(input string) (interface{}, bool) {
	var jsonData interface{}
	if err := json.Unmarshal([]byte(input), &jsonData); err != nil {
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

func routeToConverter(tool, input string) (interface{}, error) {
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
