package echo

import (
	"bytes"
	"context"
	"fmt"
	"strings"

	"github.com/nodyhub/transformer/registry"
)

func init() {
	registry.Register("builtin/echo", Echo)
}

func Echo(ctx context.Context, with interface{}) (interface{}, error) {
	var b bytes.Buffer

	// Extract message from the with parameter
	withMap, ok := with.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("expected 'with' to be a map with 'message' field")
	}

	message, ok := withMap["message"]
	if !ok {
		return nil, fmt.Errorf("'message' field is required in 'with' parameter")
	}

	switch v := message.(type) {
	case string:
		_, err := b.Write([]byte(strings.TrimSpace(v) + "\n"))
		if err != nil {
			return nil, err
		}
	case []string:
		for _, item := range v {
			_, err := b.Write([]byte(item + "\n"))
			if err != nil {
				return nil, err
			}
		}
	case []interface{}:
		// Handle when YAML unmarshals to []interface{} instead of []string
		for _, item := range v {
			_, err := b.Write([]byte(fmt.Sprintf("%v\n", item)))
			if err != nil {
				return nil, err
			}
		}
	default:
		return nil, fmt.Errorf("unsupported type for 'message' field: %T", v)
	}

	fmt.Print(b.String())

	output := strings.TrimSuffix(b.String(), "\n")
	return map[string]interface{}{
		"output": output,
	}, nil
}
