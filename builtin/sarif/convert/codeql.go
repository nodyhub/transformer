package convert

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/nodyhub/transformer/builtin/sarif/common"
)

// convertCodeQL converts CodeQL output to SARIF
// CodeQL natively outputs SARIF, but can also output JSON or CSV
func convertCodeQL(input string) (interface{}, error) {
	// Try to parse as JSON first
	var jsonData interface{}
	if err := json.Unmarshal([]byte(input), &jsonData); err != nil {
		return nil, fmt.Errorf("failed to parse CodeQL output: %w", err)
	}

	// Check if it's already SARIF format (CodeQL's default output)
	if jsonMap, ok := jsonData.(map[string]interface{}); ok {
		if version, ok := jsonMap["version"].(string); ok && strings.HasPrefix(version, "2.1.") {
			if _, ok := jsonMap["runs"]; ok {
				// Already SARIF 2.1.0, return as-is
				return jsonData, nil
			}
		}

		// Handle CodeQL's JSON result format (non-SARIF)
		// This is the format when using --format=json instead of --format=sarif
		if columns, ok := jsonMap["#select"].(map[string]interface{}); ok {
			return convertCodeQLJSON(jsonMap, columns)
		}
	}

	return nil, fmt.Errorf("unrecognized CodeQL output format")
}

// convertCodeQLJSON converts CodeQL JSON format to SARIF
func convertCodeQLJSON(output map[string]interface{}, _ map[string]interface{}) (interface{}, error) {
	results := []common.Result{}

	// CodeQL JSON has results in various formats depending on the query
	// Common structure: tuples array with data
	if tuples, ok := output["tuples"].([]interface{}); ok {
		for _, tupleRaw := range tuples {
			tuple, ok := tupleRaw.([]interface{})
			if !ok || len(tuple) == 0 {
				continue
			}

			// Extract information from tuple
			// Typical format: [location, message, severity, ...]
			var message, path string
			var startLine int
			severity := "warning"

			// Parse tuple elements
			for i, elem := range tuple {
				if elemMap, ok := elem.(map[string]interface{}); ok {
					// Check for location information
					if url, ok := elemMap["url"].(string); ok {
						path = url
						if line, ok := elemMap["startLine"].(float64); ok {
							startLine = int(line)
						}
					}
					// Check for label (message)
					if label, ok := elemMap["label"].(string); ok {
						message = label
					}
				} else if i == 0 {
					// First element often contains the message
					if str, ok := elem.(string); ok && message == "" {
						message = str
					}
				}
			}

			if message == "" {
				message = "CodeQL finding"
			}
			if startLine == 0 {
				startLine = 1
			}

			results = append(results, common.Result{
				RuleID: "codeql-query",
				Level:  severity,
				Message: common.Message{
					Text: message,
				},
				Locations: []common.Location{
					{
						PhysicalLocation: common.PhysicalLocation{
							ArtifactLocation: common.ArtifactLocation{
								URI: path,
							},
							Region: common.Region{
								StartLine: startLine,
							},
						},
					},
				},
			})
		}
	}

	return createSarifReport("CodeQL", results), nil
}
