package parse

import (
	"os"
	"testing"
)

func TestParseCommand_Basic(t *testing.T) {
	result := ParseCommand([]string{"echo", "hello"})
	if result.Err != nil {
		t.Fatalf("unexpected error: %v", result.Err)
	}
	if result.Cmd != "echo" {
		t.Errorf("expected command 'echo', got '%s'", result.Cmd)
	}
	if len(result.Args) != 1 || result.Args[0] != "hello" {
		t.Errorf("expected args ['hello'], got %v", result.Args)
	}
	if result.Redirect != nil {
		t.Errorf("expected nil redirect, got %v", result.Redirect)
	}
	if result.ErrorRedirect != nil {
		t.Errorf("expected nil errorRedirect, got %v", result.ErrorRedirect)
	}
}

func TestParseCommand_WithRedirect(t *testing.T) {
	result := ParseCommand([]string{"echo", "hello", ">", "file.txt"})
	if result.Err != nil {
		t.Fatalf("unexpected error: %v", result.Err)
	}
	if result.Cmd != "echo" {
		t.Errorf("expected command 'echo', got '%s'", result.Cmd)
	}
	if len(result.Args) != 1 || result.Args[0] != "hello" {
		t.Errorf("expected args ['hello'], got %v", result.Args)
	}
	if result.Redirect == nil {
		t.Fatalf("expected redirect, got nil")
	}
	if result.Redirect.Target != "file.txt" {
		t.Errorf("expected redirect target 'file.txt', got '%s'", result.Redirect.Target)
	}
	if result.ErrorRedirect != nil {
		t.Errorf("expected nil errorRedirect, got %v", result.ErrorRedirect)
	}
}

func TestParseCommand_MissingRedirectTarget(t *testing.T) {
	result := ParseCommand([]string{"echo", ">"})
	if result.Err == nil {
		t.Errorf("expected error for missing redirect target, got nil")
	}
}

func TestParseCommand_EmptyTokens(t *testing.T) {
	result := ParseCommand([]string{})
	if result.Err == nil {
		t.Errorf("expected error for empty tokens, got nil")
	}
}

func TestRedirect_Open_ValidFile(t *testing.T) {
	file := "test_redirect.txt"
	redirect := &Redirect{FD: 1, Target: file}
	writer, err := redirect.Open()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if writer.file == nil {
		t.Errorf("expected file to be opened, got nil")
	}
	writer.Close()
	os.Remove(file)
}

func TestRedirect_Open_InvalidFile(t *testing.T) {
	redirect := &Redirect{FD: 1, Target: "/invalid/path/to/file.txt"}
	_, err := redirect.Open()
	if err == nil {
		t.Errorf("expected error for invalid file path, got nil")
	}
}

func TestRedirectWriter_Writer_Fallback(t *testing.T) {
	w := &RedirectWriter{file: nil}
	fallback := os.Stdout
	got := w.Writer(fallback)
	if got != fallback {
		t.Errorf("expected fallback writer, got %v", got)
	}
}

func TestRedirectWriter_Close_NilFile(t *testing.T) {
	w := &RedirectWriter{file: nil}
	if err := w.Close(); err != nil {
		t.Errorf("expected nil error for nil file, got %v", err)
	}
}

func TestRedirectWriter_Close_ValidFile(t *testing.T) {
	file := "test_close.txt"
	f, _ := os.Create(file)
	w := &RedirectWriter{file: f}
	if err := w.Close(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	os.Remove(file)
}
