package main

import (
	"bytes"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/darthkoax/confluencecli/internal/config"
)

func TestMain_Help(t *testing.T) {
	binary := buildBinary(t)
	defer os.Remove(binary)

	out, err := exec.Command(binary, "help").CombinedOutput()
	if err != nil {
		t.Fatalf("help command failed: %v\n%s", err, out)
	}
	output := string(out)
	if !strings.Contains(output, "confluencecli") {
		t.Errorf("help output missing 'confluencecli': %s", output)
	}
	if !strings.Contains(output, "init") {
		t.Errorf("help output missing 'init': %s", output)
	}
	if !strings.Contains(output, "connect") {
		t.Errorf("help output missing 'connect': %s", output)
	}
}

func TestMain_Version(t *testing.T) {
	binary := buildBinary(t)
	defer os.Remove(binary)

	out, err := exec.Command(binary, "version").CombinedOutput()
	if err != nil {
		t.Fatalf("version command failed: %v\n%s", err, out)
	}
	if !strings.Contains(string(out), "confluencecli") {
		t.Errorf("version output missing 'confluencecli': %s", out)
	}
}

func TestMain_Init(t *testing.T) {
	binary := buildBinary(t)
	defer os.Remove(binary)

	tmpdir, err := os.MkdirTemp("", "confluencecli-init-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpdir)

	out, err := exec.Command(binary, "init", "--dir", tmpdir).CombinedOutput()
	if err != nil {
		t.Fatalf("init command failed: %v\n%s", err, out)
	}

	configPath := filepath.Join(tmpdir, "config.toml")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Errorf("init did not create config at %s", configPath)
	}

	output := string(out)
	if !strings.Contains(output, configPath) {
		t.Errorf("init output missing config path: %s", output)
	}
}

