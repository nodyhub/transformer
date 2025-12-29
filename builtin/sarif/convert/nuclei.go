package convert

import (
	"encoding/json"
	"strings"

	"github.com/nodyhub/transformer/builtin/sarif/common"
)

// convertNuclei converts Nuclei JSON output to SARIF
func convertNuclei(input interface{}) (interface{}, error) {
	// Nuclei outputs JSONL (one JSON object per line)
	lines := strings.Split(strings.TrimSpace(input.(string)), "\n")
	results := []common.Result{}

	for _, line := range lines {
		if line == "" {
			continue
		}

		var nucleiResult map[string]interface{}
		if err := json.Unmarshal([]byte(line), &nucleiResult); err != nil {
			continue
		}

		templateID := getStringField(nucleiResult, "template-id")
		matchedAt := getStringField(nucleiResult, "matched-at")
		name := getStringField(nucleiResult, "info", "name")
		severity := getStringField(nucleiResult, "info", "severity")

		message := name
		if message == "" {
			message = templateID
		}

		results = append(results, common.Result{
			RuleID: templateID,
			Level:  severityToLevel(severity),
			Message: common.Message{
				Text: message,
			},
			Locations: []common.Location{
				{
					PhysicalLocation: common.PhysicalLocation{
						ArtifactLocation: common.ArtifactLocation{
							URI: matchedAt,
						},
						Region: common.Region{
							StartLine: 1,
						},
					},
				},
			},
		})
	}

	return createSarifReport("Nuclei", results), nil
}
