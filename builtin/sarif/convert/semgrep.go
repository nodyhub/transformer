package convert

import (
	"encoding/json"
	"fmt"

	"github.com/nodyhub/transformer/builtin/sarif/common"
)

// convertSemgrep converts Semgrep JSON output to SARIF
func convertSemgrep(input interface{}) (interface{}, error) {
	var semgrepJson map[string]interface{}
	switch v := input.(type) {
	case string:
		if err := json.Unmarshal([]byte(v), &semgrepJson); err != nil {
			return nil, fmt.Errorf("invalid Semgrep output format: could not parse string as JSON object: %w", err)
		}
	case map[string]interface{}:
		semgrepJson = v
	default:
		return nil, fmt.Errorf("invalid Semgrep output format: expected a JSON object or string, got %T", input)
	}

	results := []common.Result{}

	// Semgrep output has "results" array
	if resultsRaw, ok := semgrepJson["results"].([]interface{}); ok {
		for _, resultRaw := range resultsRaw {
			resultMap, ok := resultRaw.(map[string]interface{})
			if !ok {
				continue
			}

			checkID := getStringField(resultMap, "check_id")
			path := getStringField(resultMap, "path")
			message := getStringField(resultMap, "extra", "message")
			severity := getStringField(resultMap, "extra", "severity")

			startLine := 1
			if start, ok := resultMap["start"].(map[string]interface{}); ok {
				if line, ok := start["line"].(float64); ok {
					startLine = int(line)
				}
			}

			results = append(results, common.Result{
				RuleID: checkID,
				Level:  severityToLevel(severity),
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

	return createSarifReport("Semgrep", results), nil
}
