package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func main() {
	if err := generateRules(); err != nil {
		fmt.Fprintf(os.Stderr, "Error generating rules: %v\n", err)
		os.Exit(1)
	}
}

func generateRules() error {
	rootDir, err := findGitRoot()
	if err != nil {
		return fmt.Errorf("finding git root: %w", err)
	}

	sharedDir := filepath.Join(rootDir, "agentrules", "shared")
	cursorSrcDir := filepath.Join(rootDir, "agentrules", "cursor")
	windsurfSrcDir := filepath.Join(rootDir, "agentrules", "windsurf")
	claudeCodeSrcDir := filepath.Join(rootDir, "agentrules", "claude-code")
	chatgptCodexSrcDir := filepath.Join(rootDir, "agentrules", "chatgpt-codex")
	cursorDst := filepath.Join(rootDir, ".cursor", "rules")
	windsurfRules := filepath.Join(rootDir, ".windsurfrules")
	claudeMD := filepath.Join(rootDir, "CLAUDE.md")
	agentsMD := filepath.Join(rootDir, "AGENTS.md")

	// Create destination directories
	if err := os.MkdirAll(cursorDst, 0755); err != nil {
		return fmt.Errorf("creating cursor directory: %w", err)
	}

	// Clear and generate Cursor rules (.mdc files) from shared + cursor-specific
	if err := os.RemoveAll(cursorDst); err != nil {
		return fmt.Errorf("clearing cursor directory: %w", err)
	}
	if err := generateCursorRules(sharedDir, cursorDst); err != nil {
		return fmt.Errorf("generating cursor rules from shared: %w", err)
	}
	if err := generateCursorRules(cursorSrcDir, cursorDst); err != nil {
		return fmt.Errorf("generating cursor rules from cursor-specific: %w", err)
	}

	// Generate Windsurf rules (concatenated .windsurfrules) from shared + windsurf-specific
	if err := generateWindsurfRules([]string{sharedDir, windsurfSrcDir}, windsurfRules); err != nil {
		return fmt.Errorf("generating windsurf rules: %w", err)
	}

	// Generate Claude rules (concatenated CLAUDE.md) from shared + claude-code-specific
	if err := generateClaudeRules([]string{sharedDir, claudeCodeSrcDir}, claudeMD); err != nil {
		return fmt.Errorf("generating claude rules: %w", err)
	}

	// Generate Agent rules (concatenated AGENTS.md) from shared + chatgpt-codex-specific
	if err := generateAgentRules([]string{sharedDir, chatgptCodexSrcDir}, agentsMD); err != nil {
		return fmt.Errorf("generating agent rules: %w", err)
	}

	// Add warning files to generated directories
	if err := addWarningFiles(cursorDst, claudeMD); err != nil {
		return fmt.Errorf("adding warning files: %w", err)
	}

	fmt.Println("Rules generation completed successfully")
	return nil
}

func findGitRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		if _, err := os.Stat(filepath.Join(dir, ".git")); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", errors.New("not in a git repository")
		}
		dir = parent
	}
}

// ensureTrailingNewline adds a trailing newline to content if it doesn't already have one
func ensureTrailingNewline(content []byte) []byte {
	if len(content) == 0 || content[len(content)-1] != '\n' {
		return append(content, '\n')
	}
	return content
}

func generateCursorRules(srcDir, dstDir string) error {
	// Ensure destination directory exists
	if err := os.MkdirAll(dstDir, 0755); err != nil {
		return err
	}

	files, err := filepath.Glob(filepath.Join(srcDir, "*.md"))
	if err != nil {
		return err
	}

	// Sort files for consistent output
	sort.Strings(files)

	for _, file := range files {
		base := strings.TrimSuffix(filepath.Base(file), ".md")
		outFile := filepath.Join(dstDir, base+".gen.mdc")

		content, err := os.ReadFile(file)
		if err != nil {
			return err
		}

		var output strings.Builder

		// Check if file already has YAML frontmatter
		if strings.HasPrefix(string(content), "---") {
			// Copy verbatim if it has YAML header
			output.Write(content)
		} else {
			// Add default YAML header
			output.WriteString("---\n")
			output.WriteString(fmt.Sprintf("description: Auto‑generated from %s.md\n", base))
			output.WriteString("globs: [\"**\"]\n")
			output.WriteString("alwaysApply: false\n")
			output.WriteString("---\n")
			output.Write(content)
		}

		if err := os.WriteFile(outFile, ensureTrailingNewline([]byte(output.String())), 0600); err != nil {
			return err
		}
		fmt.Printf("[cursor] %s\n", outFile)
	}
	return nil
}

