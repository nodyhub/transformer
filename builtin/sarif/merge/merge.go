package merge

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/nodyhub/transformer/builtin/sarif/common"
	"github.com/nodyhub/transformer/registry"
)

func init() {
	registry.Register("builtin/sarif/merge", Merge)
}

// Merge combines multiple SARIF reports into a single report
func Merge(ctx context.Context, with interface{}) (interface{}, error) {
	withMap, ok := with.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("expected 'with' to be a map")
	}

	// Get reports
	reportsRaw, ok := withMap["reports"]
	if !ok {
		return nil, fmt.Errorf("'reports' field is required")
	}

	reports, ok := reportsRaw.([]interface{})
	if !ok {
		return nil, fmt.Errorf("'reports' must be an array")
	}

	// Create merged report
	merged := common.Report{
		Version: "2.1.0",
		Schema:  "https://raw.githubusercontent.com/oasis-tcs/sarif-spec/master/Schemata/sarif-schema-2.1.0.json",
		Runs:    []common.Run{},
	}

	// Parse and merge each report
	for i, reportRaw := range reports {
		var report common.Report

		// Handle different input types
		switch v := reportRaw.(type) {
		case string:
			// Parse JSON string
			if err := json.Unmarshal([]byte(v), &report); err != nil {
				return nil, fmt.Errorf("failed to parse report %d: %w", i, err)
			}
		case map[string]interface{}:
			// Convert map to JSON and back to struct
			data, err := json.Marshal(v)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal report %d: %w", i, err)
			}
			if err := json.Unmarshal(data, &report); err != nil {
				return nil, fmt.Errorf("failed to parse report %d: %w", i, err)
			}
		default:
			return nil, fmt.Errorf("report %d has invalid type: %T", i, reportRaw)
		}

		// Add runs from this report to merged report
		merged.Runs = append(merged.Runs, report.Runs...)
	}

	// Convert merged report to JSON string
	mergedJSON, err := json.MarshalIndent(merged, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal merged report: %w", err)
	}

	return string(mergedJSON), nil
}
