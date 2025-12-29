package convert

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/nodyhub/transformer/builtin/sarif/common"
)

// convertTrufflehog converts Trufflehog JSON output to SARIF
func convertTrufflehog(input interface{}) (interface{}, error) {
	// Trufflehog outputs JSONL (one JSON object per line)
	lines := strings.Split(strings.TrimSpace(input.(string)), "\n")
	results := []common.Result{}

	for _, line := range lines {
		if line == "" {
			continue
		}

		var tfResult map[string]interface{}
		if err := json.Unmarshal([]byte(line), &tfResult); err != nil {
			continue
		}

		detectorName := getStringField(tfResult, "DetectorName")
		decoderName := getStringField(tfResult, "DecoderName")
		sourceName := getStringField(tfResult, "SourceMetadata", "Data", "Filesystem", "file")
		if sourceName == "" {
			sourceName = getStringField(tfResult, "SourceMetadata", "Data", "Git", "file")
		}
		if sourceName == "" {
			sourceName = "unknown"
		}

		// Get line number if available
		startLine := 1
		if sourceMetadata, ok := tfResult["SourceMetadata"].(map[string]interface{}); ok {
			if data, ok := sourceMetadata["Data"].(map[string]interface{}); ok {
				if git, ok := data["Git"].(map[string]interface{}); ok {
					if line, ok := git["line"].(float64); ok {
						startLine = int(line)
					}
				} else if fs, ok := data["Filesystem"].(map[string]interface{}); ok {
					if line, ok := fs["line"].(float64); ok {
						startLine = int(line)
					}
				}
			}
		}

		message := fmt.Sprintf("Secret detected: %s", detectorName)
		if decoderName != "" {
			message = fmt.Sprintf("%s (decoder: %s)", message, decoderName)
		}

		results = append(results, common.Result{
			RuleID: detectorName,
			Level:  "error",
			Message: common.Message{
				Text: message,
			},
			Locations: []common.Location{
				{
					PhysicalLocation: common.PhysicalLocation{
						ArtifactLocation: common.ArtifactLocation{
							URI: sourceName,
						},
						Region: common.Region{
							StartLine: startLine,
						},
					},
				},
			},
			Properties: map[string]interface{}{
				"detector": detectorName,
				"decoder":  decoderName,
			},
		})
	}

	return createSarifReport("Trufflehog", results), nil
}
