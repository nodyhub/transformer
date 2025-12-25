package file

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/nodyhub/transformer/registry"
)

func init() {
	registry.Register("builtin/file/write", Write)
}

// Write writes content to a file
func Write(ctx context.Context, with interface{}) (interface{}, error) {
	withMap, ok := with.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("expected 'with' to be a map")
	}

	path, err := extractPath(withMap)
	if err != nil {
		return nil, err
	}

	content, err := extractContent(withMap)
	if err != nil {
		return nil, err
	}

	mode := extractMode(withMap)
	perm := extractPermissions(withMap)

	if err := ensureDirectoryExists(path); err != nil {
		return nil, err
	}

	if err := writeContent(path, content, mode, perm); err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"path":  path,
		"bytes": len(content),
	}, nil
}

func extractPath(withMap map[string]interface{}) (string, error) {
	pathRaw, ok := withMap["path"]
	if !ok {
		return "", fmt.Errorf("'path' field is required")
	}
	path, ok := pathRaw.(string)
	if !ok {
		return "", fmt.Errorf("'path' must be a string")
	}
	return path, nil
}

func extractContent(withMap map[string]interface{}) (string, error) {
	contentRaw, ok := withMap["content"]
	if !ok {
		return "", fmt.Errorf("'content' field is required")
	}

	switch v := contentRaw.(type) {
	case string:
		return v, nil
	case []byte:
		return string(v), nil
	default:
		return fmt.Sprintf("%v", v), nil
	}
}

func extractMode(withMap map[string]interface{}) string {
	mode := "write"
	if modeRaw, ok := withMap["mode"]; ok {
		if modeStr, ok := modeRaw.(string); ok {
			mode = modeStr
		}
	}
	return mode
}

func extractPermissions(withMap map[string]interface{}) os.FileMode {
	perm := os.FileMode(0644)
	if permRaw, ok := withMap["perm"]; ok {
		switch v := permRaw.(type) {
		case int:
			if v >= 0 {
				perm = os.FileMode(uint32(v)) //nolint:gosec
			}
		case float64:
			if v >= 0 {
				perm = os.FileMode(uint32(v)) //nolint:gosec
			}
		case string:
			// Support octal string notation like "0644"
			var parsed uint32
			if _, err := fmt.Sscanf(v, "%o", &parsed); err == nil {
				perm = os.FileMode(parsed)
			}
		}
	}
	return perm
}

func ensureDirectoryExists(path string) error {
	dir := filepath.Dir(path)
	if dir != "." && dir != "" {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}
	}
	return nil
}

func writeContent(path, content, mode string, perm os.FileMode) error {
	switch mode {
	case "write", "truncate":
		return writeFile(path, content, perm)
	case "append":
		return appendFile(path, content, perm)
	default:
		return fmt.Errorf("invalid mode: %s (must be 'write' or 'append')", mode)
	}
}

func writeFile(path, content string, perm os.FileMode) error {
	if err := os.WriteFile(path, []byte(content), perm); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}
	return nil
}

func appendFile(path, content string, perm os.FileMode) error {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, perm)
	if err != nil {
		return fmt.Errorf("failed to open file for append: %w", err)
	}
	defer f.Close()

	if _, err := f.WriteString(content); err != nil {
		return fmt.Errorf("failed to append to file: %w", err)
	}
	return nil
}
