# ADR-0018: Adopt MADR as ADR Format

**Status:** accepted

**Date:** 2026-07-21

**Deciders:** @trevor-vaughan @marcusburghardt @jpower432

## Context

[ADR-0001](0001-use-adrs.md) established that ComplyTime records architectural decisions as numbered Markdown files in `docs/ADRs/`. The current template captures Context, Decision, and Consequences. That covers simple choices, but it says nothing about what alternatives were evaluated, what forces drove the decision, or what tradeoffs each option carried.

In practice, several existing ADRs already work around this gap. [ADR-0004](0004-grpc-provider-plugin-architecture.md) embeds a numbered options list in its Context section. [ADR-0006](0006-complypack-content-envelope.md) enumerates four review-identified problems that motivated the choice. These are ad-hoc adaptations of a pattern that a structured template would provide by default.

As the contributor base grows and decisions shift from retroactive capture to prospective proposals, the lack of explicit sections for drivers, alternatives, and per-option analysis creates predictable review friction. "What else did you consider?" and "Why not X?" become recurring PR comments that the template should preempt.

MADR (Markdown Architectural Decision Records) is the most widely adopted ADR template format in the open-source ecosystem. Key facts:

- Created by Oliver Kopp, Anita Armbruster, and Olaf Zimmermann. Peer-reviewed in "Markdown Architectural Decision Records: Format and Tool Support" (ZEUS 2018, CEUR-WS Vol. 2072, pp. 55–62). Presented at the Second Software Documentation Generation Challenge (DocGen2, 2020).
- 2,300+ GitHub stars, 460+ forks. Maintained under the [adr GitHub organization](https://github.com/adr), the primary hub for ADR tooling and templates. ADRs as a practice are referenced by the [Microsoft Azure Well-Architected Framework](https://learn.microsoft.com/en-us/azure/well-architected/architect-role/architecture-decision-record), [AWS Prescriptive Guidance](https://docs.aws.amazon.com/prescriptive-guidance/latest/architectural-decision-records/adr-process.html), and IEEE Software (Keeling, Vol. 39, 2022). Azure and AWS both link to the ADR GitHub organization for templates, where MADR is the most-starred project.
- Version 4.0.0 released September 17, 2024. The project captures its own design decisions as 19 self-referential ADRs (0000–0018), covering choices from file naming conventions to metadata format. Follows Semantic Versioning 2.0.0 with a maintained changelog.
- Dual-licensed MIT and CC0-1.0 (SPDX: `MIT OR CC0-1.0`). No licensing friction for any downstream use.
- Ships with markdownlint configuration and a GitHub Actions workflow for CI enforcement. Published on npm for scaffolding. Template variants (full, minimal, bare) let teams match formality to decision weight.

Alternatives evaluated:

1. Adopt MADR 4.0 (selected). Community-validated structure with explicit sections for Decision Drivers, Considered Options, per-option Pros and Cons, and a Confirmation section for validating compliance. Well-documented, externally maintained, and already familiar to contributors who have encountered it in other projects.
2. Keep the current minimal template. Low overhead, but does not enforce capture of alternatives or drivers. Review friction will increase as decision complexity grows (the trajectory from [ADR-0001](0001-use-adrs.md) through [ADR-0006](0006-complypack-content-envelope.md) already shows this).
3. Nygard's original format (the 2011 blog post that started the ADR movement). Captures Status, Context, Decision, and Consequences. Structurally identical to the current template. Does not add the sections we are missing.
4. Design a custom template. Maximum flexibility but creates an undocumented format that new contributors must learn from scratch. No community validation, no linting tooling, no external documentation. Maintenance burden falls entirely on this project.

## :stop_sign: Do I even need an ADR?

If there is only one viable alternative, then an ADR is not required but may be submitted if the information is broadly scoped and may be useful to future maintainers.

## Decision

Adopt MADR 4.0 as the ADR format for all new architectural decision records.

The project ADR template (`docs/ADRs/adr-template.md`) is replaced with a MADR 4.0-based template adapted to project conventions (ADR numbering in the title, `docs/ADRs/` directory location).

Existing accepted ADRs remain as-is. Per [ADR-0001](0001-use-adrs.md), accepted records are immutable. Format consistency across the historical record is not a goal; decision preservation is.

New ADRs must include:

- Context and Problem Statement
- Decision Drivers
- Considered Options
  - Not required for simple decisions
  - Provide as many options as is _practical_
  - Three is generally a good target with one (favored) at the top and two rejected alternatives
- Decision Outcome with a "because" clause referencing a listed driver
- Consequences with positive and negative entries
- `Pros and Cons of the Options` and `More Information` sections are recommended for complex decisions but not required.

The MADR 4.0 template uses YAML frontmatter for metadata (`status`, `date`, `deciders`, `consulted`, `informed`) rather than inline bold fields. This aligns with machine-parseable conventions and enables future tooling (index generation, status filtering) without fragile text parsing.

## Consequences

- Good, because new ADRs capture alternatives and drivers by default, reducing "what about X?" review cycles.
- Good, because the MADR format is externally documented: contributors can reference [adr.github.io/madr](https://adr.github.io/madr/) rather than relying solely on project documentation.
- Good, because template variants (full vs. minimal) allow lightweight decisions to use a shorter form without abandoning the structure.
- Good, because YAML frontmatter enables automated index generation and status tracking without custom parsing.
- Good, because markdownlint configuration from MADR can enforce formatting consistency in CI.
- Neutral, because the project takes a dependency on an external template's conventions. MADR's dual MIT/CC0 license and self-contained Markdown format mean this dependency is effectively zero-cost: no runtime dependency, no build dependency, no license encumbrance. If MADR were abandoned, the template remains usable as-is.
- Bad, because our current `Docsify` implementation strips YAML frontmatter as [prescribed by MADR](https://adr.github.io/madr/decisions/0013-use-yaml-front-matter-for-meta-data.html)
  - The remediation is to incorporate [Docisfy-mustache](https://docsify-mustache.github.io/#/) into our build system
    which may be an overall usability improvement in the long term
- Bad, because writing overhead per ADR increases. The template has more sections to fill. This is the intended tradeoff: the cost is borne once at authoring time; the benefit is realized on every subsequent read.
- Bad, because historical ADRs ([0001](0001-use-adrs.md)–[0006](0006-complypack-content-envelope.md)) use the prior format. Readers encounter two formats when browsing the full record. This is a documentation consistency cost accepted in favor of immutability.

## References

- [MADR 4.0 specification](https://adr.github.io/madr/)
- [MADR GitHub repository](https://github.com/adr/madr)
- [Kopp, Armbruster, Zimmermann, "Markdown Architectural Decision Records: Format and Tool Support" (ZEUS 2018, CEUR-WS Vol. 2072)](https://ceur-ws.org/Vol-2072/paper9.pdf)
- [Nygard, "Documenting Architecture Decisions" (2011)](https://cognitect.com/blog/2011/11/15/documenting-architecture-decisions)
- [ADR-0001: Use Architecture Decision Records](0001-use-adrs.md)
