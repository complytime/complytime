# ADR-0001: Use Architecture Decision Records

**Status:** accepted

**Date:** 2026-06-05

**Deciders:** @jpower432

## Context

ComplyTime's architectural decisions are scattered across OpenSpec specs, CEP proposals, PR discussions, and implementation code. There is no single, browsable record of what was decided and why. 
Contributors discovering the project must reconstruct rationale from multiple sources.

Fullsend and other open-source projects demonstrate that lightweight ADRs — short, immutable records of decisions — provide a navigable history without the overhead of formal enhancement proposals.

## Decision

Record architectural decisions as numbered ADR files in `docs/ADRs/`. Each ADR captures context, the decision, and consequences in one page or less.

ADRs are immutable once accepted. If a decision is reversed or superseded, a new ADR records the change and references the original.

The prior CEP process is retired. Existing CEP content is migrated into ADRs (for decisions) and problem docs or plans (for explorations and implementation details).

## Consequences

- Decisions are discoverable in one location with stable numbering.
- The bar for recording a decision is lower (one page vs. multi-page CEP).
- Problem exploration moves to `docs/problems/` where it can evolve without the formality of a proposal lifecycle.
- Historical decisions from OpenSpec specs and PRs can be retroactively captured.
- ADR acceptance requires two maintainer approvals (see [CONTRIBUTING.md](../../CONTRIBUTING.md#process) for the full approval process).
