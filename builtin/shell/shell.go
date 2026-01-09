package shell

import (
	"bytes"
	"context"
	"embed"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/exec"
	"strings"

	"github.com/nodyhub/transformer/registry"
)

//go:embed module-hook/*
var moduleHookFS embed.FS

func init() {
	registry.Register("builtin/shell", Shell, moduleHook)
}

// Shell executes a shell command specified in the 'with' parameter.
func Shell(ctx context.Context, with interface{}) (interface{}, error) {
	// Extract parameters from the with map
	withMap, ok := with.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("expected 'with' to be a map with 'command' field")
	}

	// Get the command (required)
	commandRaw, ok := withMap["command"]
	if !ok {
		return nil, fmt.Errorf("'command' field is required in 'with' parameter")
	}

	command, ok := commandRaw.(string)
	if !ok {
		return nil, fmt.Errorf("'command' field must be a string")
	}

	// Create the shell command
	cmd := exec.CommandContext(ctx, "sh", "-c", command)

	outputBuffer := bytes.Buffer{}

	// Get redirection paths
	stdoutPath, stdoutSet := withMap["stdout"].(string)
	stderrPath, stderrSet := withMap["stderr"].(string)

	// Check for circular redirection
	if stdoutSet && stderrSet && stdoutPath == "stderr" && stderrPath == "stdout" {
		return nil, fmt.Errorf("circular redirection detected: stdout->stderr and stderr->stdout")
	}

	stdoutWrapper, stderrWrapper, err := setupWriters(cmd, stdoutSet, stdoutPath, stderrSet, stderrPath, &outputBuffer)
	if err != nil {
		return nil, err
	}
	//nolint:govet
	defer func() {
		if err := stdoutWrapper.Close(); err != nil {
			slog.Error("failed to close stdout writer", "error", err)
		}
		if err := stderrWrapper.Close(); err != nil {
			slog.Error("failed to close stderr writer", "error", err)
		}
	}()

	var stderrBuffer strings.Builder
	// If not redirected, also capture stderr
	if !stderrSet {
		cmd.Stderr = &stderrBuffer
	}

	err = cmd.Run()

	result := map[string]interface{}{
		"output":    strings.TrimSuffix(outputBuffer.String(), "\n"),
		"exit_code": -1,
	}

	// Only include stderr if it was captured (not redirected)
	if stderrBuffer.Len() > 0 {
		result["stderr"] = strings.TrimSuffix(stderrBuffer.String(), "\n")
	}

	if cmd.ProcessState != nil {
		result["exit_code"] = cmd.ProcessState.ExitCode()
	}

	// Handle errors from command execution
	// ExitError is expected for non-zero exit codes and is handled above via exit_code
	if err != nil {
		var exitErr *exec.ExitError
		if !errors.As(err, &exitErr) {
			// Actual execution error (command not found, permission denied, etc.)
			return nil, fmt.Errorf("failed to execute command: %w", err)
		}
	}

	return result, nil
}

func moduleHook(ctx context.Context) error {
	slog.Info("executing shell module hook to install system utilities")
	tempDir, err := os.MkdirTemp("", "transformer-shell-*")
	if err != nil {
		return fmt.Errorf("failed to create temporary directory for module hook: %w", err)
	}

	defer func() {
		err = os.RemoveAll(tempDir)
		if err != nil {
			slog.Error("failed to remove temporary directory for module hook", "error", err)
		}
	}()

	// copy all files from moduleHook embedded files to tempDir
	files, err := moduleHookFS.ReadDir("module-hook")
	if err != nil {
		return fmt.Errorf("failed to read module hook files: %w", err)
	}

	for _, file := range files {
		data, err := moduleHookFS.ReadFile("module-hook/" + file.Name())
		if err != nil {
			return fmt.Errorf("failed to read module hook file %s: %w", file.Name(), err)
		}

		destPath := fmt.Sprintf("%s/%s", tempDir, file.Name())
		err = os.WriteFile(destPath, data, 0600)
		if err != nil {
			return fmt.Errorf("failed to write module hook file %s: %w", destPath, err)
		}
	}

	scriptPath := fmt.Sprintf("%s/install-system-utilities.sh", tempDir)

	cmd := exec.CommandContext(ctx, "sh", scriptPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to execute module hook script: %w", err)
	}

	return nil
}

// create temporary script file

// Set file permissions

func setupWriters(cmd *exec.Cmd, stdoutSet bool, stdoutPath string, stderrSet bool, stderrPath string, outputBuffer *bytes.Buffer) (writerWrapper, writerWrapper, error) {
	var stdoutWrapper, stderrWrapper writerWrapper
	var err error

	// Setup stdout writer
	if stdoutSet {
		// If stdout is set to stderr, set it later
		if stdoutPath != "stderr" {
			stdoutWrapper, err = getWriter(stdoutPath)
			if err != nil {
				return writerWrapper{}, writerWrapper{}, err
			}
		}
	} else {
		stdoutWrapper = writerWrapper{writer: outputBuffer, closer: nil}
	}
	cmd.Stdout = stdoutWrapper.writer

	// Setup stderr writer
	if stderrSet {
		// Special case: if stderr is set to "stdout", merge with stdout writer
		if stderrPath == "stdout" {
			stderrWrapper = writerWrapper{writer: cmd.Stdout, closer: nil}
		} else {
			stderrWrapper, err = getWriter(stderrPath)
			if err != nil {
				// Close stdout if it was opened
				if stdoutWrapper.closer != nil {
					stdoutWrapper.closer.Close()
				}
				return writerWrapper{}, writerWrapper{}, err
			}
		}
	} else {
		stderrWrapper = writerWrapper{writer: outputBuffer, closer: nil}
	}
	cmd.Stderr = stderrWrapper.writer

	// Handle case where stdout is set to stderr
	if stdoutSet && stdoutPath == "stderr" {
		cmd.Stdout = cmd.Stderr
	}

	return stdoutWrapper, stderrWrapper, nil
}

func getWriter(path string) (writerWrapper, error) {
	switch path {
	case "discard":
		return writerWrapper{writer: io.Discard, closer: nil}, nil
	case "stdout":
		return writerWrapper{writer: os.Stdout, closer: nil}, nil
	case "stderr":
		return writerWrapper{writer: os.Stderr, closer: nil}, nil
	default:
		f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			return writerWrapper{}, fmt.Errorf("failed to open file %s: %w", path, err)
		}
		return writerWrapper{writer: f, closer: f}, nil
	}
}

type writerWrapper struct {
	writer io.Writer
	closer io.Closer
}

func (w writerWrapper) Close() error {
	if w.closer != nil {
		return w.closer.Close()
	}

	return nil
}
