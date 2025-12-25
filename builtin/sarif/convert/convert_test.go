package convert

import (
	"context"
	"strings"
	"testing"

	"github.com/nodyhub/transformer/builtin/sarif/common"
)

func TestConvert_Trivy(t *testing.T) {
	ctx := context.Background()

	trivyOutput := `{
  "Results": [
    {
      "Target": "Dockerfile",
      "Vulnerabilities": [
        {
          "VulnerabilityID": "CVE-2023-1234",
          "PkgName": "openssl",
          "Severity": "HIGH",
          "Title": "Buffer overflow in OpenSSL"
        }
      ]
    }
  ]
}`

	result, err := Convert(ctx, map[string]interface{}{
		"tool":  "trivy",
		"input": trivyOutput,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	sarifMap, ok := result.(common.Report)
	if !ok {
		t.Fatal("expected result to be SarifReport")
	}

	if sarifMap.Version != "2.1.0" {
		t.Errorf("expected version 2.1.0, got %s", sarifMap.Version)
	}

	if len(sarifMap.Runs) != 1 {
		t.Fatalf("expected 1 run, got %d", len(sarifMap.Runs))
	}

	if sarifMap.Runs[0].Tool.Driver.Name != "Trivy" {
		t.Errorf("expected tool name 'Trivy', got %s", sarifMap.Runs[0].Tool.Driver.Name)
	}

	if len(sarifMap.Runs[0].Results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(sarifMap.Runs[0].Results))
	}

	result0 := sarifMap.Runs[0].Results[0]
	if result0.RuleID != "CVE-2023-1234" {
		t.Errorf("expected ruleID 'CVE-2023-1234', got %s", result0.RuleID)
	}
	if result0.Level != "error" {
		t.Errorf("expected level 'error', got %s", result0.Level)
	}
}

func TestConvert_Semgrep(t *testing.T) {
	ctx := context.Background()

	semgrepOutput := `{
  "results": [
    {
      "check_id": "javascript.express.security.audit.xss",
      "path": "server.js",
      "start": {
        "line": 42
      },
      "extra": {
        "message": "Potential XSS vulnerability",
        "severity": "WARNING"
      }
    }
  ]
}`

	result, err := Convert(ctx, map[string]interface{}{
		"tool":  "semgrep",
		"input": semgrepOutput,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	sarifMap, ok := result.(common.Report)
	if !ok {
		t.Fatal("expected result to be SarifReport")
	}

	if len(sarifMap.Runs[0].Results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(sarifMap.Runs[0].Results))
	}

	result0 := sarifMap.Runs[0].Results[0]
	if result0.RuleID != "javascript.express.security.audit.xss" {
		t.Errorf("expected ruleID 'javascript.express.security.audit.xss', got %s", result0.RuleID)
	}
	if result0.Locations[0].PhysicalLocation.Region.StartLine != 42 {
		t.Errorf("expected line 42, got %d", result0.Locations[0].PhysicalLocation.Region.StartLine)
	}
}

func TestConvert_Nuclei(t *testing.T) {
	ctx := context.Background()

	nucleiOutput := `{"template-id":"CVE-2023-5678","matched-at":"https://example.com","info":{"name":"SQL Injection","severity":"high"}}
{"template-id":"exposed-panel","matched-at":"https://example.com/admin","info":{"name":"Admin Panel Exposed","severity":"medium"}}`

	result, err := Convert(ctx, map[string]interface{}{
		"tool":  "nuclei",
		"input": nucleiOutput,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	sarifMap, ok := result.(common.Report)
	if !ok {
		t.Fatal("expected result to be SarifReport")
	}

	if len(sarifMap.Runs[0].Results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(sarifMap.Runs[0].Results))
	}

	result0 := sarifMap.Runs[0].Results[0]
	if result0.RuleID != "CVE-2023-5678" {
		t.Errorf("expected ruleID 'CVE-2023-5678', got %s", result0.RuleID)
	}
	if result0.Level != "error" {
		t.Errorf("expected level 'error' for high severity, got %s", result0.Level)
	}

	result1 := sarifMap.Runs[0].Results[1]
	if result1.Level != "warning" {
		t.Errorf("expected level 'warning' for medium severity, got %s", result1.Level)
	}
}

func TestConvert_AlreadySARIF(t *testing.T) {
	ctx := context.Background()

	sarifInput := `{
  "version": "2.1.0",
  "runs": [
    {
      "tool": {
        "driver": {
          "name": "TestTool"
        }
      },
      "results": []
    }
  ]
}`

	result, err := Convert(ctx, map[string]interface{}{
		"tool":  "trivy",
		"input": sarifInput,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	resultMap, ok := result.(map[string]interface{})
	if !ok {
		t.Fatal("expected result to be a map")
	}

	version, ok := resultMap["version"].(string)
	if !ok || version != "2.1.0" {
		t.Error("expected version 2.1.0 to be preserved")
	}
}

func TestConvert_CodeQL_SARIF(t *testing.T) {
	ctx := context.Background()

	// CodeQL already outputs SARIF by default
	codeqlSARIF := `{
  "version": "2.1.0",
  "$schema": "https://json.schemastore.org/sarif-2.1.0.json",
  "runs": [
    {
      "tool": {
        "driver": {
          "name": "CodeQL"
        }
      },
      "results": [
        {
          "ruleId": "js/sql-injection",
          "message": {
            "text": "SQL injection vulnerability"
          },
          "locations": [
            {
              "physicalLocation": {
                "artifactLocation": {
                  "uri": "app.js"
                },
                "region": {
                  "startLine": 25
                }
              }
            }
          ]
        }
      ]
    }
  ]
}`

	result, err := Convert(ctx, map[string]interface{}{
		"tool":  "codeql",
		"input": codeqlSARIF,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	resultMap, ok := result.(map[string]interface{})
	if !ok {
		t.Fatal("expected result to be a map")
	}

	version, ok := resultMap["version"].(string)
	if !ok || version != "2.1.0" {
		t.Error("expected version 2.1.0 to be preserved")
	}

	runs, ok := resultMap["runs"].([]interface{})
	if !ok || len(runs) == 0 {
		t.Fatal("expected runs array")
	}
}

func TestConvert_Trufflehog(t *testing.T) {
	ctx := context.Background()

	trufflehogOutput := `{"DetectorName":"AWS","DecoderName":"PLAIN","SourceMetadata":{"Data":{"Git":{"file":"config.yml","line":42}}},"Raw":"AKIAIOSFODNN7EXAMPLE"}
{"DetectorName":"GitHub","DecoderName":"PLAIN","SourceMetadata":{"Data":{"Filesystem":{"file":"secrets.txt","line":10}}},"Raw":"ghp_1234567890abcdef"}`

	result, err := Convert(ctx, map[string]interface{}{
		"tool":  "trufflehog",
		"input": trufflehogOutput,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	sarifMap, ok := result.(common.Report)
	if !ok {
		t.Fatal("expected result to be SarifReport")
	}

	if len(sarifMap.Runs[0].Results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(sarifMap.Runs[0].Results))
	}

	result0 := sarifMap.Runs[0].Results[0]
	if result0.RuleID != "AWS" {
		t.Errorf("expected ruleID 'AWS', got %s", result0.RuleID)
	}
	if result0.Level != "error" {
		t.Errorf("expected level 'error', got %s", result0.Level)
	}
	if result0.Locations[0].PhysicalLocation.Region.StartLine != 42 {
		t.Errorf("expected line 42, got %d", result0.Locations[0].PhysicalLocation.Region.StartLine)
	}
}

func TestConvert_Nmap(t *testing.T) {
	ctx := context.Background()

	nmapOutput := `{
  "nmaprun": {
    "host": [
      {
        "address": [
          {"addr": "192.168.1.1"}
        ],
        "ports": {
          "port": [
            {
              "portid": "22",
              "protocol": "tcp",
              "state": {"state": "open"},
              "service": {"name": "ssh", "product": "OpenSSH 8.0"}
            },
            {
              "portid": "80",
              "protocol": "tcp",
              "state": {"state": "open"},
              "service": {"name": "http"}
            }
          ]
        }
      }
    ]
  }
}`

	result, err := Convert(ctx, map[string]interface{}{
		"tool":  "nmap",
		"input": nmapOutput,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	sarifMap, ok := result.(common.Report)
	if !ok {
		t.Fatal("expected result to be SarifReport")
	}

	if len(sarifMap.Runs[0].Results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(sarifMap.Runs[0].Results))
	}

	result0 := sarifMap.Runs[0].Results[0]
	if result0.RuleID != "nmap-open-port-22" {
		t.Errorf("expected ruleID 'nmap-open-port-22', got %s", result0.RuleID)
	}
	if result0.Locations[0].PhysicalLocation.ArtifactLocation.URI != "192.168.1.1" {
		t.Errorf("expected URI '192.168.1.1', got %s", result0.Locations[0].PhysicalLocation.ArtifactLocation.URI)
	}
}

func TestConvert_Gosec(t *testing.T) {
	ctx := context.Background()

	gosecOutput := `{
  "Issues": [
    {
      "severity": "HIGH",
      "confidence": "HIGH",
      "rule_id": "G401",
      "details": "Use of weak cryptographic primitive",
      "file": "main.go",
      "line": "42"
    },
    {
      "severity": "MEDIUM",
      "confidence": "HIGH",
      "rule_id": "G404",
      "details": "Use of weak random number generator",
      "file": "utils.go",
      "line": "15"
    }
  ]
}`

	result, err := Convert(ctx, map[string]interface{}{
		"tool":  "gosec",
		"input": gosecOutput,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	sarifMap, ok := result.(common.Report)
	if !ok {
		t.Fatal("expected result to be Report")
	}

	if len(sarifMap.Runs[0].Results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(sarifMap.Runs[0].Results))
	}

	result0 := sarifMap.Runs[0].Results[0]
	if result0.RuleID != "G401" {
		t.Errorf("expected ruleID 'G401', got %s", result0.RuleID)
	}
	if result0.Level != "error" {
		t.Errorf("expected level 'error' for high severity, got %s", result0.Level)
	}
	if result0.Locations[0].PhysicalLocation.Region.StartLine != 42 {
		t.Errorf("expected line 42, got %d", result0.Locations[0].PhysicalLocation.Region.StartLine)
	}

	result1 := sarifMap.Runs[0].Results[1]
	if result1.Level != "warning" {
		t.Errorf("expected level 'warning' for medium severity, got %s", result1.Level)
	}
}

func TestConvert_OSV(t *testing.T) {
	ctx := context.Background()

	osvOutput := `{
  "results": [
    {
      "source": {
        "path": "package-lock.json",
        "type": "lockfile"
      },
      "packages": [
        {
          "package": {
            "name": "lodash",
            "version": "4.17.20",
            "ecosystem": "npm"
          },
          "vulnerabilities": [
            {
              "id": "GHSA-xxxx-yyyy-zzzz",
              "summary": "Prototype Pollution",
              "severity": [
                {
                  "type": "CVSS_V3",
                  "score": "7.5"
                }
              ]
            }
          ]
        },
        {
          "package": {
            "name": "axios",
            "version": "0.21.0",
            "ecosystem": "npm"
          },
          "vulnerabilities": [
            {
              "id": "GHSA-aaaa-bbbb-cccc",
              "summary": "Server-Side Request Forgery",
              "severity": [
                {
                  "type": "CVSS_V3",
                  "score": "9.1"
                }
              ]
            }
          ]
        }
      ]
    }
  ]
}`

	result, err := Convert(ctx, map[string]interface{}{
		"tool":  "osv",
		"input": osvOutput,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	sarifMap, ok := result.(common.Report)
	if !ok {
		t.Fatal("expected result to be Report")
	}

	if len(sarifMap.Runs[0].Results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(sarifMap.Runs[0].Results))
	}

	result0 := sarifMap.Runs[0].Results[0]
	if result0.RuleID != "GHSA-xxxx-yyyy-zzzz" {
		t.Errorf("expected ruleID 'GHSA-xxxx-yyyy-zzzz', got %s", result0.RuleID)
	}
	if result0.Level != "error" {
		t.Errorf("expected level 'error' for HIGH severity, got %s", result0.Level)
	}
	if !strings.Contains(result0.Message.Text, "lodash") {
		t.Errorf("expected message to contain 'lodash', got %s", result0.Message.Text)
	}

	result1 := sarifMap.Runs[0].Results[1]
	if result1.RuleID != "GHSA-aaaa-bbbb-cccc" {
		t.Errorf("expected ruleID 'GHSA-aaaa-bbbb-cccc', got %s", result1.RuleID)
	}
	if result1.Level != "error" {
		t.Errorf("expected level 'error' for CRITICAL severity, got %s", result1.Level)
	}
}

func TestConvert_MissingFields(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name        string
		with        interface{}
		expectedErr string
	}{
		{
			name:        "missing tool",
			with:        map[string]interface{}{"input": "test"},
			expectedErr: "'tool' field is required",
		},
		{
			name:        "missing input",
			with:        map[string]interface{}{"tool": "trivy"},
			expectedErr: "'input' field is required",
		},
		{
			name:        "unsupported tool",
			with:        map[string]interface{}{"tool": "unsupported", "input": "{}"},
			expectedErr: "unsupported tool",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Convert(ctx, tt.with)
			if err == nil {
				t.Fatal("expected error but got none")
			}
			if !strings.Contains(err.Error(), tt.expectedErr) {
				t.Errorf("expected error containing %q, got %q", tt.expectedErr, err.Error())
			}
		})
	}
}

func TestSeverityToLevel(t *testing.T) {
	tests := []struct {
		severity string
		expected string
	}{
		{"CRITICAL", "error"},
		{"HIGH", "error"},
		{"MEDIUM", "warning"},
		{"LOW", "note"},
		{"INFO", "note"},
		{"unknown", "warning"},
	}

	for _, tt := range tests {
		t.Run(tt.severity, func(t *testing.T) {
			result := severityToLevel(tt.severity)
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}
