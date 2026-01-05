package read

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestRead(t *testing.T) {
	ctx := context.Background()
	tmpDir := t.TempDir()

	tests := []struct {
		name        string
		setup       func(t *testing.T, path string)
		with        interface{}
		expectError bool
		validate    func(t *testing.T, result interface{})
	}{
		{
			name: "read simple text file",
			setup: func(t *testing.T, path string) {
				err := os.WriteFile(path, []byte("Hello, World!"), 0644)
				if err != nil {
					t.Fatalf("failed to write test file: %v", err)
				}
			},
			with: map[string]interface{}{
				"path": filepath.Join(tmpDir, "test1.txt"),
			},
			expectError: false,
			validate: func(t *testing.T, result interface{}) {
				resultMap, ok := result.(map[string]interface{})
				if !ok {
					t.Fatal("expected result to be a map")
				}
				content, ok := resultMap["content"].(string)
				if !ok {
					t.Fatal("expected result to have 'content' field")
				}
				if content != "Hello, World!" {
					t.Errorf("expected 'Hello, World!', got %q", content)
				}
				bytes, ok := resultMap["bytes"].(int)
				if !ok {
					t.Fatal("expected result to have 'bytes' field")
				}
				if bytes != 13 {
					t.Errorf("expected 13 bytes, got %d", bytes)
				}
			},
		},
		{
			name: "read file from nested directory",
			setup: func(t *testing.T, path string) {
				dir := filepath.Dir(path)
				err := os.MkdirAll(dir, 0755)
				if err != nil {
					t.Fatalf("failed to create directory: %v", err)
				}
				err = os.WriteFile(path, []byte("nested content"), 0644)
				if err != nil {
					t.Fatalf("failed to write test file: %v", err)
				}
			},
			with: map[string]interface{}{
				"path": filepath.Join(tmpDir, "subdir", "test2.txt"),
			},
			expectError: false,
			validate: func(t *testing.T, result interface{}) {
				resultMap, ok := result.(map[string]interface{})
				if !ok {
					t.Fatal("expected result to be a map")
				}
				content, ok := resultMap["content"].(string)
				if !ok {
					t.Fatal("expected result to have 'content' field")
				}
				if content != "nested content" {
					t.Errorf("expected 'nested content', got %q", content)
				}
			},
		},
		{
			name: "read empty file",
			setup: func(t *testing.T, path string) {
				err := os.WriteFile(path, []byte(""), 0644)
				if err != nil {
					t.Fatalf("failed to write test file: %v", err)
				}
			},
			with: map[string]interface{}{
				"path": filepath.Join(tmpDir, "empty.txt"),
			},
			expectError: false,
			validate: func(t *testing.T, result interface{}) {
				resultMap, ok := result.(map[string]interface{})
				if !ok {
					t.Fatal("expected result to be a map")
				}
				content, ok := resultMap["content"].(string)
				if !ok {
					t.Fatal("expected result to have 'content' field")
				}
				if content != "" {
					t.Errorf("expected empty string, got %q", content)
				}
				bytes, ok := resultMap["bytes"].(int)
				if !ok {
					t.Fatal("expected result to have 'bytes' field")
				}
				if bytes != 0 {
					t.Errorf("expected 0 bytes, got %d", bytes)
				}
			},
		},
		{
			name: "read multiline file",
			setup: func(t *testing.T, path string) {
				content := "line 1\nline 2\nline 3\n"
				err := os.WriteFile(path, []byte(content), 0644)
				if err != nil {
					t.Fatalf("failed to write test file: %v", err)
				}
			},
			with: map[string]interface{}{
				"path": filepath.Join(tmpDir, "multiline.txt"),
			},
			expectError: false,
			validate: func(t *testing.T, result interface{}) {
				resultMap, ok := result.(map[string]interface{})
				if !ok {
					t.Fatal("expected result to be a map")
				}
				content, ok := resultMap["content"].(string)
				if !ok {
					t.Fatal("expected result to have 'content' field")
				}
				expected := "line 1\nline 2\nline 3\n"
				if content != expected {
					t.Errorf("expected %q, got %q", expected, content)
				}
			},
		},
		{
			name:        "missing path",
			setup:       func(t *testing.T, path string) {},
			with:        map[string]interface{}{},
			expectError: true,
		},
		{
			name:  "file does not exist",
			setup: func(t *testing.T, path string) {},
			with: map[string]interface{}{
				"path": filepath.Join(tmpDir, "nonexistent.txt"),
			},
			expectError: true,
		},
		{
			name: "invalid path type",
			setup: func(t *testing.T, path string) {
			},
			with: map[string]interface{}{
				"path": 12345,
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Get the path from 'with' if it's a map
			path := ""
			if withMap, ok := tt.with.(map[string]interface{}); ok {
				if p, ok := withMap["path"].(string); ok {
					path = p
				}
			}

			// Run setup if path is valid
			if path != "" {
				tt.setup(t, path)
			}

			result, err := Read(ctx, tt.with)
			if tt.expectError {
				if err == nil {
					t.Error("expected error but got none")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			// Validate result
			if tt.validate != nil {
				tt.validate(t, result)
			}
		})
	}
}
