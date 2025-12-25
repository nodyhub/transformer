package convert

import (
	"strings"

	"github.com/nodyhub/transformer/builtin/sarif/common"
)

// createSarifReport creates a SARIF 2.1.0 report structure
func createSarifReport(toolName string, results []common.Result) common.Report {
	return common.Report{
		Schema:  "https://raw.githubusercontent.com/oasis-tcs/sarif-spec/master/Schemata/sarif-schema-2.1.0.json",
		Version: "2.1.0",
		Runs: []common.Run{
			{
				Tool: common.Tool{
					Driver: common.Driver{
						Name:    toolName,
						Version: "1.0.0",
					},
				},
				Results: results,
			},
		},
	}
}

// getStringField extracts a string field from nested maps
func getStringField(m map[string]interface{}, keys ...string) string {
	current := m
	for i, key := range keys {
		if i == len(keys)-1 {
			if val, ok := current[key].(string); ok {
				return val
			}
			return ""
		}
		if next, ok := current[key].(map[string]interface{}); ok {
			current = next
		} else {
			return ""
		}
	}
	return ""
}

// severityToLevel maps severity strings to SARIF levels
func severityToLevel(severity string) string {
	switch strings.ToUpper(severity) {
	case "CRITICAL", "HIGH":
		return "error"
	case "MEDIUM":
		return "warning"
	case "LOW", "INFO":
		return "note"
	default:
		return "warning"
	}
}