func generateWindsurfRules(srcDirs []string, windsurfFile string) error {
	var allFiles []string

	// Collect files from all source directories
	for _, srcDir := range srcDirs {
		files, err := filepath.Glob(filepath.Join(srcDir, "*.md"))
		if err != nil {
			return err
		}
		allFiles = append(allFiles, files...)
	}

	// Sort files for consistent output
	sort.Strings(allFiles)

	var output strings.Builder
	output.WriteString("# .windsurfrules\n")
	output.WriteString("<!-- generated; DO NOT EDIT. Edit files in agentrules/shared. -->\n\n")

	for _, file := range allFiles {
		func() {
			f, err := os.Open(file)
			if err != nil {
				return
			}
			defer f.Close()

			scanner := bufio.NewScanner(f)
			firstLine := true
			for scanner.Scan() {
				line := scanner.Text()
				// Skip first line if it's a heading (starts with # )
				if firstLine && strings.HasPrefix(line, "# ") {
					firstLine = false
					continue
				}
				// Skip YAML frontmatter
				if firstLine && strings.HasPrefix(line, "---") {
					for scanner.Scan() {
						if strings.HasPrefix(scanner.Text(), "---") {
							break
						}
					}
					firstLine = false
					continue
				}
				firstLine = false
				output.WriteString(line + "\n")
			}

			output.WriteString("\n---\n\n")
		}()
	}

	if err := os.WriteFile(windsurfFile, ensureTrailingNewline([]byte(output.String())), 0600); err != nil {
		return err
	}

	fmt.Println("[windsurf] rebuilt")
	return nil
}

func generateClaudeRules(srcDirs []string, claudeFile string) error {
	var allFiles []string

	// Collect files from all source directories
	for _, srcDir := range srcDirs {
		files, err := filepath.Glob(filepath.Join(srcDir, "*.md"))
		if err != nil {
			return err
		}
		allFiles = append(allFiles, files...)
	}

	// Sort files for consistent output
	sort.Strings(allFiles)

	var output strings.Builder
	output.WriteString("# CLAUDE.md\n")
	output.WriteString("<!-- generated; DO NOT EDIT. Edit files in agentrules/shared. -->\n\n")

	for _, file := range allFiles {
		func() {
			f, err := os.Open(file)
			if err != nil {
				return
			}
			defer f.Close()

			scanner := bufio.NewScanner(f)
			firstLine := true
			for scanner.Scan() {
				line := scanner.Text()
				// Skip first line if it's a heading (starts with # )
				if firstLine && strings.HasPrefix(line, "# ") {
					firstLine = false
					continue
				}
				// Skip YAML frontmatter
				if firstLine && strings.HasPrefix(line, "---") {
					for scanner.Scan() {
						if strings.HasPrefix(scanner.Text(), "---") {
							break
						}
					}
					firstLine = false
					continue
				}
				firstLine = false
				output.WriteString(line + "\n")
			}

			output.WriteString("\n---\n\n")
		}()
	}

	if err := os.WriteFile(claudeFile, ensureTrailingNewline([]byte(output.String())), 0600); err != nil {
		return err
	}

	fmt.Println("[claude] rebuilt")
	return nil
}

func generateAgentRules(srcDirs []string, agentFile string) error {
	var allFiles []string

	// Collect files from all source directories
	for _, srcDir := range srcDirs {
		files, err := filepath.Glob(filepath.Join(srcDir, "*.md"))
		if err != nil {
			return err
		}
		allFiles = append(allFiles, files...)
	}

	// Sort files for consistent output
	sort.Strings(allFiles)

	var output strings.Builder
	output.WriteString("<!-- This is used by Codex to find instructions. -->\n")

	for _, file := range allFiles {
		func() {
			f, err := os.Open(file)
			if err != nil {
				return
			}
			defer f.Close()

			scanner := bufio.NewScanner(f)
			firstLine := true
			for scanner.Scan() {
				line := scanner.Text()
				// Skip first line if it's a heading (starts with # )
				if firstLine && strings.HasPrefix(line, "# ") {
					firstLine = false
					continue
				}
				// Skip YAML frontmatter
				if firstLine && strings.HasPrefix(line, "---") {
					for scanner.Scan() {
						if strings.HasPrefix(scanner.Text(), "---") {
							break
						}
					}
					firstLine = false
					continue
				}
				firstLine = false
				output.WriteString(line + "\n")
			}

			output.WriteString("\n")
		}()
	}

	if err := os.WriteFile(agentFile, ensureTrailingNewline([]byte(output.String())), 0600); err != nil {
		return err
	}

	fmt.Println("[agents] rebuilt")
	return nil
}

func addWarningFiles(cursorDst, claudeMD string) error {
	// Add README to cursor rules directory
	cursorReadme := filepath.Join(cursorDst, "README.md")
	cursorWarning := `# Generated Files - Do Not Edit

⚠️ **WARNING**: All files in this directory are automatically generated.

**DO NOT EDIT** these files directly. Instead, edit the source files in:
- ` + "`agentrules/shared/`" + `

To regenerate these files, run:
` + "```bash" + `
go generate ./agentrules
` + "```" + `

Any changes made directly to files in this directory will be lost when the rules are regenerated.
`
	if err := os.WriteFile(cursorReadme, ensureTrailingNewline([]byte(cursorWarning)), 0600); err != nil {
		return fmt.Errorf("writing cursor README: %w", err)
	}

	fmt.Println("[warnings] added README file to generated directory")
	return nil
}

