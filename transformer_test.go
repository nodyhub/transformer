package transformer

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

// Test trySubstituteOutputDirect
func TestTrySubstituteOutputDirect(t *testing.T) {
	outputs := map[string]interface{}{
		"test_id": "output_value",
		"numeric": 42,
		"complex": []string{"a", "b", "c"},
	}

	tests := []struct {
		name     string
		input    string
		expected interface{}
		ok       bool
	}{
		{
			name:     "valid output substitution",
			input:    "${{ outputs.test_id }}",
			expected: "output_value",
			ok:       true,
		},
		{
			name:     "numeric output",
			input:    "${{ outputs.numeric }}",
			expected: 42,
			ok:       true,
		},
		{
			name:     "complex output",
			input:    "${{ outputs.complex }}",
			expected: []string{"a", "b", "c"},
			ok:       true,
		},
		{
			name:     "non-existent output",
			input:    "${{ outputs.missing }}",
			expected: nil,
			ok:       false,
		},
		{
			name:     "partial match",
			input:    "prefix ${{ outputs.test_id }} suffix",
			expected: nil,
			ok:       false,
		},
		{
			name:     "with extra spaces",
			input:    "${{  outputs.test_id  }}",
			expected: "output_value",
			ok:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, ok := trySubstituteOutputDirect(tt.input, outputs)
			if ok != tt.ok {
				t.Errorf("expected ok=%v, got ok=%v", tt.ok, ok)
			}
			if ok {
				if diff := cmp.Diff(tt.expected, result); diff != "" {
					t.Errorf("mismatch (-want +got):\n%s", diff)
				}
			}
		})
	}
}

// Test trySubstituteAllArgsDirect
func TestTrySubstituteAllArgsDirect(t *testing.T) {
	args := []string{"arg1", "arg2", "arg3"}

	tests := []struct {
		name     string
		input    string
		expected []string
		ok       bool
	}{
		{
			name:     "valid args substitution",
			input:    "${{ args }}",
			expected: args,
			ok:       true,
		},
		{
			name:     "with extra spaces",
			input:    "${{  args  }}",
			expected: args,
			ok:       true,
		},
		{
			name:     "partial match",
			input:    "prefix ${{ args }} suffix",
			expected: nil,
			ok:       false,
		},
		{
			name:     "indexed args",
			input:    "${{ args.0 }}",
			expected: nil,
			ok:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, ok := trySubstituteAllArgsDirect(tt.input, args)
			if ok != tt.ok {
				t.Errorf("expected ok=%v, got ok=%v", tt.ok, ok)
			}
			if ok {
				if diff := cmp.Diff(tt.expected, result); diff != "" {
					t.Errorf("mismatch (-want +got):\n%s", diff)
				}
			}
		})
	}
}

