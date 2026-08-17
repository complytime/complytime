# ADR-0011: Cross-Framework Mapping Principles

**Status:** proposed

**Date:** 2026-06-26

**Deciders:** ComplyTime maintainers

## Context

Organizations assessed against multiple compliance frameworks face overlapping requirements described in different vocabulary, at different granularity, and with different intent. The [cross-framework-mapping](../problems/cross-framework-mapping.md) problem doc explores why this is hard and surveys prior art (manual spreadsheets, NIST OLIR, IR 8477, OSCAL, SCF, UCF, LLM-assisted analysis).

Five decisions have crystallized from that exploration.

## Decision

1. **Confidence thresholds are configurable at every level.** Defaults per mapping type, overridable per framework, per requirement type, and per organization. The system tracks agreement rates between machine proposals and human decisions over time to calibrate thresholds.

2. **LLM reasoning traces are always captured; human accountability is required for trust.** Reasoning traces are always part of the provenance record. A human must be accountable for the final trust decision on any mapping. Whether that human accepts the LLM trace as sufficient rationale or provides their own reasoning is a policy knob tunable per organization.

3. **Divergent authority mappings are both valid; jurisdiction determines applicability.** Two authorities mapping the same requirements differently is not a conflict, because both are correct within their respective jurisdictions. Only the mapping from the authority with actual authorization over a specific system governs that system's compliance posture.

4. **Version-aware edges with targeted re-evaluation.** Framework versions coexist in the graph. When a new version arrives, re-evaluation targets only mappings referencing the superseded version. Old-version mappings remain valid for systems still assessed against them.

5. **Map as-is; never decompose or compose requirements to fit a mapping.** Requirements are mapped at the granularity the frameworks provide. Altering the original to fit a mapping risks corrupting meaning, violating requirement fidelity. This is the core principle of [Complytime Labs - CrossCodex](https://github.com/complytime-labs/crosscodex).

## Consequences

- The mapping engine must support layered threshold configuration with organizational overrides.
- Every mapping edge in the property graph carries provenance: who proposed, who ratified, confidence score, reasoning trace, framework version pair.
- The graph must support multiple concurrent framework versions with version-scoped mapping edges.
- Asymmetric mappings (one broad to many specific) produce multiple edges with appropriate IR 8477 relationship types rather than synthetic decomposition.
- Community commons mappings compound effort across organizations but require governance for quality.
