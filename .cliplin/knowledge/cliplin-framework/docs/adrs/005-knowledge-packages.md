# ADR-005: Knowledge Packages (Cliplin as Knowledge Package Manager)

## Status
Accepted

## Context

This project uses Cliplin and can depend on **knowledge packages**: external repositories that contain ADRs, TDRs/TS4, business docs, features, rules, or skills. Those packages are installed under the project and indexed in the same context store as project specs, so the AI can use them as context.

## Decision

### Command and configuration

- **CLI command**: `cliplin knowledge` with subcommands: `list`, `add`, `remove`, `update`, `show`, `install`.
- **Configuration**: Package list is declared in `cliplin.yaml` at project root under the top-level key `knowledge` (list of entries with `name`, `source`, `version`).
- **Installation**: Packages live under `.cliplin/knowledge/<name>-<source_normalized>/`. Content is obtained via git sparse checkout. The **name** may be a path (e.g. `AWS/aws-sqs`) to install a nested subfolder from a monorepo, or a top-level folder name for multi-package repos.

### Repository layout examples

To make the expected structure explicit for AI systems and humans, knowledge package repositories SHOULD follow one of these patterns:

- **Single-package repository** (one knowledge package per repo). Example: a commons library:

  ```text
  cliplin-commons/
  ├── adrs/
  ├── tdrs/
  ├── skills/
  └── business/
  ```

- **Multi-package repository with nested subpackages** (one repo, many knowledge packages). Example: provider- and service-specific packages:

  ```text
  cliplin-knowledge/
  ├── aws/
  │   ├── sqs/
  │   │   ├── adrs/
  │   │   ├── tdrs/
  │   │   ├── skills/
  │   │   └── business/
  │   └── ec2/
  │       ├── adrs/
  │       ├── tdrs/
  │       ├── skills/
  │       └── business/
  └── google-cloud/
      └── pubsub/
          ├── adrs/
          ├── tdrs/
          ├── skills/
          └── business/
  ```

  Each leaf folder (`aws/sqs`, `aws/ec2`, `google-cloud/pubsub`, …) acts as an independent knowledge package. In `cliplin.yaml` these can be declared using the path as `name` (for example `aws/sqs` or `google-cloud/pubsub`), combined with the repository `source` and `version`.

### Subcommands (summary)

- `cliplin knowledge list` — List packages declared and their install status.
- `cliplin knowledge add <name> <source> <version>` — Add a package, update config, clone, and reindex.
- `cliplin knowledge remove <name>` — Remove from config, delete directory, purge documents from context store.
- `cliplin knowledge update <name>` — Update to the configured (or given) version and reindex.
- `cliplin knowledge show <name>` — Show name, source, version, path, and file count.
- `cliplin knowledge install` — Install all packages declared in cliplin.yaml (add if missing, update if installed). Use `--force` to reinstall from scratch with the configured version.

### Context store and visibility

- Documents under `.cliplin/knowledge/**` are indexed in the same collections as project docs (e.g. adrs → business-and-architecture, tdrs → technical-decision-records, ts4 → tech-specs). The AI loads them via the Cliplin MCP (context store) when querying context.
- After add/update, the package is reindexed automatically; after remove, its documents are removed from the store.

### Full usage and conventions

For detailed usage, configuration format, multi-package vs single-package repos, and conventions for package repositories, see **docs/business/knowledge-packages.md**.

## Notes

- This ADR is created by `cliplin init` so that AI assistants and users have visibility of the knowledge package feature and know that the `cliplin knowledge` command is available in this project.
- Indexed in the context store collection `business-and-architecture`.
