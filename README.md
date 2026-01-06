# agentrules

Generates AI coding assistant rules files from markdown source files for Cursor, Windsurf, Claude Code, and ChatGPT Codex.

## Usage

```bash
go install github.com/sprucehealth/agentrules@latest
agentrules
```

Or run directly: `go run github.com/sprucehealth/agentrules@latest`

**Requirements**: Go 1.25+

## Directory Structure

The tool expects the following directory structure in your git repository:

```
your-repo/
├── agentrules/
│   ├── shared/             # Shared rules for all assistants
│   ├── cursor/             # Cursor-specific rules
│   ├── windsurf/           # Windsurf-specific rules
│   ├── claude-code/        # Claude Code-specific rules
│   ├── chatgpt-codex/      # ChatGPT Codex-specific rules
│   └── review-guidelines/  # Review guidelines (optional)
```

## Generated Files

The tool generates the following files in your repository root:

- `.cursor/rules/*.gen.mdc` - Cursor rules files (with YAML frontmatter)
- `.cursor/BUGBOT.md` - Cursor Bugbot review guidelines (from `review-guidelines/`, as-is)
- `.windsurfrules` - Windsurf rules file
- `CLAUDE.md` - Claude Code rules file
- `AGENTS.md` - ChatGPT Codex rules file (includes `## Review guidelines` section if present)

## How It Works

1. Finds the git repository root from the current working directory
2. Reads markdown files from `agentrules/` subdirectories
3. Generates formatted output files for each AI assistant:
   - Cursor: adds YAML frontmatter if not present
   - Windsurf/Claude/Codex: strips YAML frontmatter and first headings
4. If `review-guidelines/` exists:
   - Appends processed content to `AGENTS.md` under `## Review guidelines`
   - Generates `.cursor/BUGBOT.md` with raw content (no processing)
