package parse

import (
	"os"
	"testing"
)

func TestParseCommand_Basic(t *testing.T) {
	result := ParseCommand([]string{"echo", "hello"})
	if len(result.Commands) != 1 {
		t.Fatalf("expected 1 command, got %d", len(result.Commands))
	}
	cmd := result.Commands[0]
	if cmd.Err != nil {
		t.Fatalf("unexpected error: %v", cmd.Err)
	}
	if cmd.Cmd != "echo" {
		t.Errorf("expected command 'echo', got '%s'", cmd.Cmd)
	}
	if len(cmd.Args) != 1 || cmd.Args[0] != "hello" {
		t.Errorf("expected args ['hello'], got %v", cmd.Args)
	}
	if cmd.Redirect != nil {
		t.Errorf("expected nil redirect, got %v", cmd.Redirect)
	}
	if cmd.ErrorRedirect != nil {
		t.Errorf("expected nil errorRedirect, got %v", cmd.ErrorRedirect)
	}
}

func TestParseCommand_WithRedirect(t *testing.T) {
	result := ParseCommand([]string{"echo", "hello", ">", "file.txt"})
	if len(result.Commands) != 1 {
		t.Fatalf("expected 1 command, got %d", len(result.Commands))
	}
	cmd := result.Commands[0]
	if cmd.Err != nil {
		t.Fatalf("unexpected error: %v", cmd.Err)
	}
	if cmd.Cmd != "echo" {
		t.Errorf("expected command 'echo', got '%s'", cmd.Cmd)
	}
	if len(cmd.Args) != 1 || cmd.Args[0] != "hello" {
		t.Errorf("expected args ['hello'], got %v", cmd.Args)
	}
	if cmd.Redirect == nil {
		t.Fatalf("expected redirect, got nil")
	}
	if cmd.Redirect.Target != "file.txt" {
		t.Errorf("expected redirect target 'file.txt', got '%s'", cmd.Redirect.Target)
	}
	if cmd.ErrorRedirect != nil {
		t.Errorf("expected nil errorRedirect, got %v", cmd.ErrorRedirect)
	}
}

func TestParseCommand_MissingRedirectTarget(t *testing.T) {
	result := ParseCommand([]string{"echo", ">"})
	if len(result.Commands) != 1 {
		t.Fatalf("expected 1 command, got %d", len(result.Commands))
	}
	cmd := result.Commands[0]
	if cmd.Err == nil {
		t.Errorf("expected error for missing redirect target, got nil")
	}
}

func TestParseCommand_EmptyTokens(t *testing.T) {
	result := ParseCommand([]string{})
	if len(result.Commands) != 1 {
		t.Errorf("expected 1 command (empty), got %d", len(result.Commands))
	}
	if result.Commands[0].Err == nil {
		t.Errorf("expected error for empty command, got nil")
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

func TestSplitByPipe_NoPipe(t *testing.T) {
	result := SplitByPipe([]string{"echo", "hello"})
	if len(result) != 1 {
		t.Errorf("expected 1 segment, got %d", len(result))
	}
	if len(result[0]) != 2 || result[0][0] != "echo" || result[0][1] != "hello" {
		t.Errorf("expected [echo hello], got %v", result)
	}
}

func TestSplitByPipe_SinglePipe(t *testing.T) {
	result := SplitByPipe([]string{"ls", "|", "grep", ".go"})
	if len(result) != 2 {
		t.Errorf("expected 2 segments, got %d", len(result))
	}
	if len(result[0]) != 1 || result[0][0] != "ls" {
		t.Errorf("first segment expected [ls], got %v", result[0])
	}
	if len(result[1]) != 2 || result[1][0] != "grep" || result[1][1] != ".go" {
		t.Errorf("second segment expected [grep .go], got %v", result[1])
	}
}

func TestSplitByPipe_MultiplePipes(t *testing.T) {
	result := SplitByPipe([]string{"ls", "|", "grep", ".go", "|", "wc", "-l"})
	if len(result) != 3 {
		t.Errorf("expected 3 segments, got %d", len(result))
	}
	if result[0][0] != "ls" || result[1][0] != "grep" || result[2][0] != "wc" {
		t.Errorf("unexpected segments: %v", result)
	}
}

func TestSplitByPipe_PipeAtStart(t *testing.T) {
	result := SplitByPipe([]string{"|", "grep", "foo"})
	if len(result) != 2 {
		t.Errorf("expected 2 segments, got %d", len(result))
	}
	if len(result[0]) != 0 {
		t.Errorf("first segment should be empty, got %v", result[0])
	}
	if result[1][0] != "grep" {
		t.Errorf("second segment should start with grep, got %v", result[1])
	}
}

func TestSplitByPipe_PipeAtEnd(t *testing.T) {
	result := SplitByPipe([]string{"echo", "hello", "|"})
	if len(result) != 2 {
		t.Errorf("expected 2 segments, got %d", len(result))
	}
	if result[0][0] != "echo" {
		t.Errorf("first segment should be [echo], got %v", result[0])
	}
	if len(result[1]) != 0 {
		t.Errorf("second segment should be empty, got %v", result[1])
	}
}

func TestParseCommand_WithPipe(t *testing.T) {
	result := ParseCommand([]string{"ls", "|", "grep", ".go"})
	if len(result.Commands) != 2 {
		t.Fatalf("expected 2 commands, got %d", len(result.Commands))
	}
	if result.Commands[0].Cmd != "ls" {
		t.Errorf("first command should be 'ls', got '%s'", result.Commands[0].Cmd)
	}
	if result.Commands[1].Cmd != "grep" || len(result.Commands[1].Args) != 1 || result.Commands[1].Args[0] != ".go" {
		t.Errorf("second command should be 'grep .go', got %v", result.Commands[1])
	}
}

func TestParseCommand_MultiplePipes(t *testing.T) {
	result := ParseCommand([]string{"ls", "|", "grep", ".go", "|", "wc", "-l"})
	if len(result.Commands) != 3 {
		t.Fatalf("expected 3 commands, got %d", len(result.Commands))
	}
	if result.Commands[0].Cmd != "ls" {
		t.Errorf("first command should be 'ls'")
	}
	if result.Commands[1].Cmd != "grep" {
		t.Errorf("second command should be 'grep'")
	}
	if result.Commands[2].Cmd != "wc" || result.Commands[2].Args[0] != "-l" {
		t.Errorf("third command should be 'wc -l'")
	}
}
