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
	resMap := result.(map[string]interface{})
	expected := "Hello World"
	if resMap["output"] != expected {
		t.Errorf("expected '%s', got '%v'", expected, resMap["output"])
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
	resMap := result.(map[string]interface{})
	expected := "line1\nline2"
	if resMap["output"] != expected {
		t.Errorf("expected '%s', got '%v'", expected, resMap["output"])
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
	resMap := result.(map[string]interface{})
	expected := "first\nsecond\nthird"
	if resMap["output"] != expected {
		t.Errorf("expected '%s', got '%v'", expected, resMap["output"])
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
	resMap := result.(map[string]interface{})
	// Should return empty output since redirected
	if resMap["output"] != "" {
		t.Errorf("expected empty result, got '%v'", resMap["output"])
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
	resMap := result.(map[string]interface{})
	// Should return empty output since redirected
	if resMap["output"] != "" {
		t.Errorf("expected empty result, got '%v'", resMap["output"])
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
	resMap := result.(map[string]interface{})
	if resMap["output"] != "" {
		t.Errorf("expected empty result, got '%v'", resMap["output"])
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
	resMap := result.(map[string]interface{})
	if resMap["output"] != "" {
		t.Errorf("expected empty result, got '%v'", resMap["output"])
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
	resMap := result.(map[string]interface{})
	// Both should be captured together
	if !strings.Contains(resMap["output"].(string), "stdout") {
		t.Errorf("expected result to contain 'stdout', got '%v'", resMap["output"])
	}
	if !strings.Contains(resMap["output"].(string), "stderr") {
		t.Errorf("expected result to contain 'stderr', got '%v'", resMap["output"])
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
	resMap := result.(map[string]interface{})
	// stdout should be in output, stderr should be in stderr
	if !strings.Contains(resMap["output"].(string), "stdout") {
		t.Errorf("expected output to contain 'stdout', got '%v'", resMap["output"])
	}
	if !strings.Contains(resMap["stderr"].(string), "stderr") {
		t.Errorf("expected stderr to contain 'stderr', got '%v'", resMap["stderr"])
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
	resMap := result.(map[string]interface{})
	// Should return empty output since both were redirected
	if resMap["output"] != "" {
		t.Errorf("expected empty result, got '%v'", resMap["output"])
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

	result, err := Shell(ctx, with)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	resMap := result.(map[string]interface{})
	if resMap["exit_code"].(int) == 0 {
		t.Fatal("expected non-zero exit code for failed command")
	}
	if resMap["error"] == nil || resMap["error"] == "" {
		t.Errorf("expected error message in result, got '%v'", resMap["error"])
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
	resMap := result.(map[string]interface{})
	// Empty command should return empty output
	if resMap["output"] != "" {
		t.Errorf("expected empty result, got '%v'", resMap["output"])
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
	resMap := result.(map[string]interface{})
	// stdout should be in output, stderr should be in stderr
	if !strings.Contains(resMap["output"].(string), "normal") {
		t.Errorf("expected output to contain 'normal', got '%s'", resMap["output"])
	}
	if !strings.Contains(resMap["stderr"].(string), "error") {
		t.Errorf("expected stderr to contain 'error', got '%s'", resMap["stderr"])
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
