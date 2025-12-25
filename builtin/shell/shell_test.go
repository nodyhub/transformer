package shell

import (
	"context"
	"os"
	"strings"
	"testing"
)

func TestShell_SimpleCommand(t *testing.T) {
	ctx := context.Background()
	with := map[string]interface{}{
		"command": "echo 'Hello World'",
	}

	result, err := Shell(ctx, with)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "Hello World"
	if result != expected {
		t.Errorf("expected '%s', got '%v'", expected, result)
	}
}

func TestShell_CommandWithOutput(t *testing.T) {
	ctx := context.Background()
	with := map[string]interface{}{
		"command": "echo 'line1'; echo 'line2'",
	}

	result, err := Shell(ctx, with)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "line1\nline2"
	if result != expected {
		t.Errorf("expected '%s', got '%v'", expected, result)
	}
}

func TestShell_MultilineCommand(t *testing.T) {
	ctx := context.Background()
	with := map[string]interface{}{
		"command": "echo 'first'\necho 'second'\necho 'third'",
	}

	result, err := Shell(ctx, with)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "first\nsecond\nthird"
	if result != expected {
		t.Errorf("expected '%s', got '%v'", expected, result)
	}
}

func TestShell_StdoutRedirectToFile(t *testing.T) {
	ctx := context.Background()
	tmpfile := "/tmp/shell_test_stdout.log"
	defer os.Remove(tmpfile)

	with := map[string]interface{}{
		"command": "echo 'logged to file'",
		"stdout":  tmpfile,
	}

	result, err := Shell(ctx, with)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Should return empty since output was redirected
	if result != "" {
		t.Errorf("expected empty result, got '%v'", result)
	}

	// Check file contents
	content, err := os.ReadFile(tmpfile)
	if err != nil {
		t.Fatalf("failed to read log file: %v", err)
	}

	expected := "logged to file\n"
	if string(content) != expected {
		t.Errorf("expected file content '%s', got '%s'", expected, string(content))
	}
}

func TestShell_StderrRedirectToFile(t *testing.T) {
	ctx := context.Background()
	tmpfile := "/tmp/shell_test_stderr.log"
	defer os.Remove(tmpfile)

	with := map[string]interface{}{
		"command": "echo 'error message' >&2",
		"stderr":  tmpfile,
	}

	result, err := Shell(ctx, with)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Should return empty since stderr was redirected
	if result != "" {
		t.Errorf("expected empty result, got '%v'", result)
	}

	// Check file contents
	content, err := os.ReadFile(tmpfile)
	if err != nil {
		t.Fatalf("failed to read log file: %v", err)
	}

	expected := "error message\n"
	if string(content) != expected {
		t.Errorf("expected file content '%s', got '%s'", expected, string(content))
	}
}

func TestShell_StdoutDiscard(t *testing.T) {
	ctx := context.Background()
	with := map[string]interface{}{
		"command": "echo 'this is discarded'",
		"stdout":  "discard",
	}

	result, err := Shell(ctx, with)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != "" {
		t.Errorf("expected empty result, got '%v'", result)
	}
}

func TestShell_StderrDiscard(t *testing.T) {
	ctx := context.Background()
	with := map[string]interface{}{
		"command": "echo 'this is discarded' >&2",
		"stderr":  "discard",
	}

	result, err := Shell(ctx, with)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != "" {
		t.Errorf("expected empty result, got '%v'", result)
	}
}

func TestShell_MergeStderrToStdout(t *testing.T) {
	ctx := context.Background()
	with := map[string]interface{}{
		"command": "echo 'stdout'; echo 'stderr' >&2",
		"stderr":  "stdout",
	}

	result, err := Shell(ctx, with)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Both should be captured together
	if !strings.Contains(result.(string), "stdout") {
		t.Errorf("expected result to contain 'stdout', got '%v'", result)
	}
	if !strings.Contains(result.(string), "stderr") {
		t.Errorf("expected result to contain 'stderr', got '%v'", result)
	}
}

func TestShell_MergeStdoutToStderr(t *testing.T) {
	ctx := context.Background()
	with := map[string]interface{}{
		"command": "echo 'stdout'; echo 'stderr' >&2",
		"stdout":  "stderr",
	}

	result, err := Shell(ctx, with)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Both should be captured together
	if !strings.Contains(result.(string), "stdout") {
		t.Errorf("expected result to contain 'stdout', got '%v'", result)
	}
	if !strings.Contains(result.(string), "stderr") {
		t.Errorf("expected result to contain 'stderr', got '%v'", result)
	}
}