func TestMain_InitAlreadyExists(t *testing.T) {
	binary := buildBinary(t)
	defer os.Remove(binary)

	tmpdir, err := os.MkdirTemp("", "confluencecli-init-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpdir)

	_, err = exec.Command(binary, "init", "--dir", tmpdir).CombinedOutput()
	if err != nil {
		t.Fatalf("first init failed: %v", err)
	}

	out, err := exec.Command(binary, "init", "--dir", tmpdir).CombinedOutput()
	if err == nil {
		t.Error("second init should have failed, but succeeded")
	}
	output := string(out)
	if !strings.Contains(output, "already exists") {
		t.Errorf("expected 'already exists' error, got: %s", output)
	}
}

func TestMain_UnknownCommand(t *testing.T) {
	binary := buildBinary(t)
	defer os.Remove(binary)

	out, err := exec.Command(binary, "nonexistent").CombinedOutput()
	if err == nil {
		t.Error("unknown command should have failed")
	}
	output := string(out)
	if !strings.Contains(output, "Unknown command") {
		t.Errorf("expected 'Unknown command' error, got: %s", output)
	}
}

func TestMain_NoArgs(t *testing.T) {
	binary := buildBinary(t)
	defer os.Remove(binary)

	out, err := exec.Command(binary).CombinedOutput()
	if err == nil {
		t.Error("no args should have failed")
	}
	output := string(out)
	if !strings.Contains(output, "USAGE") && !strings.Contains(output, "Usage") {
		t.Errorf("expected usage output, got: %s", output)
	}
}

func TestParseFlag(t *testing.T) {
	args := []string{"--config", "/path/to/config", "--limit", "10", "--type", "page"}

	if got := parseFlag(args, "--config"); got != "/path/to/config" {
		t.Errorf("parseFlag(--config) = %v, want /path/to/config", got)
	}
	if got := parseFlag(args, "--limit"); got != "10" {
		t.Errorf("parseFlag(--limit) = %v, want 10", got)
	}
	if got := parseFlag(args, "--type"); got != "page" {
		t.Errorf("parseFlag(--type) = %v, want page", got)
	}
	if got := parseFlag(args, "--missing"); got != "" {
		t.Errorf("parseFlag(--missing) = %v, want empty", got)
	}
}

func TestParseFlagInt(t *testing.T) {
	args := []string{"--limit", "10", "--start", "0", "--invalid", "abc"}

	if got := parseFlagInt(args, "--limit", 25); got != 10 {
		t.Errorf("parseFlagInt(--limit) = %v, want 10", got)
	}
	if got := parseFlagInt(args, "--start", 25); got != 0 {
		t.Errorf("parseFlagInt(--start) = %v, want 0", got)
	}
	if got := parseFlagInt(args, "--invalid", 25); got != 25 {
		t.Errorf("parseFlagInt(--invalid) = %v, want 25 (default)", got)
	}
	if got := parseFlagInt(args, "--missing", 25); got != 25 {
		t.Errorf("parseFlagInt(--missing) = %v, want 25 (default)", got)
	}
}

func TestSplit(t *testing.T) {
	tests := []struct {
		input    string
		sep      string
		expected []string
	}{
		{"a,b,c", ",", []string{"a", "b", "c"}},
		{"single", ",", []string{"single"}},
		{"", ",", []string{""}},
		{"a|b|c", "|", []string{"a", "b", "c"}},
	}

	for _, tt := range tests {
		got := split(tt.input, tt.sep)
		if len(got) != len(tt.expected) {
			t.Errorf("split(%q, %q) = %v, want %v", tt.input, tt.sep, got, tt.expected)
			continue
		}
		for i := range got {
			if got[i] != tt.expected[i] {
				t.Errorf("split(%q, %q)[%d] = %v, want %v", tt.input, tt.sep, i, got[i], tt.expected[i])
			}
		}
	}
}

func TestIndexOf(t *testing.T) {
	tests := []struct {
		s        string
		substr   string
		expected int
	}{
		{"hello world", "world", 6},
		{"hello world", "hello", 0},
		{"hello world", "xyz", -1},
		{"", "test", -1},
		{"test", "", 0},
	}

	for _, tt := range tests {
		got := indexOf(tt.s, tt.substr)
		if got != tt.expected {
			t.Errorf("indexOf(%q, %q) = %v, want %v", tt.s, tt.substr, got, tt.expected)
		}
	}
}

func TestGetDefaultConfigPath(t *testing.T) {
	path := config.DefaultConfigPath()
	if path == "" {
		t.Error("config.DefaultConfigPath() returned empty string")
	}
	if !strings.Contains(path, "config.toml") {
		t.Errorf("config.DefaultConfigPath() = %v, should contain config.toml", path)
	}
}

func TestGetDefaultConfigPath_WithEnv(t *testing.T) {
	original := os.Getenv("CONFLUENCE_CONFIG_DIR")
	defer os.Setenv("CONFLUENCE_CONFIG_DIR", original)

	os.Setenv("CONFLUENCE_CONFIG_DIR", "/custom/path")
	path := config.DefaultConfigPath()
	if path != "/custom/path/config.toml" {
		t.Errorf("config.DefaultConfigPath() with env = %v, want /custom/path/config.toml", path)
	}
}

func TestPrintRootHelp(t *testing.T) {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	printRootHelp()

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	if !strings.Contains(output, "confluencecli") {
		t.Errorf("printRootHelp output missing 'confluencecli'")
	}
	if !strings.Contains(output, "COMMANDS") {
		t.Errorf("printRootHelp output missing 'COMMANDS'")
	}
}

func TestPrintVersion(t *testing.T) {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	printVersion()

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	if !strings.Contains(output, "confluencecli version") {
		t.Errorf("printVersion output missing 'confluencecli version'")
	}
}

func TestHandleCommandWithHelp_NoArgs(t *testing.T) {
	helpCalled := false
	printHelp := func() { helpCalled = true }
	actionHandlers := map[string]func([]string){}
	printActionHelp := func(string) {}

	handleCommandWithHelp("test", []string{}, printHelp, actionHandlers, printActionHelp)

	if !helpCalled {
		t.Error("printHelp should have been called with no args")
	}
}

func TestHandleCommandWithHelp_HelpFlag(t *testing.T) {
	helpCalled := false
	printHelp := func() { helpCalled = true }
	actionHandlers := map[string]func([]string){}
	printActionHelp := func(string) {}

	handleCommandWithHelp("test", []string{"help"}, printHelp, actionHandlers, printActionHelp)

	if !helpCalled {
		t.Error("printHelp should have been called with 'help' arg")
	}
}

func TestHandleCommandWithHelp_ValidAction(t *testing.T) {
	actionCalled := false
	printHelp := func() {}
	actionHandlers := map[string]func([]string){
		"valid": func(args []string) { actionCalled = true },
	}
	printActionHelp := func(string) {}

	handleCommandWithHelp("test", []string{"valid"}, printHelp, actionHandlers, printActionHelp)

	if !actionCalled {
		t.Error("valid action handler should have been called")
	}
}

func buildBinary(t *testing.T) string {
	t.Helper()
	tmpdir, err := os.MkdirTemp("", "confluencecli-build-*")
	if err != nil {
		t.Fatal(err)
	}
	binary := filepath.Join(tmpdir, "confluencecli")
	cmd := exec.Command("go", "build", "-o", binary, ".")
	cmd.Dir = filepath.Join(getProjectRoot(t), "cmd", "confluencecli")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("failed to build binary: %v\n%s", err, out)
	}
	return binary
}

func getProjectRoot(t *testing.T) string {
	t.Helper()
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			t.Fatal("could not find project root")
		}
		dir = parent
	}
}