// Test trySubstituteEnvDirect
func TestTrySubstituteEnvDirect(t *testing.T) {
	// Set test environment variable
	t.Setenv("TEST_VAR", "test_value")

	tests := []struct {
		name     string
		input    string
		expected string
		ok       bool
	}{
		{
			name:     "valid env substitution",
			input:    "${{ env.TEST_VAR }}",
			expected: "test_value",
			ok:       true,
		},
		{
			name:     "with extra spaces",
			input:    "${{  env.TEST_VAR  }}",
			expected: "test_value",
			ok:       true,
		},
		{
			name:     "partial match",
			input:    "prefix ${{ env.TEST_VAR }} suffix",
			expected: "",
			ok:       false,
		},
		{
			name:     "non-existent env var",
			input:    "${{ env.NONEXISTENT }}",
			expected: "",
			ok:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, ok := trySubstituteEnvDirect(tt.input)
			if ok != tt.ok {
				t.Errorf("expected ok=%v, got ok=%v", tt.ok, ok)
			}
			if ok && result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

// Test substituteOutputsInline
func TestSubstituteOutputsInline(t *testing.T) {
	outputs := map[string]interface{}{
		"id1": "value1",
		"id2": 42,
		"id3": "value3",
	}

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "single substitution",
			input:    "Result: ${{ outputs.id1 }}",
			expected: "Result: value1",
		},
		{
			name:     "multiple substitutions",
			input:    "${{ outputs.id1 }} and ${{ outputs.id3 }}",
			expected: "value1 and value3",
		},
		{
			name:     "numeric value",
			input:    "Count: ${{ outputs.id2 }}",
			expected: "Count: 42",
		},
		{
			name:     "no substitution",
			input:    "No patterns here",
			expected: "No patterns here",
		},
		{
			name:     "non-existent output",
			input:    "${{ outputs.missing }}",
			expected: "${{ outputs.missing }}",
		},
		{
			name:     "with extra spaces",
			input:    "${{  outputs.id1  }}",
			expected: "value1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := substituteOutputsInline(tt.input, outputs)
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

// Test substituteEnvInline
func TestSubstituteEnvInline(t *testing.T) {
	t.Setenv("VAR1", "value1")
	t.Setenv("VAR2", "value2")

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "single env var",
			input:    "Value: ${{ env.VAR1 }}",
			expected: "Value: value1",
		},
		{
			name:     "multiple env vars",
			input:    "${{ env.VAR1 }} and ${{ env.VAR2 }}",
			expected: "value1 and value2",
		},
		{
			name:     "non-existent env var",
			input:    "${{ env.MISSING }}",
			expected: "",
		},
		{
			name:     "no substitution",
			input:    "No patterns here",
			expected: "No patterns here",
		},
		{
			name:     "with extra spaces",
			input:    "${{  env.VAR1  }}",
			expected: "value1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := substituteEnvInline(tt.input)
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

// Test substituteAllArgsInline
func TestSubstituteAllArgsInline(t *testing.T) {
	args := []string{"arg1", "arg2", "arg3"}

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "single substitution",
			input:    "Args: ${{ args }}",
			expected: "Args: [arg1 arg2 arg3]",
		},
		{
			name:     "multiple substitutions",
			input:    "${{ args }} and ${{ args }}",
			expected: "[arg1 arg2 arg3] and [arg1 arg2 arg3]",
		},
		{
			name:     "no substitution",
			input:    "No patterns here",
			expected: "No patterns here",
		},
		{
			name:     "with extra spaces",
			input:    "${{  args  }}",
			expected: "[arg1 arg2 arg3]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := substituteAllArgsInline(tt.input, args)
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

// Test substituteIndexedArgsInline
func TestSubstituteIndexedArgsInline(t *testing.T) {
	args := []string{"first", "second", "third"}

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "first arg",
			input:    "First: ${{ args.0 }}",
			expected: "First: first",
		},
		{
			name:     "middle arg",
			input:    "Second: ${{ args.1 }}",
			expected: "Second: second",
		},
		{
			name:     "last arg",
			input:    "Third: ${{ args.2 }}",
			expected: "Third: third",
		},
		{
			name:     "multiple args",
			input:    "${{ args.0 }}, ${{ args.1 }}, ${{ args.2 }}",
			expected: "first, second, third",
		},
		{
			name:     "out of bounds",
			input:    "${{ args.10 }}",
			expected: "${{ args.10 }}",
		},
		{
			name:     "no substitution",
			input:    "No patterns here",
			expected: "No patterns here",
		},
		{
			name:     "with extra spaces",
			input:    "${{  args.0  }}",
			expected: "first",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := substituteIndexedArgsInline(tt.input, args)
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

// Test warnUnsubstitutedPatterns
func TestWarnUnsubstitutedPatterns(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "valid pattern - no warning",
			input: "This is a normal string",
		},
		{
			name:  "invalid pattern - should warn",
			input: "${{ invalid.pattern }}",
		},
		{
			name:  "multiple invalid patterns",
			input: "${{ pattern1 }} and ${{ pattern2 }}",
		},
		{
			name:  "typo in pattern",
			input: "${{ output.typo }}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(_ *testing.T) {
			// This function only logs warnings, so we just ensure it doesn't panic
			warnUnsubstitutedPatterns(tt.input)
		})
	}
}

// Test substituteInputs with complex scenarios
func TestSubstituteInputs(t *testing.T) {
	outputs := map[string]interface{}{
		"result": "success",
		"count":  3,
	}
	args := []string{"arg0", "arg1"}
	t.Setenv("TEST_ENV", "env_value")

	tests := []struct {
		name     string
		input    interface{}
		expected interface{}
	}{
		{
			name:     "string with output",
			input:    "Status: ${{ outputs.result }}",
			expected: "Status: success",
		},
		{
			name:     "direct output substitution",
			input:    "${{ outputs.count }}",
			expected: 3,
		},
		{
			name:     "direct args substitution",
			input:    "${{ args }}",
			expected: args,
		},
		{
			name:     "indexed arg",
			input:    "First arg: ${{ args.0 }}",
			expected: "First arg: arg0",
		},
		{
			name:     "env var",
			input:    "Env: ${{ env.TEST_ENV }}",
			expected: "Env: env_value",
		},
		{
			name: "map with substitutions",
			input: map[string]interface{}{
				"key1": "${{ outputs.result }}",
				"key2": "${{ args.0 }}",
			},
			expected: map[string]interface{}{
				"key1": "success",
				"key2": "arg0",
			},
		},
		{
			name: "slice with substitutions",
			input: []interface{}{
				"${{ outputs.result }}",
				"${{ args.1 }}",
			},
			expected: []interface{}{
				"success",
				"arg1",
			},
		},
		{
			name:     "mixed substitutions",
			input:    "${{ outputs.result }} ${{ args.0 }} ${{ env.TEST_ENV }}",
			expected: "success arg0 env_value",
		},
		{
			name: "nested map",
			input: map[string]interface{}{
				"outer": map[string]interface{}{
					"inner": "${{ outputs.result }}",
				},
			},
			expected: map[string]interface{}{
				"outer": map[string]interface{}{
					"inner": "success",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := substituteInputs(tt.input, outputs, args)

			if diff := cmp.Diff(tt.expected, result); diff != "" {
				t.Errorf("mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

// Test parseInput
func TestParseInput(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected Steps
		hasError bool
	}{
		{
			name:  "string YAML input",
			input: "- uses: builtin/shell\n  with:\n    command: echo hello",
			expected: Steps{
				{Uses: "builtin/shell", With: map[string]interface{}{"command": "echo hello"}},
			},
			hasError: false,
		},
		{
			name:  "string JSON input",
			input: `[{"uses": "builtin/shell", "with": {"command": "echo hello"}}]`,
			expected: Steps{
				{Uses: "builtin/shell", With: map[string]interface{}{"command": "echo hello"}},
			},
			hasError: false,
		},
		{
			name:  "single Step input",
			input: Step{Uses: "builtin/shell", With: map[string]interface{}{"command": "echo hello"}},
			expected: Steps{
				{Uses: "builtin/shell", With: map[string]interface{}{"command": "echo hello"}},
			},
			hasError: false,
		},
		{
			name: "[]Step input",
			input: []Step{
				{Uses: "builtin/shell", With: map[string]interface{}{"command": "echo hello"}},
				{Uses: "builtin/shell", With: map[string]interface{}{"command": "echo world"}},
			},
			expected: Steps{
				{Uses: "builtin/shell", With: map[string]interface{}{"command": "echo hello"}},
				{Uses: "builtin/shell", With: map[string]interface{}{"command": "echo world"}},
			},
			hasError: false,
		},
		{
			name: "Steps input",
			input: Steps{
				{Uses: "builtin/shell", With: map[string]interface{}{"command": "echo hello"}},
			},
			expected: Steps{
				{Uses: "builtin/shell", With: map[string]interface{}{"command": "echo hello"}},
			},
			hasError: false,
		},
		{
			name:     "unsupported type",
			input:    123,
			expected: nil,
			hasError: true,
		},
		{
			name:     "invalid YAML string",
			input:    "invalid: yaml: content: [unclosed",
			expected: nil,
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parseInput(tt.input)
			if tt.hasError {
				if err == nil {
					t.Errorf("expected error, but got none")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if diff := cmp.Diff(tt.expected, result); diff != "" {
					t.Errorf("mismatch (-want +got):\n%s", diff)
				}
			}
		})
	}
}
