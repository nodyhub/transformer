package suppress

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/nodyhub/transformer/builtin/sarif/common"
	"github.com/nodyhub/transformer/registry"
)

func init() {
	registry.Register("builtin/sarif/suppress", Suppress)
}

// Suppress filters out results from a SARIF report based on suppression rules
func Suppress(ctx context.Context, with interface{}) (interface{}, error) {
	withMap, ok := with.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("expected 'with' to be a map")
	}

	// Get input report
	inputRaw, ok := withMap["input"]
	if !ok {
		return nil, fmt.Errorf("'input' field is required")
	}

	// Parse input as JSON string or object
	var report common.Report
	switch v := inputRaw.(type) {
	case string:
		if err := json.Unmarshal([]byte(v), &report); err != nil {
			return nil, fmt.Errorf("failed to parse SARIF input: %w", err)
		}
	case map[string]interface{}:
		data, err := json.Marshal(v)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal input: %w", err)
		}
		if err := json.Unmarshal(data, &report); err != nil {
			return nil, fmt.Errorf("failed to parse SARIF input: %w", err)
		}
	default:
		return nil, fmt.Errorf("'input' must be a string or object")
	}

	// Extract suppression rules
	rules := extractSuppressionRules(withMap)

	// Apply suppressions
	for i := range report.Runs {
		report.Runs[i].Results = filterResults(report.Runs[i].Results, rules)
	}

	// Convert back to JSON
	resultJSON, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal suppressed report: %w", err)
	}

	return string(resultJSON), nil
}

// SuppressionRules contains all suppression criteria
type SuppressionRules struct {
	RuleIDs      []string         // Exact rule IDs to suppress (CVE-2024-1234, etc.)
	RulePatterns []*regexp.Regexp // Regex patterns for rule IDs
	Paths        []string         // Exact file paths to suppress
	PathPatterns []*regexp.Regexp // Glob/regex patterns for file paths
	Severities   []string         // Severity levels to suppress (error, warning, note)
	Tools        []string         // Tool names to suppress
}

func extractSuppressionRules(withMap map[string]interface{}) SuppressionRules {
	rules := SuppressionRules{}

	rules.RuleIDs = extractStringArray(withMap, "rule_ids")
	rules.RulePatterns = extractRegexPatterns(withMap, "rule_patterns")
	rules.Paths = extractStringArray(withMap, "paths")
	rules.PathPatterns = extractGlobPatterns(withMap, "path_patterns")
	rules.Severities = extractLowerStringArray(withMap, "severities")
	rules.Tools = extractLowerStringArray(withMap, "tools")

	return rules
}

func extractStringArray(m map[string]interface{}, key string) []string {
	result := []string{}
	if raw, ok := m[key].([]interface{}); ok {
		for _, item := range raw {
			if str, ok := item.(string); ok {
				result = append(result, str)
			}
		}
	}
	return result
}

func extractLowerStringArray(m map[string]interface{}, key string) []string {
	result := []string{}
	if raw, ok := m[key].([]interface{}); ok {
		for _, item := range raw {
			if str, ok := item.(string); ok {
				result = append(result, strings.ToLower(str))
			}
		}
	}
	return result
}

func extractRegexPatterns(m map[string]interface{}, key string) []*regexp.Regexp {
	result := []*regexp.Regexp{}
	if raw, ok := m[key].([]interface{}); ok {
		for _, item := range raw {
			if pattern, ok := item.(string); ok {
				if regex, err := regexp.Compile(pattern); err == nil {
					result = append(result, regex)
				}
			}
		}
	}
	return result
}

func extractGlobPatterns(m map[string]interface{}, key string) []*regexp.Regexp {
	result := []*regexp.Regexp{}
	if raw, ok := m[key].([]interface{}); ok {
		for _, item := range raw {
			if pattern, ok := item.(string); ok {
				regexPattern := globToRegex(pattern)
				if regex, err := regexp.Compile(regexPattern); err == nil {
					result = append(result, regex)
				}
			}
		}
	}
	return result
}

func filterResults(results []common.Result, rules SuppressionRules) []common.Result {
	filtered := []common.Result{}

	for _, result := range results {
		if shouldSuppress(result, rules) {
			continue
		}
		filtered = append(filtered, result)
	}

	return filtered
}

func shouldSuppress(result common.Result, rules SuppressionRules) bool {
	// Check rule ID exact match
	for _, ruleID := range rules.RuleIDs {
		if result.RuleID == ruleID {
			return true
		}
	}

	// Check rule ID pattern match
	for _, pattern := range rules.RulePatterns {
		if pattern.MatchString(result.RuleID) {
			return true
		}
	}

	// Check severity match
	for _, severity := range rules.Severities {
		if strings.ToLower(result.Level) == severity {
			return true
		}
	}

	// Check path and path patterns
	if len(result.Locations) > 0 {
		uri := result.Locations[0].PhysicalLocation.ArtifactLocation.URI

		// Check exact path match
		for _, path := range rules.Paths {
			if uri == path {
				return true
			}
		}

		// Check path pattern match
		for _, pattern := range rules.PathPatterns {
			if pattern.MatchString(uri) {
				return true
			}
		}
	}

	return false
}

// globToRegex converts a glob pattern to a regex pattern
func globToRegex(glob string) string {
	var result strings.Builder
	result.WriteString("^")

	i := 0
	for i < len(glob) {
		switch glob[i] {
		case '*':
			// Check for **
			if i+1 < len(glob) && glob[i+1] == '*' {
				// Check if ** is followed by /
				if i+2 < len(glob) && glob[i+2] == '/' {
					// **/ means zero or more path segments
					result.WriteString("(?:.*/)?")
					i += 3
				} else {
					// ** at the end or not followed by /
					result.WriteString(".*")
					i += 2
				}
			} else {
				result.WriteString("[^/]*")
				i++
			}
		case '?':
			result.WriteString(".")
			i++
		case '.', '+', '(', ')', '|', '[', ']', '{', '}', '^', '$', '\\':
			// Escape special regex characters
			result.WriteString("\\")
			result.WriteByte(glob[i])
			i++
		default:
			result.WriteByte(glob[i])
			i++
		}
	}

	result.WriteString("$")
	return result.String()
}
