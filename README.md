# agentrules

A tool for generating AI coding assistant rules files from markdown source files.

## Overview

This tool generates rules files for various AI coding assistants (Cursor, Windsurf, Claude Code, ChatGPT Codex) by processing markdown files from the `agentrules` directory structure in a git repository.

## Usage

Run the tool from within a git repository that has an `agentrules` directory structure:

```bash
go run github.com/sprucehealth/agentrules@latest
```

Or install it first:

```bash
go install github.com/sprucehealth/agentrules@latest
agentrules
```

## Directory Structure

The tool expects the following directory structure in your git repository:

```
your-repo/
├── agentrules/
│   ├── shared/          # Shared rules for all assistants
│   ├── cursor/          # Cursor-specific rules
│   ├── windsurf/        # Windsurf-specific rules
│   ├── claude-code/     # Claude Code-specific rules
│   └── chatgpt-codex/   # ChatGPT Codex-specific rules
```

## Generated Files

The tool generates the following files in your repository root:

- `.cursor/rules/*.gen.mdc` - Cursor rules files (with YAML frontmatter)
- `.windsurfrules` - Windsurf rules file
- `CLAUDE.md` - Claude Code rules file
- `AGENTS.md` - ChatGPT Codex rules file

## How It Works

1. Finds the git repository root from the current working directory
2. Reads markdown files from the `agentrules` subdirectories
3. Generates formatted output files for each AI assistant
4. Adds YAML frontmatter to Cursor rules files if not already present
5. Strips YAML frontmatter and first-level headings from concatenated files

## Installation

### For Team Members

Once this repository is pushed to GitHub, team members can install the tool in two ways:

#### Option 1: Install as a binary (recommended)

```bash
go install github.com/sprucehealth/agentrules@latest
```

This installs the `agentrules` binary to `$GOPATH/bin` (or `$HOME/go/bin` by default). Make sure this directory is in your `PATH`.

After installation, you can run it from anywhere:
```bash
agentrules
```

#### Option 2: Use directly without installation

You can run it directly without installing:
```bash
go run github.com/sprucehealth/agentrules@latest
```

#### Option 3: Use via go generate (in backend repo)

If you're working in the backend repository, you can use it via `go generate`:

```bash
go generate ./agentrules
```

This will automatically run the tool using the `//go:generate` directive in `agentrules/generate.go`.

### Prerequisites

- Go 1.25 or later
- Git access to `github.com/sprucehealth/agentrules` repository
- For private repositories, ensure your Git credentials are configured

