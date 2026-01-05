package read

import (
	"context"
	"fmt"
	"os"

	"github.com/nodyhub/transformer/registry"
)

func init() {
	registry.Register("builtin/file/read", Read)
}

// Read reads content from a file
func Read(ctx context.Context, with interface{}) (interface{}, error) {
	withMap, ok := with.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("expected 'with' to be a map with 'path' field")
	}

	path, err := extractPath(withMap)
	if err != nil {
		return nil, err
	}

	content, bytes, err := readContent(path)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"path":    path,
		"content": content,
		"bytes":   bytes,
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

func readContent(path string) (string, int, error) {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return "", 0, fmt.Errorf("file does not exist: %s", path)
		}
		return "", 0, fmt.Errorf("failed to stat file: %w", err)
	}

	content, err := os.ReadFile(path)
	if err != nil {
		return "", 0, fmt.Errorf("failed to read file: %w", err)
	}

	return string(content), len(content), nil
}