func TestShell_BothRedirectedToFiles(t *testing.T) {
	ctx := context.Background()
	stdoutFile := "/tmp/shell_test_stdout2.log"
	stderrFile := "/tmp/shell_test_stderr2.log"
	defer os.Remove(stdoutFile)
	defer os.Remove(stderrFile)

	with := map[string]interface{}{
		"command": "echo 'to stdout'; echo 'to stderr' >&2",
		"stdout":  stdoutFile,
		"stderr":  stderrFile,
	}

	result, err := Shell(ctx, with)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Should return empty since both were redirected
	if result != "" {
		t.Errorf("expected empty result, got '%v'", result)
	}

	// Check stdout file
	stdoutContent, err := os.ReadFile(stdoutFile)
	if err != nil {
		t.Fatalf("failed to read stdout file: %v", err)
	}
	if !strings.Contains(string(stdoutContent), "to stdout") {
		t.Errorf("stdout file should contain 'to stdout', got '%s'", string(stdoutContent))
	}

	// Check stderr file
	stderrContent, err := os.ReadFile(stderrFile)
	if err != nil {
		t.Fatalf("failed to read stderr file: %v", err)
	}
	if !strings.Contains(string(stderrContent), "to stderr") {
		t.Errorf("stderr file should contain 'to stderr', got '%s'", string(stderrContent))
	}
}

func TestShell_CommandFailure(t *testing.T) {
	ctx := context.Background()
	with := map[string]interface{}{
		"command": "exit 1",
	}

	_, err := Shell(ctx, with)
	if err == nil {
		t.Fatal("expected error for failed command")
	}

	if !strings.Contains(err.Error(), "failed to execute shell command") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestShell_MissingCommandField(t *testing.T) {
	ctx := context.Background()
	with := map[string]interface{}{}

	_, err := Shell(ctx, with)
	if err == nil {
		t.Fatal("expected error for missing command field")
	}

	expectedErr := "'command' field is required in 'with' parameter"
	if err.Error() != expectedErr {
		t.Errorf("expected error '%s', got '%s'", expectedErr, err.Error())
	}
}

func TestShell_InvalidWithType(t *testing.T) {
	ctx := context.Background()
	with := "not a map"

	_, err := Shell(ctx, with)
	if err == nil {
		t.Fatal("expected error for invalid with type")
	}

	expectedErr := "expected 'with' to be a map with 'command' field"
	if err.Error() != expectedErr {
		t.Errorf("expected error '%s', got '%s'", expectedErr, err.Error())
	}
}

func TestShell_InvalidCommandType(t *testing.T) {
	ctx := context.Background()
	with := map[string]interface{}{
		"command": 12345, // not a string
	}

	_, err := Shell(ctx, with)
	if err == nil {
		t.Fatal("expected error for invalid command type")
	}

	expectedErr := "'command' field must be a string"
	if err.Error() != expectedErr {
		t.Errorf("expected error '%s', got '%s'", expectedErr, err.Error())
	}
}

func TestShell_EmptyCommand(t *testing.T) {
	ctx := context.Background()
	with := map[string]interface{}{
		"command": "",
	}

	result, err := Shell(ctx, with)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Empty command should return empty output
	if result != "" {
		t.Errorf("expected empty result, got '%v'", result)
	}
}

func TestShell_CaptureStdoutAndStderr(t *testing.T) {
	ctx := context.Background()
	with := map[string]interface{}{
		"command": "echo 'normal'; echo 'error' >&2",
	}

	result, err := Shell(ctx, with)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Both stdout and stderr should be captured
	resultStr := result.(string)
	if !strings.Contains(resultStr, "normal") {
		t.Errorf("expected result to contain 'normal', got '%s'", resultStr)
	}
	if !strings.Contains(resultStr, "error") {
		t.Errorf("expected result to contain 'error', got '%s'", resultStr)
	}
}

func TestShell_CircularRedirection(t *testing.T) {
	ctx := context.Background()
	with := map[string]interface{}{
		"command": "echo 'test'",
		"stdout":  "stderr",
		"stderr":  "stdout",
	}

	_, err := Shell(ctx, with)
	if err == nil {
		t.Fatal("expected error for circular redirection (stdout->stderr, stderr->stdout)")
	}

	if !strings.Contains(err.Error(), "circular") {
		t.Errorf("expected error to mention 'circular', got: %v", err)
	}
}
