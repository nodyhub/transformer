package suppress

import (
	"context"
	"encoding/json"
	"regexp"
	"strings"
	"testing"

	"github.com/nodyhub/transformer/builtin/sarif/common"
)

func TestSuppress(t *testing.T) {
	ctx := context.Background()

	// Sample SARIF report with multiple findings
	report := `{
		"version": "2.1.0",
		"$schema": "https://raw.githubusercontent.com/oasis-tcs/sarif-spec/master/Schemata/sarif-schema-2.1.0.json",
		"runs": [{
			"tool": {
				"driver": {
					"name": "Trivy"
				}
			},
			"results": [
				{
					"ruleId": "CVE-2024-1234",
					"level": "error",
					"message": {"text": "High severity vulnerability"},
					"locations": [{
						"physicalLocation": {
							"artifactLocation": {"uri": "package.json"},
							"region": {"startLine": 1}
						}
					}]
				},
				{
					"ruleId": "CVE-2024-5678",
					"level": "warning",
					"message": {"text": "Medium severity vulnerability"},
					"locations": [{
						"physicalLocation": {
							"artifactLocation": {"uri": "go.mod"},
							"region": {"startLine": 1}
						}
					}]
				},
				{
					"ruleId": "CVE-2023-9999",
					"level": "note",
					"message": {"text": "Low severity vulnerability"},
					"locations": [{
						"physicalLocation": {
							"artifactLocation": {"uri": "vendor/lib.js"},
							"region": {"startLine": 1}
						}
					}]
				}
			]
		}]
	}`

	tests := []struct {
		name             string
		with             interface{}
		expectError      bool
		expectedCount    int
		shouldContain    []string
		shouldNotContain []string
	}{
		{
			name: "suppress by specific CVE",
			with: map[string]interface{}{
				"input": report,
				"rule_ids": []interface{}{
					"CVE-2024-1234",
				},
			},
			expectedCount:    2,
			shouldNotContain: []string{"CVE-2024-1234"},
			shouldContain:    []string{"CVE-2024-5678", "CVE-2023-9999"},
		},
		{
			name: "suppress by CVE pattern",
			with: map[string]interface{}{
				"input": report,
				"rule_patterns": []interface{}{
					"CVE-2024-.*",
				},
			},
			expectedCount:    1,
			shouldNotContain: []string{"CVE-2024-1234", "CVE-2024-5678"},
			shouldContain:    []string{"CVE-2023-9999"},
		},
		{
			name: "suppress by path",
			with: map[string]interface{}{
				"input": report,
				"paths": []interface{}{
					"package.json",
				},
			},
			expectedCount:    2,
			shouldNotContain: []string{"package.json"},
			shouldContain:    []string{"go.mod", "vendor/lib.js"},
		},
		{
			name: "suppress by path pattern (glob)",
			with: map[string]interface{}{
				"input": report,
				"path_patterns": []interface{}{
					"vendor/**",
				},
			},
			expectedCount:    2,
			shouldNotContain: []string{"vendor/lib.js"},
			shouldContain:    []string{"package.json", "go.mod"},
		},
		{
			name: "suppress by severity",
			with: map[string]interface{}{
				"input": report,
				"severities": []interface{}{
					"note",
				},
			},
			expectedCount:    2,
			shouldNotContain: []string{"CVE-2023-9999"},
			shouldContain:    []string{"CVE-2024-1234", "CVE-2024-5678"},
		},
		{
			name: "suppress multiple criteria",
			with: map[string]interface{}{
				"input": report,
				"rule_ids": []interface{}{
					"CVE-2024-1234",
				},
				"severities": []interface{}{
					"note",
				},
			},
			expectedCount:    1,
			shouldNotContain: []string{"CVE-2024-1234", "CVE-2023-9999"},
			shouldContain:    []string{"CVE-2024-5678"},
		},
		{
			name: "suppress all high severity",
			with: map[string]interface{}{
				"input": report,
				"severities": []interface{}{
					"error",
				},
			},
			expectedCount:    2,
			shouldNotContain: []string{"CVE-2024-1234"},
			shouldContain:    []string{"CVE-2024-5678", "CVE-2023-9999"},
		},
		{
			name: "no suppressions",
			with: map[string]interface{}{
				"input": report,
			},
			expectedCount: 3,
			shouldContain: []string{"CVE-2024-1234", "CVE-2024-5678", "CVE-2023-9999"},
		},
		{
			name:        "missing input",
			with:        map[string]interface{}{},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Suppress(ctx, tt.with)
			if tt.expectError {
				if err == nil {
					t.Error("expected error but got none")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			resultStr, ok := result.(string)
			if !ok {
				t.Fatal("expected result to be a string")
			}

			// Parse result
			var suppressedReport common.Report
			if err := json.Unmarshal([]byte(resultStr), &suppressedReport); err != nil {
				t.Fatalf("failed to parse suppressed report: %v", err)
			}

			// Check result count
			totalResults := 0
			for _, run := range suppressedReport.Runs {
				totalResults += len(run.Results)
			}
			if totalResults != tt.expectedCount {
				t.Errorf("expected %d results, got %d", tt.expectedCount, totalResults)
			}

			// Check should contain
			for _, expected := range tt.shouldContain {
				if !strings.Contains(resultStr, expected) {
					t.Errorf("expected result to contain %q", expected)
				}
			}

			// Check should not contain
			for _, notExpected := range tt.shouldNotContain {
				if strings.Contains(resultStr, notExpected) {
					t.Errorf("expected result to NOT contain %q", notExpected)
				}
			}
		})
	}
}

func TestGlobToRegex(t *testing.T) {
	tests := []struct {
		glob    string
		input   string
		matches bool
	}{
		{"*.js", "file.js", true},
		{"*.js", "file.ts", false},
		{"vendor/**", "vendor/lib.js", true},
		{"vendor/**", "vendor/sub/lib.js", true},
		{"vendor/**", "src/lib.js", false},
		{"test/*.go", "test/main.go", true},
		{"test/*.go", "test/sub/main.go", false},
		{"**/*.yaml", "config/deploy.yaml", true},
		{"**/*.yaml", "deploy.yaml", true},
		{"src/?.go", "src/a.go", true},
		{"src/?.go", "src/ab.go", false},
	}

	for _, tt := range tests {
		t.Run(tt.glob, func(t *testing.T) {
			pattern := globToRegex(tt.glob)
			regex, err := regexp.Compile(pattern)
			if err != nil {
				t.Fatalf("failed to compile pattern %q: %v", pattern, err)
			}

			matches := regex.MatchString(tt.input)
			if matches != tt.matches {
				t.Errorf("pattern %q, input %q: expected match=%v, got match=%v", tt.glob, tt.input, tt.matches, matches)
			}
		})
	}
}
