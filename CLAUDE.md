# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

A Go CLI tool that generates AI coding assistant rules files (for Cursor, Windsurf, Claude Code, and ChatGPT Codex) from markdown source files in an `agentrules/` directory within any git repository.

## Commands

- **Build**: `go build -o agentrules .`
- **Run**: `go run .`
- **Install**: `go install github.com/sprucehealth/agentrules@latest`
- **Vet**: `go vet ./...`

There are no tests or external dependencies.

## Architecture

This is a single-file Go program (`main.go`) with no dependencies beyond the standard library.

The tool finds the git root from the current working directory, then reads markdown files from `agentrules/` subdirectories (`shared/`, `cursor/`, `windsurf/`, `claude-code/`, `chatgpt-codex/`, `review-guidelines/`) and generates output files:

| Source dirs | Output | Processing |
|---|---|---|
| `shared/` + `cursor/` | `.cursor/rules/*.gen.mdc` | Adds YAML frontmatter if missing |
| `shared/` + `windsurf/` | `.windsurfrules` | Strips frontmatter/first headings, concatenates |
| `shared/` + `claude-code/` | `CLAUDE.md` | Strips frontmatter/first headings, concatenates |
| `shared/` + `chatgpt-codex/` | `AGENTS.md` | Strips frontmatter/first headings, concatenates |
| `review-guidelines/` | `.cursor/BUGBOT.md` | Raw content (no processing) |
| `review-guidelines/` | Appended to `AGENTS.md` | Under `## Review guidelines` section |

Files within each directory are sorted alphabetically for deterministic output. The `shared/` directory content is included in all outputs (except review guidelines).
