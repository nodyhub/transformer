package convert

import (
	"encoding/json"
	"fmt"

	"github.com/nodyhub/transformer/builtin/sarif/common"
)

// convertTrivy converts Trivy JSON output to SARIF
func convertTrivy(input string) (interface{}, error) {
	var trivyOutput map[string]interface{}
	if err := json.Unmarshal([]byte(input), &trivyOutput); err != nil {
		return nil, fmt.Errorf("failed to parse Trivy output: %w", err)
	}

	results := []common.Result{}

	// Trivy output has "Results" array
	if resultsRaw, ok := trivyOutput["Results"].([]interface{}); ok {
		for _, resultRaw := range resultsRaw {
			resultMap, ok := resultRaw.(map[string]interface{})
			if !ok {
				continue
			}

			target := ""
			if t, ok := resultMap["Target"].(string); ok {
				target = t
			}

			// Process vulnerabilities
			if vulnsRaw, ok := resultMap["Vulnerabilities"].([]interface{}); ok {
				for _, vulnRaw := range vulnsRaw {
					vulnMap, ok := vulnRaw.(map[string]interface{})
					if !ok {
						continue
					}

					vulnID := getStringField(vulnMap, "VulnerabilityID")
					pkgName := getStringField(vulnMap, "PkgName")
					severity := getStringField(vulnMap, "Severity")
					title := getStringField(vulnMap, "Title")

					message := fmt.Sprintf("%s: %s", vulnID, title)
					if message == ": " {
						message = vulnID
					}

					results = append(results, common.Result{
						RuleID: vulnID,
						Level:  severityToLevel(severity),
						Message: common.Message{
							Text: message,
						},
						Locations: []common.Location{
							{
								PhysicalLocation: common.PhysicalLocation{
									ArtifactLocation: common.ArtifactLocation{
										URI: target,
									},
									Region: common.Region{
										StartLine: 1,
									},
								},
							},
						},
						Properties: map[string]interface{}{
							"package":  pkgName,
							"severity": severity,
						},
					})
				}
			}
		}
	}

	return createSarifReport("Trivy", results), nil
}
