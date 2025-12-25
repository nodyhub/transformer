package file

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestWrite(t *testing.T) {
	ctx := context.Background()
	tmpDir := t.TempDir()

	tests := []struct {
		name        string
		with        interface{}
		expectError bool
		validate    func(t *testing.T, path string)
	}{
		{
			name: "write simple text",
			with: map[string]interface{}{
				"path":    filepath.Join(tmpDir, "test1.txt"),
				"content": "Hello, World!",
			},
			expectError: false,
			validate: func(t *testing.T, path string) {
				data, err := os.ReadFile(path)
				if err != nil {
					t.Fatalf("failed to read file: %v", err)
				}
				if string(data) != "Hello, World!" {
					t.Errorf("expected 'Hello, World!', got %q", string(data))
				}
			},
		},
		{
			name: "write to nested directory",
			with: map[string]interface{}{
				"path":    filepath.Join(tmpDir, "subdir", "test2.txt"),
				"content": "nested content",
			},
			expectError: false,
			validate: func(t *testing.T, path string) {
				data, err := os.ReadFile(path)
				if err != nil {
					t.Fatalf("failed to read file: %v", err)
				}
				if string(data) != "nested content" {
					t.Errorf("expected 'nested content', got %q", string(data))
				}
			},
		},
		{
			name: "append mode",
			with: map[string]interface{}{
				"path":    filepath.Join(tmpDir, "append.txt"),
				"content": "line1\n",
				"mode":    "write",
			},
			expectError: false,
			validate: func(t *testing.T, path string) {
				// Write again in append mode
				_, err := Write(ctx, map[string]interface{}{
					"path":    path,
					"content": "line2\n",
					"mode":    "append",
				})
				if err != nil {
					t.Fatalf("failed to append: %v", err)
				}
				data, err := os.ReadFile(path)
				if err != nil {
					t.Fatalf("failed to read file: %v", err)
				}
				expected := "line1\nline2\n"
				if string(data) != expected {
					t.Errorf("expected %q, got %q", expected, string(data))
				}
			},
		},
		{
			name: "missing path",
			with: map[string]interface{}{
				"content": "test",
			},
			expectError: true,
		},
		{
			name: "missing content",
			with: map[string]interface{}{
				"path": filepath.Join(tmpDir, "test.txt"),
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Write(ctx, tt.with)
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
			resultMap, ok := result.(map[string]interface{})
			if !ok {
				t.Fatal("expected result to be a map")
			}
			path, ok := resultMap["path"].(string)
			if !ok {
				t.Fatal("expected result to have 'path' field")
			}

			// Run custom validation
			if tt.validate != nil {
				tt.validate(t, path)
			}
		})
	}
}

func TestExtractPermissions(t *testing.T) {
	tests := []struct {
		name     string
		input    map[string]interface{}
		expected os.FileMode
	}{
		{
			name:     "default permissions when not provided",
			input:    map[string]interface{}{},
			expected: os.FileMode(0644),
		},
		{
			name: "int type positive value",
			input: map[string]interface{}{
				"perm": 420, // decimal for 0644
			},
			expected: os.FileMode(420),
		},
		{
			name: "int type with octal value",
			input: map[string]interface{}{
				"perm": 0755,
			},
			expected: os.FileMode(0755),
		},
		{
			name: "int type negative value uses default",
			input: map[string]interface{}{
				"perm": -1,
			},
			expected: os.FileMode(0644),
		},
		{
			name: "float64 type positive value",
			input: map[string]interface{}{
				"perm": 420.0,
			},
			expected: os.FileMode(420),
		},
		{
			name: "float64 type negative value uses default",
			input: map[string]interface{}{
				"perm": -1.0,
			},
			expected: os.FileMode(0644),
		},
		{
			name: "string octal notation 0644",
			input: map[string]interface{}{
				"perm": "0644",
			},
			expected: os.FileMode(0644),
		},
		{
			name: "string octal notation 0755",
			input: map[string]interface{}{
				"perm": "0755",
			},
			expected: os.FileMode(0755),
		},
		{
			name: "string octal notation 0600",
			input: map[string]interface{}{
				"perm": "0600",
			},
			expected: os.FileMode(0600),
		},
		{
			name: "string octal notation without leading zero",
			input: map[string]interface{}{
				"perm": "644",
			},
			expected: os.FileMode(0644),
		},
		{
			name: "string invalid format uses default",
			input: map[string]interface{}{
				"perm": "invalid",
			},
			expected: os.FileMode(0644),
		},
		{
			name: "string empty uses default",
			input: map[string]interface{}{
				"perm": "",
			},
			expected: os.FileMode(0644),
		},
		{
			name: "unsupported type uses default",
			input: map[string]interface{}{
				"perm": []int{644},
			},
			expected: os.FileMode(0644),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractPermissions(tt.input)
			if result != tt.expected {
				t.Errorf("expected %o (%d), got %o (%d)", tt.expected, tt.expected, result, result)
			}
		})
	}
}
