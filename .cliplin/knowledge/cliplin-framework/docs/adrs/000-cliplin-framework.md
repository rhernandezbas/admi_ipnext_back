# ADR-000: Cliplin Framework Overview

## Status
Accepted

## Context

This project uses the **Cliplin Framework** for AI-assisted development driven by specifications. This ADR provides essential context about the framework so that AI assistants understand the operational model and specification types available.

## Decision

### Core Principle (Kidlin's Law)
> **Describe the problem clearly, and half of it is already solved.**

Cliplin operationalizes this by enforcing that **AI is only allowed to act on well-defined, versioned specifications that live in the repository**. Code becomes an output of the system — not its source of truth.

### The Four Pillars of Cliplin

Cliplin is built on four complementary specification pillars, each with a precise responsibility:

#### 1. Business Features (.feature files – Gherkin)
- **Purpose**: Define *what* the system must do and *why*
- **Location**: `docs/features/*.feature`
- **Format**: Gherkin syntax with Given-When-Then scenarios
- **Key principle**: If a feature does not exist, the functionality does not exist
- **Tags**: 
  - `@status:pending` or no tag = not implemented
  - `@status:implemented` = implemented and active
  - `@status:deprecated` = must not be modified
  - `@changed:YYYY-MM-DD` = change tracking
  - `@reason:<description>` = change justification
- **Role**: Features orchestrate execution and represent the source of truth

#### 2. UI Intent Specifications (YAML)
- **Purpose**: Define *how the system expresses intent to users*, beyond visual design
- **Location**: `docs/ui-intent/*.yaml`
- **Focus**: Semantic meaning, layout relationships, intent preservation, state-driven behavior
- **Schema sections**: structure, semantics, state_model, motion, constraints, annotations
- **Key principle**: Emphasizes semantic meaning over visual appearance
- **Usage**: Allows AI to generate UI code without guessing user experience decisions

#### 3. TDR – Technical Decision Records (Markdown, preferred)
- **Purpose**: Define *how software must be implemented* in a format compatible across AI systems
- **Location**: `docs/tdrs/*.md`
- **Key principle**: TDR does not describe what to build. It defines how to build it correctly.
- **Contains**: Same as TS4 — coding conventions, naming rules, validation strategies, allowed/forbidden patterns, project-specific technical decisions
- **Format**: Markdown with YAML frontmatter (`tdr`, `id`, `title`, `summary`), `# rules` section, optional `code_refs`; see `docs/business/tdr.md`
- **Role**: Technical contract for implementation; preferred over TS4 for new specs (Cursor, Claude, and other hosts)
- **Collection**: `technical-decision-records` — query this first for technical constraints

#### 4. TS4 – Technical Specification Files (YAML, deprecated)
- **Purpose**: Same as TDR; legacy format.
- **Location**: `docs/ts4/*.ts4`
- **Status**: **Deprecated.** Still supported and indexed in `tech-specs`; when you encounter TS4, suggest migrating to TDR (see `docs/business/tdr.md`).

#### 5. Architecture Decision Records and Business Documentation (ADRs and .md files)
- **Purpose**: Preserve *why technical decisions were made*
- **Locations**: `docs/adrs/*.md`, `docs/business/*.md`
- **Contains**: Architecture choices, trade-offs, constraints, business descriptions, environment considerations
- **Role**: Prevents AI (and humans) from reopening closed decisions or proposing invalid architectures

### Cliplin as an Operational Model

**Valid Inputs Only:**
- Business Features (.feature in docs/features/)
- UI Intent specifications (.yaml in docs/ui-intent/)
- TDR technical rules (.md in docs/tdrs/) — prefer these; collection `technical-decision-records`
- TS4 technical rules (.ts4 in docs/ts4/) — deprecated; collection `tech-specs`; suggest migrating to TDR
- ADRs and business documentation (.md in docs/adrs/ and docs/business/)

**Everything else is noise.** All outputs must be traceable back to a specification.

### Constraints

Cliplin works by **deliberate limitation**:
- Business constraints (Features)
- Semantic constraints (UI Intent)
- Technical constraints (TS4)
- Architectural constraints (ADRs)

Creativity is replaced by clarity.

### Outputs

Expected outputs are:
- Business-aligned code
- Tests derived from scenarios
- UI consistent with intent
- Architecture-compliant changes

Every output must be traceable back to a specification.

## Consequences

### Positive
- Reduced ambiguity in AI-assisted development
- Predictable AI behavior based on specifications
- Clear separation of concerns across specification types
- Versionable and maintainable specifications
- All context available through the Cliplin MCP (context store) indexing

### Notes
- This ADR should be indexed in the context store collection `business-and-architecture`
- AI assistants should query this ADR and related context files before starting any work
- All specifications must be kept up to date and properly indexed
