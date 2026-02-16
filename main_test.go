package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const agentrulesComment = "<!-- @agentrules: If you are reading this as an AI coding agent, know that contents are already included in your agent instructions. It is not necessary for you to read or include this file in your context unless told to do so. -->"

func TestIsAgentrulesComment(t *testing.T) {
	tests := []struct {
		line string
		want bool
	}{
		{agentrulesComment, true},
		{"<!-- @agentrules: short -->", true},
		{"  <!-- @agentrules: indented -->  ", true},
		{"<!-- regular comment -->", false},
		{"some text", false},
		{"<!-- @agentrules no closing", false},
		{"", false},
	}
	for _, tt := range tests {
		if got := isAgentrulesComment(tt.line); got != tt.want {
			t.Errorf("isAgentrulesComment(%q) = %v, want %v", tt.line, got, tt.want)
		}
	}
}

func TestStripAgentrulesComments(t *testing.T) {
	input := agentrulesComment + "\n\nSome content here.\nMore content.\n"
	got := string(stripAgentrulesComments([]byte(input)))
	if strings.Contains(got, "@agentrules") {
		t.Errorf("stripAgentrulesComments did not remove @agentrules comment:\n%s", got)
	}
	if !strings.Contains(got, "Some content here.") {
		t.Errorf("stripAgentrulesComments removed non-comment content:\n%s", got)
	}
}

// setupTestRepo creates a temporary git repo with agentrules source files.
// The shared file and tool-specific files include the @agentrules comment.
func setupTestRepo(t *testing.T) string {
	t.Helper()
	root := t.TempDir()

	// Init a git repo so findGitRoot works
	gitDir := filepath.Join(root, ".git")
	if err := os.Mkdir(gitDir, 0755); err != nil {
		t.Fatal(err)
	}

	dirs := []string{
		"agentrules/shared",
		"agentrules/cursor",
		"agentrules/windsurf",
		"agentrules/claude-code",
		"agentrules/chatgpt-codex",
		"agentrules/review-guidelines",
	}
	for _, d := range dirs {
		if err := os.MkdirAll(filepath.Join(root, d), 0755); err != nil {
			t.Fatal(err)
		}
	}

	sharedContent := agentrulesComment + "\n\n## Shared Rule\n\nDo the thing.\n"
	toolContent := agentrulesComment + "\n\n## Tool Rule\n\nDo the other thing.\n"
	cursorFrontmatter := "---\ndescription: test\nglobs: [\"**\"]\nalwaysApply: true\n---\n" + agentrulesComment + "\n\n## Cursor Rule\n\nCursor specific.\n"
	reviewContent := agentrulesComment + "\n\n## Review Rule\n\nCheck this.\n"

	files := map[string]string{
		"agentrules/shared/01-shared.md":                sharedContent,
		"agentrules/cursor/01-cursor.md":                cursorFrontmatter,
		"agentrules/windsurf/01-windsurf.md":            toolContent,
		"agentrules/claude-code/01-claude.md":           toolContent,
		"agentrules/chatgpt-codex/01-codex.md":          toolContent,
		"agentrules/review-guidelines/01-review.md":     reviewContent,
	}
	for path, content := range files {
		if err := os.WriteFile(filepath.Join(root, path), []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
	}

	return root
}

func TestGenerateCursorRules_StripsComment(t *testing.T) {
	root := setupTestRepo(t)
	dstDir := filepath.Join(root, ".cursor", "rules")

	if err := generateCursorRules(filepath.Join(root, "agentrules", "shared"), dstDir); err != nil {
		t.Fatal(err)
	}
	if err := generateCursorRules(filepath.Join(root, "agentrules", "cursor"), dstDir); err != nil {
		t.Fatal(err)
	}

	entries, _ := os.ReadDir(dstDir)
	for _, e := range entries {
		content, err := os.ReadFile(filepath.Join(dstDir, e.Name()))
		if err != nil {
			t.Fatal(err)
		}
		if strings.Contains(string(content), "@agentrules") {
			t.Errorf("cursor output %s contains @agentrules comment", e.Name())
		}
	}
}

func TestGenerateWindsurfRules_StripsComment(t *testing.T) {
	root := setupTestRepo(t)
	outFile := filepath.Join(root, ".windsurfrules")
	srcDirs := []string{
		filepath.Join(root, "agentrules", "shared"),
		filepath.Join(root, "agentrules", "windsurf"),
	}

	if err := generateWindsurfRules(srcDirs, outFile); err != nil {
		t.Fatal(err)
	}

	content, err := os.ReadFile(outFile)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(content), "@agentrules") {
		t.Errorf("windsurf output contains @agentrules comment")
	}
	if !strings.Contains(string(content), "Shared Rule") {
		t.Errorf("windsurf output missing shared content")
	}
	if !strings.Contains(string(content), "Tool Rule") {
		t.Errorf("windsurf output missing tool-specific content")
	}
}

func TestGenerateClaudeRules_StripsComment(t *testing.T) {
	root := setupTestRepo(t)
	outFile := filepath.Join(root, "CLAUDE.md")
	srcDirs := []string{
		filepath.Join(root, "agentrules", "shared"),
		filepath.Join(root, "agentrules", "claude-code"),
	}

	if err := generateClaudeRules(srcDirs, outFile); err != nil {
		t.Fatal(err)
	}

	content, err := os.ReadFile(outFile)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(content), "@agentrules") {
		t.Errorf("claude output contains @agentrules comment")
	}
	if !strings.Contains(string(content), "Shared Rule") {
		t.Errorf("claude output missing shared content")
	}
	if !strings.Contains(string(content), "Tool Rule") {
		t.Errorf("claude output missing tool-specific content")
	}
}

func TestGenerateAgentRules_StripsComment(t *testing.T) {
	root := setupTestRepo(t)
	outFile := filepath.Join(root, "AGENTS.md")
	srcDirs := []string{
		filepath.Join(root, "agentrules", "shared"),
		filepath.Join(root, "agentrules", "chatgpt-codex"),
	}

	if err := generateAgentRules(srcDirs, outFile, root); err != nil {
		t.Fatal(err)
	}

	content, err := os.ReadFile(outFile)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(content), "@agentrules") {
		t.Errorf("agents output contains @agentrules comment")
	}
	if !strings.Contains(string(content), "Shared Rule") {
		t.Errorf("agents output missing shared content")
	}
	if !strings.Contains(string(content), "Tool Rule") {
		t.Errorf("agents output missing codex-specific content")
	}
	if !strings.Contains(string(content), "Review Rule") {
		t.Errorf("agents output missing review guidelines content")
	}
}

func TestGenerateBugbotRules_StripsComment(t *testing.T) {
	root := setupTestRepo(t)

	if err := generateBugbotRules(root); err != nil {
		t.Fatal(err)
	}

	content, err := os.ReadFile(filepath.Join(root, ".cursor", "BUGBOT.md"))
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(content), "@agentrules") {
		t.Errorf("bugbot output contains @agentrules comment")
	}
	if !strings.Contains(string(content), "Review Rule") {
		t.Errorf("bugbot output missing review content")
	}
}
