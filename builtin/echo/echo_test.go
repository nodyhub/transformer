package echo

import (
	"bytes"
	"context"
	"io"
	"os"
	"testing"
)

func captureStdout(f func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String()
}

func getOutputFromResult(t *testing.T, result interface{}) string {
	resultMap, ok := result.(map[string]interface{})
	if !ok {
		t.Fatal("expected result to be a map")
	}
	output, ok := resultMap["output"].(string)
	if !ok {
		t.Fatal("expected result to have 'output' field")
	}
	return output
}

func TestPrint_SimpleString(t *testing.T) {
	ctx := context.Background()
	with := map[string]interface{}{
		"message": "Hello World",
	}

	output := captureStdout(func() {
		result, err := Echo(ctx, with)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		resultMap, ok := result.(map[string]interface{})
		if !ok {
			t.Fatal("expected result to be a map")
		}

		output, ok := resultMap["output"].(string)
		if !ok {
			t.Fatal("expected result to have 'output' field")
		}

		if output != "Hello World" {
			t.Errorf("expected 'Hello World', got '%v'", output)
		}
	})

	expected := "Hello World\n"
	if output != expected {
		t.Errorf("expected stdout '%s', got '%s'", expected, output)
	}
}

func TestPrint_StringWithWhitespace(t *testing.T) {
	ctx := context.Background()
	with := map[string]interface{}{
		"message": "  Hello World  ",
	}

	output := captureStdout(func() {
		result, err := Echo(ctx, with)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		output := getOutputFromResult(t, result)
		if output != "Hello World" {
			t.Errorf("expected 'Hello World', got '%v'", output)
		}
	})

	expected := "Hello World\n"
	if output != expected {
		t.Errorf("expected stdout '%s', got '%s'", expected, output)
	}
}

func TestPrint_MultilineString(t *testing.T) {
	ctx := context.Background()
	with := map[string]interface{}{
		"message": "Line 1\nLine 2\nLine 3",
	}

	output := captureStdout(func() {
		result, err := Echo(ctx, with)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		output := getOutputFromResult(t, result)
		expected := "Line 1\nLine 2\nLine 3"
		if output != expected {
			t.Errorf("expected '%s', got '%v'", expected, output)
		}
	})

	expected := "Line 1\nLine 2\nLine 3\n"
	if output != expected {
		t.Errorf("expected stdout '%s', got '%s'", expected, output)
	}
}

func TestPrint_ListOfStrings(t *testing.T) {
	ctx := context.Background()
	with := map[string]interface{}{
		"message": []string{"First", "Second", "Third"},
	}

	output := captureStdout(func() {
		result, err := Echo(ctx, with)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		output := getOutputFromResult(t, result)
		expected := "First\nSecond\nThird"
		if output != expected {
			t.Errorf("expected '%s', got '%v'", expected, output)
		}
	})

	expected := "First\nSecond\nThird\n"
	if output != expected {
		t.Errorf("expected stdout '%s', got '%s'", expected, output)
	}
}

func TestPrint_ListOfInterfaces(t *testing.T) {
	ctx := context.Background()
	with := map[string]interface{}{
		"message": []interface{}{"String", 123, true},
	}

	output := captureStdout(func() {
		result, err := Echo(ctx, with)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		output := getOutputFromResult(t, result)
		expected := "String\n123\ntrue"
		if output != expected {
			t.Errorf("expected '%s', got '%v'", expected, output)
		}
	})

	expected := "String\n123\ntrue\n"
	if output != expected {
		t.Errorf("expected stdout '%s', got '%s'", expected, output)
	}
}

func TestPrint_EmptyString(t *testing.T) {
	ctx := context.Background()
	with := map[string]interface{}{
		"message": "",
	}

	output := captureStdout(func() {
		result, err := Echo(ctx, with)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		output := getOutputFromResult(t, result)
		if output != "" {
			t.Errorf("expected empty string, got '%v'", output)
		}
	})

	expected := "\n"
	if output != expected {
		t.Errorf("expected stdout '%s', got '%s'", expected, output)
	}
}

func TestPrint_EmptyList(t *testing.T) {
	ctx := context.Background()
	with := map[string]interface{}{
		"message": []string{},
	}

	output := captureStdout(func() {
		result, err := Echo(ctx, with)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		output := getOutputFromResult(t, result)
		if output != "" {
			t.Errorf("expected empty string, got '%v'", output)
		}
	})

	if output != "" {
		t.Errorf("expected no stdout, got '%s'", output)
	}
}

func TestPrint_MissingMessageField(t *testing.T) {
	ctx := context.Background()
	with := map[string]interface{}{}

	_, err := Echo(ctx, with)
	if err == nil {
		t.Fatal("expected error for missing message field")
	}

	expectedErr := "'message' field is required in 'with' parameter"
	if err.Error() != expectedErr {
		t.Errorf("expected error '%s', got '%s'", expectedErr, err.Error())
	}
}

func TestPrint_InvalidWithType(t *testing.T) {
	ctx := context.Background()
	with := "not a map"

	_, err := Echo(ctx, with)
	if err == nil {
		t.Fatal("expected error for invalid with type")
	}

	expectedErr := "expected 'with' to be a map with 'message' field"
	if err.Error() != expectedErr {
		t.Errorf("expected error '%s', got '%s'", expectedErr, err.Error())
	}
}

func TestPrint_UnsupportedMessageType(t *testing.T) {
	ctx := context.Background()
	with := map[string]interface{}{
		"message": 12345, // number not supported
	}

	_, err := Echo(ctx, with)
	if err == nil {
		t.Fatal("expected error for unsupported message type")
	}

	if err.Error() != "unsupported type for 'message' field: int" {
		t.Errorf("unexpected error: %v", err)
	}
}
