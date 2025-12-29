package nmap

import (
	"context"
	"fmt"
	"strings"

	"github.com/nodyhub/transformer/registry"
)

func init() {
	registry.Register("builtin/tools/nmap", Nmap)
}

// Nmap executes Nmap network scanner
// Supported inputs:
//   - target: target host(s), IP address(es), or network(s) (required)
//   - ports: port specification (e.g., "22,80,443" or "1-1000")
//   - scan-type: scan type (-sS, -sT, -sU, -sV, -sC, -A)
//   - output: output file path (base name, extensions added automatically)
//   - format: output format (normal, xml, grepable, all)
//   - timing: timing template (0-5, or paranoid/sneaky/polite/normal/aggressive/insane)
//   - script: NSE script(s) to run
//   - os-detection: enable OS detection (-O)
//   - version-detection: enable version detection (-sV)
//   - aggressive: enable aggressive scan (-A)
//   - additional-args: array of additional arguments
func Nmap(ctx context.Context, with interface{}) (interface{}, error) {
	withMap, ok := with.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("expected 'with' to be a map")
	}

	args, err := buildNmapArgs(withMap)
	if err != nil {
		return nil, err
	}

	cmdStr := strings.Join(args, " ")
	shellInput := map[string]interface{}{
		"command": cmdStr,
	}
	if v, ok := withMap["stdout"]; ok {
		shellInput["stdout"] = v
	}
	if v, ok := withMap["stderr"]; ok {
		shellInput["stderr"] = v
	}
	return registry.Call(ctx, "builtin/shell", shellInput)
}

func buildNmapArgs(m map[string]interface{}) ([]string, error) {
	target, ok := m["target"].(string)
	if !ok || target == "" {
		return nil, fmt.Errorf("'target' field is required")
	}

	args := []string{"nmap"}

	// Scan type
	if scanType, ok := m["scan-type"].(string); ok && scanType != "" {
		args = append(args, scanType)
	}

	// Ports
	if ports, ok := m["ports"].(string); ok && ports != "" {
		args = append(args, "-p", ports)
	}

	// Timing
	if timing, ok := m["timing"].(string); ok && timing != "" {
		args = append(args, "-T", timing)
	} else if timing, ok := m["timing"].(int); ok {
		args = append(args, "-T", fmt.Sprintf("%d", timing))
	}

	// OS detection
	if osDetection, ok := m["os-detection"].(bool); ok && osDetection {
		args = append(args, "-O")
	}

	// Version detection
	if versionDetection, ok := m["version-detection"].(bool); ok && versionDetection {
		args = append(args, "-sV")
	}

	// Aggressive scan
	if aggressive, ok := m["aggressive"].(bool); ok && aggressive {
		args = append(args, "-A")
	}

	// NSE scripts
	if script, ok := m["script"].(string); ok && script != "" {
		args = append(args, "--script", script)
	}

	// Output format and file
	if output, ok := m["output"].(string); ok && output != "" {
		format := "normal"
		if f, ok := m["format"].(string); ok && f != "" {
			format = f
		}
		switch format {
		case "normal":
			args = append(args, "-oN", output)
		case "xml":
			args = append(args, "-oX", output)
		case "grepable":
			args = append(args, "-oG", output)
		case "all":
			args = append(args, "-oA", output)
		}
	}

	// Additional arguments
	if additionalArgs, ok := m["additional-args"].([]interface{}); ok {
		for _, arg := range additionalArgs {
			if argStr, ok := arg.(string); ok {
				args = append(args, argStr)
			}
		}
	}

	// Add target
	args = append(args, target)
	return args, nil
}
