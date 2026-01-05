package transformer

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"regexp"
	"strconv"
	"strings"

	// Import built-in transformers to register them
	_ "github.com/nodyhub/transformer/builtin"

	"github.com/nodyhub/transformer/registry"
	"gopkg.in/yaml.v3"

	"github.com/nodyhub/goperf"
)

var (
	outputsMatcher       = regexp.MustCompile(`\$\{\{\s*outputs\.([a-zA-Z0-9_-]+(?:\.[a-zA-Z0-9_-]+)*)\s*\}\}`)
	allArgsMatcher       = regexp.MustCompile(`\$\{\{\s*args\s*\}\}`)
	argsMatcher          = regexp.MustCompile(`\$\{\{\s*args\.([0-9]+)\s*\}\}`)
	envMatcher           = regexp.MustCompile(`\$\{\{\s*env\.([a-zA-Z0-9_]+)\s*\}\}`)
	unsubstitutedMatcher = regexp.MustCompile(`\$\{\{[^}]+\}\}`)
)

// Run executes the steps defined in the input.
// Input can be:
// - string: YAML/JSON definition of steps
// - step: single step struct
// - []step: slice of step structs
func Run(ctx context.Context, input interface{}, args []string) error {
	perf := goperf.Start()

	s, err := parseInput(input)
	if err != nil {
		return err
	}

	outputs := make(map[string]interface{})

	for _, step := range s {
		fn, ok := registry.Get(step.Uses)
		if ok {
			out, err := fn(ctx, substituteInputs(step.With, outputs, args))
			if err != nil {
				return fmt.Errorf("error executing step \"%s\": %w", step.Name, err)
			}
			slog.Debug("step executed", slog.String("name", step.Name), slog.Any("output", out))

			if step.ID != "" {
				outputs[step.ID] = out
				slog.Debug("output stored", slog.String("id", step.ID), slog.Any("output", out))
			}
		} else {
			return fmt.Errorf("unknown step function (id: %s): %s", step.ID, step.Uses)
		}

		perf.Mark("step " + step.Name)

	}

	slog.Debug("all steps completed", slog.Duration("duration", perf.GetMarks()["from_start"]))

	return nil
}

type Steps []Step

type Step struct {
	ID       string      `json:"id,omitempty" yaml:"id,omitempty"`
	LogLevel string      `json:"logLevel,omitempty" yaml:"log_level,omitempty"`
	Name     string      `json:"name,omitempty" yaml:"name,omitempty"`
	Uses     string      `json:"uses" yaml:"uses"`
	With     interface{} `json:"with,omitempty" yaml:"with,omitempty"`
}

func parseInput(input interface{}) (Steps, error) {
	switch v := input.(type) {
	case string:
		// YAML/JSON string definition
		return parseSteps(v)
	case Step:
		// Single step
		return Steps{v}, nil
	case []Step:
		// Slice of steps
		return v, nil
	case Steps:
		// Already parsed steps
		return v, nil
	default:
		return nil, fmt.Errorf("unsupported input type: %T (expected string, step, or []step)", input)
	}
}

func parseSteps(data string) (Steps, error) {
	var s Steps

	// Try JSON first if data looks like JSON
	trimmed := strings.TrimSpace(data)
	if strings.HasPrefix(trimmed, "[") || strings.HasPrefix(trimmed, "{") {
		if err := json.Unmarshal([]byte(data), &s); err == nil {
			return s, nil
		}
	}

	// Fall back to YAML (which also handles JSON as JSON is valid YAML)
	if err := yaml.Unmarshal([]byte(data), &s); err != nil {
		return nil, err
	}

	return s, nil
}

func substituteInputs(input interface{}, outputs map[string]interface{}, args []string) interface{} {
	switch v := input.(type) {
	case map[string]interface{}:
		for key, val := range v {
			v[key] = substituteInputs(val, outputs, args)
		}
		return v

	case []interface{}:
		for i, val := range v {
			v[i] = substituteInputs(val, outputs, args)
		}
		return v

	case string:
		// Check if entire string is a single substitution pattern (return directly)
		if result, ok := trySubstituteOutputDirect(v, outputs); ok {
			return result
		}
		if result, ok := trySubstituteAllArgsDirect(v, args); ok {
			return result
		}
		if result, ok := trySubstituteEnvDirect(v); ok {
			return result
		}

		// Perform inline substitutions
		result := substituteOutputsInline(v, outputs)
		result = substituteEnvInline(result)
		result = substituteAllArgsInline(result, args)
		result = substituteIndexedArgsInline(result, args)
		warnUnsubstitutedPatterns(result)

		return result
	}

	return input
}

