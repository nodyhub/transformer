package convert

import (
	"encoding/json"
	"fmt"

	"github.com/nodyhub/transformer/builtin/sarif/common"
)

// convertSemgrep converts Semgrep JSON output to SARIF
func convertSemgrep(input string) (interface{}, error) {
	var semgrepOutput map[string]interface{}
	if err := json.Unmarshal([]byte(input), &semgrepOutput); err != nil {
		return nil, fmt.Errorf("failed to parse Semgrep output: %w", err)
	}

	results := []common.Result{}

	// Semgrep output has "results" array
	if resultsRaw, ok := semgrepOutput["results"].([]interface{}); ok {
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
