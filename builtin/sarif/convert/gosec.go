package convert

import (
	"encoding/json"
	"fmt"

	"github.com/nodyhub/transformer/builtin/sarif/common"
)

// convertGosec converts Gosec JSON output to SARIF
func convertGosec(input interface{}) (interface{}, error) {
	var gosecJson map[string]interface{}
	switch v := input.(type) {
	case string:
		if err := json.Unmarshal([]byte(v), &gosecJson); err != nil {
			return nil, fmt.Errorf("invalid Gosec output format: could not parse string as JSON object: %w", err)
		}
	case map[string]interface{}:
		gosecJson = v
	default:
		return nil, fmt.Errorf("invalid Gosec output format: expected a JSON object or string, got %T", input)
	}

	results := []common.Result{}

	// Gosec output has "Issues" array
	if issuesRaw, ok := gosecJson["Issues"].([]interface{}); ok {
		for _, issueRaw := range issuesRaw {
			issueMap, ok := issueRaw.(map[string]interface{})
			if !ok {
				continue
			}

			ruleID := getStringField(issueMap, "rule_id")
			severity := getStringField(issueMap, "severity")
			confidence := getStringField(issueMap, "confidence")
			details := getStringField(issueMap, "details")
			file := getStringField(issueMap, "file")

			startLine := 1
			if line, ok := issueMap["line"].(string); ok {
				_, _ = fmt.Sscanf(line, "%d", &startLine)
			}

			message := details
			if message == "" {
				message = ruleID
			}

			results = append(results, common.Result{
				RuleID: ruleID,
				Level:  severityToLevel(severity),
				Message: common.Message{
					Text: message,
				},
				Locations: []common.Location{
					{
						PhysicalLocation: common.PhysicalLocation{
							ArtifactLocation: common.ArtifactLocation{
								URI: file,
							},
							Region: common.Region{
								StartLine: startLine,
							},
						},
					},
				},
				Properties: map[string]interface{}{
					"severity":   severity,
					"confidence": confidence,
				},
			})
		}
	}

	return createSarifReport("Gosec", results), nil
}
