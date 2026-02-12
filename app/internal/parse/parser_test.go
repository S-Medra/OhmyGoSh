package parse

import (
	"os"
	"testing"
)

func TestParseCommand_Basic(t *testing.T) {
	cmd, args, redirect, errorRedirect, err := ParseCommand([]string{"echo", "hello"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cmd != "echo" {
		t.Errorf("expected command 'echo', got '%s'", cmd)
	}
	if len(args) != 1 || args[0] != "hello" {
		t.Errorf("expected args ['hello'], got %v", args)
	}
	if redirect != nil {
		t.Errorf("expected nil redirect, got %v", redirect)
	}
	if errorRedirect != nil {
		t.Errorf("expected nil errorRedirect, got %v", errorRedirect)
	}
}

func TestParseCommand_WithRedirect(t *testing.T) {
	cmd, args, redirect, errorRedirect, err := ParseCommand([]string{"echo", "hello", ">", "file.txt"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cmd != "echo" {
		t.Errorf("expected command 'echo', got '%s'", cmd)
	}
	if len(args) != 1 || args[0] != "hello" {
		t.Errorf("expected args ['hello'], got %v", args)
	}
	if redirect == nil {
		t.Fatalf("expected redirect, got nil")
	}
	if redirect.Target != "file.txt" {
		t.Errorf("expected redirect target 'file.txt', got '%s'", redirect.Target)
	}
	if errorRedirect != nil {
		t.Errorf("expected nil errorRedirect, got %v", errorRedirect)
	}
}

func TestParseCommand_MissingRedirectTarget(t *testing.T) {
	_, _, _, _, err := ParseCommand([]string{"echo", ">"})
	if err == nil {
		t.Errorf("expected error for missing redirect target, got nil")
	}
}

func TestParseCommand_EmptyTokens(t *testing.T) {
	_, _, _, _, err := ParseCommand([]string{})
	if err == nil {
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
