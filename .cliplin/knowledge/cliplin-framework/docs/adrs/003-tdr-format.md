# ADR-003: TDR (Technical Decision Record) Format and Usage

## Status
Accepted

## Context

TDR is the preferred format for technical decision records in Cliplin. It uses standard Markdown with YAML frontmatter instead of a custom YAML-only format (TS4), improving compatibility across AI systems (Cursor, Claude, etc.) and avoiding a Cliplin-specific syntax. This ADR explains the TDR format so that AI assistants can understand and work with TDR files correctly.

## Decision

### What is TDR?

TDR (Technical Decision Record) is a Markdown-based format with the same conceptual model as TS4: technical rules, code references, and a clear structure. Each TDR file contains implementation rules and optional code references in a format that is widely supported by AI tools.

### TDR File Structure

A TDR file is a Markdown file (`.md`) with:

1. **YAML frontmatter** (between `---` lines): `tdr`, `id` (kebab-case), `title`, `summary`
2. **Body**: A `# rules` section with optional `##` subsections; rules as bullet lists or prose
3. **Optional `code_refs`** at the end (YAML block or list of file paths)

Example:

```markdown
---
tdr: "1.0"
id: "chromadb-library"
title: "ChromaDB as Library Usage"
summary: "Rules for using ChromaDB as the context store library."
---

# rules

## ChromaDB client usage
- Use `chromadb.PersistentClient(path=...)` for the project context store.
- ALWAYS pass an absolute, resolved path (critical for Windows).

code_refs:
  - "src/cliplin/utils/chromadb.py"
  - "docs/adrs/002-chromadb-rag-context-base.md"
```

### Field Descriptions

- **tdr**: Format version (e.g. `"1.0"`)
- **id**: Unique identifier in kebab-case (lowercase words separated by hyphens)
- **title**: Short title of the technical decision
- **summary**: Brief description for indexing and retrieval
- **Body**: Markdown with `# rules` and optional `##` headings; bullets or prose
- **code_refs**: Optional list of file paths or patterns related to this specification

### Key Principles

1. **TDR does not describe what to build. It defines how to build it correctly.**
2. TDR files act as a **technical contract** for implementation, like TS4
3. Each TDR file should focus on a specific technical decision or set of related rules
4. Use standard Markdown so that any AI host or tool can parse and display it

### Usage

- TDR files are located in `docs/tdrs/` directory
- They are indexed in the context store collection `technical-decision-records`
- AI assistants should query `technical-decision-records` first for technical constraints; use `tech-specs` (TS4) as fallback
- TDR files complement ADRs: ADRs explain *why*, TDRs define *how*
- If the project still has TS4 files in `docs/ts4/`, suggest migrating them to TDR (see project `docs/business/tdr.md`)

## Consequences

### Positive
- Compatible with standard Markdown tooling and AI systems
- Same conceptual model as TS4 (rules, code_refs) with better portability
- Optimized for AI context retrieval; single collection `technical-decision-records` for technical rules

### Notes
- This ADR is indexed in the context store collection `business-and-architecture`
- When creating new technical specs, prefer TDR in `docs/tdrs/` over TS4 in `docs/ts4/`