func trySubstituteOutputDirect(s string, outputs map[string]interface{}) (interface{}, bool) {
	matches := outputsMatcher.FindStringSubmatch(s)
	if len(matches) == 2 && matches[0] == s {
		path := matches[1]
		value, found := getNestedOutputValue(path, outputs)
		if found {
			slog.Debug("substituting output", slog.String("path", path), slog.Any("value", value))
			return value, true
		}
	}
	return nil, false
}

func trySubstituteAllArgsDirect(s string, args []string) (interface{}, bool) {
	if allArgsMatcher.MatchString(s) && allArgsMatcher.FindString(s) == s {
		slog.Debug("substituting all args", slog.Any("value", args))
		return args, true
	}
	return nil, false
}

func trySubstituteEnvDirect(s string) (interface{}, bool) {
	matches := envMatcher.FindStringSubmatch(s)
	if len(matches) == 2 && matches[0] == s {
		envVar := matches[1]
		envValue := os.Getenv(envVar)
		slog.Debug("substituting env var", slog.String("var", envVar), slog.String("value", envValue))
		return envValue, true
	}
	return nil, false
}

func substituteOutputsInline(s string, outputs map[string]interface{}) string {
	return outputsMatcher.ReplaceAllStringFunc(s, func(m string) string {
		matches := outputsMatcher.FindStringSubmatch(m)
		if len(matches) == 2 {
			path := matches[1]
			value, found := getNestedOutputValue(path, outputs)
			if found {
				slog.Debug("substituting output inline", slog.String("path", path), slog.Any("value", value))
				return fmt.Sprintf("%v", value)
			}
		}
		return m
	})
}

func substituteEnvInline(s string) string {
	return envMatcher.ReplaceAllStringFunc(s, func(m string) string {
		matches := envMatcher.FindStringSubmatch(m)
		if len(matches) == 2 {
			envVar := matches[1]
			envValue := os.Getenv(envVar)
			if envValue == "" {
				slog.Warn("environment variable not set", slog.String("var", envVar))
			}
			slog.Debug("substituting env var", slog.String("var", envVar), slog.String("value", envValue))
			return envValue
		}
		return m
	})
}

func substituteAllArgsInline(s string, args []string) string {
	return allArgsMatcher.ReplaceAllStringFunc(s, func(_ string) string {
		slog.Debug("substituting all args", slog.Any("value", args))
		return fmt.Sprintf("%v", args)
	})
}

func substituteIndexedArgsInline(s string, args []string) string {
	return argsMatcher.ReplaceAllStringFunc(s, func(m string) string {
		matches := argsMatcher.FindStringSubmatch(m)
		if len(matches) == 2 {
			n, err := strconv.Atoi(matches[1])
			if err == nil && n >= 0 && n < len(args) {
				slog.Debug("substituting arg", slog.Int("index", n), slog.String("value", args[n]))
				return fmt.Sprintf("%v", args[n])
			}
		}
		return m
	})
}

// getNestedOutputValue traverses a path like "step_id.property.nested" in the outputs map
func getNestedOutputValue(path string, outputs map[string]interface{}) (interface{}, bool) {
	parts := strings.Split(path, ".")
	if len(parts) == 0 {
		return nil, false
	}

	// Get the step ID (first part)
	stepID := parts[0]
	value, ok := outputs[stepID]
	if !ok {
		return nil, false
	}

	// If there are no nested properties, return the value as-is
	if len(parts) == 1 {
		return value, true
	}

	// Traverse nested properties
	for _, key := range parts[1:] {
		switch v := value.(type) {
		case map[string]interface{}:
			if nested, found := v[key]; found {
				value = nested
			} else {
				return nil, false
			}
		default:
			// Cannot traverse further - value is not a map
			return nil, false
		}
	}

	return value, true
}

func warnUnsubstitutedPatterns(s string) {
	if unsubstituted := unsubstitutedMatcher.FindAllString(s, -1); len(unsubstituted) > 0 {
		for _, pattern := range unsubstituted {
			slog.Warn("pattern missmatch found", slog.String("pattern", pattern))
		}
	}
}
