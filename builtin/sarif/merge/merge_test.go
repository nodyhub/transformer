package merge

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/nodyhub/transformer/builtin/sarif/common"
)

func TestMerge(t *testing.T) {
	ctx := context.Background()

	report1 := `{
		"version": "2.1.0",
		"$schema": "https://raw.githubusercontent.com/oasis-tcs/sarif-spec/master/Schemata/sarif-schema-2.1.0.json",
		"runs": [{
			"tool": {
				"driver": {
					"name": "Trivy",
					"version": "0.50.0"
				}
			},
			"results": [{
				"ruleId": "CVE-2024-1234",
				"level": "error",
				"message": {
					"text": "High severity vulnerability"
				}
			}]
		}]
	}`

	report2 := `{
		"version": "2.1.0",
		"$schema": "https://raw.githubusercontent.com/oasis-tcs/sarif-spec/master/Schemata/sarif-schema-2.1.0.json",
		"runs": [{
			"tool": {
				"driver": {
					"name": "Semgrep",
					"version": "1.50.0"
				}
			},
			"results": [{
				"ruleId": "go.lang.security.audit.xss",
				"level": "warning",
				"message": {
					"text": "Potential XSS vulnerability"
				}
			}]
		}]
	}`

	tests := []struct {
		name        string
		with        interface{}
		expectError bool
		validate    func(t *testing.T, result interface{})
	}{
		{
			name: "merge two reports",
			with: map[string]interface{}{
				"reports": []interface{}{report1, report2},
			},
			expectError: false,
			validate: func(t *testing.T, result interface{}) {
				resultStr, ok := result.(string)
				if !ok {
					t.Fatal("expected result to be a string")
				}

				var merged common.Report
				if err := json.Unmarshal([]byte(resultStr), &merged); err != nil {
					t.Fatalf("failed to parse merged report: %v", err)
				}

				if len(merged.Runs) != 2 {
					t.Errorf("expected 2 runs, got %d", len(merged.Runs))
				}

				if merged.Runs[0].Tool.Driver.Name != "Trivy" {
					t.Errorf("expected first tool to be Trivy, got %s", merged.Runs[0].Tool.Driver.Name)
				}

				if merged.Runs[1].Tool.Driver.Name != "Semgrep" {
					t.Errorf("expected second tool to be Semgrep, got %s", merged.Runs[1].Tool.Driver.Name)
				}
			},
		},
		{
			name: "merge single report",
			with: map[string]interface{}{
				"reports": []interface{}{report1},
			},
			expectError: false,
			validate: func(t *testing.T, result interface{}) {
				resultStr, ok := result.(string)
				if !ok {
					t.Fatal("expected result to be a string")
				}

				var merged common.Report
				if err := json.Unmarshal([]byte(resultStr), &merged); err != nil {
					t.Fatalf("failed to parse merged report: %v", err)
				}

				if len(merged.Runs) != 1 {
					t.Errorf("expected 1 run, got %d", len(merged.Runs))
				}
			},
		},
		{
			name: "invalid JSON",
			with: map[string]interface{}{
				"reports": []interface{}{"invalid json"},
			},
			expectError: true,
		},
		{
			name:        "missing reports field",
			with:        map[string]interface{}{},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Merge(ctx, tt.with)
			if tt.expectError {
				if err == nil {
					t.Error("expected error but got none")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if tt.validate != nil {
				tt.validate(t, result)
			}
		})
	}
}

func TestMergeEmptyReports(t *testing.T) {
	ctx := context.Background()

	result, err := Merge(ctx, map[string]interface{}{
		"reports": []interface{}{},
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	resultStr := result.(string)
	if !strings.Contains(resultStr, `"runs": []`) {
		t.Error("expected empty runs array")
	}
}
