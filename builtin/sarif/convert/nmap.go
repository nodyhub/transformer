package convert

import (
	"encoding/json"
	"fmt"

	"github.com/nodyhub/transformer/builtin/sarif/common"
)

// convertNmap converts Nmap XML or JSON output to SARIF
func convertNmap(input string) (interface{}, error) {
	var nmapOutput map[string]interface{}
	if err := json.Unmarshal([]byte(input), &nmapOutput); err != nil {
		return nil, fmt.Errorf("failed to parse Nmap output (XML not yet supported): %w", err)
	}

	results := parseNmapRun(nmapOutput)
	return createSarifReport("Nmap", results), nil
}

func parseNmapRun(nmapOutput map[string]interface{}) []common.Result {
	results := []common.Result{}

	nmaprun, ok := nmapOutput["nmaprun"].(map[string]interface{})
	if !ok {
		return results
	}

	hosts, ok := nmaprun["host"].([]interface{})
	if !ok {
		return results
	}

	for _, hostRaw := range hosts {
		hostMap, ok := hostRaw.(map[string]interface{})
		if !ok {
			continue
		}

		hostAddr := extractHostAddress(hostMap)
		hostResults := parseHostPorts(hostMap, hostAddr)
		results = append(results, hostResults...)
	}

	return results
}

func extractHostAddress(hostMap map[string]interface{}) string {
	addresses, ok := hostMap["address"].([]interface{})
	if !ok || len(addresses) == 0 {
		return "unknown"
	}

	addr, ok := addresses[0].(map[string]interface{})
	if !ok {
		return "unknown"
	}

	addrStr, ok := addr["addr"].(string)
	if !ok {
		return "unknown"
	}

	return addrStr
}

func parseHostPorts(hostMap map[string]interface{}, hostAddr string) []common.Result {
	results := []common.Result{}

	ports, ok := hostMap["ports"].(map[string]interface{})
	if !ok {
		return results
	}

	portList, ok := ports["port"].([]interface{})
	if !ok {
		return results
	}

	for _, portRaw := range portList {
		portMap, ok := portRaw.(map[string]interface{})
		if !ok {
			continue
		}

		if result, ok := createNmapResult(portMap, hostAddr); ok {
			results = append(results, result)
		}
	}

	return results
}

func createNmapResult(portMap map[string]interface{}, hostAddr string) (common.Result, bool) {
	state := extractPortState(portMap)
	if state != "open" {
		return common.Result{}, false
	}

	portID := extractPortID(portMap)
	protocol := getStringField(portMap, "protocol")
	serviceName, serviceProduct := extractServiceInfo(portMap)

	message := buildNmapMessage(portID, protocol, serviceName, serviceProduct)
	ruleID := fmt.Sprintf("nmap-open-port-%s", portID)

	return common.Result{
		RuleID: ruleID,
		Level:  "note",
		Message: common.Message{
			Text: message,
		},
		Locations: []common.Location{
			{
				PhysicalLocation: common.PhysicalLocation{
					ArtifactLocation: common.ArtifactLocation{
						URI: hostAddr,
					},
					Region: common.Region{
						StartLine: 1,
					},
				},
			},
		},
		Properties: map[string]interface{}{
			"port":     portID,
			"protocol": protocol,
			"service":  serviceName,
			"product":  serviceProduct,
		},
	}, true
}

func extractPortID(portMap map[string]interface{}) string {
	if pid, ok := portMap["portid"].(string); ok {
		return pid
	}
	if pid, ok := portMap["portid"].(float64); ok {
		return fmt.Sprintf("%.0f", pid)
	}
	return ""
}

func extractPortState(portMap map[string]interface{}) string {
	stateMap, ok := portMap["state"].(map[string]interface{})
	if !ok {
		return ""
	}
	return getStringField(stateMap, "state")
}

func extractServiceInfo(portMap map[string]interface{}) (string, string) {
	service, ok := portMap["service"].(map[string]interface{})
	if !ok {
		return "", ""
	}
	serviceName := getStringField(service, "name")
	serviceProduct := getStringField(service, "product")
	return serviceName, serviceProduct
}

func buildNmapMessage(portID, protocol, serviceName, serviceProduct string) string {
	message := fmt.Sprintf("Open port %s/%s", portID, protocol)
	if serviceName != "" {
		message = fmt.Sprintf("%s - %s", message, serviceName)
	}
	if serviceProduct != "" {
		message = fmt.Sprintf("%s (%s)", message, serviceProduct)
	}
	return message
}
