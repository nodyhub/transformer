package convert

import (
	"encoding/json"
	"fmt"

	"github.com/nodyhub/transformer/builtin/sarif/common"
)

// convertOSV converts OSV Scanner JSON output to SARIF
func convertOSV(input interface{}) (interface{}, error) {
	osvResult, ok := input.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid OSV output format: expected a JSON object, got %T", input)
	}

	osvOutputString, ok := osvResult["output"].(string)
	if !ok {
		return nil, fmt.Errorf("failed to extract OSV output, missing 'output' field or invalid format: %T", osvResult["output"])
	}

	var osvOutput map[string]interface{}
	if err := json.Unmarshal([]byte(osvOutputString), &osvOutput); err != nil {
		return nil, fmt.Errorf("failed to parse OSV output JSON: %w", err)
	}

	results := []common.Result{}

	// OSV output has "results" array
	if resultsRaw, ok := osvOutput["results"].([]interface{}); ok {
		for _, resultRaw := range resultsRaw {
			resultMap, ok := resultRaw.(map[string]interface{})
			if !ok {
				continue
			}

			sourcePath := extractOSVSourcePath(resultMap)
			packageResults := processOSVPackages(resultMap, sourcePath)
			results = append(results, packageResults...)
		}
	}

	return createSarifReport("OSV-Scanner", results), nil
}

func extractOSVSourcePath(resultMap map[string]interface{}) string {
	if source, ok := resultMap["source"].(map[string]interface{}); ok {
		return getStringField(source, "path")
	}
	return ""
}

func processOSVPackages(resultMap map[string]interface{}, sourcePath string) []common.Result {
	results := []common.Result{}
	packagesRaw, ok := resultMap["packages"].([]interface{})
	if !ok {
		return results
	}

	for _, pkgRaw := range packagesRaw {
		pkgMap, ok := pkgRaw.(map[string]interface{})
		if !ok {
			continue
		}

		pkgInfo := extractOSVPackageInfo(pkgMap)
		vulnResults := processOSVVulnerabilities(pkgMap, pkgInfo, sourcePath)
		results = append(results, vulnResults...)
	}

	return results
}

type osvPackageInfo struct {
	name      string
	version   string
	ecosystem string
}

func extractOSVPackageInfo(pkgMap map[string]interface{}) osvPackageInfo {
	info := osvPackageInfo{}
	if pkg, ok := pkgMap["package"].(map[string]interface{}); ok {
		info.name = getStringField(pkg, "name")
		info.version = getStringField(pkg, "version")
		info.ecosystem = getStringField(pkg, "ecosystem")
	}
	return info
}

func processOSVVulnerabilities(pkgMap map[string]interface{}, pkgInfo osvPackageInfo, sourcePath string) []common.Result {
	results := []common.Result{}
	vulnsRaw, ok := pkgMap["vulnerabilities"].([]interface{})
	if !ok {
		return results
	}

	for _, vulnRaw := range vulnsRaw {
		vulnMap, ok := vulnRaw.(map[string]interface{})
		if !ok {
			continue
		}

		vulnID := getStringField(vulnMap, "id")
		summary := getStringField(vulnMap, "summary")
		severity := determineOSVSeverity(vulnMap)

		message := fmt.Sprintf("%s: %s in %s@%s", vulnID, summary, pkgInfo.name, pkgInfo.version)
		if summary == "" {
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
							URI: sourcePath,
						},
						Region: common.Region{
							StartLine: 1,
						},
					},
				},
			},
			Properties: map[string]interface{}{
				"package":   pkgInfo.name,
				"version":   pkgInfo.version,
				"ecosystem": pkgInfo.ecosystem,
				"severity":  severity,
			},
		})
	}

	return results
}

func determineOSVSeverity(vulnMap map[string]interface{}) string {
	severityRaw, ok := vulnMap["severity"].([]interface{})
	if !ok || len(severityRaw) == 0 {
		return "MEDIUM"
	}

	sevMap, ok := severityRaw[0].(map[string]interface{})
	if !ok {
		return "MEDIUM"
	}

	score, ok := sevMap["score"].(string)
	if !ok {
		return "MEDIUM"
	}

	var cvssScore float64
	if _, err := fmt.Sscanf(score, "%f", &cvssScore); err != nil {
		return "MEDIUM"
	}

	switch {
	case cvssScore >= 9.0:
		return "CRITICAL"
	case cvssScore >= 7.0:
		return "HIGH"
	case cvssScore >= 4.0:
		return "MEDIUM"
	default:
		return "LOW"
	}
}
